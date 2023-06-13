# build stage
FROM golang:1.18 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

# final stage
FROM debian:11
WORKDIR /app
COPY --from=builder /app/start.sh /app
COPY --from=builder /app/main /app
COPY --from=builder /go/bin/goose /app
COPY --from=builder /app/dbScripts/migration /app
COPY --from=builder /app/app.env /app

# run stage
EXPOSE 8080
CMD ["/app/main"]
ENTRYPOINT [ "/app/start.sh" ]