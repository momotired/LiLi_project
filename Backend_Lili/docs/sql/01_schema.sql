-- 数据库建表脚本（MySQL 8.0）
-- 说明：请根据自己的数据库名、字符集按需修改。执行顺序：01_schema.sql -> 02_indexes.sql -> 03_seed.sql

-- 用户与偏好
CREATE TABLE IF NOT EXISTS users (
  id INT PRIMARY KEY AUTO_INCREMENT,
  openid VARCHAR(100) NOT NULL UNIQUE,
  unionid VARCHAR(100) NULL,
  session_key VARCHAR(100) NULL,
  nickname VARCHAR(100) NULL,
  avatar VARCHAR(500) NULL,
  gender INT NOT NULL DEFAULT 0,
  city VARCHAR(100) NULL,
  province VARCHAR(100) NULL,
  country VARCHAR(100) NULL,
  language VARCHAR(50) NULL,
  status INT NOT NULL DEFAULT 1,
  last_login_at DATETIME NULL,
  created_at DATETIME NOT NULL,
  updated_at DATETIME NOT NULL,
  deleted_at DATETIME NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS user_preferences (
  id INT PRIMARY KEY AUTO_INCREMENT,
  user_id INT NOT NULL UNIQUE,
  notification_enabled TINYINT(1) NOT NULL DEFAULT 0,
  price_alert_enabled TINYINT(1) NOT NULL DEFAULT 0,
  warranty_reminder_enabled TINYINT(1) NOT NULL DEFAULT 0,
  theme VARCHAR(20) NULL,
  language VARCHAR(20) NULL,
  created_at DATETIME NOT NULL,
  updated_at DATETIME NOT NULL,
  CONSTRAINT fk_user_preferences_user FOREIGN KEY (user_id) REFERENCES users(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 标签
CREATE TABLE IF NOT EXISTS tags (
  id INT PRIMARY KEY AUTO_INCREMENT,
  name VARCHAR(100) NOT NULL,
  description VARCHAR(500) NULL,
  category VARCHAR(100) NULL,
  color VARCHAR(20) NULL,
  icon VARCHAR(200) NULL,
  type VARCHAR(20) NOT NULL DEFAULT 'system',
  active TINYINT(1) NOT NULL DEFAULT 1,
  usage_count INT NOT NULL DEFAULT 0,
  owner_id INT NULL,
  created_at DATETIME NOT NULL,
  updated_at DATETIME NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS user_tags (
  id INT PRIMARY KEY AUTO_INCREMENT,
  user_id INT NOT NULL,
  tag_id INT NOT NULL,
  CONSTRAINT fk_user_tags_user FOREIGN KEY (user_id) REFERENCES users(id),
  CONSTRAINT fk_user_tags_tag FOREIGN KEY (tag_id) REFERENCES tags(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 设备、分类、模板
CREATE TABLE IF NOT EXISTS categories (
  id INT PRIMARY KEY AUTO_INCREMENT,
  name VARCHAR(100) NOT NULL,
  description TEXT NULL,
  parent_id INT NULL DEFAULT 0,
  icon VARCHAR(500) NULL,
  color VARCHAR(20) NULL,
  sort_order INT NOT NULL DEFAULT 0,
  type VARCHAR(20) NOT NULL DEFAULT 'system',
  user_id INT NULL,
  is_active TINYINT(1) NOT NULL DEFAULT 1,
  created_at DATETIME NOT NULL,
  updated_at DATETIME NOT NULL,
  deleted_at DATETIME NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS device_templates (
  id INT PRIMARY KEY AUTO_INCREMENT,
  name VARCHAR(100) NOT NULL,
  category_id INT NULL,
  description TEXT NULL,
  icon VARCHAR(500) NULL,
  fields JSON NOT NULL,
  is_active TINYINT(1) NOT NULL DEFAULT 1,
  use_count INT NOT NULL DEFAULT 0,
  created_at DATETIME NOT NULL,
  updated_at DATETIME NOT NULL,
  deleted_at DATETIME NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS devices (
  id INT PRIMARY KEY AUTO_INCREMENT,
  user_id INT NULL,
  template_id INT NULL,
  category_id INT NULL,
  name VARCHAR(200) NOT NULL,
  brand VARCHAR(100) NOT NULL,
  model VARCHAR(100) NOT NULL,
  serial_number VARCHAR(200) NULL,
  color VARCHAR(50) NULL,
  storage VARCHAR(50) NULL,
  memory VARCHAR(50) NULL,
  processor VARCHAR(100) NULL,
  screen_size VARCHAR(50) NULL,
  purchase_price DECIMAL(10,2) NOT NULL,
  current_value DECIMAL(10,2) NOT NULL DEFAULT 0,
  purchase_date DATE NOT NULL,
  warranty_date DATE NULL,
  `condition` VARCHAR(20) NOT NULL DEFAULT 'new',
  status VARCHAR(20) NOT NULL DEFAULT 'active',
  sale_price DECIMAL(10,2) NULL,
  sale_date DATE NULL,
  notes TEXT NULL,
  specifications JSON NULL,
  created_at DATETIME NOT NULL,
  updated_at DATETIME NOT NULL,
  deleted_at DATETIME NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS device_images (
  id INT PRIMARY KEY AUTO_INCREMENT,
  device_id INT NOT NULL,
  image_url VARCHAR(500) NOT NULL,
  image_type VARCHAR(20) NOT NULL DEFAULT 'normal',
  sort_order INT NOT NULL DEFAULT 0,
  created_at DATETIME NOT NULL,
  CONSTRAINT fk_device_images_device FOREIGN KEY (device_id) REFERENCES devices(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 价格
CREATE TABLE IF NOT EXISTS prices (
  id INT PRIMARY KEY AUTO_INCREMENT,
  device_id INT NOT NULL,
  user_id INT NOT NULL,
  current_price DECIMAL(10,2) NOT NULL,
  market_price DECIMAL(10,2) NULL,
  average_price DECIMAL(10,2) NULL,
  min_price DECIMAL(10,2) NULL,
  max_price DECIMAL(10,2) NULL,
  price_change DECIMAL(10,2) NOT NULL DEFAULT 0,
  change_rate DECIMAL(5,2) NOT NULL DEFAULT 0,
  trend_status VARCHAR(20) NOT NULL DEFAULT 'stable',
  last_update_at DATETIME NULL,
  created_at DATETIME NOT NULL,
  updated_at DATETIME NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS price_histories (
  id INT PRIMARY KEY AUTO_INCREMENT,
  device_id INT NOT NULL,
  user_id INT NOT NULL,
  source VARCHAR(50) NOT NULL,
  source_id VARCHAR(100) NULL,
  platform VARCHAR(50) NULL,
  price DECIMAL(10,2) NOT NULL,
  `condition` VARCHAR(20) NOT NULL,
  description VARCHAR(500) NULL,
  url VARCHAR(1000) NULL,
  record_date DATE NOT NULL,
  created_at DATETIME NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS price_alerts (
  id INT PRIMARY KEY AUTO_INCREMENT,
  device_id INT NOT NULL,
  user_id INT NOT NULL,
  alert_type VARCHAR(20) NOT NULL,
  threshold DECIMAL(10,2) NOT NULL,
  threshold_type VARCHAR(20) NOT NULL,
  enabled TINYINT(1) NOT NULL DEFAULT 1,
  notification_methods VARCHAR(200) NULL,
  last_triggered_at DATETIME NULL,
  trigger_count INT NOT NULL DEFAULT 0,
  status VARCHAR(20) NOT NULL DEFAULT 'active',
  created_at DATETIME NOT NULL,
  updated_at DATETIME NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS price_sources (
  id INT PRIMARY KEY AUTO_INCREMENT,
  name VARCHAR(100) NOT NULL,
  platform VARCHAR(50) NOT NULL,
  base_url VARCHAR(500) NULL,
  api_endpoint VARCHAR(500) NULL,
  status VARCHAR(20) NOT NULL DEFAULT 'active',
  reliability DECIMAL(3,2) NOT NULL DEFAULT 1.00,
  update_freq INT NOT NULL DEFAULT 24,
  last_sync DATETIME NULL,
  `config` JSON NULL,
  created_at DATETIME NOT NULL,
  updated_at DATETIME NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS price_predictions (
  id INT PRIMARY KEY AUTO_INCREMENT,
  device_id INT NOT NULL,
  user_id INT NOT NULL,
  prediction_type VARCHAR(20) NOT NULL,
  predicted_price DECIMAL(10,2) NOT NULL,
  confidence DECIMAL(3,2) NOT NULL,
  algorithm VARCHAR(50) NOT NULL,
  factors JSON NULL,
  valid_until DATETIME NOT NULL,
  actual_price DECIMAL(10,2) NULL,
  accuracy DECIMAL(3,2) NULL,
  created_at DATETIME NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 认证（黑名单与会话）
CREATE TABLE IF NOT EXISTS token_blacklist (
  id INT PRIMARY KEY AUTO_INCREMENT,
  token VARCHAR(1000) NOT NULL,
  expires_at DATETIME NOT NULL,
  created_at DATETIME NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS user_session (
  id INT PRIMARY KEY AUTO_INCREMENT,
  user_id INT NOT NULL,
  openid VARCHAR(100) NULL,
  session_key VARCHAR(100) NULL,
  access_token VARCHAR(1000) NULL,
  refresh_token VARCHAR(1000) NULL,
  expires_at DATETIME NULL,
  last_login_at DATETIME NULL,
  created_at DATETIME NOT NULL,
  updated_at DATETIME NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
