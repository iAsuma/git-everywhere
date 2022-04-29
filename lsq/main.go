package main

import (
	_ "git-auto/internal/packed"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/glog"
)

func main() {
	ctx := gctx.New()
	defer func() {
		if exception := recover(); exception != nil {
			if err, ok := exception.(error); ok {
				glog.Print(ctx, err.Error())
			} else {
				panic(exception)
			}
		}
	}()
}
