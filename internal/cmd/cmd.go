package cmd

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/gogf/gf/v2/util/gtag"
)

var LSQ = cLSQ{}

type cLSQ struct {
	g.Meta `name:"lsq" ad:"{cLSQAd}"`
}

const (
	cLSQAd = `lsq COMMAND is making `
)

func init() {
	gtag.Sets(g.MapStrStr{
		"cLSQAd": cLSQAd,
	})
}

type cLSQInput struct {
	g.Meta `name:"lsq"`
}

type cLSQOutput struct{}

func (c cLSQ) Index(ctx context.Context, in cLSQInput) (out *cLSQOutput, err error) {
	fmt.Println("** this first Command of Asuma **")
	gcmd.CommandFromCtx(ctx).Print()
	return
}
