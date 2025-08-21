package repository

import (
	"Backend_Lili/internal/device/model"
	"errors"
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type DeviceRepository struct{}

func NewDeviceRepository() *DeviceRepository {
	return &DeviceRepository{}
}

// GetDevicesList 获取设备列表   键为字符串类型 值为任意类型的映射 处理结构不固定 类型不确定的数据（动态数据）
func (r *DeviceRepository) GetDevicesList(userID int, params map[string]interface{}) ([]*model.Device, int64, error) {
	o := orm.NewOrm()
	qs := o.QueryTable("devices").Filter("user_id", userID).Filter("deleted_at__isnull", true)

	// 应用筛选条件
	if categoryID, ok := params["category_id"].(int); ok && categoryID > 0 {
		qs = qs.Filter("category_id", categoryID)
	}
	if status, ok := params["status"].(string); ok && status != "" {
		qs = qs.Filter("status", status)
	}
	if search, ok := params["search"].(string); ok && search != "" {
		qs = qs.Filter("name__icontains", search).Filter("brand__icontains", search).Filter("model__icontains", search)
	}

	// 应用排序
	if sort, ok := params["sort"].(string); ok && sort != "" {
		order := "desc"
		if orderParam, ok := params["order"].(string); ok && orderParam == "asc" {
			order = "asc"
		}
		if order == "asc" {
			qs = qs.OrderBy(sort)
		} else {
			qs = qs.OrderBy("-" + sort)
		}
	} else {
		qs = qs.OrderBy("-created_at")
	}

	// 获取总数
	total, err := qs.Count()
	if err != nil {
		return nil, 0, err
	}

	// 分页
	if page, ok := params["page"].(int); ok && page > 0 {
		limit := 10
		if limitParam, ok := params["limit"].(int); ok && limitParam > 0 {
			limit = limitParam
		}
		offset := (page - 1) * limit
		qs = qs.Limit(limit, offset)
	}

	var devices []*model.Device
	_, err = qs.RelatedSel().All(&devices)
	return devices, total, err
}

// GetDeviceByID 根据ID获取设备详情
func (r *DeviceRepository) GetDeviceByID(deviceID, userID int) (*model.Device, error) {
	o := orm.NewOrm()
	device := &model.Device{}
	err := o.QueryTable("devices").
		Filter("id", deviceID).
		Filter("user_id", userID).
		Filter("deleted_at__isnull").
		RelatedSel().
		One(device)

	if err == orm.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	// 加载设备图片
	var images []*model.DeviceImage
	_, err = o.QueryTable("device_images").
		Filter("device_id", deviceID).
		OrderBy("sort_order", "created_at").
		All(&images)
	if err == nil {
		device.Images = images
	}

	return device, nil
}

// CreateDevice 创建设备
func (r *DeviceRepository) CreateDevice(device *model.Device) error {
	o := orm.NewOrm()
	device.CreatedAt = time.Now()
	device.UpdatedAt = time.Now()
	_, err := o.Insert(device)
	return err
}

// UpdateDevice 更新设备
func (r *DeviceRepository) UpdateDevice(device *model.Device) error {
	o := orm.NewOrm()
	device.UpdatedAt = time.Now()
	_, err := o.Update(device)
	return err
}

// SoftDeleteDevice 软删除设备
func (r *DeviceRepository) SoftDeleteDevice(deviceID, userID int) error {
	o := orm.NewOrm()
	_, err := o.QueryTable("devices").
		Filter("id", deviceID).
		Filter("user_id", userID).
		Update(orm.Params{
			"deleted_at": time.Now(),
			"updated_at": time.Now(),
		})
	return err
}

// UpdateDeviceStatus 更新设备状态
func (r *DeviceRepository) UpdateDeviceStatus(deviceID, userID int, status string, salePrice *float64, saleDate *time.Time, notes string) error {
	o := orm.NewOrm()
	params := orm.Params{
		"status":     status,
		"updated_at": time.Now(),
	}

	if notes != "" {
		params["notes"] = notes
	}
	if salePrice != nil {
		params["sale_price"] = *salePrice
	}
	if saleDate != nil {
		params["sale_date"] = *saleDate
	}

	_, err := o.QueryTable("devices").
		Filter("id", deviceID).
		Filter("user_id", userID).
		Update(params)
	return err
}

// GetDeviceImages 获取设备图片
func (r *DeviceRepository) GetDeviceImages(deviceID, userID int) ([]*model.DeviceImage, error) {
	o := orm.NewOrm()

	// 先验证设备归属
	exists := o.QueryTable("devices").
		Filter("id", deviceID).
		Filter("user_id", userID).
		Filter("deleted_at__isnull", true).
		Exist()
	if !exists {
		return nil, errors.New("设备不存在或无权限")
	}

	var images []*model.DeviceImage
	_, err := o.QueryTable("device_images").
		Filter("device_id", deviceID).
		OrderBy("sort_order", "created_at").
		All(&images)
	return images, err
}

// AddDeviceImage 添加设备图片
func (r *DeviceRepository) AddDeviceImage(image *model.DeviceImage) error {
	o := orm.NewOrm()
	image.CreatedAt = time.Now()
	_, err := o.Insert(image)
	return err
}

// DeleteDeviceImage 删除设备图片
func (r *DeviceRepository) DeleteDeviceImage(imageID, deviceID, userID int) error {
	o := orm.NewOrm()

	// 验证权限：先检查设备是否属于该用户
	deviceExists := o.QueryTable("devices").
		Filter("id", deviceID).
		Filter("user_id", userID).
		Filter("deleted_at__isnull", true).
		Exist()
	if !deviceExists {
		return errors.New("设备不存在或无权限")
	}

	// 验证图片是否存在且属于该设备
	imageExists := o.QueryTable("device_images").
		Filter("id", imageID).
		Filter("device_id", deviceID).
		Exist()
	if !imageExists {
		return errors.New("图片不存在")
	}

	_, err := o.QueryTable("device_images").Filter("id", imageID).Delete()
	return err
}

// GetPriceHistory 获取设备价格历史
// 注意：价格历史功能应该在价格模块中实现
func (r *DeviceRepository) GetPriceHistory(deviceID, userID int, limit int) ([]interface{}, error) {
	// 这个方法应该调用价格模块的服务
	// 暂时返回空结果，避免循环依赖
	return []interface{}{}, nil
}

// AddPriceHistory 添加价格历史记录
// 注意：价格历史功能应该在价格模块中实现
func (r *DeviceRepository) AddPriceHistory(history interface{}) error {
	// 这个方法应该调用价格模块的服务
	// 暂时返回nil，避免循环依赖
	return nil
}

// UpdateDeviceCurrentValue 更新设备当前估值
func (r *DeviceRepository) UpdateDeviceCurrentValue(deviceID int, currentValue float64) error {
	o := orm.NewOrm()
	_, err := o.QueryTable("devices").
		Filter("id", deviceID).
		Update(orm.Params{
			"current_value": currentValue,
			"updated_at":    time.Now(),
		})
	return err
}

// BatchImportDevices 批量导入设备
func (r *DeviceRepository) BatchImportDevices(devices []*model.Device) (int, error) {
	o := orm.NewOrm()
	successCount := 0

	for _, device := range devices {
		device.CreatedAt = time.Now()
		device.UpdatedAt = time.Now()
		_, err := o.Insert(device)
		if err != nil {
			continue // 跳过失败的记录
		}
		successCount++
	}

	return successCount, nil
}
