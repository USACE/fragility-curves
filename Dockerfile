FROM osgeo/gdal:alpine-normal-3.2.1 as dev

COPY --from=golang:1.18-alpine3.14 /usr/local/go/ /usr/local/go/

RUN apk add --no-cache \
    pkgconfig \
    gcc \
    libc-dev \
    git

ENV GOROOT=/usr/local/go
ENV GOPATH=/go
ENV GO111MODULE="on"
ENV PATH=$PATH:$GOROOT/bin:$GOPATH/bin

# Hot-Reloader
RUN go install github.com/githubnemo/CompileDaemon@v1.4.0

COPY ./ /app
WORKDIR /app

RUN go mod download

RUN go build main.go
ENTRYPOINT /go/bin/CompileDaemon --build="go build main.go" --command="./main -payload=payload.yml"

# TODO: add prod build
# FROM osgeo/gdal:alpine-normal-3.2.1 as prod
# Production container
FROM golang:1.18-alpine3.14 AS prod
WORKDIR /app
COPY --from=dev /app/main .
CMD [ "./main" ]