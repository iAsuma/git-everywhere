# 人无远虑，必有近忧

### 将你的github仓库代码自动同步到码云等其他仓库

#### 前置问题

需要在你的个人主机或服务器配置SSH公钥

若使用`docker`，需将你的ssh钥匙目录映射到容器内，否则会提示权限问题

**重要：**

**1. 目标仓库名必须要跟来源仓库名一致！** 比如：github的仓库名叫git-everywhere，在gitee里也必须叫git-everywhere

**2. 目标仓库必须为空仓库或者同步双方的仓库提交历史必须一致，** 建议目标仓库为空仓库

#### 使用docker

同步单个目标主机

```shell
$ docker run -d \
    --name git-everywhere \
    -v ~/.ssh:/root/.ssh \
    -e CI_GIT_FROM_ORIGIN="github.com/your-name" \
    -e CI_GIT_TO_ORIGIN="gitee.com/your-name" \
    -e CI_GIT_DELAY=10 \
    -e CI_GIT_ORIGIN_REPO="your-project1,your-project2" \
    iasuma/git-everywhere
```

同步到多个目标主机

``` shell
$ docker run -d \
    --name git-everywhere \
    -v ~/.ssh:/root/.ssh \
    -e CI_GIT_FROM_ORIGIN="github.com/your-name" \
    -e CI_GIT_TO_ORIGIN=[\"gitee.com/your-name\",\"xxxx.net:8022/your-name\"] \
    -e CI_GIT_DELAY=20 \
    -e CI_GIT_ORIGIN_REPO="your-project1,your-project2" \
    iasuma/git-everywhere
```

若是数组值，引号需要使用反斜杠（\）转义


#### 使用docker-compose

##### 同步单个目标主机

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
      CI_GIT_FROM_ORIGIN: "github.com/your-name"
      CI_GIT_TO_ORIGIN: "xxxlab.com/your-name"
      CI_GIT_DELAY: 15
      CI_GIT_ORIGIN_REPO: "your-project1,your-project2"

```

##### 同步到多个目标主机

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
      CI_GIT_FROM_ORIGIN: "github.com/your-name"
      CI_GIT_TO_ORIGIN: "[\"gitee.com/your-name\"]"
      # CI_GIT_TO_ORIGIN: "[\"gitee.com/your-name\", \"xxxlab.com/your-name\"]"
      CI_GIT_DELAY: 15
      CI_GIT_ORIGIN_REPO: "your-project1,your-project2"
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
    from: "github.com/your-name"
    to:
      - "gitee.com/your-name"
    #  - "xxxxlab.com/your-name"
    repo: "your-project1,your-project2"
```

git-repo-url.lsq

```
project1
project2
you-name.github.io
```



`config.yaml`中的`ci.git.repo`优先级高于`git-repo-url.lsq`列表，且只会其中一个生效。配了`ci.git.repo`就不用再配置`git-repo-url.lsq`了，反之亦然

docker-compose.yaml

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
    -fr, --from    copy from who's rep, example: github.com/iasuma
    -to, --to      copy to who's repo, example: gitee.com/iasuma
    -rp, --repo    the repositories which you want sync, example: your-project1,your-project2
    -h, --help     more information about this command
~~~


#### 环境变量

##### `CI_GIT_FROM_ORIGIN`

来源仓库地址，示例：`github.com/iasuma`

##### `CI_GIT_TO_ORIGIN`

目标仓库地址，目标仓库名必须与来源仓库名一致，示例：`gitee.com/iasuma`

##### `CI_GIT_DELAY`

同步间隔时间，默认10s，示例：`20`

##### `CI_GIT_ORIGIN_REPO`

需要同步的项目仓库，示例：`git-everywhere`

若同步的是github的个人主页项目`your-name/your-name.github.io`，目标主机是`gitee`，将会同步到`your-name/your-name`仓库。比如：iAsuma/iAsuma.github.io会同步到gitee的iAsuma/iAsuma仓库

#### 常见问题

![image-20220503202636583](https://img-oss2.udzan.com/md/202205032026952.png)

请检查git仓库的ssh权限，或者ssh公钥映射问题