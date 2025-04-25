package components

import "fyne.io/fyne/v2"

// CustomVBoxLayout 自定义垂直布局，支持自适应宽度、固定组件高度和间隙
type CustomVBoxLayout struct {
	componentHeight float32 // 组件固定高度
	gap             float32 // 固定空白高度
}

// NewCustomVBoxLayout 创建新的自定义垂直布局
// componentHeight: 组件的固定高度
// gap: 组件之间的固定空白高度
func NewCustomVBoxLayout(componentHeight float32, gap float32) *CustomVBoxLayout {
	return &CustomVBoxLayout{
		componentHeight: componentHeight,
		gap:             gap,
	}
}

// Add 向容器添加子组件
func (l *CustomVBoxLayout) Add(container *fyne.Container, child fyne.CanvasObject) {
	if container == nil || child == nil {
		return
	}
	container.Objects = append(container.Objects, child)
	container.Refresh() // 触发布局更新
}

// Layout 实现组件布局
func (l *CustomVBoxLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	y := float32(0)
	// 左右边距各 12 像素，组件宽度自适应
	margin := float32(12)
	componentWidth := size.Width - 2*margin // 自适应容器宽度，减去边距

	for i, obj := range objects {
		// 设置组件大小：宽度自适应，高度固定
		obj.Resize(fyne.NewSize(componentWidth, l.componentHeight))
		// 左侧对齐，偏移 12 像素
		obj.Move(fyne.NewPos(margin, y))

		// 累加高度
		y += l.componentHeight

		// 添加固定空白（最后一个组件后不加）
		if i < len(objects)-1 {
			y += l.gap
		}
	}
}

// MinSize 计算布局最小尺寸
func (l *CustomVBoxLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	height := float32(0)
	for i := range objects {
		// 累加高度
		height += l.componentHeight
		// 添加固定空白（最后一个组件后不加）
		if i < len(objects)-1 {
			height += l.gap
		}
	}
	// 宽度设为 0，允许容器自适应
	return fyne.NewSize(0, height)
}
