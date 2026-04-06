package common

import (
	"context"
	"fmt"
	"os"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

func RunHertz(serviceName, instanceName, port string, echoAll bool) {
	h := server.Default(
		server.WithHostPorts(fmt.Sprintf(":%s", port)),
	)

	h.GET("/ping", func(c context.Context, ctx *app.RequestContext) {
		ctx.JSON(consts.StatusOK, map[string]string{
			"service":  serviceName,
			"instance": instanceName,
			"message":  "pong",
		})
	})

	h.GET("/info", func(c context.Context, ctx *app.RequestContext) {
		ctx.JSON(consts.StatusOK, buildHertzResponse(serviceName, instanceName, ctx))
	})

	if echoAll {
		h.Any("/*path", func(c context.Context, ctx *app.RequestContext) {
			ctx.JSON(consts.StatusOK, buildHertzResponse(serviceName, instanceName, ctx))
		})
	}

	h.Spin()
}

func EnvOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}

func buildHertzResponse(service, instance string, ctx *app.RequestContext) InfoResponse {
	return BuildRequestInfoResponse(
		service,
		instance,
		string(ctx.Method()),
		string(ctx.Path()),
		string(ctx.Host()),
		ctx.ClientIP(),
		map[string]string{
			"X-Forwarded-For":   string(ctx.GetHeader("X-Forwarded-For")),
			"X-Forwarded-Host":  string(ctx.GetHeader("X-Forwarded-Host")),
			"X-Forwarded-Proto": string(ctx.GetHeader("X-Forwarded-Proto")),
		},
	)
}
