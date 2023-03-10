package router

import (
	"github.com/gin-gonic/gin"
	"github.com/hex-techs/blade/pkg/util/storage"
	"github.com/hex-techs/blade/pkg/view/authentication"
)

func InstallAuthn(r *gin.Engine, s *storage.Engine) {
	api := authentication.NewAuthn(s)
	group := r.Group("/api/v1/auth")
	{
		group.POST("/register", api.Register)
		group.POST("/login", api.Login)
		group.POST("/restpasswordrequest", api.ResetPasswordRequest)
		group.PUT("/resetpassword/:token", api.ResetPassword)
		group.PUT("/changepassword", api.ResetPassword)
	}
}
