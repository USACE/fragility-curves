FROM golang:1.19-alpine3.16 AS builder
RUN apk update && apk add --no-cache git && apk add --no-cache ca-certificates && update-ca-certificates
WORKDIR /go/src/app
COPY . .

# Install Dependencies
RUN go get -d -v &&\ 
    GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-w -s" -o /go/bin/main

# SCRATCH IMAGE
FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/bin/main /go/bin/fragilitycurve
ENTRYPOINT ["/go/bin/fragilitycurve"]