package web

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// 方法类型
type Method string

const (
	Create Method = "create"
	Delete Method = "delete"
	Update Method = "update"
	Get    Method = "get"
	List   Method = "list"
)

// RestController restful风格的控制器
type RestController interface {
	// Create is the method for router.POST
	Create(c *gin.Context)
	// Delete is the method for router.DELETE
	Delete(c *gin.Context)
	// Update is the method for router.PUT
	Update(c *gin.Context)
	// Get is the method for router.GET
	Get(c *gin.Context)
	// List is the method for router.GET with query parameters
	List(c *gin.Context)
	// 当前api版本号
	Version() string
	// 当前资源名称
	Name() string
	// return the relation
	RelationObject() map[Method][]gin.HandlerFunc
}

// basicAPIGroup is the basic api group
func basicAPIGroup(e *gin.Engine) *gin.RouterGroup {
	return e.Group("/api")
}

// RestfulAPI restful api struct
type RestfulAPI struct {
	// 处理路由时，path是中间态，longpath是最终的路由
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
	for method, ops := range rc.RelationObject() {
		r.option(method, versionAPIGroup, ops)
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

func (r *RestfulAPI) option(method Method, group *gin.RouterGroup, funcs []gin.HandlerFunc) {
	switch method {
	case Get:
		r.get(group, funcs)
	case List:
		r.list(group, funcs)
	case Create:
		r.create(group, funcs)
	case Delete:
		r.delete(group, funcs)
	case Update:
		r.update(group, funcs)
	}
}
func (r *RestfulAPI) create(group *gin.RouterGroup, funcs []gin.HandlerFunc) {
	for _, ops := range funcs {
		group.POST(r.longpath, ops)
	}
}
func (r *RestfulAPI) delete(group *gin.RouterGroup, funcs []gin.HandlerFunc) {
	for _, ops := range funcs {
		group.DELETE(r.longpath, ops)
	}
}
func (r *RestfulAPI) update(group *gin.RouterGroup, funcs []gin.HandlerFunc) {
	for _, ops := range funcs {
		group.PUT(r.longpath, ops)
	}
}
func (r *RestfulAPI) get(group *gin.RouterGroup, funcs []gin.HandlerFunc) {
	for _, ops := range funcs {
		group.GET(r.longpath, ops)
	}
}
func (r *RestfulAPI) list(group *gin.RouterGroup, funcs []gin.HandlerFunc) {
	for _, ops := range funcs {
		group.GET(r.path, ops)
	}
}
