user-rpc-dev:
	@make -f deploy/mk/user-rpc.mk release-test

# 构建镜像
release-test: user-rpc-dev

# 部署服务
install-server:
	cd ./deploy/script && chmod +x release-test.sh && ./release-test.sh