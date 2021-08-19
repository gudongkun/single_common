package jaeger_log

import (
	"context"
	"github.com/micro/go-micro/v2/server"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

//todo 实现zap接口，zap方式格式日志
//logToSpan 把日志写入jaeger的span中
func logToSpan(ctx context.Context, level string, msg string, fields ...log.Field) {
	span := opentracing.SpanFromContext(ctx)
	fa := make([]log.Field, 0, 2+len(fields))

	fa = append(fa, log.String("event", msg))
	fa = append(fa, log.String("level", level))
	for _, field := range fields {
		fa = append(fa,   field)
	}
	span.LogFields(fa...)
}

func Info(ctx context.Context, msg string, fields ...log.Field) {
	logToSpan(ctx, "info", msg, fields...)
}

func Debug(ctx context.Context, msg string, fields ...log.Field) {
	logToSpan(ctx, "debug", msg, fields...)
}

func Error(ctx context.Context, msg string, fields ...log.Field) {
	logToSpan(ctx, "error", msg, fields...)
}
//NewLogWrapper  类似于go micro的middleware 参见：https://github.com/microhq/go-plugins/tree/master/wrapper
func NewLogWrapper(fn server.HandlerFunc) server.HandlerFunc {
	return func(ctx context.Context, req server.Request, rsp interface{}) error {
		Info(ctx,req.Endpoint()+" begin:", log.Object("request",req.Body()))
		err := fn(ctx, req, rsp)
		if err != nil {
			Error(ctx,req.Endpoint()+" end_err:", log.Object("error",err))
		} else {
			Info(ctx,req.Endpoint()+" end:", log.Object("response",rsp))
		}

		return err
	}
}