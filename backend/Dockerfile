FROM golang:1.22-alpine AS builder
WORKDIR /src
COPY go.mod .
COPY ./cmd ./cmd
COPY ./internal ./internal
COPY ./pkg ./pkg
RUN apk add --no-cache git
RUN go build -o /bin/api ./cmd/api

FROM alpine:3.18
RUN apk add --no-cache ca-certificates
COPY --from=builder /bin/api /bin/api
EXPOSE 8080
ENTRYPOINT ["/bin/api"]
