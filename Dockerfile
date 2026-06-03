FROM golang:latest AS builder
WORKDIR /src
COPY . .
RUN go build -o /bin/app ./cmd/app

FROM golang:latest
WORKDIR /app
COPY --from=builder /bin/app .
CMD ["./app"]
