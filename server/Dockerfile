FROM golang:alpine AS builder
WORKDIR /app
COPY go.* .
COPY server ./server
COPY persistence ./persistence
COPY share ./share
COPY game ./game

RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o main server/main.go

FROM scratch
WORKDIR /app
COPY --from=builder /app/main .
COPY database .
CMD ["./main"]
