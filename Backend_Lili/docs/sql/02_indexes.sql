-- 索引与约束（在建表之后执行）
-- 说明：为兼容更多 MySQL 版本，取消 IF NOT EXISTS，首次执行即可。

-- users
ALTER TABLE users ADD INDEX idx_users_openid (openid);
ALTER TABLE users ADD INDEX idx_users_status (status);

-- user_preferences（外键已在 01_schema.sql 中创建，这里不再重复添加）

-- tags
ALTER TABLE tags ADD INDEX idx_tags_type_active (type, active);
ALTER TABLE tags ADD INDEX idx_tags_category (category);
ALTER TABLE tags ADD INDEX idx_tags_owner (owner_id);

-- user_tags
ALTER TABLE user_tags ADD INDEX idx_user_tags_user_tag (user_id, tag_id);

-- categories
ALTER TABLE categories ADD INDEX idx_categories_type (type);
ALTER TABLE categories ADD INDEX idx_categories_parent (parent_id);

-- device_templates
ALTER TABLE device_templates ADD INDEX idx_device_templates_category (category_id);

-- devices
ALTER TABLE devices ADD INDEX idx_devices_user (user_id);
ALTER TABLE devices ADD INDEX idx_devices_category (category_id);
ALTER TABLE devices ADD INDEX idx_devices_brand (brand);
ALTER TABLE devices ADD INDEX idx_devices_status (status);
ALTER TABLE devices ADD INDEX idx_devices_purchase_date (purchase_date);

-- device_images
ALTER TABLE device_images ADD INDEX idx_device_images_device (device_id);

-- prices
ALTER TABLE prices ADD INDEX idx_prices_user_device (user_id, device_id);

-- price_histories
ALTER TABLE price_histories ADD INDEX idx_price_histories_user_device_date (user_id, device_id, record_date);
ALTER TABLE price_histories ADD INDEX idx_price_histories_platform (platform);

-- price_alerts
ALTER TABLE price_alerts ADD INDEX idx_price_alerts_user_device (user_id, device_id);
ALTER TABLE price_alerts ADD INDEX idx_price_alerts_status (status);

-- price_sources
ALTER TABLE price_sources ADD INDEX idx_price_sources_status (status);

-- price_predictions
ALTER TABLE price_predictions ADD INDEX idx_price_predictions_user_device (user_id, device_id);

-- token_blacklist
ALTER TABLE token_blacklist ADD INDEX idx_token_blacklist_token (token);

-- user_session
ALTER TABLE user_session ADD INDEX idx_user_session_user (user_id);
