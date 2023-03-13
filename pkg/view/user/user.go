package user

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hex-techs/blade/pkg/models"
	"github.com/hex-techs/blade/pkg/util/log"
	"github.com/hex-techs/blade/pkg/util/storage"
	"github.com/hex-techs/blade/pkg/util/web"
	"github.com/hex-techs/blade/pkg/view"
)

// UserController user controller
type UserController struct {
	Store *storage.Engine
}

// NewUserController return a new user controller
func NewUserController(s *storage.Engine) web.RestController {
	return &UserController{
		Store: s,
	}
}

// 版本号
func (*UserController) Version() string {
	return "v1"
}

// 资源名
func (*UserController) Name() string {
	return "user"
}

// Create 创建新用户，只有管理员才能创建
func (uc *UserController) Create(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, web.ExceptResponse(errorMap[ErrInvalidParam], err))
		return
	}
	log.Debugf("create user: %v", user)
	if err := uc.Store.Create(context.TODO(), &user); err != nil {
		c.JSON(http.StatusOK, web.ExceptResponse(errorMap[ErrCreateUserFailed], err))
		return
	}
	c.JSON(http.StatusOK, web.OkResponse())
}

// Delete 删除用户，只有管理员才能删除
func (uc *UserController) Delete(c *gin.Context) {
	id, err := view.GetID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, web.ExceptResponse(errorMap[ErrInvalidParam], err))
		return
	}
	u := web.GetCurrentUser(c)
	// 不能删除自己
	if u.ID == id {
		c.JSON(http.StatusOK, web.ExceptResponse(errorMap[ErrDeleteSelf], ErrDeleteSelf))
		return
	}
	if err := uc.Store.Delete(context.TODO(), id, "", &models.User{}); err != nil {
		c.JSON(http.StatusOK, web.ExceptResponse(errorMap[ErrDeleteUserFailed], err))
		return
	}
	c.JSON(http.StatusOK, web.OkResponse())
}

// Update 管理员普通用户都可以
func (uc *UserController) Update(c *gin.Context) {
	id, err := view.GetID(c)
	if err != nil {
		c.JSON(http.StatusOK, web.ExceptResponse(errorMap[ErrID], err))
		return
	}
	var user models.User
	c.ShouldBindJSON(&user)
	u := web.GetCurrentUser(c)
	log.Debugw("update user", "id", id, "user", user, "current user", u)
	if !u.Admin && u.ID != id {
		c.JSON(http.StatusOK, web.ExceptResponse(errorMap[ErrUpdateOther], ErrUpdateOther))
		return
	}
	if err := uc.Store.Update(context.TODO(), id, "", &models.User{}, &user); err != nil {
		c.JSON(http.StatusOK, web.ExceptResponse(errorMap[ErrUpdateUserFailed], err))
		return
	}
	// flush user info and token
	if err := user.GenUser(); err != nil {
		c.JSON(http.StatusOK, web.ExceptResponse(errorMap[ErrGenerateUserToken], err))
		return
	}
	user.TruncatePassword()
	c.JSON(http.StatusOK, web.DataResponse(user))
}

// Get 获取用户详情
func (uc *UserController) Get(c *gin.Context) {
	id, err := view.GetID(c)
	if err != nil {
		c.JSON(http.StatusOK, web.ExceptResponse(errorMap[ErrID], err))
		return
	}
	u := web.GetCurrentUser(c)
	if !u.Admin && u.ID != id {
		c.JSON(http.StatusOK, web.ExceptResponse(errorMap[ErrGetOther], ErrGetOther))
		return
	}
	var user models.User
	if err := uc.Store.Get(context.TODO(), id, "", &user); err != nil {
		c.JSON(http.StatusOK, web.ExceptResponse(errorMap[ErrGetUserFailed], err))
		return
	}
	user.TruncatePassword()
	c.JSON(http.StatusOK, web.DataResponse(user))
}

// List 获取用户列表
func (uc *UserController) List(c *gin.Context) {
	var req web.Request
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, web.ExceptResponse(errorMap[ErrInvalidParam], err))
		return
	}
	req.Default()
	log.Debugf("list user: %v", req)
	var users []models.User
	total, err := uc.Store.List(context.TODO(), req.Limit, req.Page, "", &users)
	if err != nil {
		c.JSON(http.StatusOK, web.ExceptResponse(errorMap[ErrGetUserListFailed], err))
		return
	}
	for i, u := range users {
		u.TruncatePassword()
		users[i] = u
	}
	c.JSON(http.StatusOK, web.ListResponse(int(total), users))
}

// RelationObject returns the action and func correspondence
func (uc *UserController) RelationObject() map[web.Method]web.HandlerFunc {
	return map[web.Method]web.HandlerFunc{
		web.Create: {Funcs: []gin.HandlerFunc{uc.Create}, Admin: true},
		web.Delete: {Funcs: []gin.HandlerFunc{uc.Delete}, Admin: true},
		web.Update: {Funcs: []gin.HandlerFunc{uc.Update}},
		web.Get:    {Funcs: []gin.HandlerFunc{uc.Get}},
		web.List:   {Funcs: []gin.HandlerFunc{uc.List}, Admin: true},
	}
}
