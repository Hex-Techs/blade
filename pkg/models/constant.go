package models

import (
	"time"

	"gorm.io/gorm"
)

const salt = "blade"

// 基类
type Base struct {
	// 默认主键为ID，类型为uint
	ID uint `gorm:"primary_key,index,unique" json:"id"`
	// 创建时间
	CreatedAt time.Time `json:"createdAt"`
	// 更新时间
	UpdatedAt time.Time `json:"updatedAt"`
	// 删除时间
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// BeforeInsert 在插入数据之前设置创建时间
func (b *Base) BeforeInsert() {
	b.CreatedAt = time.Now()
}
