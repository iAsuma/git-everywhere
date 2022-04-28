package main

import (
	_ "git-auto/internal/packed"

	"github.com/gogf/gf/v2/os/gctx"

	"git-auto/internal/cmd"
)

func main() {
	cmd.Main.Run(gctx.New())
}
