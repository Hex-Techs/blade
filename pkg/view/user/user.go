package user

import (
	"context"
	"net/http"

	"github.com/fize/go-ext/log"
	"github.com/gin-gonic/gin"
	"github.com/hex-techs/blade/pkg/models"
	"github.com/hex-techs/blade/pkg/utils/storage"
	"github.com/hex-techs/blade/pkg/utils/web"
	"github.com/hex-techs/blade/pkg/view"
)

// UserController user controller
type UserController struct {
	web.DefaultController
	Store *storage.Engine
}

// NewUserController return a new user controller
func NewUserController(s *storage.Engine) web.RestController {
	return &UserController{
		Store: s,
	}
}

// 资源名
func (*UserController) Name() string {
	return "user"
}

// Create 创建新用户，只有管理员才能创建
func (uc *UserController) Create() (gin.HandlerFunc, error) {
	return func(c *gin.Context) {
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
	}, nil
}

// Delete 删除用户，只有管理员才能删除
func (uc *UserController) Delete() (gin.HandlerFunc, error) {
	return func(c *gin.Context) {
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
	}, nil
}

// Update 管理员普通用户都可以
func (uc *UserController) Update() (gin.HandlerFunc, error) {
	return func(c *gin.Context) {
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
		// 刷新用户信息
		if err := user.GenUser(); err != nil {
			c.JSON(http.StatusOK, web.ExceptResponse(errorMap[ErrGenerateUserToken], err))
			return
		}
		user.TruncatePassword()
		c.JSON(http.StatusOK, web.DataResponse(user))
	}, nil
}

// Get 获取用户详情
func (uc *UserController) Get() (gin.HandlerFunc, error) {
	return func(c *gin.Context) {
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
	}, nil
}

// List 获取用户列表
func (uc *UserController) List() (gin.HandlerFunc, error) {
	return func(c *gin.Context) {
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
	}, nil
}
