package config

import (
	"sync"

	"github.com/kkyr/fig"
)

// default configuration
const (
	// 默认的数据库类型
	_defaultDBType = "sqlite"
	// 默认的数据库文件
	_defaultDB = "./blade.db"
	// 默认端口
	_defaultServerPort = 8080
	// 默认token超时时间24小时
	_defaultTokenExpired = 3600 * 24
	// 默认url超时时间10分钟
	_defaultURLExpired = 600
	// 默认管理员密码
	_defaultPassword = "admin"
	// 默认日志文件路径
	_defaultLogPath = "./blade.log"
	// 默认日志级别
	_defaultLogLevel = "info"
	// 默认日志
	_defaultLogMaxSize = 10
	// 默认日志备份数
	_defaultLogMaxBackups = 5
	// 默认日志最大保存时间
	_defaultLogMaxAge = 30
	// 默认日志格式
	_defaultLogFormat = "string"
)

// dbType 数据库类型
type dbType string

const (
	// mysql数据库
	Mysql dbType = "mysql"
	// sqlite数据库
	Sqlite3 dbType = "sqlite"
)

// 数据库配置
type DBConfig struct {
	// 数据库类型，只支持mysql 和 sqlite，默认 sqllite
	Type dbType `fig:"type"`
	// 数据库地址，此部分包含了端口号 127.0.0.1:3306
	Host string `fig:"host"`
	// 用户名
	User string `fig:"user"`
	// 密码
	Password string `fig:"password"`
	// 数据库名
	DB string `fig:"db"`
	// 连接池配置
	MaxIdleConns int `fig:"maxIdleConns"`
	MaxOpenConns int `fig:"maxOpenConns"`
	// 开启详细的sql日志
	SqlDebug bool `fig:"sqlDebug"`
}

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

// 邮件配置
type Email struct {
	// 是否开启邮件发送服务
	Enabled bool `fig:"enabled"`
	// 邮件账号用户名
	Account string `fig:"account"`
	// smtp地址
	SMTP string `fig:"smtp"`
	// smtp端口
	Port int `fig:"port"`
	// 密码
	Password string `fig:"password"`
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

// 日志配置
type Log struct {
	// 日志文件路径
	Filename string `fig:"filename"`
	// 日志文件最大大小，单位MB
	MaxSize int `fig:"maxSize"`
	// 日志文件最大备份数量
	MaxBackups int `fig:"maxBackups"`
	// 日志文件最大保存时间，单位天
	MaxAge int `fig:"maxAge"`
	// 是否压缩
	Compress bool `fig:"compress"`
	// 日志级别
	Level string `fig:"level"`
	// 日志格式
	Format string `fig:"format"`
	// 输出方式
	Output string `fig:"output"`
}

// 全局配置
type Config struct {
	// 服务配置
	Service *ServiceConfig `fig:"service"`
	// 数据库配置
	DB *DBConfig `fig:"db"`
	// 邮件配置
	Email *Email `fig:"email"`
	// ldap配置
	Ldap *Ldap `fig:"ldap"`
	// 日志配置
	Log *Log `fig:"log"`
}

// 配置内容
var (
	lock   *sync.RWMutex
	config *Config
)

// Load 加载配置，只在程序启动时加载一次
func Load(path string) error {
	lock = new(sync.RWMutex)
	lock.Lock()
	defer lock.Unlock()
	config = new(Config)
	err := fig.Load(config, fig.File(path))
	if err != nil {
		return err
	}
	// 设置默认数据库类型为sqlite
	if config.DB == nil {
		config.DB = new(DBConfig)
	}
	if config.DB.Type != Mysql && config.DB.Type != Sqlite3 {
		config.DB.Type = _defaultDBType
	}
	// 当使用sqlite时设置默认数据库文件
	if config.DB.Type == Sqlite3 && config.DB.DB == "" {
		config.DB.DB = _defaultDB
	}
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
	// 设置默认日志配置
	if config.Log == nil {
		config.Log = new(Log)
	}
	if config.Log.Filename == "" {
		config.Log.Filename = _defaultLogPath
	}
	if config.Log.MaxSize == 0 {
		config.Log.MaxSize = _defaultLogMaxSize
	}
	if config.Log.MaxBackups == 0 {
		config.Log.MaxBackups = _defaultLogMaxBackups
	}
	if config.Log.MaxAge == 0 {
		config.Log.MaxAge = _defaultLogMaxAge
	}
	if config.Log.Level == "" {
		config.Log.Level = _defaultLogLevel
	}
	if config.Log.Format != "json" {
		config.Log.Format = _defaultLogFormat
	}
	return nil
}

// Read 读取配置
func Read() *Config {
	return config
}
