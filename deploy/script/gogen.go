package main

import (
	"path/filepath"
	"runtime"

	"zeroIM/apps/user/models"

	"gorm.io/gen"
)

func main() {
	// 获取当前文件所在目录
	_, filename, _, _ := runtime.Caller(0)
	root := filepath.Dir(filepath.Dir(filepath.Dir(filename))) // zeroIM/

	// 指定输出目录：apps/user/rpc/internal/dao
	outputDir := filepath.Join(root, "apps", "user", "rpc", "internal", "dao")

	// 初始化生成器
	g := gen.NewGenerator(gen.Config{
		OutPath:       outputDir, // 输出目录
		Mode:          gen.WithDefaultQuery | gen.WithQueryInterface,
		FieldNullable: true,
	})

	// 使用结构体生成查询代码
	// ApplyBasic 指定需要生成的模型
	g.ApplyBasic(models.User{})

	// 执行生成
	g.Execute()
}
