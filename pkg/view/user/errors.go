package user

const (
	// 创建用户失败
	ErrCreateUserFailed = "create user failed"
	// 删除用户失败
	ErrDeleteUserFailed = "delete user failed"
	// 不能删除自身
	ErrDeleteSelf = "can not delete yourself"
	// 不能修改其他用户
	ErrUpdateOther = "can not update other user"
	// id错误
	ErrID = "id error"
	// 更新失败
	ErrUpdateUserFailed = "update user failed"
	// 不能获取其他用户信息
	ErrGetOther = "can not get other user info"
	// 获取用户信息失败
	ErrGetUserFailed = "get user failed"
	// 获取用户列表失败
	ErrGetUserListFailed = "get user list failed"
	// 无效的参数
	ErrInvalidParam = "invalid param"
	// 生成token失败
	ErrGenerateUserToken = "generate user token failed"
)

var errorMap = map[string]int{
	ErrCreateUserFailed:  20001,
	ErrDeleteUserFailed:  20002,
	ErrDeleteSelf:        20003,
	ErrUpdateOther:       20004,
	ErrID:                20005,
	ErrUpdateUserFailed:  20006,
	ErrGetOther:          20007,
	ErrGetUserFailed:     20008,
	ErrGetUserListFailed: 20009,
	ErrInvalidParam:      20010,
	ErrGenerateUserToken: 20011,
}
