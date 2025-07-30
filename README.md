# 🚀 Dự án API với Go Gin, GORM, PostgreSQL và Redis

Dự án này là một boilerplate RESTful API được xây dựng bằng ngôn ngữ Go, sử dụng các công nghệ phổ biến:

- [Gin](https://github.com/gin-gonic/gin): Web framework hiệu năng cao
- [GORM](https://gorm.io/): ORM (Object Relational Mapping) cho Go
- [PostgreSQL](https://www.postgresql.org/): Hệ quản trị cơ sở dữ liệu quan hệ mã nguồn mở
- [Redis](https://redis.io/): Hệ thống lưu trữ key-value trong bộ nhớ (cache/session)

---

## 📦 Thư viện Go đã sử dụng

| Thư viện                              | Chức năng                           |
| ------------------------------------- | ----------------------------------- |
| `github.com/gin-gonic/gin`            | Framework HTTP                      |
| `gorm.io/gorm`                        | ORM                                 |
| `gorm.io/driver/postgres`             | Driver PostgreSQL cho GORM          |
| `github.com/joho/godotenv`            | Load biến môi trường từ file `.env` |
| `github.com/go-redis/redis/v9`        | Redis client cho Go                 |
| `github.com/google/uuid`              | Sinh UUID                           |
| `gopkg.in/go-playground/validator.v9` | Validation                          |

---

## 🐳 Khởi động nhanh bằng Docker Compose

Dự án đã cấu hình sẵn Docker Compose để bạn có thể khởi động nhanh môi trường PostgreSQL và Redis.

### ✅ Yêu cầu:

- Cài đặt sẵn [Docker](https://www.docker.com/)

**Sao chép file `local.env` thành `.env`** (nếu chưa có):

```bash
cp local.env .env
```

### ▶️ Chạy lệnh sau để khởi động môi trường:

```bash
docker-compose up --build
```


---

## 📝 Viết file `.http` để test API

Sau khi hoàn thành mỗi API, hãy tạo một file `.http` trong thư mục dự án để mô tả và kiểm thử các endpoint.

**Ví dụ về file `user.http`:**

```http
### Đăng ký tài khoản mới
POST http://localhost:8080/api/v1/users/register
Content-Type: application/json

{
    "username": "testuser",
    "password": "123456"
}

### Đăng nhập
POST http://localhost:8080/api/v1/users/login
Content-Type: application/json

{
    "username": "testuser",
    "password": "123456"
}
```

**Lưu ý:**  
- Ghi chú rõ ràng chức năng của từng API ngay trong file `.http`.
- Có thể tạo nhiều file `.http` cho từng nhóm chức năng (user, product, ...).
- Sử dụng các biến môi trường nếu cần để dễ dàng chuyển đổi giữa các môi trường (local, staging, production).