package global

import (
	"cdbot/src/models/config"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	DATABASE *gorm.DB
	CONFIG   config.ServerConfig
	VIPER    *viper.Viper
	LOGGER   *zap.Logger
)
