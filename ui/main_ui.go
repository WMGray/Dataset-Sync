package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// MainUI 主界面结构
type MainUI struct {
	app         fyne.App
	window      fyne.Window
	contentArea *fyne.Container // 右侧内容区容器
	currentPage string          // 当前页面标识
}

// NewMainUI 创建新的主界面
func NewMainUI() *MainUI {
	mainUI := &MainUI{
		app:    app.New(),
		window: nil,
	}
	mainUI.init()
	return mainUI
}

// init 初始化主界面
func (m *MainUI) init() {
	m.window = m.app.NewWindow("Dataset Sync")
	m.window.Resize(fyne.NewSize(1024, 768))

	// 初始化右侧内容区
	m.contentArea = container.NewMax(widget.NewLabel("欢迎使用 Dataset Sync"))

	m.currentPage = "home"

	// 创建左侧导航栏
	navBar := m.createNavBar()

	// 创建主布局
	mainLayout := container.NewBorder(nil, nil, navBar, nil, m.contentArea)
	m.window.SetContent(mainLayout)
}

// createNavBar 创建左侧导航栏
func (m *MainUI) createNavBar() fyne.CanvasObject {
	// 功能按钮
	topFuncBtns := container.NewVBox(
		widget.NewButtonWithIcon("数据集", theme.StorageIcon(), func() {
			m.switchPage("dataset")
		}),
		widget.NewButtonWithIcon("上传", theme.UploadIcon(), func() {
			m.switchPage("upload")
		}),
	)

	//	// 下半功能区 -- 账户、设置
	btmFuncBtns := container.NewVBox(
		widget.NewButtonWithIcon("账户", theme.AccountIcon(), func() {
			m.switchPage("account")
		}),
		widget.NewButtonWithIcon("设置", theme.SettingsIcon(), func() {
			m.switchPage("settings")
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

// switchPage 切换页面
func (m *MainUI) switchPage(page string) {
	var newContent fyne.CanvasObject

	switch page {
	case "dataset":
		newContent = CreateDatasetContainer()
	case "upload":
		newContent = NewUploadPage(m.window).Container()
	case "account":
		newContent = widget.NewLabel("账户功能区")
	case "settings":
		newContent = widget.NewLabel("设置功能区")
	default:
		newContent = widget.NewLabel("未知页面")
	}

	// 更新右侧内容区
	m.currentPage = page
	m.contentArea.Objects = []fyne.CanvasObject{newContent}
	m.contentArea.Refresh()
}

// Run 运行主界面
func (m *MainUI) Run() {
	m.window.ShowAndRun()
}
