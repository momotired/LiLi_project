package model

import (
	"time"
)

// Device 设备信息表
type Device struct {
	ID             int       `orm:"column(id);auto;pk" json:"id"`
	UserID         int       `orm:"column(user_id)" json:"user_id"`
	TemplateID     *int      `orm:"column(template_id);null" json:"template_id"`
	CategoryID     *int      `orm:"column(category_id);null" json:"category_id"`
	Name           string    `orm:"column(name);size(200)" json:"name"`
	Brand          string    `orm:"column(brand);size(100)" json:"brand"`
	Model          string    `orm:"column(model);size(100)" json:"model"`
	SerialNumber   string    `orm:"column(serial_number);size(200);null" json:"serial_number"`
	Color          string    `orm:"column(color);size(50);null" json:"color"`
	Storage        string    `orm:"column(storage);size(50);null" json:"storage"`
	Memory         string    `orm:"column(memory);size(50);null" json:"memory"`
	Processor      string    `orm:"column(processor);size(100);null" json:"processor"`
	ScreenSize     string    `orm:"column(screen_size);size(50);null" json:"screen_size"`
	PurchasePrice  float64   `orm:"column(purchase_price);digits(10);decimals(2)" json:"purchase_price"`
	CurrentValue   float64   `orm:"column(current_value);digits(10);decimals(2);default(0)" json:"current_value"`
	PurchaseDate   time.Time `orm:"column(purchase_date);type(date)" json:"purchase_date"`
	WarrantyDate   time.Time `orm:"column(warranty_date);type(date);null" json:"warranty_date"`
	Condition      string    `orm:"column(condition);size(20);default(new)" json:"condition"` // new/good/fair/poor
	Status         string    `orm:"column(status);size(20);default(active)" json:"status"`    // active/sold/broken/lost
	SalePrice      float64   `orm:"column(sale_price);digits(10);decimals(2);null" json:"sale_price"`
	SaleDate       time.Time `orm:"column(sale_date);type(date);null" json:"sale_date"`
	Notes          string    `orm:"column(notes);type(text);null" json:"notes"`
	Specifications string    `orm:"column(specifications);type(json);null" json:"specifications"` // JSON格式存储其他规格
	CreatedAt      time.Time `orm:"column(created_at);auto_now_add;type(datetime)" json:"created_at"`
	UpdatedAt      time.Time `orm:"column(updated_at);auto_now;type(datetime)" json:"updated_at"`
	DeletedAt      time.Time `orm:"column(deleted_at);null;type(datetime)" json:"-"`

	// 关联字段 (不使用ORM自动关联，在代码中手动加载)
	Images []*DeviceImage `orm:"-" json:"images,omitempty"`
}

func (d *Device) TableName() string {
	return "devices"
}

// DeviceTemplate 设备模板表
type DeviceTemplate struct {
	ID          int       `orm:"column(id);auto;pk" json:"id"`
	Name        string    `orm:"column(name);size(100)" json:"name"`
	Description string    `orm:"column(description);type(text);null" json:"description"`
	Icon        string    `orm:"column(icon);size(500);null" json:"icon"`
	Fields      string    `orm:"column(fields);type(json)" json:"fields"` // JSON格式定义字段模板
	CategoryID  *int      `orm:"column(category_id);null" json:"category_id"`
	IsActive    bool      `orm:"column(is_active);default(true)" json:"is_active"`
	UseCount    int       `orm:"column(use_count);default(0)" json:"use_count"` // 使用次数，用于热门模板统计
	CreatedAt   time.Time `orm:"column(created_at);auto_now_add;type(datetime)" json:"created_at"`
	UpdatedAt   time.Time `orm:"column(updated_at);auto_now;type(datetime)" json:"updated_at"`
	DeletedAt   time.Time `orm:"column(deleted_at);null;type(datetime)" json:"-"`


}

func (dt *DeviceTemplate) TableName() string {
	return "device_templates"
}

// TemplateField 模板字段定义结构
type TemplateField struct {
	FieldName       string   `json:"field_name"`       // 字段名称
	FieldLabel      string   `json:"field_label"`      // 字段显示名称
	FieldType       string   `json:"field_type"`       // 字段类型：text/number/select/date/textarea
	Required        bool     `json:"required"`         // 是否必填
	DefaultValue    string   `json:"default_value"`    // 默认值
	ValidationRules string   `json:"validation_rules"` // 验证规则
	Options         []string `json:"options"`          // 选项列表(仅select类型)
	Placeholder     string   `json:"placeholder"`      // 占位符
	HelpText        string   `json:"help_text"`        // 帮助文本
}

// Category 设备分类表
type Category struct {
	ID          int       `orm:"column(id);auto;pk" json:"id"`
	Name        string    `orm:"column(name);size(100)" json:"name"`
	Description string    `orm:"column(description);type(text);null" json:"description"`
	ParentID    int       `orm:"column(parent_id);null;default(0)" json:"parent_id"` // 0表示顶级分类
	Icon        string    `orm:"column(icon);size(500);null" json:"icon"`
	Color       string    `orm:"column(color);size(20);null" json:"color"`
	SortOrder   int       `orm:"column(sort_order);default(0)" json:"sort_order"`
	Type        string    `orm:"column(type);size(20);default(system)" json:"type"` // system/custom
	UserID      int       `orm:"column(user_id);null" json:"user_id"`               // 自定义分类的用户ID
	IsActive    bool      `orm:"column(is_active);default(true)" json:"is_active"`
	DeviceCount int       `orm:"-" json:"device_count,omitempty"` // 设备数量，不存储在数据库
	CreatedAt   time.Time `orm:"column(created_at);auto_now_add;type(datetime)" json:"created_at"`
	UpdatedAt   time.Time `orm:"column(updated_at);auto_now;type(datetime)" json:"updated_at"`
	DeletedAt   time.Time `orm:"column(deleted_at);null;type(datetime)" json:"-"`

	// 关联字段
	Children []*Category `orm:"-" json:"children,omitempty"` // 子分类，不存储在数据库
}

func (c *Category) TableName() string {
	return "categories"
}

// DeviceImage 设备图片表
type DeviceImage struct {
	ID        int       `orm:"column(id);auto;pk" json:"id"`
	DeviceID  int       `orm:"column(device_id)" json:"device_id"`
	ImageURL  string    `orm:"column(image_url);size(500)" json:"image_url"`
	ImageType string    `orm:"column(image_type);size(20);default(normal)" json:"image_type"` // normal/cover
	SortOrder int       `orm:"column(sort_order);default(0)" json:"sort_order"`
	CreatedAt time.Time `orm:"column(created_at);auto_now_add;type(datetime)" json:"created_at"`


}

func (di *DeviceImage) TableName() string {
	return "device_images"
}


