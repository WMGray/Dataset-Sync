package ui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"image/color"

	"dataset-sync/database"
	"dataset-sync/ui/components"
)

func CreateUpload(w fyne.Window) *fyne.Container {
	// 创建搜索输入框
	searchEntry := createSearchEntry(w)

	// 拖放区域
	content := createUploadPrompt(w)
	dropArea := components.CreateDropArea(w, container.NewCenter(content))
	dropArea.Resize(fyne.NewSize(500, 200))
	dropArea.Move(fyne.NewPos(0, 30)) // 往下偏移一下，避免被 label 挡住

	// 上方区域  -- 搜索框和拖放区域
	topSection := container.NewBorder(
		searchEntry, nil, nil, nil,
		dropArea,
	)

	// 下方区域 -- 上传历史记录
	historyList := createHistoryList()
	bottomSection := container.NewVScroll(historyList)

	split := container.NewVSplit(topSection, bottomSection)
	split.Offset = 0.4 // 上下比例，上面占 40%，下面占 60%

	return container.NewStack(split)
}

// createSearchEntry 创建搜索输入框
func createSearchEntry(window fyne.Window) *fyne.Container {
	searchEntry := widget.NewEntry()
	searchEntry.SetPlaceHolder("请粘贴图片的链接")
	searchEntry.Wrapping = fyne.TextWrapOff

	searchButton := widget.NewButtonWithIcon("搜索", theme.SearchIcon(), func() {
		keyword := searchEntry.Text
		if keyword == "" {
			dialog.ShowInformation("提示", "请输入图片链接", window)
			return
		}
		fmt.Printf("Captured URL: %s\n", keyword)
		dialog.ShowInformation("URL 已捕获", fmt.Sprintf("已捕获 URL: %s", keyword), window)
	})
	searchButton.Importance = widget.HighImportance

	searchEntryContainer := container.NewBorder(
		nil, nil, nil, searchButton,
		searchEntry,
	)
	return searchEntryContainer
}

// createUploadPrompt 创建拖放区域提示
func createUploadPrompt(window fyne.Window) *fyne.Container {
	icon := canvas.NewImageFromResource(theme.UploadIcon())
	icon.FillMode = canvas.ImageFillContain
	icon.SetMinSize(fyne.NewSize(24, 24))

	promptText := widget.NewLabel("将图片放到此处")
	promptText.Alignment = fyne.TextAlignCenter

	uploadLink := widget.NewHyperlink("上传文件", nil)
	uploadLink.OnTapped = func() {
		dialog.ShowFileOpen(func(file fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, window)
				return
			}
			if file == nil {
				return
			}
			defer file.Close()

			filePath := file.URI().Path()
			fmt.Printf("Captured file: %s\n", filePath)
			dialog.ShowInformation("上传成功", fmt.Sprintf("已捕获文件: %s", filePath), window)
		}, window)
	}

	promptContainer := container.NewVBox(
		container.NewCenter(container.NewHBox(
			icon,
			promptText,
		)),
		container.NewCenter(widget.NewLabel("或")),
		container.NewCenter(uploadLink),
	)

	return promptContainer
}

// createHistoryList 创建上传历史记录列表
func createHistoryList() *fyne.Container {
	// 上传历史记录 -- 可滚动
	historyList := container.NewVBox()
	// 表头
	header := container.NewGridWithColumns(7,
		// TODO 当前字体不支持加粗
		widget.NewLabelWithStyle("数据集名称", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("图片名称", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("图片路径", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("图片大小", fyne.TextAlignTrailing, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("上传进度", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("上传状态", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("上传时间", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
	)
	historyList.Add(header) // 添加表头

	// 获取上传历史记录
	topHistory := database.GetUploadHistory()
	// 排版输出
	for i, record := range topHistory {
		// 上传进度条
		progressBar := components.NewCustomProgressBar()
		progressBar.SetValue(getProgressValue(record.UploadStatus))

		// 根据状态设置进度条颜色
		switch record.UploadStatus {
		case "成功":
			progressBar.SetBarColor(color.RGBA{R: 0, G: 128, B: 0, A: 255}) // 绿色
		case "失败":
			progressBar.SetBarColor(color.RGBA{R: 255, G: 0, B: 0, A: 255}) // 红色
		default:
			progressBar.SetBarColor(color.RGBA{R: 0, G: 122, B: 255, A: 255}) // 蓝色（进行中）
		}

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

		// 创建行数据
		rowContent := container.NewGridWithColumns(7,
			widget.NewLabelWithStyle(record.DatasetName, fyne.TextAlignLeading, fyne.TextStyle{Bold: false}),
			widget.NewLabelWithStyle(record.ImageName, fyne.TextAlignLeading, fyne.TextStyle{Bold: false}),
			widget.NewLabelWithStyle(record.ImagePath, fyne.TextAlignLeading, fyne.TextStyle{Bold: false}),
			widget.NewLabelWithStyle(record.ImageSize, fyne.TextAlignTrailing, fyne.TextStyle{Bold: false}),
			container.NewCenter(progressBar),
			statusText,
			widget.NewLabelWithStyle(record.UploadTime, fyne.TextAlignLeading, fyne.TextStyle{Bold: false}),
		)

		// 添加交替背景色
		var background *canvas.Rectangle
		if i%2 == 0 {
			background = canvas.NewRectangle(color.White) // 偶数行（包括 0）为白色
		} else {
			background = canvas.NewRectangle(color.NRGBA{R: 240, G: 240, B: 240, A: 255}) // 奇数行为浅灰色
		}
		background.FillColor = background.FillColor
		background.SetMinSize(fyne.NewSize(600, 20)) // 假设总宽度 600 像素，行高 20 像素，可调整

		// 使用 VBox 组合背景和行内容
		row := container.NewVBox(
			container.NewStack(
				background,
				rowContent,
			),
		)
		historyList.Add(row)
	}
	return historyList
}

// getProgressValue 根据上传状态返回进度值
func getProgressValue(status string) float64 {
	switch status {
	case "成功":
		return 1.0
	case "失败":
		return 0.0
	default:
		return 0.0 // 模拟进行中状态
	}
}
