-- 设备模块初始化数据脚本
-- 用于在开发和测试环境中插入基础的分类和模板数据

-- 清理现有数据（开发环境使用）
-- DELETE FROM device_templates WHERE id > 0;
-- DELETE FROM categories WHERE id > 0;

-- 插入设备分类数据
INSERT INTO categories (id, name, parent_id, icon, sort_order, is_active) VALUES
-- 一级分类
(1, '手机', 0, 'phone', 1, true),
(2, '电脑', 0, 'laptop', 2, true),
(3, '平板', 0, 'tablet', 3, true),
(4, '穿戴设备', 0, 'watch', 4, true),
(5, '音频设备', 0, 'headphone', 5, true),
(6, '摄影设备', 0, 'camera', 6, true),
(7, '游戏设备', 0, 'gamepad', 7, true),

-- 二级分类 - 手机
(101, 'iPhone', 1, 'iphone', 1, true),
(102, 'Android手机', 1, 'android', 2, true),

-- 二级分类 - 电脑
(201, 'MacBook', 2, 'macbook', 1, true),
(202, 'Windows笔记本', 2, 'windows-laptop', 2, true),
(203, '台式机', 2, 'desktop', 3, true),

-- 二级分类 - 平板
(301, 'iPad', 3, 'ipad', 1, true),
(302, 'Android平板', 3, 'android-tablet', 2, true),

-- 二级分类 - 穿戴设备
(401, 'Apple Watch', 4, 'apple-watch', 1, true),
(402, '智能手环', 4, 'smart-band', 2, true),

-- 二级分类 - 音频设备
(501, 'AirPods', 5, 'airpods', 1, true),
(502, '头戴式耳机', 5, 'headphone', 2, true),
(503, '音箱', 5, 'speaker', 3, true),

-- 二级分类 - 摄影设备
(601, '数码相机', 6, 'camera', 1, true),
(602, '运动相机', 6, 'action-camera', 2, true),
(603, '镜头', 6, 'lens', 3, true);

-- 插入设备模板数据
INSERT INTO device_templates (id, name, category_id, fields, description, is_active) VALUES
-- iPhone模板
(1, 'iPhone标准模板', 101, JSON_OBJECT(
    'required_fields', JSON_ARRAY('brand', 'model', 'color', 'storage', 'condition'),
    'optional_fields', JSON_ARRAY('memory', 'serial_number', 'warranty_date'),
    'field_options', JSON_OBJECT(
        'color', JSON_ARRAY('深空灰', '银色', '金色', '玫瑰金', '蓝色', '紫色', '红色'),
        'storage', JSON_ARRAY('64GB', '128GB', '256GB', '512GB', '1TB'),
        'condition', JSON_ARRAY('全新', '99新', '95新', '9成新', '8成新', '7成新')
    )
), 'iPhone系列设备模板，包含常见配置选项', true),

-- MacBook模板
(2, 'MacBook标准模板', 201, JSON_OBJECT(
    'required_fields', JSON_ARRAY('brand', 'model', 'screen_size', 'processor', 'memory', 'storage'),
    'optional_fields', JSON_ARRAY('color', 'serial_number', 'warranty_date'),
    'field_options', JSON_OBJECT(
        'screen_size', JSON_ARRAY('13.3英寸', '14.2英寸', '15.3英寸', '16.2英寸'),
        'processor', JSON_ARRAY('M1', 'M1 Pro', 'M1 Max', 'M2', 'M2 Pro', 'M2 Max', 'M3', 'M3 Pro', 'M3 Max'),
        'memory', JSON_ARRAY('8GB', '16GB', '32GB', '64GB', '128GB'),
        'storage', JSON_ARRAY('256GB', '512GB', '1TB', '2TB', '4TB', '8TB'),
        'color', JSON_ARRAY('深空灰', '银色')
    )
), 'MacBook系列设备模板', true),

-- iPad模板
(3, 'iPad标准模板', 301, JSON_OBJECT(
    'required_fields', JSON_ARRAY('brand', 'model', 'screen_size', 'storage'),
    'optional_fields', JSON_ARRAY('color', 'memory', 'cellular', 'serial_number', 'warranty_date'),
    'field_options', JSON_OBJECT(
        'screen_size', JSON_ARRAY('7.9英寸', '8.3英寸', '10.2英寸', '10.9英寸', '11英寸', '12.9英寸'),
        'storage', JSON_ARRAY('32GB', '64GB', '128GB', '256GB', '512GB', '1TB', '2TB'),
        'color', JSON_ARRAY('深空灰', '银色', '金色', '玫瑰金', '蓝色', '绿色', '紫色'),
        'cellular', JSON_ARRAY('WiFi版', 'WiFi+蜂窝版')
    )
), 'iPad系列设备模板', true),

-- Apple Watch模板
(4, 'Apple Watch标准模板', 401, JSON_OBJECT(
    'required_fields', JSON_ARRAY('brand', 'model', 'screen_size', 'case_material'),
    'optional_fields', JSON_ARRAY('color', 'band_type', 'cellular', 'serial_number', 'warranty_date'),
    'field_options', JSON_OBJECT(
        'screen_size', JSON_ARRAY('38mm', '40mm', '41mm', '42mm', '44mm', '45mm', '49mm'),
        'case_material', JSON_ARRAY('铝金属', '不锈钢', '钛金属', '陶瓷'),
        'color', JSON_ARRAY('银色', '深空灰', '金色', '石墨色', '蓝色', '红色'),
        'band_type', JSON_ARRAY('运动型表带', '运动回环表带', '编织单圈表带', '皮制链式表带', '米兰尼斯表带'),
        'cellular', JSON_ARRAY('GPS版', 'GPS+蜂窝版')
    )
), 'Apple Watch系列设备模板', true),

-- 数码相机模板
(5, '数码相机标准模板', 601, JSON_OBJECT(
    'required_fields', JSON_ARRAY('brand', 'model', 'sensor_type', 'resolution'),
    'optional_fields', JSON_ARRAY('lens_mount', 'video_capability', 'serial_number', 'warranty_date'),
    'field_options', JSON_OBJECT(
        'sensor_type', JSON_ARRAY('全画幅', 'APS-C', '4/3', '1英寸', '更小传感器'),
        'resolution', JSON_ARRAY('1200万像素', '1600万像素', '2000万像素', '2400万像素', '3000万像素', '4000万像素+'),
        'lens_mount', JSON_ARRAY('佳能EF', '佳能RF', '尼康F', '尼康Z', '索尼E', '富士X', '奥林巴斯M4/3'),
        'video_capability', JSON_ARRAY('1080p', '4K 30fps', '4K 60fps', '8K')
    )
), '数码相机设备模板', true);

-- 重置自增ID
ALTER TABLE categories AUTO_INCREMENT = 700;
ALTER TABLE device_templates AUTO_INCREMENT = 100; 