FROM golang:1.17.4-alpine3.15 as builder

WORKDIR /app
COPY . .
RUN go mod tidy

ENV GOARCH=amd64

RUN go build \
    -ldflags "-X main.buildcommit=`git rev-parse --short HEAD` \
    -X main.buildtime=`date "+%Y-%m-%dT%H:%M:%S%Z:00"`" \
    -o /go/bin/app

## Deploy
FROM alpine:3.15  
WORKDIR /

COPY --from=builder /go/bin/app /app
COPY ./config ./config
EXPOSE 8080

CMD ["/app"] 