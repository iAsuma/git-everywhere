package cmd

import (
	"github.com/gogf/gf/v2/os/gproc"
	"github.com/gogf/gf/v2/text/gregex"
	"github.com/gogf/gf/v2/text/gstr"
	"strings"

	//"bufio"
	"context"
	"fmt"
	"git-auto/utility"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/util/gtag"
)

const (
	cCiGitUsage = "lsq ci git"
	cCiGitBrief = "lss"
)

func init() {
	gtag.Sets(g.MapStrStr{
		"cCiGitUsage": cCiGitUsage,
		"cCiGitBrief": cCiGitBrief,
	})
}

type (
	cCiGitInput struct {
		g.Meta `name:"git" usage:"{cCiGitUsage}" brief:"{cCiGitBrief}"`
	}
	cCiGitOutput struct{}
)

func (c cCi) Git(ctx context.Context, in cCiGitInput) (out cCiGitOutput, err error) {
	fmt.Println(g.Cfg().Get(ctx, "git"))

	rootPath := utility.GetPwd()
	fileName := fmt.Sprintf("%s%s", rootPath, "/res/git-url.lsq")

	target, err := c.getTargetPath(fileName)
	for _, v := range target {
		fullP := "data/" + v
		if !gfile.Exists(fullP) {
			err = gfile.Mkdir(fullP)
			if err != nil {
				panic(err)
			}
			err = gproc.ShellRun("git init " + fullP)
			if err != nil {
				panic(err)
			}

		}
	}
	fmt.Println(target)

	return
}

func (c cCi) getTargetPath(fileName string) ([]string, error) {
	var targetArr []string
	err := gfile.ReadLines(fileName, func(text string) error {
		if text == "" {
			return nil
		}

		lineInfo := gstr.Split(text, "/")
		infoLen := len(lineInfo)
		if infoLen == 0 {
			return nil
		}

		name, err := gregex.ReplaceString(".git$", "", lineInfo[infoLen-1])
		if err != nil {
			return nil
		}

		targetArr = append(targetArr, strings.ToLower(name))
		return nil
	})

	return targetArr, err
}
