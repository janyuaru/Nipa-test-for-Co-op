CREATE TABLE ticket (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    contact VARCHAR(255) NOT NULL,
    status VARCHAR(255) NOT NULL, CHECK (status IN ('pending', 'accepted', 'resolved', 'rejected')),
    created_ticket_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    lastest_updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- สร้าง function สำหรับอัพเดท updated_at โดยอัตโนมัติ
CREATE OR REPLACE FUNCTION update_modified_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.lastest_updated_at = now();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- สร้าง trigger สำหรับอัพเดท updated_at
CREATE TRIGGER update_users_modtime
    BEFORE UPDATE ON ticket
    FOR EACH ROW
    EXECUTE FUNCTION update_modified_column();

-- เพิ่มข้อมูลตัวอย่าง
INSERT INTO ticket (title, description, contact, status) VALUES 
('ระบบล่ม', 'ไม่สามารถเข้าหน้า Dashboard ได้', 'user1@example.com', 'pending'),
('ระบบล่ม', 'ไม่สามารถเข้าหน้า Dashboard ได้', 'user2@example.com', 'pending');
