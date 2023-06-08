package router

import (
	"github.com/gin-gonic/gin"
	"github.com/hex-techs/blade/pkg/utils/storage"
	"github.com/hex-techs/blade/pkg/utils/web"
	"github.com/hex-techs/blade/pkg/view/authentication"
	"github.com/hex-techs/blade/pkg/view/module"
	"github.com/hex-techs/blade/pkg/view/user"
)

func InstallAPI(r *gin.Engine, s *storage.Engine) {
	installAuthn(r, s)
	installUserAPI(r, s)
	installModuleAPI(r, s)
}

func installAuthn(r *gin.Engine, s *storage.Engine) {
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

func installUserAPI(r *gin.Engine, s *storage.Engine) {
	u := web.RestfulAPI{
		PostParameter: "/:id",
	}
	u.Install(r, user.NewUserController(s))
}

func installModuleAPI(r *gin.Engine, s *storage.Engine) {
	u := web.RestfulAPI{
		PostParameter: "/:id",
	}
	u.Install(r, module.NewModuleController(s))
}
