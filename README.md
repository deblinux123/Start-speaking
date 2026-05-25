# Start-speaking

# 🚀 StartSpeek Backend API

یک بک‌اند با زبان Go (Gin) برای مدیریت احراز هویت کاربران، JWT، Redis و Docker.  
این پروژه برای اتصال به اپلیکیشن موبایل (Android Kotlin) طراحی شده است.

---

# 🧱 تکنولوژی‌ها

- Go (Gin Framework)
- JWT Authentication
- Redis (مدیریت سشن و logout)
- SQLite / GORM
- Docker + Docker Compose

---

# ✨ قابلیت‌ها

- ثبت‌نام کاربر (Register)
- ورود کاربر (Login)
- تولید Access Token (کوتاه‌مدت)
- تولید Refresh Token (بلندمدت)
- خروج امن (Logout)
- مدیریت سشن با Redis
- مسیرهای محافظت‌شده (/api/me)
- Refresh Token API
- کاملاً قابل اجرا با Docker

---

# 📁 ساختار پروژه

```

.
├── main.go
├── controllers/
├── middleware/
├── utils/
├── db/
├── routes/
├── Dockerfile
├── docker-compose.yml
├── go.mod
├── go.sum
└── README.md

````

---

# 🚀 اجرای پروژه

## 1. کلون کردن پروژه

```bash
git clone <repo-url>
cd startSpeek
````

---

## 2. اجرای پروژه با Docker (پیشنهادی)

```bash
sudo docker-compose up --build
```

---

## 3. آدرس API

```
http://localhost:8080
```

---

# 🔐 سیستم احراز هویت

## 1. ثبت‌نام

```
POST /register
```

---

## 2. ورود

```
POST /login
```

### پاسخ:

```json
{
  "access_token": "...",
  "refresh_token": "...",
  "user": {
    "id": 1,
    "name": "Ali",
    "email": "ali@gmail.com"
  }
}
```

---

## 3. گرفتن اطلاعات کاربر

```
GET /api/me
Authorization: Bearer <access_token>
```

---

## 4. Refresh Token

```
POST /refresh
```

```json
{
  "refresh_token": "..."
}
```

---

## 5. خروج از حساب

```
POST /api/logout
Authorization: Bearer <access_token>
```

---

# 🧠 سیستم توکن‌ها

* Access Token → کوتاه‌مدت (مثلاً 15 دقیقه)
* Refresh Token → بلندمدت (مثلاً 7 روز)
* Redis برای:

  * مدیریت سشن
  * جلوگیری از استفاده بعد از logout

---

# 🐳 Docker

این پروژه شامل دو سرویس است:

* app → Go API روی پورت 8080
* redis → Redis server روی پورت 6379

---

## اجرای کامل پروژه

```bash
sudo docker-compose up --build
```

---

## توقف پروژه

```bash
sudo docker-compose down
```

---

# 📡 اتصال به موبایل (Android)

* لوکال:

```
http://localhost:8080
```

* با ngrok:

```
https://xxxx.ngrok-free.app
```

---

# 🚀 وضعیت پروژه

✔ آماده اتصال به موبایل
✔ آماده توسعه AI
✔ دارای سیستم auth کامل
✔ Dockerized و قابل deploy

---

# 🔥 مرحله بعد

* اضافه کردن AI Chat (Python + Gemma)
* حافظه چت با Redis
* تبدیل متن به گفتار (TTS)
* سیستم یادگیری زبان انگلیسی

```

---
