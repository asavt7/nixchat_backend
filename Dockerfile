FROM golang:1.16 as builder
ADD . /

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -v -o /bin/main /cmd/chat/main.go

FROM scratch as app
COPY --from=builder /bin/* /
COPY --from=builder /configs /configs

CMD ["/main"]