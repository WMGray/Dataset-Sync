package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
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

var ui *MainUI // 全局变量，存储主界面实例

func NewMainUI(window fyne.Window) *MainUI {
	ui = &MainUI{
		window: window,
	}
	ui.window.Resize(fyne.NewSize(1024, 768))

	// 创建每个功能模块的容器
	ui.dataset = CreateDatasets()
	ui.currentContent = ui.dataset

	// 创建左侧导航栏
	navBar := ui.createNavBar()

	// 创建主内容区域
	mainContent := container.NewStack(ui.currentContent)

	// 创建水平分割布局，左侧导航栏占 20%
	split := container.NewHSplit(navBar, mainContent)
	split.Offset = 0.2

	// 设置窗口内容
	window.SetContent(split)

	return ui
}

// createNavBar 创建左侧导航栏
func (m *MainUI) createNavBar() fyne.CanvasObject {
	// 功能按钮
	topFuncBtns := container.NewVBox(
		widget.NewButtonWithIcon("数据集", theme.StorageIcon(), func() {
			ui.showContent(ui.dataset)
		}),
		widget.NewButtonWithIcon("上传", theme.UploadIcon(), func() {
			ui.showContent(ui.upload)
		}),
	)

	//	// 下半功能区 -- 账户、设置
	btmFuncBtns := container.NewVBox(
		widget.NewButtonWithIcon("账户", theme.AccountIcon(), func() {
			ui.showContent(ui.validation)
		}),
		widget.NewButtonWithIcon("设置", theme.SettingsIcon(), func() {
			ui.showContent(ui.settings)
		}),
	)

	// 导航栏容器
	return container.NewVBox(
		topFuncBtns,
		layout.NewSpacer(),    // 添加弹性空间
		widget.NewSeparator(), // 动态分割线
		btmFuncBtns,
	)
}

// CreateDatasets 创建数据集管理界面
func (ui *MainUI) showContent(content *fyne.Container) {
	if ui.currentContent != nil {
		ui.currentContent.Hide()
	}
	content.Show()
	ui.currentContent = content
}
