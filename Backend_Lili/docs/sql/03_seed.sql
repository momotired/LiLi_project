-- 初始化数据（可选）

-- 预置部分系统标签
INSERT INTO tags (name, description, category, type, active, usage_count, created_at, updated_at)
VALUES
 ('手机', '手机相关', 'device_type', 'system', 1, 0, NOW(), NOW()),
 ('相机', '相机相关', 'device_type', 'system', 1, 0, NOW(), NOW()),
 ('电脑', '电脑相关', 'device_type', 'system', 1, 0, NOW(), NOW())
ON DUPLICATE KEY UPDATE updated_at = VALUES(updated_at);

-- 预置分类（示例）
INSERT INTO categories (name, description, parent_id, type, is_active, created_at, updated_at)
VALUES
 ('手机', '手机类设备', 0, 'system', 1, NOW(), NOW()),
 ('相机', '相机类设备', 0, 'system', 1, NOW(), NOW()),
 ('电脑', '电脑类设备', 0, 'system', 1, NOW(), NOW())
ON DUPLICATE KEY UPDATE updated_at = VALUES(updated_at);
