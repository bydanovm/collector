FROM golang:alpine AS builder
RUN apk add --no-cache git
WORKDIR /go/src/miniapp
COPY . .
RUN go get -d -v ./...
RUN go build -o /go/bin/miniapp -v ./cmd/miniapp/main.go

FROM alpine:latest
VOLUME [ "/var/logs" ]
RUN apk --no-cache add ca-certificates tzdata
COPY --from=builder /go/bin/miniapp /miniapp
ENTRYPOINT /miniapp
EXPOSE 8080