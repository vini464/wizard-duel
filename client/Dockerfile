FROM golang:alpine AS builder
WORKDIR /app
COPY go.mod  .
COPY client  ./client
COPY share  ./share
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o main client/main.go

FROM scratch
WORKDIR /app
COPY --from=builder /app/main .
CMD ["./main"]
