package authentication

// 登录表单
type LoginForm struct {
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// 修改密码表单
type ChangePasswordForm struct {
	// 在重置密码时，不需要旧密码，所以这里不使用required
	OldPassword        string `json:"oldPassword"`
	NewPassword        string `json:"newPassword" binding:"required"`
	NewPasswordConfirm string `json:"newPasswordConfirm" binding:"eqfield=Password"`
	// 在重置密码时，需要token，所以这里不使用required
	Token string `json:"token"`
}

// 注册表单
type RegisterForm struct {
	Email     string `json:"email" binding:"required,email"`
	Name      string `json:"name" binding:"required"`
	Password  string `json:"password" binding:"required"`
	Password2 string `json:"password2" binding:"eqfield=Password"`
	CnName    string `json:"cnName" binding:"required"`
	Phone     string `json:"phone"`
	IM        string `json:"im"`
}
