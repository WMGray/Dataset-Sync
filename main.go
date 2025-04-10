package main

import (
	"dataset-sync/conf"
	"dataset-sync/ui"
	"fmt"
)

func main() {
	// 启动软件
	fmt.Println("软件启动中...")
	// 读取参数
	fmt.Println("读取参数...")
	// 读取配置文件
	if err := conf.Init("conf/config.yaml"); err != nil {
		fmt.Printf("读取配置文件失败, err:%v\n", err)
		return
	}
	fmt.Println("读取配置文件成功")
	// 启动
	fmt.Println("启动软件....")
	// 显示主界面
	ui.NewMainUI().Run()
}
