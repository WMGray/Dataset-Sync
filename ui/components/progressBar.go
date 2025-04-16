package components

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
	"image/color"
)

// CustomProgressBar 自定义进度条，支持设置颜色
type CustomProgressBar struct {
	widget.ProgressBar
	BarColor color.Color // 进度条填充颜色
}

// NewCustomProgressBar 创建新的自定义进度条
func NewCustomProgressBar() *CustomProgressBar {
	bar := &CustomProgressBar{
		BarColor: color.RGBA{R: 0, G: 122, B: 255, A: 255}, // 默认蓝色
	}
	bar.ExtendBaseWidget(bar)
	bar.Resize(fyne.NewSize(100, 10)) // 默认大小
	return bar
}

// SetBarColor 设置进度条颜色
func (bar *CustomProgressBar) SetBarColor(c color.Color) {
	bar.BarColor = c
	bar.Refresh()
}

// CreateRenderer 重写渲染器，自定义进度条样式
func (bar *CustomProgressBar) CreateRenderer() fyne.WidgetRenderer {
	// 背景矩形（灰色）
	background := canvas.NewRectangle(color.NRGBA{R: 200, G: 200, B: 200, A: 255})
	background.CornerRadius = 2

	// 填充矩形（使用自定义颜色）
	fill := canvas.NewRectangle(bar.BarColor)
	fill.CornerRadius = 2

	return &customProgressBarRenderer{
		bar:        bar,
		background: background,
		fill:       fill,
		objects:    []fyne.CanvasObject{background, fill},
	}
}

// customProgressBarRenderer 自定义进度条渲染器
type customProgressBarRenderer struct {
	bar        *CustomProgressBar
	background *canvas.Rectangle
	fill       *canvas.Rectangle
	objects    []fyne.CanvasObject
}

// MinSize 返回最小尺寸
func (r *customProgressBarRenderer) MinSize() fyne.Size {
	return fyne.NewSize(100, 10)
}

// Layout 布局组件
func (r *customProgressBarRenderer) Layout(size fyne.Size) {
	r.background.Resize(size)
	fillWidth := size.Width * float32(r.bar.Value)
	r.fill.Resize(fyne.NewSize(fillWidth, size.Height))
}

// Refresh 刷新渲染
func (r *customProgressBarRenderer) Refresh() {
	r.fill.FillColor = r.bar.BarColor
	fillWidth := r.background.Size().Width * float32(r.bar.Value)
	r.fill.Resize(fyne.NewSize(fillWidth, r.background.Size().Height))
	r.background.Refresh()
	r.fill.Refresh()
}

// Objects 返回渲染对象
func (r *customProgressBarRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

// Destroy 清理资源
func (r *customProgressBarRenderer) Destroy() {}
