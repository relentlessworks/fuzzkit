FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o fuzzkit ./cmd/fuzzkit

FROM scratch
COPY --from=builder /app/fuzzkit /fuzzkit
EXPOSE 8427
ENTRYPOINT ["/fuzzkit"]
