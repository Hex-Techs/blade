package module

import (
	"context"
	"fmt"
	"net/http"

	"github.com/fize/go-ext/log"
	"github.com/gin-gonic/gin"
	"github.com/hex-techs/blade/pkg/models"
	"github.com/hex-techs/blade/pkg/utils/storage"
	"github.com/hex-techs/blade/pkg/utils/web"
	"github.com/hex-techs/blade/pkg/view"
)

// ModuleController module controller
type ModuleController struct {
	web.DefaultController
	Store *storage.Engine
}

// NewModuleController return a new module controller
func NewModuleController(s *storage.Engine) web.RestController {
	return &ModuleController{
		Store: s,
	}
}

// 资源名
func (*ModuleController) Name() string {
	return "module"
}

// Create 创建新模块
func (mc *ModuleController) Create() (gin.HandlerFunc, error) {
	return func(c *gin.Context) {
		var module models.Module
		if err := c.ShouldBindJSON(&module); err != nil {
			c.JSON(http.StatusBadRequest, web.ExceptResponse(errorMap[ErrInvalidParam], err))
			return
		}
		log.Debugf("create module: %v", module)
		if err := mc.Store.Create(context.TODO(), &module); err != nil {
			c.JSON(http.StatusOK, web.ExceptResponse(errorMap[ErrCreateModuleFailed], err))
			return
		}
		c.JSON(http.StatusOK, web.OkResponse())
	}, nil
}

// Delete 删除模块
func (mc *ModuleController) Delete() (gin.HandlerFunc, error) {
	return func(c *gin.Context) {
		id, err := view.GetID(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, web.ExceptResponse(errorMap[ErrID], err))
			return
		}
		if err := mc.delete(id); err != nil {
			c.JSON(http.StatusOK, web.ExceptResponse(errorMap[ErrDeleteModuleFailed], err))
			return
		}
		c.JSON(http.StatusOK, web.OkResponse())
	}, nil
}

// 递归删除parent_id为id的所有模块
func (mc *ModuleController) delete(id uint) error {
	var modules []models.Module
	_, err := mc.Store.List(context.TODO(), -1, 1, "parent_id = "+fmt.Sprint(id), &modules)
	if err != nil {
		return err
	}
	for _, module := range modules {
		if err := mc.delete(module.ID); err != nil {
			return err
		}
	}
	return mc.Store.ForceDelete(context.TODO(), id, "", &models.Module{})
}

// Update 更新模块信息
func (mc *ModuleController) Update() (gin.HandlerFunc, error) {
	return func(c *gin.Context) {
		id, err := view.GetID(c)
		if err != nil {
			c.JSON(http.StatusOK, web.ExceptResponse(errorMap[ErrID], err))
			return
		}
		var (
			new models.Module
			old models.Module
		)
		c.ShouldBindJSON(&new)
		if err := mc.Store.Get(context.TODO(), id, "", &old); err != nil {
			c.JSON(http.StatusOK, web.ExceptResponse(errorMap[ErrUpdateModuleFailed], err))
			return
		}
		if old.Description == new.Description {
			log.Debugf("module %d description not changed", id)
			c.JSON(http.StatusOK, web.OkResponse())
			return
		}
		log.Debugw("update module", "old", old.Description, "new", new.Description)
		old.Description = new.Description
		if err := mc.Store.Update(context.TODO(), id, "", &models.Module{}, &old); err != nil {
			c.JSON(http.StatusOK, web.ExceptResponse(errorMap[ErrUpdateModuleFailed], err))
			return
		}
		c.JSON(http.StatusOK, web.OkResponse())
	}, nil
}

// Get 获取模块详情
func (mc *ModuleController) Get() (gin.HandlerFunc, error) {
	return func(c *gin.Context) {
		id, err := view.GetID(c)
		if err != nil {
			c.JSON(http.StatusOK, web.ExceptResponse(errorMap[ErrID], err))
			return
		}
		var module models.Module
		if err := mc.Store.Get(context.TODO(), id, "", &module); err != nil {
			c.JSON(http.StatusOK, web.ExceptResponse(errorMap[ErrGetModuleFailed], err))
			return
		}
		c.JSON(http.StatusOK, web.DataResponse(module))
	}, nil
}

// List 获取模块列表，可根据父模块id和level进行过滤
func (mc *ModuleController) List() (gin.HandlerFunc, error) {
	return func(c *gin.Context) {
		var req web.Request
		c.ShouldBindQuery(&req)
		req.Default()
		log.Debugf("list modules: %+v", req)
		if req.Level == 0 && req.ParentID == 0 {
			c.JSON(http.StatusBadRequest, web.ExceptResponse(errorMap[ErrInvalidParam],
				"level and parentID can't be 0 at the same time"))
			return
		}
		var condition string
		if req.Level != 0 {
			condition = "level = " + fmt.Sprint(req.Level)
			log.Debugw("list modules by level", "condition", condition)
		}
		// 查询条件中，parent_id比level有更高的优先级
		if req.ParentID != 0 {
			condition = "parent_id = " + fmt.Sprint(req.ParentID)
			log.Debugw("list modules by parentID", "condition", condition)
		}
		var modules []models.Module
		total, err := mc.Store.List(context.TODO(), req.Limit, req.Page, condition, &modules)
		if err != nil {
			c.JSON(http.StatusOK, web.ExceptResponse(errorMap[ErrGetModuleListFailed], err))
			return
		}
		c.JSON(http.StatusOK, web.ListResponse(int(total), modules))
	}, nil
}

func (uc *ModuleController) Middlewares() []web.MiddlewaresObject {
	return []web.MiddlewaresObject{
		{
			Methods:     []string{web.CREATE, web.DELETE},
			Middlewares: []gin.HandlerFunc{web.LoginRequired(), web.AdminRequired()},
		},
		{
			Methods:     []string{web.UPDATE, web.GET, web.LIST},
			Middlewares: []gin.HandlerFunc{web.LoginRequired()},
		},
	}
}
