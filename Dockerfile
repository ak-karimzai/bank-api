FROM golang:1.19-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN go build -o /bank-api ./cmd/web/main.go

RUN apk add --no-cache curl
RUN curl -L https://github.com/pressly/goose/releases/download/v3.13.1/goose_linux_x86_64 --output goose
RUN chmod +x goose

FROM alpine:3.18.2 AS final
WORKDIR /app
COPY .env .
COPY db/migration ./migration

COPY --from=builder /app/goose .
COPY --from=builder /bank-api .

COPY start.sh .
COPY wait-for.sh .

EXPOSE 8080
CMD [ "./bank-api" ]
ENTRYPOINT [ "/app/start.sh" ]