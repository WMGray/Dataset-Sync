package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// DatasetPage 数据集页面组件
type DatasetPage struct {
	container *fyne.Container
	window    fyne.Window
}

// NewDatasetPage 创建新的数据集页面
func NewDatasetPage(window fyne.Window) *DatasetPage {
	page := &DatasetPage{
		window: window,
	}
	page.init()
	return page
}

// init 初始化数据集页面
func (p *DatasetPage) init() {
	// 创建搜索栏
	searchEntry := widget.NewEntry()
	searchEntry.SetPlaceHolder("搜索数据集...")

	// 创建数据集网格
	grid := container.NewGridWrap(fyne.NewSize(200, 200))

	// 创建主容器
	p.container = container.NewBorder(
		searchEntry, // 顶部搜索栏
		nil,         // 底部
		nil,         // 左侧
		nil,         // 右侧
		grid,        // 中间内容区域
	)
}

// Container 返回页面容器
func (p *DatasetPage) Container() fyne.CanvasObject {
	return p.container
}
