package components

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"image/color"
	"time"
)

// DropZone 自定义拖放区域组件
type DropZone struct {
	widget.BaseWidget
	background *canvas.Rectangle
	content    *fyne.Container
	onDropped  func([]fyne.URI)
	window     fyne.Window
}

// NewDropZone 创建新的拖放区域
func NewDropZone(window fyne.Window, content *fyne.Container, onDropped func([]fyne.URI)) *DropZone {
	d := &DropZone{
		background: canvas.NewRectangle(color.NRGBA{R: 240, G: 240, B: 240, A: 255}),
		content:    content,
		onDropped:  onDropped,
		window:     window,
	}
	d.ExtendBaseWidget(d)
	return d
}

// CreateRenderer 实现自定义渲染
func (d *DropZone) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(container.NewStack(d.background, d.content))
}

// isInBounds 检查点是否在组件边界内
func (d *DropZone) isInBounds(pos fyne.Position) bool {
	x, y := d.Position().X, d.Position().Y
	w, h := d.Size().Width, d.Size().Height
	return pos.X >= x && pos.X <= x+w && pos.Y >= y && pos.Y <= y+h
}

// Highlight 提供视觉反馈 - 高亮
func (d *DropZone) Highlight() {
	d.background.FillColor = color.NRGBA{R: 220, G: 220, B: 250, A: 255}
	d.background.Refresh()

	// 短暂高亮后恢复默认颜色
	go func() {
		time.Sleep(500 * time.Millisecond) // 高亮 0.5 秒
		d.ResetHighlight()
	}()
}

// ResetHighlight 恢复默认背景颜色
func (d *DropZone) ResetHighlight() {
	d.background.FillColor = color.NRGBA{R: 240, G: 240, B: 240, A: 255}
	d.background.Refresh()
}

// CreateDropArea 创建拖放区域
func CreateDropArea(window fyne.Window, content *fyne.Container) *DropZone {
	// 创建拖放区域
	dropZone := NewDropZone(window, container.NewCenter(content), func(uris []fyne.URI) {
		fmt.Println("Dropped files:")
		for _, uri := range uris {
			fmt.Printf("Dropped file: %s\n", uri.Path())
		}
	})

	// 绑定窗口的拖放事件
	window.SetOnDropped(func(pos fyne.Position, uris []fyne.URI) {
		if dropZone.isInBounds(pos) {
			dropZone.Highlight() // 文件放下时高亮
			dropZone.onDropped(uris)
		}
	})

	return dropZone
}
