FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/server ./cmd/api

FROM gcr.io/distroless/static-debian12

COPY --from=builder /bin/server /server

EXPOSE 8080

ENTRYPOINT ["/server"]
