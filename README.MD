# 人无远虑，必有近忧

### 将你的github仓库代码自动同步到码云等其他仓库

#### 前置问题

需要在你的个人主机或VPS配置SSH公钥

若使用`docker`，需将你的ssh钥匙目录映射到容器内，否则会提示权限问题

#### 使用docker

推送单个仓库

```shell
$ docker run -d \
    --name git-everywhere \
    -v ~/.ssh:/root/.ssh \
    -e CI_GIT_FROM_ORIGIN="https://github.com/your-name" \
    -e CI_GIT_TO_ORIGIN="https://gitee.com/your-name" \
    -e CI_GIT_DELAY=10 \
    -e CI_GIT_ORIGIN_REPO="git@github.com:your-name/your-project.git" \
    iasuma/git-everywhere
```



推送多个仓库

``` shell
$ docker run -d \
    --name git-everywhere \
    -v ~/.ssh:/root/.ssh \
    -e CI_GIT_FROM_ORIGIN="https://github.com/your-name" \
    -e CI_GIT_TO_ORIGIN=[\"https://gitee.com/your-name\",\"https://xxxx.com/your-name\"] \
    -e CI_GIT_DELAY=20 \
    -e CI_GIT_ORIGIN_REPO=[\"git@github.com:your-name/your-project1.git\",\"your-name/your-project2.git\"] \
    iasuma/git-everywhere
```

若是数组值，引号需要使用反斜杠（\）转义


#### 使用docker-compose

##### 同步多个仓库

```yaml
version: "3.1"
services:
  git-everywhere:
    image: iasuma/git-everywhere
    container_name: git-everywhere
    restart: always
    deploy:
      restart_policy:
        condition: on-failure
    volumes:
      - ~/.ssh/:/root/.ssh/:ro
      - ./data/:/var/work/data/
    environment:
      TZ: Asia/Shanghai
      CI_GIT_FROM_ORIGIN: "https://github.com/your-name"
      CI_GIT_TO_ORIGIN: "https://gitee.com/your-name"
      CI_GIT_DELAY: 15
      CI_GIT_ORIGIN_REPO: "git@github.com:your-name/your-project.git"

```

##### 同步多个仓库

```yaml
version: "3.1"
services:
  git-everywhere:
    image: iasuma/git-everywhere
    container_name: git-everywhere
    restart: always
    deploy:
      restart_policy:
        condition: on-failure
    volumes:
      - ~/.ssh/:/root/.ssh/:ro
      - ./data/:/var/work/data/
    environment:
      TZ: Asia/Shanghai
      CI_GIT_FROM_ORIGIN: "https://github.com/your-name"
      CI_GIT_TO_ORIGIN: "[\"https://gitee.com/your-name\"]"
      CI_GIT_DELAY: 15
      CI_GIT_ORIGIN_REPO: "[\"https://github.com/your-name/project1.git\",\"git@github.com:your-name/project2.git\"]"
```

#### USAGE

~~~shell
$ lsq-ci git -h
USAGE
    lsq-ci git

OPTION
    -dl, --delay   time interval, default 10 seconds
    -fr, --from    copy from who's rep, example: https://github.com/iasuma/
    -to, --to      copy to who's repo, example: https://gitee.com/iasuma
    -rp, --repo    the repositories which you want sync, example: git@github.com:iAsuma/your-project.git
    -h, --help     more information about this command
~~~


#### 环境变量

##### `CI_GIT_FROM_ORIGIN`

来源仓库地址，示例：`https://github.com/iasuma`

##### `CI_GIT_TO_ORIGIN`

目标仓库地址，示例：`https://gitee.com/iasuma`

##### `CI_GIT_DELAY`

定时时间，默认10s，示例：`20`

##### `CI_GIT_ORIGIN_REPO`

需要同步的仓库，仓库必须来自配置的`CI_GIT_FROM_ORIGIN`地址，示例：`git@github.com:iAsuma/git-everywhere.git`

#### 常见问题

![image-20220503202636583](http://img-oss2.udzan.com/md/202205032026952.png)

请检查git仓库的ssh权限，或者ssh公钥映射问题