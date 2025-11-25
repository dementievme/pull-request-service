FROM golang AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server ./cmd/server/main.go

FROM scratch
WORKDIR /app
COPY --from=builder /app/server .
COPY migrations ./migrations

EXPOSE 8080
CMD ["./server"]
