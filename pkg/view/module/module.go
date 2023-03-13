package module

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hex-techs/blade/pkg/models"
	"github.com/hex-techs/blade/pkg/util/log"
	"github.com/hex-techs/blade/pkg/util/storage"
	"github.com/hex-techs/blade/pkg/util/web"
	"github.com/hex-techs/blade/pkg/view"
)

// ModuleController module controller
type ModuleController struct {
	Store *storage.Engine
}

// NewModuleController return a new module controller
func NewModuleController(s *storage.Engine) web.RestController {
	return &ModuleController{
		Store: s,
	}
}

// 版本号
func (*ModuleController) Version() string {
	return "v1"
}

// 资源名
func (*ModuleController) Name() string {
	return "module"
}

// Create 创建新模块
func (uc *ModuleController) Create(c *gin.Context) {
	var module models.Module
	if err := c.ShouldBindJSON(&module); err != nil {
		c.JSON(http.StatusBadRequest, web.ExceptResponse(errorMap[ErrInvalidParam], err))
		return
	}
	log.Debugf("create module: %v", module)
	if err := uc.Store.Create(context.TODO(), &module); err != nil {
		c.JSON(http.StatusOK, web.ExceptResponse(errorMap[ErrCreateModuleFailed], err))
		return
	}
	c.JSON(http.StatusOK, web.OkResponse())
}

// Delete 删除模块
func (uc *ModuleController) Delete(c *gin.Context) {
	id, err := view.GetID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, web.ExceptResponse(errorMap[ErrID], err))
		return
	}
	if err := uc.Store.Delete(context.TODO(), id, "", &models.Module{}); err != nil {
		c.JSON(http.StatusOK, web.ExceptResponse(errorMap[ErrDeleteModuleFailed], err))
		return
	}
	c.JSON(http.StatusOK, web.OkResponse())
}

// Update 更新模块信息
func (uc *ModuleController) Update(c *gin.Context) {
	id, err := view.GetID(c)
	if err != nil {
		c.JSON(http.StatusOK, web.ExceptResponse(errorMap[ErrID], err))
		return
	}
	var module models.Module
	c.ShouldBindJSON(&module)
	log.Debugf("update module: %v", module)
	if err := uc.Store.Update(context.TODO(), id, "", &models.Module{}, &module); err != nil {
		c.JSON(http.StatusOK, web.ExceptResponse(errorMap[ErrUpdateModuleFailed], err))
		return
	}
	c.JSON(http.StatusOK, web.DataResponse(module))
}

// Get 获取模块详情
func (uc *ModuleController) Get(c *gin.Context) {
	id, err := view.GetID(c)
	if err != nil {
		c.JSON(http.StatusOK, web.ExceptResponse(errorMap[ErrID], err))
		return
	}
	var module models.Module
	if err := uc.Store.Get(context.TODO(), id, "", &module); err != nil {
		c.JSON(http.StatusOK, web.ExceptResponse(errorMap[ErrGetModuleFailed], err))
		return
	}
	c.JSON(http.StatusOK, web.DataResponse(module))
}

// List 获取模块列表，可根据父模块id和level进行过滤
func (uc *ModuleController) List(c *gin.Context) {
	var req web.Request
	c.ShouldBindQuery(&req)
	req.Default()
	log.Debugf("list modules: %+v", req)
	var condition string
	if req.ParentID != 0 {
		condition = "parent_id = " + fmt.Sprint(req.ParentID)
		log.Debugw("list modules by parent_id", "condition", condition)
	}
	var modules []models.Module
	total, err := uc.Store.List(context.TODO(), req.Limit, req.Page, condition, &modules)
	if err != nil {
		c.JSON(http.StatusOK, web.ExceptResponse(errorMap[ErrGetModuleListFailed], err))
		return
	}
	c.JSON(http.StatusOK, web.ListResponse(int(total), modules))
}

// RelationObject returns the action and func correspondence
func (uc *ModuleController) RelationObject() map[web.Method]web.HandlerFunc {
	return map[web.Method]web.HandlerFunc{
		web.Create: {Funcs: []gin.HandlerFunc{uc.Create}, Admin: true},
		web.Delete: {Funcs: []gin.HandlerFunc{uc.Delete}, Admin: true},
		web.Update: {Funcs: []gin.HandlerFunc{uc.Update}, Admin: true},
		web.Get:    {Funcs: []gin.HandlerFunc{uc.Get}, Admin: true},
		web.List:   {Funcs: []gin.HandlerFunc{uc.List}, Admin: true},
	}
}
