FROM alpine:latest

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories

RUN apk update \
    && apk upgrade \
    && apk add openssh-client \
    && apk add git

## 设置 操作系统时区
RUN rm -rf /etc/localtime \
    && ln -s /usr/share/zoneinfo/Asia/Shanghai /etc/localtime

ENV APP_NAME lsq-ci
ENV WORKDIR /var/work
ENV GF_GCFG_PATH ${WORKDIR}/config
ENV LSQ_CI_RES_DIR ${WORKDIR}/res
ENV LSQ_CI_DATA_DIR ${WORKDIR}/data

COPY bin/lsq-ci /usr/bin/
COPY config/config.yaml.example ${WORKDIR}/config/config.yaml
COPY res/git-repo-url.lsq.exmaple ${WORKDIR}/res/git-repo-url.lsq

WORKDIR ${WORKDIR}

CMD lsq-ci git