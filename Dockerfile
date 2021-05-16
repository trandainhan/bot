FROM golang:1.16-alpine
LABEL maintainer "Nhan <nhan@fiahub.com>"

RUN apk update
RUN apk add git

ENV PROJECT="/go/src/gitlab.com/fiahub/bot"

WORKDIR ${PROJECT}
COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN go install ./...
