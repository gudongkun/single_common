package custom_gorm

import (
	"context"
	log "github.com/micro/go-micro/v2/logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Eng  *gorm.DB

func InitEngine( dsn string)(  err error) {
	Eng, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	_ = Eng.Use(&OpentracingPlugin{})
	return

}

func Db(ctx context.Context)(*gorm.DB)   {
	session := Eng.WithContext(ctx)
	return session
}