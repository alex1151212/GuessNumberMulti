# Build stage
FROM golang:1.17.0-alpine3.14 AS build
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o app

# Final stage
FROM alpine:3.14
WORKDIR /app
COPY --from=build /app/app .
CMD ["./app"]