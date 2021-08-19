package custom_xorm

import (
	"fmt"
	"github.com/opentracing/opentracing-go"
	tracerLog "github.com/opentracing/opentracing-go/log"
	"github.com/uber/jaeger-client-go/log/zap"
	xormLog "xorm.io/xorm/log"
)

/*
 参考代码https://gitee.com/avtion/xormWithTracing/blob/master/custom_log.go
	下面都是实现自定义ContextLogger的部分，这里使用OpenTracing自带的Zap日志（语法糖裁剪版）
	也可以使用Zap、logrus、原生自带的log实现xorm的ContextLogger接口
*/

type CustomCtxLogger struct {
	logger  *zap.Logger
	level   xormLog.LogLevel
	showSQL bool
	span    opentracing.Span
}

// BeforeSQL implements ContextLogger
func (l *CustomCtxLogger) BeforeSQL(ctx xormLog.LogContext) {
	// ----> 重头戏在这里，需要从Context上下文中创建一个新的Span来对SQL执行进行链路监控
	l.span, _ = opentracing.StartSpanFromContext(ctx.Ctx, "XORM SQL Execute")
}

// AfterSQL implements ContextLogger
func (l *CustomCtxLogger) AfterSQL(ctx xormLog.LogContext) {
	// defer结束掉span
	defer l.span.Finish()
	fa := make([]tracerLog.Field, 0, 3)
	// 原本的SimpleLogger里面会获取一次SessionId
	var sessionPart string
	v := ctx.Ctx.Value("__xorm_session_id")
	if key, ok := v.(string); ok {
		sessionPart = fmt.Sprintf(" [%s]", key)
		fa  =  append(fa,tracerLog.String("session_id", sessionPart))
	}
	fa  =  append(fa,tracerLog.String("SQL", ctx.SQL))
	fa  =  append(fa,tracerLog.Object("args", ctx.Args))

	// 将Ctx中全部的信息写入到Span中
	l.span.LogFields(fa...)

	l.span.SetTag("execute_time", ctx.ExecuteTime)

	//if ctx.ExecuteTime > 0 {
	//	l.logger.Infof("[SQL]%s %s %v - %v", sessionPart, ctx.SQL, ctx.Args, ctx.ExecuteTime)
	//} else {
	//	l.logger.Infof("[SQL]%s %s %v", sessionPart, ctx.SQL, ctx.Args)
	//}
}

// Errorf implement ILogger
func (l *CustomCtxLogger) Errorf(format string, v ...interface{}) {
	l.logger.Error(fmt.Sprintf(format, v...))
	return
}

// Debugf implement ILogger
func (l *CustomCtxLogger) Debugf(format string, v ...interface{}) {
	l.logger.Debugf(format, v...)
	return
}

// Infof implement ILogger
func (l *CustomCtxLogger) Infof(format string, v ...interface{}) {
	l.logger.Infof(format, v...)
	return
}

// Warnf implement ILogger ---> 这里偷懒了，直接用Info代替
func (l *CustomCtxLogger) Warnf(format string, v ...interface{}) {
	l.logger.Infof(format, v...)
	return
}

// Level implement ILogger
func (l *CustomCtxLogger) Level() xormLog.LogLevel {
	return l.level
}

// SetLevel implement ILogger
func (l *CustomCtxLogger) SetLevel(lv xormLog.LogLevel) {
	l.level = lv
	return
}

// ShowSQL implement ILogger
func (l *CustomCtxLogger) ShowSQL(show ...bool) {
	if len(show) == 0 {
		l.showSQL = true
		return
	}
	l.showSQL = show[0]
}

// IsShowSQL implement ILogger
func (l *CustomCtxLogger) IsShowSQL() bool {
	return l.showSQL
}