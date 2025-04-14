package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type MainUI struct {
	window         fyne.Window
	dataset        *fyne.Container
	upload         *fyne.Container
	validation     *fyne.Container
	storage        *fyne.Container
	export         *fyne.Container
	settings       *fyne.Container
	currentContent *fyne.Container
}

func NewMainUI(window fyne.Window) *MainUI {
	ui := &MainUI{
		window: window,
	}

	// 创建每个功能模块的容器
	ui.dataset = CreateDatasetContainer()
	ui.currentContent = ui.dataset

	// 创建左侧导航栏
	sidebar := container.NewVBox(
		widget.NewLabel("图片数据集管理工具"),
		widget.NewButtonWithIcon("数据集", theme.StorageIcon(), func() {
			ui.showContent(ui.dataset)
		}),
		widget.NewButtonWithIcon("上传", theme.UploadIcon(), func() {
			ui.showContent(ui.upload)
		}),
		widget.NewButtonWithIcon("验证", theme.ConfirmIcon(), func() {
			ui.showContent(ui.validation)
		}),
		widget.NewButtonWithIcon("存储", theme.StorageIcon(), func() {
			ui.showContent(ui.storage)
		}),
		widget.NewButtonWithIcon("导出", theme.DownloadIcon(), func() {
			ui.showContent(ui.export)
		}),
		widget.NewButtonWithIcon("设置", theme.SettingsIcon(), func() {
			ui.showContent(ui.settings)
		}),
	)

	// 创建主内容区域
	mainContent := container.NewMax(ui.currentContent)

	// 创建水平分割布局，左侧导航栏占 20%
	split := container.NewHSplit(sidebar, mainContent)
	split.Offset = 0.2

	// 设置窗口内容
	window.SetContent(split)

	return ui
}

func (ui *MainUI) showContent(content *fyne.Container) {
	if ui.currentContent != nil {
		ui.currentContent.Hide()
	}
	content.Show()
	ui.currentContent = content
}
