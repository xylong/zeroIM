# 生成rpc
goctl rpc protoc ./app/user/rpc/user.proto --go_out=./app/user/rpc/ --go-grpc_out=./app/user/rpc/ --zrpc_out=./app/user/rpc/ -style goZero
