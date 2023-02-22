FROM golang:latest

ENV GOPATH=/

COPY ./ ./

RUN go mod download
RUN go mod tidy