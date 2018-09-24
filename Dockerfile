FROM golang:alpine as builder
WORKDIR /go/src/github.com/SimePel/travis-builds-status
COPY . .
RUN go build main.go

FROM alpine:latest
COPY --from=builder /go/src/github.com/SimePel/travis-builds-status .
ENTRYPOINT [ "./main" ]
