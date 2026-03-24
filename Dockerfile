# Stage 1: Builder
FROM golang:1.25-alpine AS builder
WORKDIR /app
# Copy go.mod dan go.sum dulu agar caching dependency efisien
COPY go.mod go.sum ./
RUN go mod download
# Copy semua kode source
COPY . .
# Build aplikasi menjadi binary bernama 'main'
RUN go build -o main .

# Stage 2: Runner (Hanya mengambil binary-nya saja agar ringan)
FROM alpine:latest
# Mengatur zona waktu (Sangat penting untuk transaksi & jadwal salat!)
RUN apk add --no-cache tzdata

WORKDIR /app
# Copy binary dari stage builder
COPY --from=builder /app/main .

# Copy folder docs (jika kamu pakai Swagger/Redoc)
# Hapus baris ini jika eTAKMIR tidak punya folder docs
COPY --from=builder /app/docs ./docs 

EXPOSE 3000

CMD ["./main"]