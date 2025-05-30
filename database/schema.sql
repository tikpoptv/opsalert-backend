-- ลบตารางและฟังก์ชันเดิม (ปลอดภัยไม่ error ถ้าไม่มี)
DROP TABLE IF EXISTS
  audit_logs,
  api_logs,
  staff_oa_permissions,
  staff_accounts,
  messages,
  line_users,
  system_oa_permissions,
  external_systems,
  line_official_accounts,
  api_tokens
CASCADE;

DROP FUNCTION IF EXISTS log_audit() CASCADE;


-- LINE OA ที่ระบบควบคุม
CREATE TABLE line_official_accounts (
  id SERIAL PRIMARY KEY,
  name VARCHAR(100) NOT NULL,                     -- ชื่อ OA เช่น "Main OA"
  channel_id VARCHAR(100) NOT NULL,               -- จาก LINE Developer Console
  channel_secret TEXT NOT NULL,
  channel_access_token TEXT NOT NULL,
  webhook_url TEXT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- ระบบภายนอก (เช่น ร้านค้า, ERP ฯลฯ) ที่ยิง API เข้ามา
CREATE TABLE external_systems (
  id SERIAL PRIMARY KEY,
  name VARCHAR(100) NOT NULL,
  api_key VARCHAR(100) UNIQUE NOT NULL,           -- ใช้ยืนยัน API
  description TEXT,
  is_active BOOLEAN DEFAULT TRUE,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- สิทธิ์ว่า system ไหนยิง OA ไหนได้บ้าง
CREATE TABLE system_oa_permissions (
  id SERIAL PRIMARY KEY,
  system_id INTEGER REFERENCES external_systems(id) ON DELETE CASCADE,
  oa_id INTEGER REFERENCES line_official_accounts(id) ON DELETE CASCADE,
  permission_level VARCHAR(20) DEFAULT 'send'
    CHECK (permission_level IN ('send', 'admin')) -- จำกัดสิทธิ์เฉพาะคำที่กำหนด
);

-- ผู้ใช้ LINE (แยกตาม OA)
CREATE TABLE line_users (
  id SERIAL PRIMARY KEY,
  line_user_id VARCHAR(64) NOT NULL,
  display_name VARCHAR(255),
  oa_id INTEGER REFERENCES line_official_accounts(id) ON DELETE CASCADE,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  UNIQUE (line_user_id, oa_id)                    -- ป้องกัน user ซ้ำใน OA เดียวกัน
);

-- ประวัติข้อความที่ส่งออกจากระบบ
CREATE TABLE messages (
  id SERIAL PRIMARY KEY,
  system_id INTEGER REFERENCES external_systems(id),
  oa_id INTEGER REFERENCES line_official_accounts(id),
  line_user_id INTEGER REFERENCES line_users(id),
  content TEXT NOT NULL,
  sent_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  status VARCHAR(20) DEFAULT 'success'
    CHECK (status IN ('success', 'failed')),      -- จำกัดแค่ success/failed
  error_message TEXT
);

-- พนักงาน/ผู้ดูแลระบบหลังบ้าน
CREATE TABLE staff_accounts (
  id SERIAL PRIMARY KEY,
  username VARCHAR(50) UNIQUE NOT NULL,
  password_hash TEXT NOT NULL,                    -- เข้ารหัสด้วย bcrypt
  full_name VARCHAR(100),
  role VARCHAR(20) DEFAULT 'staff'
    CHECK (role IN ('admin', 'staff')),           -- จำกัด role
  is_active BOOLEAN DEFAULT TRUE,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- สิทธิ์ว่าแอดมินเข้าถึง OA ใดได้บ้าง
CREATE TABLE staff_oa_permissions (
  id SERIAL PRIMARY KEY,
  staff_id INTEGER REFERENCES staff_accounts(id) ON DELETE CASCADE,
  oa_id INTEGER REFERENCES line_official_accounts(id) ON DELETE CASCADE,
  permission_level VARCHAR(20) DEFAULT 'manage'
    CHECK (permission_level IN ('view', 'manage')) -- จำกัดสิทธิ์
);

-- API Token สำหรับ user
CREATE TABLE api_tokens (
  id SERIAL PRIMARY KEY,
  user_id INTEGER REFERENCES staff_accounts(id) ON DELETE CASCADE,
  token VARCHAR(100) UNIQUE NOT NULL,           -- token สำหรับยืนยัน API
  name VARCHAR(100) NOT NULL,                   -- ชื่อ token (เช่น "Development", "Production")
  is_active BOOLEAN DEFAULT TRUE,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  last_used_at TIMESTAMP                        -- เก็บเวลาที่ใช้ token ล่าสุด
);

-- Log API Request จาก token
CREATE TABLE api_logs (
  id SERIAL PRIMARY KEY,
  token_id INTEGER REFERENCES api_tokens(id),
  endpoint TEXT,
  method VARCHAR(10)
    CHECK (method IN ('GET', 'POST', 'PUT', 'DELETE', 'PATCH')),
  request_body TEXT,
  response_status INTEGER
    CHECK (response_status BETWEEN 100 AND 599),  -- จำกัด HTTP status code
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Audit Log สำหรับติดตามทุกการเปลี่ยนแปลงในระบบ
CREATE TABLE audit_logs (
  id SERIAL PRIMARY KEY,
  table_name TEXT NOT NULL,
  action TEXT NOT NULL
    CHECK (action IN ('INSERT', 'UPDATE', 'DELETE')),
  record_id TEXT,
  changed_data JSONB,
  changed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  changed_by TEXT                                   -- รองรับระบบ auth ภายหลัง
);


CREATE OR REPLACE FUNCTION log_audit() RETURNS TRIGGER AS $$
DECLARE
  record_json JSONB;
BEGIN
  IF (TG_OP = 'INSERT') THEN
    record_json := to_jsonb(NEW);
    INSERT INTO audit_logs(table_name, action, record_id, changed_data)
    VALUES (TG_TABLE_NAME, 'INSERT', NEW.id::TEXT, record_json);

  ELSIF (TG_OP = 'UPDATE') THEN
    record_json := jsonb_build_object('old', to_jsonb(OLD), 'new', to_jsonb(NEW));
    INSERT INTO audit_logs(table_name, action, record_id, changed_data)
    VALUES (TG_TABLE_NAME, 'UPDATE', NEW.id::TEXT, record_json);

  ELSIF (TG_OP = 'DELETE') THEN
    record_json := to_jsonb(OLD);
    INSERT INTO audit_logs(table_name, action, record_id, changed_data)
    VALUES (TG_TABLE_NAME, 'DELETE', OLD.id::TEXT, record_json);
  END IF;
  RETURN NULL;
END;
$$ LANGUAGE plpgsql;


CREATE TRIGGER trg_audit_line_official_accounts
AFTER INSERT OR UPDATE OR DELETE ON line_official_accounts
FOR EACH ROW EXECUTE FUNCTION log_audit();

CREATE TRIGGER trg_audit_external_systems
AFTER INSERT OR UPDATE OR DELETE ON external_systems
FOR EACH ROW EXECUTE FUNCTION log_audit();

CREATE TRIGGER trg_audit_system_oa_permissions
AFTER INSERT OR UPDATE OR DELETE ON system_oa_permissions
FOR EACH ROW EXECUTE FUNCTION log_audit();

CREATE TRIGGER trg_audit_line_users
AFTER INSERT OR UPDATE OR DELETE ON line_users
FOR EACH ROW EXECUTE FUNCTION log_audit();

CREATE TRIGGER trg_audit_messages
AFTER INSERT OR UPDATE OR DELETE ON messages
FOR EACH ROW EXECUTE FUNCTION log_audit();

CREATE TRIGGER trg_audit_staff_accounts
AFTER INSERT OR UPDATE OR DELETE ON staff_accounts
FOR EACH ROW EXECUTE FUNCTION log_audit();

CREATE TRIGGER trg_audit_staff_oa_permissions
AFTER INSERT OR UPDATE OR DELETE ON staff_oa_permissions
FOR EACH ROW EXECUTE FUNCTION log_audit();

CREATE TRIGGER trg_audit_api_tokens
AFTER INSERT OR UPDATE OR DELETE ON api_tokens
FOR EACH ROW EXECUTE FUNCTION log_audit();

CREATE TRIGGER trg_audit_api_logs
AFTER INSERT OR UPDATE OR DELETE ON api_logs
FOR EACH ROW EXECUTE FUNCTION log_audit();
