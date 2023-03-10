package cmd

import (
	"context"

	"github.com/hex-techs/blade/pkg/models"
	"github.com/hex-techs/blade/pkg/util/config"
	"github.com/hex-techs/blade/pkg/util/log"
	"github.com/hex-techs/blade/pkg/util/storage"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func initDB(s *storage.Engine) {
	log.Info("initializing database...")
	db := s.Client().(*gorm.DB)
	// 设置gorm日志模式
	if config.Read().DB.SqlDebug {
		db.Config.Logger = logger.Default.LogMode(logger.Info)
	} else {
		db.Config.Logger = logger.Default.LogMode(logger.Warn)
	}
	// 设置连接池
	if config.Read().DB.MaxIdleConns != 0 {
		c, _ := db.DB()
		c.SetMaxIdleConns(config.Read().DB.MaxIdleConns)
	}
	if config.Read().DB.MaxOpenConns != 0 {
		c, _ := db.DB()
		c.SetMaxOpenConns(config.Read().DB.MaxOpenConns)
	}
	// 自动迁移
	if err := db.AutoMigrate(&models.User{}); err != nil {
		log.Fatalf("auto migrate table error: %v", err)
		return
	}
	log.Info("initialize database ok!")
}

func initAdmin(s *storage.Engine) error {
	u := &models.User{
		Name:     "admin",
		CnName:   "管理员",
		Password: config.Read().Service.AdminPassword,
		Email:    "admin@example.com",
	}
	u.EncodePasswd()
	if err := s.Create(context.TODO(), u); err != nil {
		if err.Error() != "object exist" {
			return err
		}
		return nil
	} else {
		log.Info("initialize administrator user ok!")
	}
	return nil
}
