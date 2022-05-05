# 人无远虑，必有近忧

### 将你的github仓库代码自动同步到码云等其他仓库

#### 前置问题

需要在你的个人主机或服务器配置SSH公钥

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

##### 同步单个仓库

```yaml
version: "3.1"
services:
  git-everywhere:
    image: iasuma/git-everywhere
    container_name: git-everywhere
    restart: always
    volumes:
      - ~/.ssh/:/root/.ssh/:ro
      - ./data/:/var/work/data/
    environment:
      CI_GIT_FROM_ORIGIN: "https://github.com/your-name"
      CI_GIT_TO_ORIGIN: "https://xxxlab.com/your-name"
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
    volumes:
      - ~/.ssh/:/root/.ssh/:ro
      - ./data/:/var/work/data/
    environment:
      CI_GIT_FROM_ORIGIN: "https://github.com/your-name"
      CI_GIT_TO_ORIGIN: "[\"https://gitee.com/your-name\"]"
      # CI_GIT_TO_ORIGIN: "[\"https://gitee.com/your-name\", \"https://xxxlab.com/your-name\"]"
      CI_GIT_DELAY: 15
      CI_GIT_ORIGIN_REPO: "[\"https://github.com/your-name/project1.git\",\"git@github.com:your-name/project2.git\"]"
```

#### 使用配置项目（推荐）

```
git-everywhere  启动根目录
├─config           	  配置目录
│  ├─config.yaml                配置文件
│  └─config.yaml.example        示例配置文件
├─data                数据目录
│  ├─project1          
│  ├─project2          
│  ...
├─res                 资源目录
│  ├─git-repo-url.lsq          	项目仓库地址文件
│  └─git-repo-url.lsq.example  	示例项目仓库地址文件
│
├─docker-compose.yml  docker启动文件  
```

config.yaml

```yaml
ci:
  git:
    delay: 20
    from: "https://github.com/your-name"
    to:
      - "https://gitee.com/your-name"
      - "https://xxxxlab.com/your-name"
    repo:
      - "git@github.com:your-name/project1.git"
      - "git@github.com:your-name/project2.git"
      - "git@github.com:your-name/your-name.github.io.git"
```

git-repo-url.lsq

```
git@github.com:your-name/project1.git
git@github.com:your-name/project2.git
https://github.com/your-name/you-name.github.io.git
```



`config.yaml`中的`ci.git.repo`优先级高于`git-repo-url.lsq`列表，且只会其中一个生效。配了`ci.git.repo`就不用再配置`git-repo-url.lsq`了，反之亦然

Docker-compose.yaml

```yaml
version: "3.1"
services:
  git-everywhere:
    image: iasuma/git-everywhere
    container_name: git-everywhere
    restart: always
    volumes:
      - ~/.ssh/:/root/.ssh/:ro
      - ./config/:/var/work/config/:rw
      - ./data/:/var/work/data/:rw
      - ./res/:/var/work/res/:rw
```

示例：[示例项目地址(github)](https://github.com/iAsuma/git-everywhere-demo)  [示例项目地址(gitee)](https://gitee.com/iAsuma/git-everywhere-demo)

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

同步间隔时间，默认10s，示例：`20`

##### `CI_GIT_ORIGIN_REPO`

需要同步的项目仓库，项目仓库必须与配置的`CI_GIT_FROM_ORIGIN`主机地址一致，示例：`git@github.com:iAsuma/git-everywhere.git`

若同步的是github的个人主页项目`your-name/your-name.github.io`，目标主机是`gitee`，将会同步到`your-name/your-name`仓库。比如：iAsuma/iAsuma.github.io会同步到gitee的iAsuma/iAsuma仓库

#### 常见问题

![image-20220503202636583](https://img-oss2.udzan.com/md/202205032026952.png)

请检查git仓库的ssh权限，或者ssh公钥映射问题

![image-20220503202636583](http://img-oss2.udzan.com/md/202205032026952.png)