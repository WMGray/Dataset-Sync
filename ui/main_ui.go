// 3. 右侧： 上半部分 -- 一些查找、搜索操作（一个横条）， 下面部分 -- 部分数据集展示
package ui

// 主界面UI布局
// 1. 上方横条 软件信息
// 2. 左侧： 上半部分 -- 功能模块， 下面部分 -- 账户、设置
// 3. 右侧： 上半部分 -- 一些查找、搜索操作（一个横条）， 下面部分 -- 部分数据集展示
import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// NavItem 定义导航项的结构
type NavItem struct {
	Label string
	Icon  fyne.Resource
	OnTap func()
}

// NavSection 定义导航分区的结构
type NavSection struct {
	Items []NavItem
}

// MainUI 主界面结构
type MainUI struct {
	app    fyne.App
	window fyne.Window
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

	// 创建左侧导航按钮
	datasetsBtn := widget.NewButtonWithIcon("数据集", theme.FolderIcon(), nil)
	uploadBtn := widget.NewButtonWithIcon("上传", theme.UploadIcon(), nil)

	// 创建左侧导航容器
	navContainer := container.NewVBox(
		datasetsBtn,
		uploadBtn,
	)

	// 创建右侧容器
	rightSide := container.NewBorder(nil, nil, nil, nil, nil)

	// 创建数据集页面
	datasetPage := NewDatasetPage(m.window)

	// 创建上传页面
	uploadPage := NewUploadPage(m.window)

	// 设置导航事件
	datasetsBtn.OnTapped = func() {
		rightSide.Objects = []fyne.CanvasObject{datasetPage.Container()}
		rightSide.Refresh()
	}

	uploadBtn.OnTapped = func() {
		rightSide.Objects = []fyne.CanvasObject{uploadPage.Container()}
		rightSide.Refresh()
	}

	// 创建分割容器
	split := container.NewHSplit(navContainer, rightSide)
	split.Offset = 0.2

	// 设置窗口内容
	m.window.SetContent(split)

	// 默认选择数据集页面
	datasetsBtn.OnTapped()
}

// Run 运行主界面
func (m *MainUI) Run() {
	m.window.ShowAndRun()
}
