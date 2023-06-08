package web

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
)

// RestController restful风格的控制器
type RestController interface {
	// Create is the method for router.POST
	Create() (gin.HandlerFunc, error)
	// Delete is the method for router.DELETE
	Delete() (gin.HandlerFunc, error)
	// Update is the method for router.PUT
	Update() (gin.HandlerFunc, error)
	// Patch is the method for router.PATCH
	Patch() (gin.HandlerFunc, error)
	// Get is the method for router.GET
	Get() (gin.HandlerFunc, error)
	// List is the method for router.GET with query parameters
	List() (gin.HandlerFunc, error)
	// 当前api版本号
	Version() string
	// 当前资源名称
	Name() string
}

// basicAPIGroup is the basic api group
func basicAPIGroup(e *gin.Engine) *gin.RouterGroup {
	return e.Group("/api")
}

// RestfulAPI restful api struct
type RestfulAPI struct {
	// the path and longpath for current resource
	path     string
	longpath string
	// 路由前缀
	PreParameter string
	// 路由后缀
	PostParameter string
}

// Install 装载api
func (r *RestfulAPI) Install(e *gin.Engine, rc RestController) {
	versionAPIGroup := basicAPIGroup(e).Group("/" + rc.Version())
	r.handleParameter(rc)
	if post, err := rc.Create(); err == nil {
		versionAPIGroup.POST(r.path, post)
	}
	if del, err := rc.Delete(); err == nil {
		versionAPIGroup.DELETE(r.longpath, del)
	}
	if put, err := rc.Update(); err == nil {
		versionAPIGroup.PUT(r.longpath, put)
	}
	if patch, err := rc.Patch(); err == nil {
		versionAPIGroup.PATCH(r.longpath, patch)
	}
	if get, err := rc.Get(); err == nil {
		versionAPIGroup.GET(r.longpath, get)
	}
	if list, err := rc.List(); err == nil {
		versionAPIGroup.GET(r.path, list)
	}
}

func (r *RestfulAPI) handleParameter(rc RestController) {
	if r.PreParameter != "" {
		r.path = fmt.Sprintf("/%s/%s", r.PreParameter, rc.Name())
	} else {
		r.path = fmt.Sprintf("/%s", rc.Name())
	}
	if r.PostParameter != "" {
		r.longpath = fmt.Sprintf("%s/%s", r.path, r.PostParameter)
	} else {
		r.longpath = r.path
	}
}

// ErrUnimplemented is the error for unimplemented method
var ErrUnimplemented error = errors.New("Unimplemented")

// DefaultController is the default interface for restful api.
// You can use it to composite your own interface.
type DefaultController struct{}

// Create is the method for router.POST
func (d *DefaultController) Create() (gin.HandlerFunc, error) {
	return nil, ErrUnimplemented
}

// Delete is the method for router.DELETE
func (d *DefaultController) Delete() (gin.HandlerFunc, error) {
	return nil, ErrUnimplemented
}

// Update is the method for router.PUT
func (d *DefaultController) Update() (gin.HandlerFunc, error) {
	return nil, ErrUnimplemented
}

// Patch is the method for router.PATCH
func (d *DefaultController) Patch() (gin.HandlerFunc, error) {
	return nil, ErrUnimplemented
}

// Get is the method for router.GET
func (d *DefaultController) Get() (gin.HandlerFunc, error) {
	return nil, ErrUnimplemented
}

// List is the method for router.GET with query parameters
func (d *DefaultController) List() (gin.HandlerFunc, error) {
	return nil, ErrUnimplemented
}

// Version return the restful API version
func (d *DefaultController) Version() string {
	return "v1"
}

// Name return the restful API name
func (d *DefaultController) Name() string {
	return "blade"
}
