FROM golang:1.19-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN go build -o /bank-api ./cmd/web/main.go


FROM alpine:3.18.2 AS final
WORKDIR /app
COPY .env .
COPY db/migration ./db/migration

COPY --from=builder /bank-api .

COPY start.sh .
COPY wait-for.sh .

EXPOSE 8080
CMD [ "./bank-api" ]
ENTRYPOINT [ "/app/start.sh" ]