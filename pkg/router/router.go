package router

import (
	"github.com/gin-gonic/gin"
	"github.com/hex-techs/blade/pkg/util/storage"
	"github.com/hex-techs/blade/pkg/util/web"
	"github.com/hex-techs/blade/pkg/view/authentication"
	"github.com/hex-techs/blade/pkg/view/user"
)

func InstallAuthn(r *gin.Engine, s *storage.Engine) {
	api := authentication.NewAuthn(s)
	group := r.Group("/api/v1/auth")
	{
		group.POST("/register", api.Register)
		group.POST("/login", api.Login)
		group.POST("/restpasswordrequest", api.ResetPasswordRequest)
		group.PUT("/resetpassword/:token", api.ResetPassword)
		group.PUT("/changepassword", web.LoginRequired(), api.ChangePassword)
	}
}

func InstallUserAPI(r *gin.Engine, s *storage.Engine) {
	u := web.RestfulAPI{
		PostParameter: "/:id",
	}
	u.Install(r, user.NewUserController(s))
}
