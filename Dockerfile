# ---- Build stage ----
FROM golang:1.24.4-alpine AS builder
WORKDIR /src

# copy mod trước để tận dụng cache
COPY go.mod go.sum ./
RUN go mod download

# copy toàn bộ source
COPY . .

# build binary
RUN go build -o /go-app .

# Xóa tất cả file *.env để đảm bảo không còn trong image
# RUN rm -f /src/.env

# ---- Final stage ----
FROM alpine:latest

# copy binary từ builder
COPY --from=builder /go-app .

# # expose port Gin dùng
EXPOSE 8080

# # chạy binary
ENTRYPOINT ["./go-app"]
# CMD ["sleep", "infinity"]


