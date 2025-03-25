# Gunakan image base Golang versi terbaru
FROM golang:1.24.1 AS builder

# Set direktori kerja di dalam container
WORKDIR /app

# Copy semua file dari project ke dalam container
COPY . .

# Download dependensi dan build aplikasi (statically linked binary)
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o myapp main.go

# Buat stage baru dengan image lebih ringan
FROM alpine:latest

# Set direktori kerja di container
WORKDIR /root/

# Install dependencies
RUN apk add --no-cache ca-certificates

# Copy file hasil build dari stage sebelumnya dan beri izin eksekusi
COPY --from=builder /app/myapp .
RUN chmod +x myapp

# Tentukan port yang digunakan
EXPOSE 8080

# Jalankan aplikasi saat container start
CMD ["./myapp"]
