package config

import (
	"fmt"
	"sync"

	ext "github.com/fize/go-ext/config"
	"github.com/kkyr/fig"
)

// default configuration
const (
	// 默认端口
	_defaultServerPort = 8080
	// 默认token超时时间24小时
	_defaultTokenExpired = 3600 * 24
	// 默认url超时时间10分钟
	_defaultURLExpired = 600
	// 默认管理员密码
	_defaultPassword = "admin"
)

// 服务配置
type ServiceConfig struct {
	// 服务端口
	ServerPort int `fig:"serverPort"`
	// 域名
	Domain string `fig:"domain"`
	// token过期时间
	TokenExpired int64 `fig:"tokenExpired"`
	// url过期时间
	URLExpired int64 `fig:"urlExpired"`
	// 重置密码的path
	ResetPath string `fig:"resetPath"`
	// 管理员密码
	AdminPassword string `fig:"adminPassword"`
	// 跨域相关配置
	Cors        bool     `fig:"cors"`
	AllowOrigin []string `fig:"allowOrigin"`
	// 公司邮箱后缀，只允许该后缀的邮箱注册
	Company string `fig:"company"`
	// 是否开启文档服务
	APIDoc bool `fig:"apiDoc"`
}

// ldap配置
type Ldap struct {
	// 是否开启ldap
	Enabled bool `fig:"enabled"`
	// ldap 地址 127.0.0.1:389
	Host string `fig:"host"`
	// baseDN ou=Users,dc=blade,dc=cn
	BaseDN string `fig:"baseDN"`
	// cn=admin,dc=blade,dc=cn
	User string `fig:"user"`
	// 密码
	Password string `fig:"password"`
}

// 全局配置
type Config struct {
	ext.Config
	// 服务配置
	Service *ServiceConfig `fig:"service"`
	// ldap配置
	Ldap *Ldap `fig:"ldap"`
}

// 配置内容
var (
	lock   *sync.RWMutex
	config *Config
)

// Load 加载配置，只在程序启动时加载一次
func Load(dir, name string) error {
	lock = new(sync.RWMutex)
	lock.Lock()
	defer lock.Unlock()
	config = new(Config)

	if err := ext.Load(dir, name); err != nil {
		return fmt.Errorf("load base config error: %s", err)
	}

	if err := fig.Load(config, fig.Dirs(dir), fig.File(name)); err != nil {
		return err
	}

	config.DB = ext.Read().DB
	config.Email = ext.Read().Email

	config.Ldap = new(Ldap)

	// 设置默认端口
	if config.Service == nil {
		config.Service = new(ServiceConfig)
	}
	if config.Service.ServerPort == 0 {
		config.Service.ServerPort = _defaultServerPort
	}
	// 设置默认token超时时间
	if config.Service.TokenExpired == 0 {
		config.Service.TokenExpired = _defaultTokenExpired
	}
	// 设置默认url超时时间
	if config.Service.URLExpired == 0 {
		config.Service.URLExpired = _defaultURLExpired
	}
	// 设置默认管理员密码
	if config.Service.AdminPassword == "" {
		config.Service.AdminPassword = _defaultPassword
	}
	return nil
}

// Read 读取配置
func Read() *Config {
	return config
}
