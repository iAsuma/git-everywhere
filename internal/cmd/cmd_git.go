package cmd

import (
	"context"
	"fmt"
	"git-everywhere/utility/qcmd"
	"git-everywhere/utility/qlog"
	"git-everywhere/utility/qslice"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/genv"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gproc"
	"github.com/gogf/gf/v2/os/gtimer"
	"github.com/gogf/gf/v2/text/gregex"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/gtag"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

var proc = make(chan os.Signal)

const (
	cCiGitUsage = "lsq-ci git"
	cCiGitBrief = "make git repo sync to other repos"
)

const (
	TimeInterval = 1
	TimeDelay    = 10
)

const (
	EnvCiFromKey  = "ci.git.from.origin"
	EnvCiToKey    = "ci.git.to.origin"
	EnvCiDelayKey = "ci.git.delay"
	EnvCiRepoKey  = "ci.git.origin.repo"
)

func init() {
	gtag.Sets(g.MapStrStr{
		"cCiGitUsage": cCiGitUsage,
		"cCiGitBrief": cCiGitBrief,
	})
}

var Git = cGit{}

type cGit struct {
	g.Meta `name:"git" brief:"Continuous Integration"`
}

type (
	cCiGitInput struct {
		g.Meta `name:"git" usage:"{cCiGitUsage}" brief:"{cCiGitBrief}" config:"ci.git"`
		Delay  int      `name:"delay" short:"dl" brief:"time interval, default 10 seconds"`
		From   string   `name:"from" short:"fr"  brief:"copy from who's rep, example: github.com/iasuma"`
		To     []string `name:"to" short:"to"  brief:"copy to who's repo, example: gitee.com/iasuma"`
		Repo   string   `name:"repo" short:"rp" brief:"the repositories which you want sync, example: your-project1,your-project2"`
	}
	cCiGitOutput struct{}

	repoEntity struct {
		Origin     string
		Account    string
		Addr       string
		PersonHome string
	}
)

func (c cGit) Git(ctx context.Context, in cCiGitInput) (out cCiGitOutput, err error) {
	err = c.defaultInput(&in)
	from, err := c.dealInputFrom(in.From)
	if err != nil {
		qlog.Echo("FROM url something wrong")
		return
	}

	// github个人主页项目名
	from.PersonHome = strings.ToLower(fmt.Sprintf("%s.%s.io", from.Account, from.Origin))

	target, err := c.getTargetPath(in)

	dataDir := genv.GetWithCmd(EnvDataDirKey, DataDir).String()
	if !gfile.Exists(dataDir) {
		err = gfile.Mkdir(dataDir)
		if err != nil {
			qlog.Echo(fmt.Sprintf("文件夹%s生成错误，%s", dataDir, err.Error()))
			return
		}
	}

	gtimer.AddSingleton(ctx, TimeInterval*time.Second, func(ctx context.Context) {
		qlog.Echo("===========================================================================")
		qlog.Echo("repo sync start ...")
		err = c.syncGitRepo(in, from, target, dataDir)
		if in.Delay == 0 {
			in.Delay = TimeDelay
		}
		time.Sleep(time.Second * time.Duration(in.Delay-TimeInterval))
		return
	})

	handleProcess()
	return
}

func (c cGit) syncGitRepo(in cCiGitInput, from repoEntity, target [][2]string, dataDir string) (err error) {
	var bufStr string
	for _, v := range target {
		qlog.Echo("---------------------------------------------------------------------------")
		currentProject := v[0]
		toProject := v[1]
		fullP := gstr.TrimRight(dataDir, "/") + string(os.PathSeparator) + currentProject

		if !gfile.Exists(fullP) {
			// git clone
			cloneShell := fmt.Sprintf("git clone -o %s %s/%s.git", from.Origin, from.Addr, currentProject)
			qlog.Echo(currentProject, cloneShell)
			err = qcmd.ShellRun(dataDir, cloneShell)
			if err != nil {
				qlog.Echo(fmt.Sprintf("%s，git clone -o error，%s", currentProject, err.Error()))
				continue
			}

			for _, t := range in.To {
				to, err := c.dealInputFrom(t)
				if err != nil {
					qlog.Echo(fmt.Sprintf("%s，the TO's url something wrong，%s", currentProject, err.Error()))
					continue
				}

				if gstr.Equal(currentProject, from.PersonHome) {
					if to.Origin == "gitee" {
						toProject = to.Account
					} else {
						toProject = strings.ToLower(fmt.Sprintf("%s.%s.io", to.Account, to.Origin))
					}
				}

				// git remote add
				remoteAddShell := fmt.Sprintf("git remote add %s %s/%s.git", to.Origin, to.Addr, toProject)
				qlog.Echo(currentProject, remoteAddShell)
				err = qcmd.ShellRun(fullP, remoteAddShell)
				if err != nil {
					qlog.Echo(fmt.Sprintf("%s，git remote add error，%s", v, err.Error()))
					continue
				}
			}
		}

		// git fetch
		fetchShell := fmt.Sprintf("git fetch %s", from.Origin)
		qlog.Echo(currentProject, fetchShell)
		err = qcmd.ShellRun(fullP, fetchShell)

		// git branch
		branchLShell := "git branch"
		qlog.Echo(currentProject, branchLShell)
		bufStr = qcmd.MustShellExec(fullP, branchLShell)
		bufStr = strings.ReplaceAll(bufStr, "*", "")
		localB := gstr.SplitAndTrim(bufStr, "\n")

		// git branch -r
		branchRShell := "git branch -r"
		qlog.Echo(currentProject, branchRShell)
		bufStr = qcmd.MustShellExec(fullP, branchRShell)
		remoteB := gstr.SplitAndTrim(bufStr, "\n")

		for _, r := range remoteB {
			if strings.Contains(r, "->") {
				continue
			}

			if strings.Contains(r, from.Origin) {
				newB := strings.ReplaceAll(r, from.Origin+"/", "")
				if qslice.ContainsInSliceString(localB, newB) {
					// git checkout & git merge
					checkoutShell := fmt.Sprintf("git checkout %s;git merge %s", newB, r)
					qlog.Echo(currentProject, checkoutShell)
					err = qcmd.ShellRun(fullP, checkoutShell)
					if err != nil {
						qlog.Echo("git merge error", err.Error())
						continue
					}
				} else {
					// git branch
					branchNShell := fmt.Sprintf("git branch %s %s", newB, r)
					qlog.Echo(currentProject, branchNShell)
					err = qcmd.ShellRun(fullP, branchNShell)
				}
			}
		}

		for _, t := range in.To {
			to, err := c.dealInputFrom(t)
			if err != nil {
				qlog.Echo(fmt.Sprintf("同步%s发生错误，TO url something wrong，%s", v, err.Error()))
				continue
			}

			// git push
			pushShell := fmt.Sprintf("git push %s --all;git push %s --tags", to.Origin, to.Origin)
			qlog.Echo(currentProject, pushShell)
			err = qcmd.ShellRun(fullP, pushShell)
			if err != nil {
				qlog.Echo(fmt.Sprintf("同步%s发生错误，git push error，%s", v, err.Error()))
				continue
			}
		}
	}

	return
}

func (c cGit) defaultInput(in *cCiGitInput) (err error) {
	if in.From == "" {
		in.From = genv.GetWithCmd(EnvCiFromKey).String()
	}

	if len(in.To) == 0 {
		in.To = gconv.SliceStr(genv.GetWithCmd(EnvCiToKey).Slice())
	}

	if in.Delay == 0 {
		in.Delay = genv.GetWithCmd(EnvCiDelayKey).Int()
	}

	if in.Repo == "" {
		in.Repo = genv.GetWithCmd(EnvCiRepoKey).String()
	}

	return
}

// 替换git拉取方式为ssh
func (c cGit) dealInputFrom(repoUrl string) (repo repoEntity, err error) {
	repoUrl = gstr.TrimRight(repoUrl, "/")
	repoUrl, err = gregex.ReplaceString("https://|http://", "", repoUrl)

	fromUrlArr := gstr.Split(repoUrl, "/")

	if len(fromUrlArr) != 2 {
		qlog.Print("from url is wrong")
		return
	}

	originUrl := gstr.Split(fromUrlArr[0], ":")
	if len(originUrl) > 2 {
		qlog.Print("from url is wrong")
		return
	}

	origin := gstr.Split(originUrl[0], ".")
	orLen := len(origin)

	repo.Origin = origin[orLen-2]
	repo.Account = fromUrlArr[1]
	repo.Addr = fmt.Sprintf("git@%s:%s", fromUrlArr[0], fromUrlArr[1])

	return
}

// 替换git提交方式为ssh
func (c cGit) getTargetPath(in cCiGitInput) (targetArr [][2]string, err error) {
	if in.Repo != "" {
		repoList := gstr.SplitAndTrim(in.Repo, ",")
		for _, v := range repoList {
			lineInfo := gstr.Split(v, ":")
			var t [2]string
			t[0] = lineInfo[0]
			if len(lineInfo) > 1 {
				t[1] = lineInfo[1]
			} else {
				t[1] = lineInfo[0]
			}

			targetArr = append(targetArr, t)
		}
	} else {
		resDir := genv.GetWithCmd(EnvResDirKey, ResDir).String()
		fileName := fmt.Sprintf("%s/%s", resDir, "git-repo-url.lsq")
		err = gfile.ReadLines(fileName, func(text string) error {
			if text == "" {
				return nil
			}

			lineInfo := gstr.Split(text, ":")
			var t [2]string
			t[0] = lineInfo[0]
			if len(lineInfo) > 1 {
				t[1] = lineInfo[1]
			} else {
				t[1] = lineInfo[0]
			}

			targetArr = append(targetArr, t)
			return nil
		})
	}

	return targetArr, err
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
