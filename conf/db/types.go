package db

import (
	"github.com/irisnet/rainbow-sync/constant"
	"github.com/irisnet/rainbow-sync/logger"
	"os"
)

var (
	Addrs    = "localhost:27018"
	User     = "iris"
	Passwd   = "irispassword"
	Database = "sync-iris"
)

// get value of env var
func init() {
	addrs, found := os.LookupEnv(constant.EnvNameDbAddr)
	if found {
		Addrs = addrs
	}

	user, found := os.LookupEnv(constant.EnvNameDbUser)
	if found {
		User = user
	}

	passwd, found := os.LookupEnv(constant.EnvNameDbPassWd)
	if found {
		Passwd = passwd
	}

	database, found := os.LookupEnv(constant.EnvNameDbDataBase)
	if found {
		Database = database
	}

	logger.Debug("init db config", logger.String("addrs", Addrs),
		logger.Bool("userIsEmpty", User == ""), logger.Bool("passwdIsEmpty", Passwd == ""),
		logger.String("database", Database))
}
