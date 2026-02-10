VERSION=latest

# 应用名称
APP_NAME=im
# 服务名称
SERVER_NAME=user
# 服务类型
SERVER_TYPE=api

# 测试环境配置
# docker的镜像发布地址
DOCKER_REPO_TEST=registry.cn-hangzhou.aliyuncs.com/0rz/${APP_NAME}-${SERVER_NAME}-${SERVER_TYPE}-dev
# 测试版本
VERSION_TEST=$(VERSION)
# 编译的程序名称
APP_NAME_TEST=zero-${APP_NAME}-${SERVER_NAME}-${SERVER_TYPE}-test

# 测试下的编译文件
DOCKER_FILE_TEST=./deploy/dockerfile/Dockerfile_${SERVER_NAME}_${SERVER_TYPE}_dev

# 测试环境的编译发布
build-test:

	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bin/${SERVER_NAME}-${SERVER_TYPE} ./apps/${SERVER_NAME}/${SERVER_TYPE}/${SERVER_NAME}.go
	docker build . -f ${DOCKER_FILE_TEST} --no-cache -t ${APP_NAME_TEST}

# 镜像的测试标签
tag-test:

	@echo 'create tag ${VERSION_TEST}'
	docker tag ${APP_NAME_TEST} ${DOCKER_REPO_TEST}:${VERSION_TEST}

publish-test:

	@echo 'publish ${VERSION_TEST} to ${DOCKER_REPO_TEST}'
	docker push $(DOCKER_REPO_TEST):${VERSION_TEST}

# release-test会执行上面命令，会将user服务编译成二进制文件，然后根据编译文件构建镜像，给镜像打标签，最后推送到阿里云镜像仓库中
release-test: build-test tag-test publish-test
