package models

import "gorm.io/gorm"

type Project struct {
	Base
	// 用户名，默认英文，唯一，不可为空
	Name string `gorm:"unique,index,size:64,not null" json:"name" binding:"required"`
	// 中文名称
	CnName string `gorm:"size:64;not null;unique" json:"cnName" binding:"required"`
	// 描述
	Description string `gorm:"size:1024" json:"description"`
	// 开发语言
	Language string `gorm:"size:32" json:"language"`
	// 开发框架
	Framework string `gorm:"size:32" json:"framework"`
	// 负责人
	Owner string `gorm:"size:256" json:"owner"`
	// 产品负责人
	ProductOwner string `gorm:"size:256" json:"productOwner"`
	// 测试负责人
	TestOwner string `gorm:"size:256" json:"testOwner"`
	// 所属模块 id
	ModuleID uint `gorm:"size:256;not null" json:"moduleID"`
	// 所属模块
	Module string `json:"module"`
}

func (p *Project) AfterFind(tx *gorm.DB) error {
	var md Module
	r := tx.Model(p).Where("id = ?", p.ModuleID).First(&md)
	if r.Error != nil {
		return r.Error
	}
	p.Module = md.FullName
	return nil
}
