FROM golang

LABEL maintainer="George Abumere <abumere@gmail.com>"

ENV GO111MODULE=on

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build

EXPOSE 8080
ENTRYPOINT ["/app/detector"]


