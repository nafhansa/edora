FROM golang:1.25-alpine AS builder

WORKDIR /src

# Install git dulu (penting buat download modul)
RUN apk add --no-cache git

# Copy kedua file ini DULUAN agar caching Docker bekerja optimal
COPY go.mod go.sum ./

# Download dependencies di layer terpisah
RUN go mod download

# Baru copy source code sisanya
COPY ./cmd ./cmd
COPY ./internal ./internal
COPY ./pkg ./pkg

# Build binary
RUN go build -o /bin/api ./cmd/api

# Stage 2: Image Kecil untuk Production
FROM alpine:3.18
RUN apk add --no-cache ca-certificates
COPY --from=builder /bin/api /bin/api
EXPOSE 8080
ENTRYPOINT ["/bin/api"]