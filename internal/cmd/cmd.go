package cmd

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/gogf/gf/v2/util/gtag"
)

const (
	DataDir = "data"
	ResDir  = "res"
)

var LSQ = cLSQ{}

type cLSQ struct {
	g.Meta `name:"lsq-ci" ad:"{cLsqCiAd}"`
}

const (
	cLsqCiAd = `lsq-ci COMMAND is making ci/cd for project`
)

func init() {
	gtag.Sets(g.MapStrStr{
		"cLsqCiAd": cLsqCiAd,
	})
}

type cLsqCiInput struct {
	g.Meta `name:"lsq-ci"`
}

type cLsqCiOutput struct{}

func (c cLSQ) Index(ctx context.Context, in cLsqCiInput) (out *cLsqCiOutput, err error) {
	fmt.Println("** lsq-cli Command here **")
	gcmd.CommandFromCtx(ctx).Print()
	return
}
