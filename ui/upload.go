package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"image/color"

	"dataset-sync/database"
	"dataset-sync/ui/components"
)

func CreateUpload(w fyne.Window) *fyne.Container {
	// 拖放区域
	dropArea := components.CreateDropArea(w, "拖放文件到这里")
	dropArea.Resize(fyne.NewSize(500, 200))
	dropArea.Move(fyne.NewPos(0, 30)) // 往下偏移一下，避免被 label 挡住

	// 4. 布局组合
	topSection := container.NewStack(
		dropArea,
	)

	// 上传历史记录 -- 可滚动
	historyList := createHistoryList()
	bottomSection := container.NewVScroll(historyList)

	split := container.NewVSplit(topSection, bottomSection)
	split.Offset = 0.4 // 上下比例，上面占 40%，下面占 60%

	return container.NewStack(split)

}

// centeredLabel 创建一个居中对齐的标签
func centeredLabel(text string) *widget.Label {
	label := widget.NewLabel(text)
	label.Alignment = fyne.TextAlignCenter
	return label
}

// rightedLabel 创建一个右对齐的标签
func rightedLabel(text string) *widget.Label {
	label := widget.NewLabel(text)
	label.Alignment = fyne.TextAlignTrailing
	return label
}

// createHistoryList 创建上传历史记录列表
func createHistoryList() *fyne.Container {
	// 上传历史记录 -- 可滚动
	historyList := container.NewVBox()
	// 表头
	header := container.NewGridWithColumns(6,
		widget.NewLabelWithStyle("数据集名称", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("图片名称", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("图片路径", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("图片大小", fyne.TextAlignTrailing, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("上传状态", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("上传时间", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
	)
	historyList.Add(header) // 添加表头

	// 获取上传历史记录
	topHistory := database.GetUploadHistory()
	// 排版输出
	for _, record := range topHistory {
		var statusColor color.Color
		switch record.UploadStatus {
		case "成功":
			statusColor = color.RGBA{R: 0, G: 128, B: 0, A: 255} // 绿色
		case "失败":
			statusColor = color.RGBA{R: 255, G: 0, B: 0, A: 255} // 红色
		default:
			statusColor = color.Black // 默认黑色
		}

		statusText := canvas.NewText(record.UploadStatus, statusColor)
		statusText.Alignment = fyne.TextAlignCenter

		row := container.NewGridWithColumns(6,
			widget.NewLabelWithStyle(record.DatasetName, fyne.TextAlignLeading, fyne.TextStyle{Bold: false}),
			widget.NewLabelWithStyle(record.ImageName, fyne.TextAlignLeading, fyne.TextStyle{Bold: false}),
			widget.NewLabelWithStyle(record.ImagePath, fyne.TextAlignLeading, fyne.TextStyle{Bold: false}),
			widget.NewLabelWithStyle(record.ImageSize, fyne.TextAlignTrailing, fyne.TextStyle{Bold: false}),
			statusText,
			widget.NewLabelWithStyle(record.UploadTime, fyne.TextAlignLeading, fyne.TextStyle{Bold: false}),
		)
		historyList.Add(row)
	}
	return historyList
}
