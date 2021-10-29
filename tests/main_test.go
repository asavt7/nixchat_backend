package tests

import (
	"fmt"
	"github.com/asavt7/nixchat_backend/internal/app"
	"github.com/asavt7/nixchat_backend/internal/config"
	"github.com/asavt7/nixchat_backend/internal/repos"
	"github.com/go-redis/redis/v8"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	"net"
	"runtime"
	"testing"
	"time"
)

const (
	containersExpirationTimeMs = 60
	postgresContainerName      = "postgres"
	redisContainerName         = "redis"
	pathToMigrations           = "../migrations/"
)

type authCredentials struct {
	accessToken string
}

type MainTestSuite struct {
	suite.Suite

	cfg *config.Config

	pool        *dockertest.Pool
	pgDB        *sqlx.DB
	redisClient *redis.Client

	chatApp *app.ChatApp

	credentials authCredentials
}

func TestMainTestSuite(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	suite.Run(t, new(MainTestSuite))
}

func (m *MainTestSuite) SetupSuite() {
	log.Debug("SetupSuite")

	m.initConfigs()

	pool, err := dockertest.NewPool("")
	if err != nil {
		log.WithError(err).Fatal("Could not connect to docker")
	}
	pool.MaxWait = 30 * time.Second
	m.pool = pool

	m.initPostgresDbContainer()
	m.initRedisContainer()

	m.initApp()

	// todo https://github.com/asavt7/nixchat_backend/issues/11  wait while app initialize -- use readiness probe
	time.Sleep(400 * time.Millisecond)

	m.registerUser()
	m.signIn()
}

func (m *MainTestSuite) TearDownSuite() {
	log.Info("TearDownSuite")

	if err := m.pool.RemoveContainerByName(postgresContainerName); err != nil {
		log.Warning(err)
	}
	if err := m.pool.RemoveContainerByName(redisContainerName); err != nil {
		log.Warning(err)
	}
}

func (m *MainTestSuite) initConfigs() {
	m.cfg = &config.Config{
		Logger: config.LoggerConfig{
			Level: "info",
		},
		Auth: config.AuthConfig{
			AccessTokenTTL:  120 * time.Minute,
			RefreshTokenTTL: 120 * time.Minute,
			AccessSecret:    "some",
			RefreshSecret:   "secret",
			AutoLogoffTime:  120 * time.Minute,
		},
		Postgres: config.PostgresConfig{
			Host:     "CHANGE_ME",
			Port:     "5432",
			Username: "postgres",
			Password: "postgres",
			DBName:   "postgres",
			SSLMode:  "disable",
		},
		HTTP: config.HTTPConfig{
			Host:         "localhost",
			Port:         "8080",
			ReadTimeout:  2 * time.Second,
			WriteTimeout: 2 * time.Second,
		},
		Redis: config.RedisConfig{
			Host: "CHANGE_ME",
			Port: "6379",
		},
	}
}

func (m *MainTestSuite) initRedisContainer() {

	runOpts := dockertest.RunOptions{
		Name:       redisContainerName,
		Repository: "redis",
		Tag:        "alpine",
	}
	resource, err := m.pool.RunWithOptions(&runOpts, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{
			Name: "no",
		}
	})
	if err != nil {
		log.Fatalf("Could start postgres container %s", err)
	}
	if err := resource.Expire(containersExpirationTimeMs); err != nil {
		log.Fatal(err)
	}

	m.cfg.Redis.Host = getContainerHost(resource)

	// try to connect
	err = m.pool.Retry(func() error {
		redisClient := repos.InitRedisClient(repos.InitRedisOpts(&m.cfg.Redis))

		m.redisClient = redisClient

		return nil
	})
	if err != nil {
		log.Fatalf("cannot connect to DB %s", err)
	}

}

func (m *MainTestSuite) initPostgresDbContainer() {

	postgresRunOpts := dockertest.RunOptions{
		Name:       postgresContainerName,
		Repository: "postgres",
		Tag:        "latest",
		Env: []string{
			"POSTGRES_USER=" + m.cfg.Postgres.Username,
			"POSTGRES_PASSWORD=" + m.cfg.Postgres.Password,
			"POSTGRES_DB=" + m.cfg.Postgres.DBName,
		},
	}
	resource, err := m.pool.RunWithOptions(&postgresRunOpts, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{
			Name: "no",
		}
	})
	if err != nil {
		log.Fatalf("Could start postgres container %s", err)
	}
	if err := resource.Expire(containersExpirationTimeMs); err != nil {
		log.Fatal(err)
	}

	m.cfg.Postgres.Host = getContainerHost(resource)

	// try to connect
	err = m.pool.Retry(func() error {
		pgDB, err := sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s", m.cfg.Postgres.Host, m.cfg.Postgres.Port, m.cfg.Postgres.Username, m.cfg.Postgres.DBName, m.cfg.Postgres.Password, m.cfg.Postgres.SSLMode))
		if err != nil {
			return err
		}

		err = pgDB.Ping()
		if err != nil {
			return err
		}
		//init db
		m.pgDB = pgDB
		return nil
	})
	if err != nil {
		log.Fatalf("cannot connect to DB %s", err)
	}

	//up migrations
	migrator, err := migrate.New(
		"file://"+pathToMigrations,
		fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
			m.cfg.Postgres.Username,
			m.cfg.Postgres.Password,
			m.cfg.Postgres.Host,
			m.cfg.Postgres.Port,
			m.cfg.Postgres.DBName,
			m.cfg.Postgres.SSLMode),
	)
	if err != nil {
		log.Errorf("error happened when migration %v", err)
	}
	if err := migrator.Up(); err != nil && err != migrate.ErrNoChange {
		log.Errorf("error when migration up: %v", err)
	}
	log.Println("migration completed!")
}

func getContainerHost(resource *dockertest.Resource) string {
	// Docker layer network is different on Mac
	if runtime.GOOS == "darwin" {
		return net.JoinHostPort(resource.GetBoundIP("5432/tcp"), resource.GetPort("5432/tcp"))
	}
	return resource.Container.NetworkSettings.IPAddress

}

func (m *MainTestSuite) initApp() {
	chatApp := app.NewChatApp(m.cfg)
	go chatApp.Run()
}
