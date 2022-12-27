FROM golang:1.19.2-alpine3.15 AS dev

RUN apk add --no-cache \
    build-base \
    gcc \
    git
    