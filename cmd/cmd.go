package cmd

import (
	"github.com/gin-gonic/gin"
	"github.com/hex-techs/blade/pkg/router"
	"github.com/hex-techs/blade/pkg/util/config"
	"github.com/hex-techs/blade/pkg/util/log"
	"github.com/hex-techs/blade/pkg/util/storage"
)

func Run() *gin.Engine {
	if err := config.Load("config.yaml"); err != nil {
		panic(err)
	}
	logger := log.InitLogger()
	defer logger.Sync()

	s := storage.NewEngine(config.Read().DB.Host, config.Read().DB.DB, config.Read().DB.User, config.Read().DB.Password)
	initDB(s)
	if err := initAdmin(s); err != nil {
		log.Fatalf("initialize administrator user error: %v", err)
	}

	r := gin.Default()
	router.InstallAPI(r, s)
	return r
}
