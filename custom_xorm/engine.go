package custom_xorm

import (
	"context"
	_ "github.com/go-sql-driver/mysql"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/uber/jaeger-client-go/log/zap"
	zap2 "go.uber.org/zap"
	"xorm.io/xorm"
	xormLog "xorm.io/xorm/log"
)

var Engine  *xorm.Engine

func InitEngine( sdn string) (  err error) {
	// XORM创建引擎
	Engine, err = xorm.NewEngine("mysql", sdn)
	if err != nil {
		log.Fatal(err)
		return
	}

	// 创建自定义的日志实例
	_l, err := zap2.NewDevelopment()
	if err != nil {
		return
	}

	// 将日志实例设置到XORM的引擎中
	Engine.SetLogger(&CustomCtxLogger{
		logger:  zap.NewLogger(_l),
		level:   xormLog.LOG_DEBUG,
		showSQL: true,
	})
	return
}

func Db(ctx context.Context)(*xorm.Session)   {
	session := Engine.Context(ctx)
	return session
}
