# 📡 Line Notification Gateway

A centralized system for sending messages through LINE Official Accounts (OA), supporting multiple OAs and external systems.  
Includes admin login system, access control, audit logging, and message dashboard.

---

## 📦 Key Features

- ✅ Support for multiple LINE OA accounts
- ✅ External systems can send messages via API
- ✅ Permission control for each system's OA access
- ✅ Dashboard for message history and usage logs
- ✅ Staff login system (Admin / Staff)
- ✅ Comprehensive audit logging
- ✅ PostgreSQL integration with Triggers and CHECK constraints

---

## 🧱 Database Structure

| Table Name             | Description                                   |
|------------------------|-----------------------------------------------|
| `line_official_accounts` | Stores OA information (channel_id, token, etc.) |
| `external_systems`     | External systems with API access (api_key)    |
| `system_oa_permissions`| Defines which system can access which OA      |
| `line_users`           | LINE users (separated by OA)                  |
| `messages`             | Message sending history                       |
| `staff_accounts`       | Staff accounts for Dashboard access           |
| `staff_oa_permissions` | Admin permissions for each OA                 |
| `api_logs`             | Logs of all API calls from external systems   |
| `audit_logs`           | Records all data changes (Insert/Update/Delete) |

---

## 🚀 Installation (Database)

1. Install PostgreSQL (>= v12)
2. Create a database named `line_gateway`
3. Run the `schema.sql` file containing all tables

```bash
psql -U postgres -d line_gateway -f schema.sql
```

## 🛠️ Development Setup

1. Clone the repository
```bash
git clone https://github.com/tikpoptv/opsalert-backend.git
cd opsalert-backend
```

2. Install dependencies
```bash
go mod download
```

3. Configure environment
```bash
cp example.env .env
# Edit .env with your configuration
```

4. Run the application
```bash
go run cmd/main.go
```

## 🔒 Security Features

- API key authentication for external systems
- Role-based access control for staff
- Comprehensive audit logging
- Secure password hashing
- Rate limiting for API endpoints

## 📊 API Documentation

The API documentation is available at `/swagger` when running in development mode.

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

The MIT License is a permissive license that is short and to the point. It lets people do anything they want with your code as long as they provide attribution back to you and don't hold you liable.

Key permissions:
- ✅ Commercial use
- ✅ Modification
- ✅ Distribution
- ✅ Private use

Key limitations:
- ⚠️ Liability
- ⚠️ Warranty

## 📞 Support

For support, please open an issue in the GitHub repository or contact the development team. 