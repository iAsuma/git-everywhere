package main

import (
	"lsq-cli/internal/cmd"
	_ "lsq-cli/internal/packed"

	"github.com/gogf/gf/v2/os/gcmd"
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

	command, err := gcmd.NewFromObject(cmd.LSQ)
	if err != nil {
		panic(err)
	}

	err = command.AddObject(
		cmd.Ci,
	)

	if err != nil {
		panic(err)
	}

	command.Run(ctx)
}
