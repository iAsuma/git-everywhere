package cmd

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/genv"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gproc"
	"github.com/gogf/gf/v2/os/gtimer"
	"github.com/gogf/gf/v2/text/gregex"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gtag"
	"lsq-cli/utility"
	"lsq-cli/utility/qlog"
	"lsq-cli/utility/qstr"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

var proc = make(chan os.Signal)

const (
	cCiGitUsage = "lsq ci git"
	cCiGitBrief = "make git repo sync to other repos"
)

const DataDir = "data"
const TimeInterval = 1

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
	}
	cCiGitOutput struct{}

	repoEntity struct {
		Origin     string
		Account    string
		Addr       string
		PersonHome string
	}
)

func (c cCi) Git(ctx context.Context, in cCiGitInput) (out cCiGitOutput, err error) {
	from, err := c.dealInputFrom(in.From)
	if err != nil {
		qlog.Echo("from url something wrong")
		return
	}

	// github个人主页项目名
	from.PersonHome = strings.ToLower(fmt.Sprintf("%s.%s.io", from.Account, from.Origin))

	rootPath := utility.GetPwd()
	fileName := fmt.Sprintf("%s%s", rootPath, "/res/git-repo-url.lsq")
	target, err := c.getTargetPath(fileName)

	dataDir := genv.GetWithCmd("lsq.cli.data.dir", DataDir).String()
	if !gfile.Exists(dataDir) {
		err = gfile.Mkdir(dataDir)
		if err != nil {
			panic(err)
		}
	}

	gtimer.AddSingleton(ctx, TimeInterval*time.Second, func(ctx context.Context) {
		qlog.Echo("本次同步开始")
		err = c.syncGitRepo(in, from, target, dataDir)
		qlog.Echo("本次同步结束")
		time.Sleep(time.Second * (10 - TimeInterval))
		return
	})

	handleProcess()
	return
}

func (c cCi) syncGitRepo(in cCiGitInput, from repoEntity, target []string, dataDir string) (err error) {
	for _, v := range target {
		fullP := gstr.TrimRight(dataDir, "/") + string(os.PathSeparator) + v

		if !gfile.Exists(fullP) {
			str, err := gproc.ShellExec(fmt.Sprintf("cd %s;git clone %s/%s.git", dataDir, from.Addr, v))
			if err != nil {
				panic(err)
			}
			qlog.Echo(str)

			str, err = gproc.ShellExec(fmt.Sprintf("cd %s;git remote rename origin %s", fullP, from.Origin))
			if err != nil {
				panic(err)
			}
			qlog.Echo(qstr.ReplaceN(str))

			for _, t := range in.To {
				to, err := c.dealInputFrom(t)
				if err != nil {
					panic(err)
				}

				if gstr.Equal(v, from.PersonHome) {
					if to.Origin == "gitee" {
						v = to.Account
					} else {
						v = strings.ToLower(fmt.Sprintf("%s.%s.io", to.Account, to.Origin))
					}
				}

				str, err = gproc.ShellExec(fmt.Sprintf("cd %s;git remote add %s %s/%s.git", fullP, to.Origin, to.Addr, v))
				if err != nil {
					panic(err)
				}
				qlog.Echo(qstr.ReplaceN(str))
			}
		}

		for _, t := range in.To {
			to, err := c.dealInputFrom(t)
			if err != nil {
				panic(err)
			}
			str, err := gproc.ShellExec(fmt.Sprintf("cd %s;git pull %s master -f", fullP, from.Origin))
			if err != nil {
				panic(err)
			}
			qlog.Echo(qstr.ReplaceN(str))

			str, err = gproc.ShellExec(fmt.Sprintf("cd %s;git push %s master -f", fullP, to.Origin))
			if err != nil {
				panic(err)
			}
			qlog.Echo(qstr.ReplaceN(str))
		}
	}

	return
}

func handleProcess() {
	var sig os.Signal
	signal.Notify(
		proc,
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGKILL,
		syscall.SIGTERM,
		syscall.SIGABRT,
	)

	for {
		sig = <-proc
		sigString := sig.String()
		pid := gproc.Pid()
		switch sig {
		case syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGABRT:
			qlog.Printf("%d: Server shutting down by signal, %s", pid, sigString)
			return
		case syscall.SIGTERM:
			qlog.Printf("%d: server gracefully shutting down by signal: %s", pid, sigString)
		default:

		}
	}
}

// 替换git拉取方式为ssh
func (c cCi) dealInputFrom(repoUrl string) (repo repoEntity, err error) {
	repoUrl = gstr.TrimRight(repoUrl, "/")
	repoUrl, err = gregex.ReplaceString("https://|http://|/|:|^git@", "#", repoUrl)

	fromUrlArr := gstr.Split(repoUrl, "#")

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
