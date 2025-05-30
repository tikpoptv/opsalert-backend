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

#### Set Staff Permissions (Admin Only)
```
POST /staff/permissions
```
กำหนดสิทธิ์การเข้าถึง OA ให้กับ staff

**Headers**
```
Authorization: Bearer <token>
```

**Request Body**
```json
{
    "staff_id": "integer",
    "permissions": [
        {
            "oa_id": "integer",
            "permission_level": "string" // "view" หรือ "manage"
        }
    ]
}
```

**Response**
```json
{
    "message": "staff permissions updated successfully"
}
```

#### Get Staff Permissions
```
GET /staff/permissions/:staff_id
```
ดูรายการ OA ที่ staff มีสิทธิ์เข้าถึง

**Headers**
```
Authorization: Bearer <token>
```

**Response**
```json
{
    "data": [
        {
            "oa_id": "integer",
            "oa_name": "string",
            "permission_level": "string" // "view" หรือ "manage"
        }
    ]
}
```

**Error Responses**
```json
{
    "error": "staff not found"
}
```
หรือ
```json
{
    "error": "insufficient permissions"
}
```

#### Delete Staff Permission (Admin Only)
```
DELETE /staff/permissions/:id?oa_id={oa_id}
```
ลบสิทธิ์การเข้าถึง OA ที่ระบุของ staff

**Headers**
```
Authorization: Bearer <token>
```

**Query Parameters**
- oa_id: ID ของ OA ที่ต้องการลบสิทธิ์ (required)

**Response**
```json
{
    "message": "staff permission deleted successfully"
}
```

**Error Responses**
```json
{
    "error": "invalid staff id"
}
```
หรือ
```json
{
    "error": "oa_id is required"
}
```
หรือ
```json
{
    "error": "invalid oa_id"
}
```
หรือ
```json
{
    "error": "staff not found"
}
```
หรือ
```json
{
    "error": "OA not found"
}
```
หรือ
```json
{
    "error": "staff does not have permission for this OA"
}
```
หรือ
```json
{
    "error": "cannot delete permissions for admin"
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

#### Update OA (Admin or Staff with Manage Permission)
```
PUT /oa/:id
```
อัพเดทข้อมูล LINE Official Account
- Admin สามารถแก้ไขได้ทุก OA
- Staff ต้องมีสิทธิ์ manage ถึงจะแก้ไขได้

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

**Error Responses**
```json
{
    "error": "insufficient permissions to update this OA"
}
```
หรือ
```json
{
    "error": "line official account not found"
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
ดูรายการ LINE Official Account

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

**Note:**
- ถ้าเป็น admin จะเห็น OA ทั้งหมด
- ถ้าเป็น staff จะเห็นเฉพาะ OA ที่ได้รับสิทธิ์เข้าถึงเท่านั้น

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

### ดูสิทธิ์ของ Staff

**Endpoint:** `GET /api/v1/staff/permissions/:staff_id`

**Method:** GET

**Headers:**
- Authorization: Bearer {token} (ต้องเป็น admin)

**Response:**
```json
{
    "data": [
        {
            "oa_id": 1,
            "oa_name": "OA Name",
            "permission_level": "manage"
        }
    ]
}
```

### ลบสิทธิ์ของ Staff

**Endpoint:** `DELETE /api/v1/staff/permissions/:id`

**Method:** DELETE

**Headers:**
- Authorization: Bearer {token} (ต้องเป็น admin)

**Response:**
```json
{
    "message": "staff permissions deleted successfully"
}
```

**Error Responses:**
- 400 Bad Request: `{"error": "invalid staff id"}` หรือ `{"error": "cannot delete permissions for admin"}`
- 404 Not Found: `{"error": "staff not found"}`
- 500 Internal Server Error: `{"error": "internal server error"}` 