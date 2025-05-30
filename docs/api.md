# OpsAlert API Documentation

## Base URL
```
http://localhost:8080/api/v1
```

## Authentication
API ใช้ JWT Token สำหรับการยืนยันตัวตน โดยต้องส่ง token ในรูปแบบ:
```
Authorization: Bearer <token>
```

## Endpoints

### Health Check
```
GET /health
```
ตรวจสอบสถานะการทำงานของระบบ

**Response**
```json
{
    "status": "ok"
}
```

### Staff Management

#### Login
```
POST /staff/login
```
เข้าสู่ระบบ

**Request Body**
```json
{
    "username": "string",
    "password": "string"
}
```

**Response**
```json
{
    "token": "string",
    "user": {
        "id": "integer",
        "username": "string",
        "role": "string"
    }
}
```

#### Get Profile
```
GET /staff/me
```
ดูข้อมูลโปรไฟล์ของตัวเอง

**Headers**
```
Authorization: Bearer <token>
```

**Response**
```json
{
    "id": "integer",
    "username": "string",
    "role": "string"
}
```

#### Get All Staff Accounts (Admin Only)
```
GET /staff/accounts
```
ดูรายการบัญชีพนักงานทั้งหมด

**Headers**
```
Authorization: Bearer <token>
```

**Response**
```json
{
    "data": [
        {
            "id": "integer",
            "username": "string",
            "role": "string"
        }
    ]
}
```

#### Get Staff Account by ID (Admin Only)
```
GET /staff/accounts/:id
```
ดูข้อมูลบัญชีพนักงานตาม ID

**Headers**
```
Authorization: Bearer <token>
```

**Response**
```json
{
    "id": "integer",
    "username": "string",
    "role": "string"
}
```

#### Update Staff (Admin Only)
```
PUT /staff/accounts/:id
```
อัพเดทข้อมูลพนักงาน

**Headers**
```
Authorization: Bearer <token>
```

**Request Body**
```json
{
    "username": "string",
    "password": "string",
    "role": "string"
}
```

**Response**
```json
{
    "message": "staff updated successfully"
}
```

#### Register New Staff (Admin Only)
```
POST /staff/register
```
สร้างบัญชีพนักงานใหม่

**Headers**
```
Authorization: Bearer <token>
```

**Request Body**
```json
{
    "username": "string",
    "password": "string",
    "role": "string"
}
```

**Response**
```json
{
    "message": "staff registered successfully"
}
```

### LINE Official Account Management

#### Create OA (Admin Only)
```
POST /oa
```
สร้าง LINE Official Account ใหม่

**Headers**
```
Authorization: Bearer <token>
```

**Request Body**
```json
{
    "name": "string",
    "channel_id": "string",
    "channel_secret": "string",
    "channel_access_token": "string"
}
```

**Response**
```json
{
    "message": "line official account created successfully"
}
```

#### Update OA (Admin Only)
```
PUT /oa/:id
```
อัพเดทข้อมูล LINE Official Account

**Headers**
```
Authorization: Bearer <token>
```

**Request Body**
```json
{
    "name": "string",
    "channel_id": "string",
    "channel_secret": "string",
    "channel_access_token": "string"
}
```

**Response**
```json
{
    "message": "line official account updated successfully"
}
```

#### Delete OA (Admin Only)
```
DELETE /oa/:id
```
ลบ LINE Official Account

**Headers**
```
Authorization: Bearer <token>
```

**Response**
```json
{
    "message": "line official account deleted successfully"
}
```

#### List All OAs
```
GET /oa
```
ดูรายการ LINE Official Account ทั้งหมด

**Headers**
```
Authorization: Bearer <token>
```

**Response**
```json
{
    "data": [
        {
            "id": "integer",
            "name": "string",
            "channel_id": "string",
            "channel_secret": "string",
            "channel_access_token": "string",
            "webhook_url": "string",
            "created_at": "datetime"
        }
    ]
}
```

## Error Responses

### 400 Bad Request
```json
{
    "error": "invalid request data"
}
```

### 401 Unauthorized
```json
{
    "error": "unauthorized"
}
```

### 403 Forbidden
```json
{
    "error": "forbidden"
}
```

### 404 Not Found
```json
{
    "error": "not found"
}
```

### 500 Internal Server Error
```json
{
    "error": "internal server error"
}
``` 