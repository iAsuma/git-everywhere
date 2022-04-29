package cmd

import (
	"context"
	"fmt"
	"git-auto/internal/controller"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"

	"github.com/gogf/gf/v2/os/gcmd"
)

var (
	Main = &gcmd.Command{
		Name:  "main",
		Usage: "main COMMAND [OPTION]",
		Brief: "start http server",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			s := g.Server()
			s.Group("/", func(group *ghttp.RouterGroup) {
				group.Middleware(ghttp.MiddlewareHandlerResponse)
				group.Bind(
					controller.Hello,
				)
			})
			s.Run()
			return nil
		},
	}
	Http = &gcmd.Command{
		Name:          "http",
		Description:   "",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			return nil
		},
	}
	Cron = &gcmd.Command{
		Name:          "cron",
		Usage: "cron COMMAND [OPTION]",
		Description:   "",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			fmt.Println("mihasp")
			return nil
		},
	}
)
