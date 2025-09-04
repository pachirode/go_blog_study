package main

import (
	"os"

	"github.com/pachirode/go_blog_study/cmd/miniblog/app"
	// 导入 automaxprocs 包，可以在程序启动时自动设置 GOMAXPROCS 配置，
	_ "go.uber.org/automaxprocs"
)

// Go 程序的默认入口函数。阅读项目代码的入口函数.
func main() {
	// 创建迷你博客命令
	command := app.NewMiniBlogCommand()

	// 执行命令并处理错误
	if err := command.Execute(); err != nil {
		// 如果发生错误，则退出程序
		// 返回退出码，可以使其他程序（例如 bash 脚本）根据退出码来判断服务运行状态
		os.Exit(1)
	}
}
