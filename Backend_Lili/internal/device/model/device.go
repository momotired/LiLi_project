package model

import (
	"time"
)

// Device 设备信息表
type Device struct {
	ID             int       `orm:"column(id);auto;pk" json:"id"`
	UserID         *int      `orm:"column(user_id);rel(fk);null" json:"user_id"`
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

	// 关联字段
	Template *DeviceTemplate `orm:"rel(fk);null;on_delete(set_null)" json:"template,omitempty"`
	Category *Category       `orm:"rel(fk);null;on_delete(set_null)" json:"category,omitempty"`
	Images   []*DeviceImage  `orm:"reverse(many)" json:"images,omitempty"`
}

func (d *Device) TableName() string {
	return "devices"
}

// DeviceTemplate 设备模板表
type DeviceTemplate struct {
	ID          int       `orm:"column(id);auto;pk" json:"id"`
	Name        string    `orm:"column(name);size(100)" json:"name"`
	Fields      string    `orm:"column(fields);type(json)" json:"fields"` // JSON格式定义字段模板
	Description string    `orm:"column(description);type(text);null" json:"description"`
	IsActive    bool      `orm:"column(is_active);default(true)" json:"is_active"`
	CreatedAt   time.Time `orm:"column(created_at);auto_now_add;type(datetime)" json:"created_at"`
	UpdatedAt   time.Time `orm:"column(updated_at);auto_now;type(datetime)" json:"updated_at"`

	// 关联字段
	Category *Category `orm:"rel(fk);null;on_delete(cascade)" json:"category,omitempty"`
}

func (dt *DeviceTemplate) TableName() string {
	return "device_templates"
}

// Category 设备分类表
type Category struct {
	ID        int       `orm:"column(id);auto;pk" json:"id"`
	Name      string    `orm:"column(name);size(100)" json:"name"`
	ParentID  int       `orm:"column(parent_id);null;default(0)" json:"parent_id"` // 0表示顶级分类
	Icon      string    `orm:"column(icon);size(200);null" json:"icon"`
	SortOrder int       `orm:"column(sort_order);default(0)" json:"sort_order"`
	IsActive  bool      `orm:"column(is_active);default(true)" json:"is_active"`
	CreatedAt time.Time `orm:"column(created_at);auto_now_add;type(datetime)" json:"created_at"`
	UpdatedAt time.Time `orm:"column(updated_at);auto_now;type(datetime)" json:"updated_at"`
}

func (c *Category) TableName() string {
	return "categories"
}

// DeviceImage 设备图片表
type DeviceImage struct {
	ID        int       `orm:"column(id);auto;pk" json:"id"`
	ImageURL  string    `orm:"column(image_url);size(500)" json:"image_url"`
	ImageType string    `orm:"column(image_type);size(20);default(normal)" json:"image_type"` // normal/cover
	SortOrder int       `orm:"column(sort_order);default(0)" json:"sort_order"`
	CreatedAt time.Time `orm:"column(created_at);auto_now_add;type(datetime)" json:"created_at"`

	// 关联字段
	Device *Device `orm:"rel(fk);on_delete(cascade)" json:"device,omitempty"`
}

func (di *DeviceImage) TableName() string {
	return "device_images"
}

// PriceHistory 价格历史表
type PriceHistory struct {
	ID int `orm:"column(id);auto;pk" json:"id"`

	Source      string    `orm:"column(source);size(50)" json:"source"`          // 价格来源：manual/market_api
	Platform    string    `orm:"column(platform);size(50);null" json:"platform"` // 平台名称：如闲鱼、转转等
	Price       float64   `orm:"column(price);digits(10);decimals(2)" json:"price"`
	Condition   string    `orm:"column(condition);size(20)" json:"condition"`
	Description string    `orm:"column(description);size(500);null" json:"description"`
	RecordDate  time.Time `orm:"column(record_date);type(date)" json:"record_date"`
	CreatedAt   time.Time `orm:"column(created_at);auto_now_add;type(datetime)" json:"created_at"`

	// 关联字段
	Device *Device `orm:"rel(fk);on_delete(cascade)" json:"device,omitempty"`
}

func (ph *PriceHistory) TableName() string {
	return "price_histories"
}
