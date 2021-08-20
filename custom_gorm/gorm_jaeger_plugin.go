package custom_gorm

import (
	"github.com/opentracing/opentracing-go"
	tracerLog "github.com/opentracing/opentracing-go/log"
	"gorm.io/gorm"
)

const (
	gormSpanKey        = "__gorm_span"
	callBackBeforeName = "opentracing:before"
	callBackAfterName  = "opentracing:after"
)

func beforeSql(db *gorm.DB) {
	// 先从父级spans生成子span
	span, _ := opentracing.StartSpanFromContext(db.Statement.Context, "gorm")
	// 利用db实例去传递span
	db.InstanceSet(gormSpanKey, span)
	return
}

func afterSql(db *gorm.DB) {
	// 从GORM的DB实例中取出span
	_span, isExist := db.InstanceGet(gormSpanKey)
	if !isExist {
		return
	}

	// 断言进行类型转换
	span, ok := _span.(opentracing.Span)
	if !ok {
		return
	}
	defer span.Finish()

	// Error
	if db.Error != nil {
		span.LogFields(tracerLog.Error(db.Error))
	}

	// sql
	span.LogFields(tracerLog.String("sql", db.Dialector.Explain(db.Statement.SQL.String(), db.Statement.Vars...)))
	return
}

type OpentracingPlugin struct{}

func (op *OpentracingPlugin) Name() string {
	return "opentracingPlugin"
}

func (op *OpentracingPlugin) Initialize(db *gorm.DB) (err error) {
	// 开始前
	db.Callback().Create().Before("gorm:before_create").Register(callBackBeforeName, beforeSql)
	db.Callback().Query().Before("gorm:query").Register(callBackBeforeName, beforeSql)
	db.Callback().Delete().Before("gorm:before_delete").Register(callBackBeforeName, beforeSql)
	db.Callback().Update().Before("gorm:setup_reflect_value").Register(callBackBeforeName, beforeSql)
	db.Callback().Row().Before("gorm:row").Register(callBackBeforeName, beforeSql)
	db.Callback().Raw().Before("gorm:raw").Register(callBackBeforeName, beforeSql)

	// 结束后
	db.Callback().Create().After("gorm:after_create").Register(callBackAfterName, afterSql)
	db.Callback().Query().After("gorm:after_query").Register(callBackAfterName, afterSql)
	db.Callback().Delete().After("gorm:after_delete").Register(callBackAfterName, afterSql)
	db.Callback().Update().After("gorm:after_update").Register(callBackAfterName, afterSql)
	db.Callback().Row().After("gorm:row").Register(callBackAfterName, afterSql)
	db.Callback().Raw().After("gorm:raw").Register(callBackAfterName, afterSql)
	return
}

var _ gorm.Plugin = &OpentracingPlugin{}
