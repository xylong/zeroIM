### 目录结构
```
.
├── apps    应用目录
│   └── user
│       ├── api
│       ├── bin     存放脚本执行命令
│       └── rpc
├── bin     编译后的程序存储目录
├── deploy  项目部署相关信息，如部署时候的一些配置，sql或dockerfile
│   ├── dockerfile
│   │   └── Dockerfile_user_rpc_dev
│   ├── mk
│   │   └── user-rpc.mk
│   └── script
├── docker-compose.yml
├── go.mod
├── go.sum
├── Makefile    项目编译脚本
├── pkg     项目的公共工具目录
└── README.md
```

### 用户服务
1. apps/user/rpc下创建user.proto
2. 创建exec.sh，执行脚本，根目录下运行其中命令
3. 根目录下执行go mod tidy，下载服务所需要的依赖包
4. deploy/dockerfile下创建Dockerfile_user_rpc_dev，构建容器镜像
5. Makefile编译项目

### 部署方式
1. 先将程序编译成二进制可执行的文件
2. 然后根据二进制文件构建成镜像文件
3. 再修改构建的镜像标签
4. 然后推送到阿里云
5. 在部署的时候拉取下来构建容器运行即可