package ui

import (
	"dataset-sync/conf"
	"dataset-sync/ui/components"
	"dataset-sync/utils"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// createSettingsView 创建设置界面
func createSettingsView() *fyne.Container {
	// 创建设置项容器 -- 长为窗口宽度，内容宽度为 1000
	vBoxLayout := components.NewCustomVBoxLayout(40, 20)

	content := container.New(vBoxLayout)

	// 自动重命名设置
	autoRename := widget.NewCheck("", nil)
	autoRename.SetChecked(conf.Conf.DatasetConfig.AutoRename) // 设置初始值
	// 设置回调函数
	autoRename.OnChanged = func(value bool) {
		go func() {
			err := utils.ChangeSettings(conf.Conf.DatasetConfig, "AutoRename", value)
			if err != nil {
				fmt.Println("修改设置失败:", err)
			} else {
				fmt.Println("修改设置成功:", value)
			}
		}()
	}
	autoRenameItem := components.NewSettingItem("自动重命名", autoRename)

	vBoxLayout.Add(content, autoRenameItem)

	return container.NewBorder(nil, nil, nil, nil, container.NewScroll(content))
}
