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
	content := widget.NewLabel("拖放文件到这里")
	content.Alignment = fyne.TextAlignCenter
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

// showContent 显示指定的内容
func createUploadButton(window fyne.Window) *widget.Button {
	uploadButton := widget.NewButtonWithIcon("上传", theme.UploadIcon(), func() {
		// 调用系统文件选择器
		dialog.ShowFileOpen(func(file fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, window)
				return
			}
			if file == nil {
				// 用户取消选择
				return
			}
			defer file.Close()

			// 捕获文件路径
			filePath := file.URI().Path()
			fmt.Printf("Captured file: %s\n", filePath)

			// 显示上传成功的提示
			dialog.ShowInformation("上传成功", fmt.Sprintf("已捕获文件: %s", filePath), window)
		}, window)
	})
	uploadButton.Importance = widget.HighImportance // 突出按钮样式
	return uploadButton
}

// createHistoryList 创建上传历史记录列表
func createHistoryList() *fyne.Container {
	// 上传历史记录 -- 可滚动
	historyList := container.NewVBox()
	// 表头
	header := container.NewGridWithColumns(6,
		// TODO 当前字体不支持加粗
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
