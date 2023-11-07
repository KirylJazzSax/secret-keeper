FROM golang:1.20.5-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o /secret-keeper

ENTRYPOINT [ "/secret-keeper" ]
