user-rpc-dev:
	@make -f deploy/make/user-rpc.mk release-test

release-test: user-rpc-dev