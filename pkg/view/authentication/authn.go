package authentication

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/fize/go-ext/log"
	"github.com/fize/go-ext/sendmail"
	"github.com/gin-gonic/gin"
	"github.com/hex-techs/blade/pkg/models"
	"github.com/hex-techs/blade/pkg/utils/config"
	"github.com/hex-techs/blade/pkg/utils/storage"
	"github.com/hex-techs/blade/pkg/utils/token"
	"github.com/hex-techs/blade/pkg/utils/web"
)

// Authn 认证结构体
type Authn struct {
	Store *storage.Engine
}

func NewAuthn(s *storage.Engine) *Authn {
	return &Authn{Store: s}
}

// 登录
func (a *Authn) Login(c *gin.Context) {
	var f LoginForm
	var user models.User
	if err := c.ShouldBindJSON(&f); err != nil {
		c.JSON(http.StatusBadRequest, web.ExceptResponse(errorMap[ErrInvalidParam], err))
		return
	}
	log.Debugf("login user: %v", f)
	if err := a.Store.Get(context.TODO(), 0, f.Name, &user); err != nil {
		if err.Error() != "record not found" {
			c.JSON(http.StatusOK, web.ExceptResponse(errorMap[ErrOther], err))
			return
		}
		c.JSON(http.StatusOK, web.ExceptResponse(errorMap[ErrAccountNotFound], ErrAccountNotFound))
		return
	}
	if user.ValidatePassword(f.Password) {
		if err := user.GenUser(); err != nil {
			c.JSON(http.StatusOK, web.ExceptResponse(errorMap[ErrGenerateToken], err))
			return
		}
	} else {
		c.JSON(http.StatusOK, web.ExceptResponse(errorMap[ErrPasswordInvalid], ErrPasswordInvalid))
		return
	}
	user.TruncatePassword()
	c.JSON(http.StatusOK, web.DataResponse(user))
}

// 注册用户
func (a *Authn) Register(c *gin.Context) {
	var f RegisterForm
	if err := c.ShouldBindJSON(&f); err != nil {
		c.JSON(http.StatusBadRequest, web.ExceptResponse(errorMap[ErrInvalidParam], err))
		return
	}
	log.Debugf("register user: %v", f)
	// 判断邮箱是否可以注册
	l := strings.Split(f.Email, "@")
	company := l[1]
	if config.Read().Service.Company != "" {
		if company != config.Read().Service.Company {
			c.JSON(http.StatusOK, web.ExceptResponse(errorMap[ErrEmailNotAllowed], ErrEmailNotAllowed))
			return
		}
	}
	user := models.User{
		Name:     f.Name,
		Email:    f.Email,
		CnName:   f.CnName,
		Password: f.Password,
		Phone:    f.Phone,
		IM:       f.IM,
	}
	user.EncodePasswd()
	if err := a.Store.Create(context.TODO(), &user); err != nil {
		c.JSON(http.StatusOK, web.ExceptResponse(errorMap[ErrRegisterFailed], err))
	} else {
		c.JSON(http.StatusOK, web.OkResponse())
	}
}

// 修改密码
func (a *Authn) ChangePassword(c *gin.Context) {
	var f ChangePasswordForm
	var user models.User
	if err := c.ShouldBindJSON(&f); err != nil {
		c.JSON(http.StatusBadRequest, web.ExceptResponse(errorMap[ErrInvalidParam], err))
		return
	}
	u := web.GetCurrentUser(c)
	log.Debugw("change password", "user", u, "form", f)
	if err := a.Store.Get(context.TODO(), u.ID, u.Name, &user); err != nil {
		c.JSON(http.StatusOK, web.ExceptResponse(errorMap[ErrChangePasswordFailed], err))
		return
	}
	if user.ValidatePassword(f.OldPassword) {
		user.Password = f.NewPassword
		user.EncodePasswd()
		if err := a.Store.Update(context.TODO(), user.ID, user.Name, &user, &user); err != nil {
			c.JSON(http.StatusOK, web.ExceptResponse(errorMap[ErrChangePasswordFailed], err))
			return
		}
	} else {
		c.JSON(http.StatusOK, web.ExceptResponse(errorMap[ErrPasswordInvalid], ErrPasswordInvalid))
		return
	}
	c.JSON(http.StatusOK, web.OkResponse())
}

// 请求重置密码
func (a *Authn) ResetPasswordRequest(c *gin.Context) {
	var f ForgetPasswordForm
	var user models.User
	if err := c.ShouldBindJSON(&f); err != nil {
		c.JSON(http.StatusBadRequest, web.ExceptResponse(errorMap[ErrInvalidParam], err))
		return
	}
	log.Debugw("reset password request", "user", f)
	if err := a.Store.Get(context.TODO(), 0, f.Name, &user); err != nil {
		c.JSON(http.StatusOK, web.ExceptResponse(errorMap[ErrAccountNotFound], err))
		return
	}
	mc := config.Read().Email
	// if mc != nil && mc.Enabled {
	body := fmt.Sprintf("重置密码链接： %s%s/%s，有效时间%d秒。",
		config.Read().Service.Domain, config.Read().Service.ResetPath,
		token.GenerateCustomToken(f.Name, int(config.Read().Service.URLExpired)), config.Read().Service.URLExpired)
	if err := sendmail.SendEmail("blade", mc.SMTP, mc.Account, mc.Password, user.Email, "Reset Password", body, mc.Port); err != nil {
		c.JSON(http.StatusOK, web.ExceptResponse(errorMap[ErrSendResetEmailFailed], err))
		return
	}
	// } else {
	// 	c.JSON(http.StatusOK, web.DataResponse("未开启邮件服务，请联系管理员重置密码！"))
	// 	return
	// }
	c.JSON(http.StatusOK, web.OkResponse())
}

// 重置密码
func (a *Authn) ResetPassword(c *gin.Context) {
	var f ResetPasswordForm
	if err := c.ShouldBindJSON(&f); err != nil {
		c.JSON(http.StatusBadRequest, web.ExceptResponse(errorMap[ErrInvalidParam], err))
		return
	}
	log.Debugf("reset password: %v", f)
	if f.Token == "" {
		c.JSON(http.StatusOK, web.ExceptResponse(errorMap[ErrResetTokenInvalid], ErrResetTokenInvalid))
		return
	}
	name, err := token.ParseCustomToken(f.Token)
	if err != nil {
		c.JSON(http.StatusOK, web.ExceptResponse(errorMap[ErrResetTokenInvalid], err))
		return
	}
	var u models.User
	if err := a.Store.Get(context.TODO(), 0, name, &u); err != nil {
		c.JSON(http.StatusOK, web.ExceptResponse(errorMap[ErrResetPasswordFailed], err))
		return
	}
	u.Password = f.Password
	u.EncodePasswd()
	if err := a.Store.Update(context.TODO(), u.ID, u.Name, &u, &u); err != nil {
		c.JSON(http.StatusOK, web.ExceptResponse(errorMap[ErrResetPasswordFailed], err))
		return
	}
	c.JSON(http.StatusOK, web.OkResponse())
}
