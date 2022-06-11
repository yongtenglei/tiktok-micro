package mysql

import (
	"fmt"
	usermodel "yunyandz.com/tiktok/user-part/model"
	"yunyandz.com/tiktok/user-part/settings"
	videomodel "yunyandz.com/tiktok/video-part/model"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB
var err error

func Init() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		settings.ServiceConf.MySQLConf.User,
		settings.ServiceConf.MySQLConf.Password,
		settings.ServiceConf.MySQLConf.Host,
		settings.ServiceConf.MySQLConf.Port,
		settings.ServiceConf.MySQLConf.DbName)
	println(dsn)

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		zap.S().Errorw("Open MySQL failed", "err", err.Error())
		panic(err)
	}

	err = DB.AutoMigrate(&usermodel.User{}, &videomodel.Comment{}, &videomodel.Video{})
	if err != nil {
		zap.S().Errorw("AutoMigrate model failed", "err", err.Error())
		panic(err)
	}
}
