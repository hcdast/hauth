FROM golang:1.9.2
# 指定制作我们的镜像的联系人信息（镜像创建者）
MAINTAINER hc <1984146116@qq.com>

# 设置时区
RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime  && echo 'Asia/Shanghai' >/etc/timezone

# 安装 golint
ENV GOPATH /go
ENV PATH ${GOPATH}/bin:$PATH
ENV GO111MODULE=on
ENV CGO_ENABLED 0
RUN export GOPROXY=https://goproxy.cn/

# RUN go get -u github.com/golang/lint/golint
WORKDIR ${GOPATH}/src/app/
COPY . ${GOPATH}/src/app/

# RUN go get github.com/astaxie/beego && go get github.com/astaxie/beego/orm && go get github.com/astaxie/beego/toolbox
# RUN go get github.com/astaxie/beedb

# RUN go mod init example-hauth && go mod tidy && go install && go mod vendor
EXPOSE 80
# 容器启动时执行的命令，类似npm run start
CMD ["go", "run", "main.go"]
