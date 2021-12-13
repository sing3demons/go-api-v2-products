FROM golang:1.17.4-alpine3.15 as builder

RUN apk --no-cache add gcc g++ make
RUN apk add git
WORKDIR /app
COPY . .
RUN go mod tidy

RUN go build -a -installsuffix cgo \
    -ldflags "-X main.buildcommit=`git rev-parse --short HEAD` \
    -X main.buildtime=`date "+%Y-%m-%dT%H:%M:%S%Z:00"`" \
    -o /go/bin/app

FROM alpine:latest  
RUN apk --no-cache add ca-certificates
WORKDIR /

COPY --from=builder /go/bin/app /app
COPY ./config ./config
EXPOSE 8080


CMD ["/app"] 