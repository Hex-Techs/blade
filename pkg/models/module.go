// models 服务模块
package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Module struct {
	Base
	// 名称
	Name string `gorm:"size:128;not null;unique;index" json:"name" binding:"required"`
	// 中文名称
	CnName string `gorm:"size:128;not null;unique" json:"cnName" binding:"required"`
	// 描述
	Description string `gorm:"size:1024" json:"description"`
	// 父模块id
	ParentID uint `gorm:"index" json:"parentID"`
	// 级别
	Level uint `gorm:"index" json:"level"`
}

func (m *Module) BeforeCreate(tx *gorm.DB) error {
	if m.ParentID == 0 {
		m.Level = 1
		m.CreatedAt = time.Now()
		return nil
	}
	var parent Module
	r := tx.Model(m).Where("id = ?", m.ParentID).First(&parent)
	if r.Error != nil {
		return fmt.Errorf("parent module error: %v", r.Error)
	}
	m.Level = parent.Level + 1
	return nil
}
