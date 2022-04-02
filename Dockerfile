FROM golang:latest AS build
RUN apt-get update
RUN apt-get install -y git bash
WORKDIR /testdir
COPY . .
RUN go mod tidy

ENTRYPOINT ["go", "test", "-v", "./...", "-coverprofile", "cover.out"]