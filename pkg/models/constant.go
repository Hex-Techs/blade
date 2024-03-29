package models

import (
	"time"

	"gorm.io/gorm"
)

const salt = "blade"

// 基类
type Base struct {
	// 默认主键为ID，类型为uint
	ID uint `gorm:"primary_key" json:"id"`
	// 创建时间
	CreatedAt time.Time `json:"createdAt"`
	// 更新时间
	UpdatedAt time.Time `json:"updatedAt"`
	// 删除时间
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// BeforeInsert 在插入数据之前设置创建时间
func (b *Base) BeforeCreate(tx *gorm.DB) error {
	b.CreatedAt = time.Now()
	return nil
}

func (b *Base) BeforeUpdate(tx *gorm.DB) error {
	b.UpdatedAt = time.Now()
	return nil
}
