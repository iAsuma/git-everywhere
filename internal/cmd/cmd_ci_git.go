package cmd

import (
	"git-auto/utility"
	"git-auto/utility/qlog"
	"github.com/gogf/gf/v2/os/gproc"
	"github.com/gogf/gf/v2/text/gregex"
	"github.com/gogf/gf/v2/text/gstr"
	"strings"

	//"bufio"
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/util/gtag"
)

const (
	cCiGitUsage = "lsq ci git"
	cCiGitBrief = "make git repo sync to other repos"
)

func init() {
	gtag.Sets(g.MapStrStr{
		"cCiGitUsage": cCiGitUsage,
		"cCiGitBrief": cCiGitBrief,
	})
}

type (
	cCiGitInput struct {
		g.Meta `name:"git" usage:"{cCiGitUsage}" brief:"{cCiGitBrief}" config:"ci.git"`
		From   string   `name:"from" short:"fr"  brief:"copy from who's rep, example: https://github.com/iasuma/"`
		To     []string `name:"to" short:"to"  brief:"copy to who's repo, , example: https://gitee.com/iasuma/"`
		Daemon bool     `name:"daemon" short:"d" brief:"run as a daemon"`
	}
	cCiGitOutput struct{}

	repoEntity struct {
		Origin  string
		Account string
		Addr    string
	}
)

func (c cCi) Git(ctx context.Context, in cCiGitInput) (out cCiGitOutput, err error) {
	from, err := c.dealInputFrom(in.From)
	if err != nil {
		qlog.Print("from url something wrong")
		return
	}

	// github个人主页项目名
	personHome := strings.ToLower(fmt.Sprintf("%s.%s.io", from.Account, from.Origin))

	rootPath := utility.GetPwd()
	fileName := fmt.Sprintf("%s%s", rootPath, "/res/git-repo-url.lsq")

	target, err := c.getTargetPath(fileName)
	for _, v := range target {
		fullP := "../data/" + v

		if !gfile.Exists(fullP) {
			err = gfile.Mkdir(fullP)
			if err != nil {
				panic(err)
			}

			err = gproc.ShellRun(fmt.Sprintf("cd ../data/;git clone %s/%s.git", from.Addr, v))
			if err != nil {
				panic(err)
			}

			err = gproc.ShellRun(fmt.Sprintf("cd %s;git remote rename origin %s", fullP, from.Origin))
			if err != nil {
				panic(err)
			}

			for _, t := range in.To {
				to, err := c.dealInputFrom(t)
				if err != nil {
					panic(err)
				}

				if gstr.Equal(v, personHome) {
					if to.Origin == "gitee" {
						v = to.Account
					} else {
						v = strings.ToLower(fmt.Sprintf("%s.%s.io", to.Account, to.Origin))
					}
				}

				err = gproc.ShellRun(fmt.Sprintf("cd %s;git remote add %s %s/%s.git", fullP, to.Origin, to.Addr, v))
				if err != nil {
					panic(err)
				}

			}
		}

		for _, t := range in.To {
			to, err := c.dealInputFrom(t)
			if err != nil {
				panic(err)
			}
			err = gproc.ShellRun(fmt.Sprintf("cd %s;git pull %s master -f", fullP, from.Origin))
			if err != nil {
				panic(err)
			}

			err = gproc.ShellRun(fmt.Sprintf("cd %s;git push %s master -f", fullP, to.Origin))
			if err != nil {
				panic(err)
			}
		}
	}

	return
}

// 替换git拉取方式为ssh
func (c cCi) dealInputFrom(repoUrl string) (repo repoEntity, err error) {
	repoUrl = gstr.TrimRight(repoUrl, "/")
	repoUrl, err = gregex.ReplaceString("https://|http://|/|:|^git@", "#", repoUrl)

	fromUrlArr := gstr.Split(repoUrl, "#")
	fmt.Println("fromUrlArr", fromUrlArr)

	if len(fromUrlArr) != 3 {
		qlog.Print("from url is wrong")
		return
	}

	origin := gstr.Split(fromUrlArr[1], ".")
	orLen := len(origin)

	repo.Origin = origin[orLen-2]
	repo.Account = fromUrlArr[2]
	repo.Addr = fmt.Sprintf("git@%s:%s", fromUrlArr[1], fromUrlArr[2])

	return
}

// 替换git提交方式为ssh
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
