package ui

import (
	"dataset-sync/conf"
	"dataset-sync/ui/components"
	"dataset-sync/utils"
	"errors"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	fyneDialog "fyne.io/fyne/v2/dialog" // 重命名以区分
	"fyne.io/fyne/v2/widget"
	"github.com/sqweek/dialog"
)

// createSettingsView 创建设置界面
func createSettingsView() *fyne.Container {
	// 创建设置项容器 -- 长为窗口宽度，内容宽度为 1000
	vBoxLayout := components.NewCustomVBoxLayout(40, 20)

	content := container.New(vBoxLayout)

	// 自动重命名设置
	autoRename := widget.NewCheck("", nil)
	autoRename.SetChecked(conf.Conf.DatasetConfig.AutoRename) // 设置初始值
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
	autoRenameItem := components.NewSettingItem(widget.NewLabel("自动重命名"), autoRename)

	// 文件存放目录设置
	saveDir := conf.Conf.DatasetConfig.SaveDir
	saveLabel := widget.NewLabel("文件存放目录: " + saveDir)
	saveSettingBtn := widget.NewButton("更改", func() {
		go func() {
			// 使用system dialog库打开文件夹选择对话框
			dir, err := dialog.Directory().Title("选择文件存放目录").Browse()
			if err != nil {
				// 处理错误，包括用户取消
				if !errors.Is(err, dialog.ErrCancelled) {
					fyneDialog.ShowError(err, ui.window)
				}
				return
			}

			// 成功选择了目录
			go func() {
				if err := utils.ChangeSettings(conf.Conf.DatasetConfig, "SaveDir", dir); err != nil {
					fyneDialog.ShowError(err, ui.window)
					return
				}
				// 更新UI显示
				saveLabel.SetText("文件存放目录: " + dir)
				fmt.Println("修改文件存放目录成功:", dir)
			}()
		}()
	})
	saveSettingItem := components.NewSettingItem(saveLabel, saveSettingBtn)

	// 缓存目录更改设置
	cacheDir := conf.Conf.DatasetConfig.TmpDir
	cacheLabel := widget.NewLabel("缓存目录: " + cacheDir)

	cacheSettingBtn := widget.NewButton("更改", func() {
		go func() {
			// 使用system dialog库打开文件夹选择对话框
			dir, err := dialog.Directory().Title("选择缓存目录").Browse()
			if err != nil {
				// 处理错误，包括用户取消
				if !errors.Is(err, dialog.ErrCancelled) {
					fyneDialog.ShowError(err, ui.window)
				}
				return
			}

			// 成功选择了目录
			go func() {
				if err := utils.ChangeSettings(conf.Conf.DatasetConfig, "TmpDir", dir); err != nil {
					fyneDialog.ShowError(err, ui.window)
					return
				}
				// 更新UI显示
				cacheLabel.SetText("缓存目录: " + dir)
				fmt.Println("修改缓存目录成功:", dir)
			}()
		}()
	})
	cacheSettingItem := components.NewSettingItem(cacheLabel, cacheSettingBtn)

	vBoxLayout.Add(content, autoRenameItem)
	vBoxLayout.Add(content, saveSettingItem)
	vBoxLayout.Add(content, cacheSettingItem)

	return container.NewBorder(nil, nil, nil, nil, container.NewScroll(content))
}
