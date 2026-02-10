# 生成grpc
goctl rpc protoc ./apps/user/rpc/user.proto --go_out=./apps/user/rpc/ --go-grpc_out=./apps/user/rpc/ --zrpc_out=./apps/user/rpc/ -style goZero

# 生成API
goctl api go -api apps/user/api/user.api -dir apps/user/api -style goZero
