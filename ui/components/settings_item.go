package components

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// SettingItem 单条设置项组件，包含标题和右侧组件，包裹在圆角矩形框中
type SettingItem struct {
	widget.BaseWidget
	title     *widget.Label     // 标题标签
	control   fyne.CanvasObject // 右侧组件
	container *fyne.Container   // 整体容器
}

// NewSettingItem 创建新的设置项
func NewSettingItem(title *widget.Label, control fyne.CanvasObject) *SettingItem {
	item := &SettingItem{
		title:   title,
		control: control,
	}
	item.ExtendBaseWidget(item)

	// 设置标题样式
	item.title.TextStyle = fyne.TextStyle{Bold: true}
	item.title.Alignment = fyne.TextAlignLeading // 左对齐，匹配图中样式

	// 创建内容布局：标题在左，控件在右
	content := container.NewBorder(
		nil,          // top
		nil,          // bottom
		item.title,   // left
		item.control, // right
		nil,          // center
	)

	// 创建圆角矩形框
	rect := canvas.NewRectangle(theme.BackgroundColor())
	rect.CornerRadius = 8 // 圆角半径
	// 固定

	// 使用 Stack 组合圆角框和内容
	item.container = container.NewStack(
		rect,
		content, // 直接使用 content，不再居中
	)

	return item
}

// CreateRenderer 实现渲染
func (item *SettingItem) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(item.container)
}
