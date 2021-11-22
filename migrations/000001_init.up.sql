CREATE SCHEMA IF NOT EXISTS nix;

CREATE TYPE nix.user_to_chat_relations AS ENUM ('userBlockedByChat','chatBlockedByUser', 'joined', 'visited', 'owner','userToUser');
CREATE TYPE nix.chat_types AS ENUM ('private','userToUser','public');

CREATE TYPE nix.user_relations_types AS ENUM ('friend', 'blocked');



CREATE TABLE IF NOT EXISTS nix.users
(
    id            uuid         NOT NULL unique PRIMARY KEY DEFAULT gen_random_uuid(),
    username      varchar(255) NOT NULL unique,
    email         varchar(255) NOT NULL unique
        CONSTRAINT proper_email CHECK (email ~* '^[A-Za-z0-9._%-]+@[A-Za-z0-9.-]+[.][A-Za-z]+$'),
    password_hash varchar(255),
    avatar_url    varchar      NOT NULL                    DEFAULT ''
);


CREATE TABLE IF NOT EXISTS nix.chats
(
    id          uuid           NOT NULL UNIQUE PRIMARY KEY DEFAULT gen_random_uuid(),
    name        varchar(255)   NOT NULL,
    title       varchar(255)   NOT NULL,
    description varchar(1000)  NOT NULL,
    type        nix.chat_types NOT NULL

);


CREATE TABLE IF NOT EXISTS nix.messages
(
    id        uuid                     NOT NULL unique PRIMARY KEY DEFAULT gen_random_uuid(),
    userid    uuid                     NOT NULL,
    chatid    uuid                     NOT NULL,
    text      varchar(5000)            NOT NULL,
    timestamp timestamp with time zone NOT NULL,
    FOREIGN KEY (userid) REFERENCES nix.users (id) ON DELETE NO ACTION,
    FOREIGN KEY (chatid) REFERENCES nix.chats (id) ON DELETE NO ACTION
);

CREATE INDEX IF NOT EXISTS massage_chat_index ON nix.messages (chatid);



CREATE TABLE IF NOT EXISTS nix.user_chats
(
    id       uuid NOT NULL UNIQUE PRIMARY KEY DEFAULT gen_random_uuid(),
    userid   uuid NOT NULL,
    chatid   uuid NOT NULL,
    relation nix.user_to_chat_relations,

    FOREIGN KEY (userid) REFERENCES nix.users (id) ON DELETE NO ACTION,
    FOREIGN KEY (chatid) REFERENCES nix.chats (id) ON DELETE NO ACTION
);

CREATE INDEX IF NOT EXISTS user_chats_userid_index ON nix.user_chats (userid);
CREATE INDEX IF NOT EXISTS user_chats_chatid_index ON nix.user_chats (chatid);


CREATE TABLE IF NOT EXISTS nix.users_relations
(
    id        uuid                     NOT NULL UNIQUE PRIMARY KEY DEFAULT gen_random_uuid(),
    userid    uuid                     NOT NULL,
    to_userid uuid                     NOT NULL,
    type      nix.user_relations_types NOT NULL
);

CREATE INDEX IF NOT EXISTS users_relations_user_to_user ON nix.users_relations (userid, to_userid);
