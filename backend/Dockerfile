# Build the Go binary (the Vue dist/ and config.json are embedded at compile time).
FROM golang:1.26-alpine AS build
WORKDIR /src
COPY go.mod ./
COPY . .
RUN CGO_ENABLED=0 go build -o /server .

# Minimal runtime image.
FROM alpine:latest
WORKDIR /app
COPY --from=build /server /app/server
# config.json on disk lets content be edited without a rebuild; falls back to the
# embedded copy if absent.
COPY config.json ./
# Render injects $PORT at runtime; main.go reads it (defaults to 8080 locally).
CMD ["/app/server"]
