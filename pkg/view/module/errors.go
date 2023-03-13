package module

const (
	// 创建模块失败
	ErrCreateModuleFailed = "create module failed"
	// 删除模块失败
	ErrDeleteModuleFailed = "delete module failed"
	// id错误
	ErrID = "id error"
	// 更新模块失败
	ErrUpdateModuleFailed = "update module failed"
	// 获取模块信息失败
	ErrGetModuleFailed = "get module failed"
	// 获取模块列表失败
	ErrGetModuleListFailed = "get module list failed"
	// 无效的参数
	ErrInvalidParam = "invalid param"
)

var errorMap = map[string]int{
	ErrCreateModuleFailed:  30001,
	ErrDeleteModuleFailed:  30002,
	ErrUpdateModuleFailed:  30003,
	ErrGetModuleFailed:     30004,
	ErrID:                  30005,
	ErrGetModuleListFailed: 30006,
	ErrInvalidParam:        30007,
}
