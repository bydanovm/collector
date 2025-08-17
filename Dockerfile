FROM golang:alpine AS builder
RUN apk add --no-cache git
WORKDIR /go/src/app
COPY . .
RUN go get -d -v ./...
RUN go build -o /go/bin/app -v ./cmd/app/main.go

FROM alpine:latest
VOLUME [ "/var/logs" ]
RUN apk --no-cache add ca-certificates tzdata
COPY --from=builder /go/bin/app /app
ENTRYPOINT /app