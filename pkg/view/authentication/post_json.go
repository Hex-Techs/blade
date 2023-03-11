package authentication

// 登录表单
type LoginForm struct {
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// 忘记密码表单
type ForgetPasswordForm struct {
	Name string `json:"name" binding:"required"`
}

// 修改密码表单
type ChangePasswordForm struct {
	OldPassword        string `json:"oldPassword" binding:"required"`
	NewPassword        string `json:"newPassword" binding:"required,nefield=OldPassword,min=6"`
	NewPasswordConfirm string `json:"newPasswordConfirm" binding:"eqfield=NewPassword"`
}

// 重置密码表单
type ResetPasswordForm struct {
	Password        string `json:"password" binding:"required,min=6"`
	PasswordConfirm string `json:"passwordConfirm" binding:"eqfield=Password"`
	Token           string `json:"token" binding:"required"`
}

// 注册表单
type RegisterForm struct {
	Email     string `json:"email" binding:"required,email"`
	Name      string `json:"name" binding:"required"`
	Password  string `json:"password" binding:"required,min=6"`
	Password2 string `json:"password2" binding:"eqfield=Password"`
	CnName    string `json:"cnName" binding:"required"`
	Phone     string `json:"phone"`
	IM        string `json:"im"`
}
