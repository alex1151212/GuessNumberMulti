# Build stage
FROM golang:1.20.7-alpine3.18 AS build
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o app

# Final stage
FROM alpine:3.18
WORKDIR /app
COPY --from=build /app/app .
CMD ["./app"]