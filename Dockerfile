FROM golang:1.25.5 AS builder
WORKDIR /src
COPY . .
RUN go build -o /bin/app ./cmd/app

FROM golang:1.25.5
WORKDIR /app
COPY --from=builder /bin/app .
CMD ["./app"]
