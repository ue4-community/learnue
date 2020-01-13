## Start from golang v1.12 base image
#FROM golang:1.12
#
#WORKDIR /app/studygolang
#
#COPY .. /app/studygolang
#
#RUN make build
#
#CMD ["bin/studygolang"]


# 构建阶段
FROM golang:1.13 AS build-env
#国内镜像;启用gomodules;静态编译c的依赖;64位;linux平台
ENV GOPROXY=https://goproxy.cn  GO111MODULE=on  CGO_ENABLED=0 GOARCH=amd64 GOOS=linux
WORKDIR /src
ADD . .
RUN go mod vendor && go build -o goapp -mod vendor github.com/studygolang/studygolang/goapp

# 最终阶段
FROM golang:1.13-alpine
WORKDIR /app
ENV LEARNUE_HOME=/app
COPY --from=build-env /src/goapp/goapp .
COPY --from=build-env /src/config ./config
EXPOSE 8088
CMD ./goapp