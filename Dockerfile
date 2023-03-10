FROM golang:1.19.5-alpine3.17 AS builder
WORKDIR /app
COPY . .
RUN go build -o secret-keeper main.go

FROM alpine:3.17.1
WORKDIR /app
COPY app.env .
COPY --from=builder /app/secret-keeper .

EXPOSE 8000
CMD ["/app/secret-keeper"]