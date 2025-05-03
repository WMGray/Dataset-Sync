package ui

import (
	"dataset-sync/database"
	"dataset-sync/models"
	"fyne.io/fyne/v2/dialog"
	"sort"
	"strings"

	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"

	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

var cardSize = fyne.NewSize(240, 500)                           // 每张卡片尺寸
var Datasets []*models.Dataset                                  // 全局数据集列表
var grid *fyne.Container = container.NewGridWrap(cardSize, nil) // 全局网格容器grid
var curDatasets []*models.Dataset                               // 当前数据集列表
var datasetCards = make(map[string]fyne.CanvasObject)           // 存储数据集卡片的映射

func CreateDatasets() *fyne.Container {
	gridScroll := container.NewScroll(grid) // 创建滚动容器
	// 获取数据集列表
	Datasets = database.GetDatasets()                      // 从SQL中获取Dataset切片
	curDatasets = append([]*models.Dataset{}, Datasets...) // 初始化当前数据集列表
	// 初始化卡片
	initDatasetCards()

	// 初始化网格
	updateGrid()

	// 创建新增数据集按钮
	addDatasetButton := widget.NewButtonWithIcon("", theme.ContentAddIcon(), func() {
		// TODO: 弹出新增数据集窗口吧
		// 弹出一个提示框用于捕获数据集信息
		// 这里可以使用一个简单的对话框来捕获数据集名称 并 输出
		nameEntry := widget.NewEntry()
		nameEntry.SetPlaceHolder("请输入要新创建的数据集名称")
		nameEntry.Resize(fyne.NewSize(400, 100))
		// 构造一个 FormItem 数组（这里只有一个字段）
		formItems := []*widget.FormItem{
			widget.NewFormItem("名称", nameEntry),
		} // 弹出 Form 对话框：标题“新增数据集”，确认按钮“添加”，取消按钮“取消”
		formDialog := dialog.NewForm("新增数据集", "添加", "取消", formItems, func(confirmed bool) {
			if !confirmed {
				return
			}
			name := nameEntry.Text
			if name == "" {
				dialog.ShowError(fmt.Errorf("名称不能为空"), ui.window)
				return
			}
			dialog.ShowCustom("数据集添加成功！", "确定",
				widget.NewLabel("新数据集名称: "+name), ui.window)
			fmt.Println("新数据集：", name)
		}, ui.window)

		// 设置对话框的大小
		formDialog.Resize(fyne.NewSize(400, 200))
		formDialog.Show()
	})

	// 创建搜索组件
	searchEntry := widget.NewEntry()           // 初始化搜索输入框
	searchEntry.Resize(fyne.NewSize(400, 400)) // 增加搜索框的宽度和高度
	searchEntry.OnChanged = func(keyword string) {
		searchDatasets(keyword)
	}
	searchEntry.Hide()

	// 创建搜索按钮
	searchButton := widget.NewButtonWithIcon("搜索", theme.SearchIcon(), func() {
		if searchEntry.Hidden {
			searchEntry.Show()
			searchEntry.Refresh()
		} else {
			searchEntry.Hide()
			searchEntry.SetText("") // 清空搜索内容
			searchDatasets("")      // 重置数据集显示
		}
	})
	searchButton.Resize(fyne.NewSize(40, 40)) // 调整搜索按钮大小

	// 创建排序复选菜单 -- 两个单选框组合
	sortTypeMenu := widget.NewSelect([]string{"名称", "数量", "日期"}, func(selected string) {
		sortDatasets(selected)
	})
	sortTypeMenu.PlaceHolder = "排序"

	// 控件尺寸设置
	searchEntry.Resize(fyne.NewSize(400, 50)) // 只让 Entry 大
	searchButton.Resize(fyne.NewSize(40, 40)) // 也可以不设置
	sortTypeMenu.Resize(fyne.NewSize(100, 40))

	// 用 Max 包住 Entry，防止布局压缩它
	searchEntryWrapper := container.NewStack(searchEntry)

	// 顶部导航栏布局（3列）
	topNav := container.NewBorder(
		nil, nil,
		addDatasetButton,
		container.NewGridWithColumns(3,
			searchEntryWrapper, // 用 wrapper 控制大小
			searchButton,
			sortTypeMenu,
		),
	)

	// 使用点击事件隐藏搜索框
	searchEntry.OnChanged = func(keyword string) {
		searchDatasets(keyword)
		if keyword == "" {
			searchEntry.Hide()
		}
	}

	// 整体布局：顶部按钮 + 滚动容器
	return container.NewBorder(topNav, nil, nil, nil, gridScroll)
}

// searchDatasets 搜索数据集
func searchDatasets(keyword string) {
	if keyword == "" {
		curDatasets = append([]*models.Dataset{}, Datasets...) // 重置当前数据集列表
		updateGrid()
		return
	}
	// 搜索过滤
	var filtered []*models.Dataset
	for _, item := range Datasets {
		if strings.Contains(strings.ToLower(item.Name), strings.ToLower(keyword)) {
			filtered = append(filtered, item)
		}
	}
	// 更新当前数据集列表
	curDatasets = filtered
	updateGrid()
}

// sortDatasets 对数据集进行排序
func sortDatasets(option string) {
	sorted := make([]*models.Dataset, len(curDatasets))
	copy(sorted, curDatasets)

	// 根据选择的排序类型和顺序进行排序
	switch option {
	case "名称":
		sort.Slice(sorted, func(i, j int) bool { return sorted[i].Name < sorted[j].Name })
	case "数量":
		sort.Slice(sorted, func(i, j int) bool { return sorted[i].ImageCount < sorted[j].ImageCount })
	case "日期":
		sort.Slice(sorted, func(i, j int) bool { return sorted[i].UpdatedAt.Before(sorted[j].UpdatedAt) })
	}
	// 更新全局数据集并刷新网格显示
	curDatasets = sorted
	updateGrid()
}

// updateGrid 更新网格显示
func updateGrid() {
	grid.Objects = nil // 清空网格对象
	if len(curDatasets) == 0 {
		grid.Add(widget.NewLabel("无该数据集"))
		return
	}

	var cards []fyne.CanvasObject
	for _, item := range curDatasets {
		cards = append(cards, datasetCards[item.Name])
	}
	grid.Objects = cards
	grid.Refresh()
}

// initDatasetCards 初始化数据集卡片
func initDatasetCards() {
	// 初始化数据集卡片
	for _, ds := range Datasets {
		card := createDatasetCard(ds)
		datasetCards[ds.Name] = card
	}
}

// createDatasetCard 创建数据集卡片
func createDatasetCard(ds *models.Dataset) fyne.CanvasObject {
	// 获取封面图
	thumbnail := getCover(ds.Cover)

	// 卡牌内容
	content := container.NewVBox(
		// 使用 container.NewPadded 为封面图片添加内边距
		container.NewPadded(thumbnail),
		widget.NewLabelWithStyle(ds.Name, fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabel(fmt.Sprintf("图片数量: %d", ds.ImageCount)),
		widget.NewLabel(fmt.Sprintf("最后更新日期: %s", ds.UpdatedAt.Format("2006-01-02"))),
		// 状态 0 -- 未更新 1 -- 已更新
		widget.NewLabel(fmt.Sprintf("状态: %s", func() string {
			if ds.Status == 0 {
				return "未更新"
			}
			return "已更新"
		}())),
	)

	// 使用 container.NewStack 确保背景和内容完全重叠
	cardContainer := container.NewStack(content)
	cardContainer.Resize(fyne.NewSize(240, 300)) // 设置卡片大小

	return cardContainer
}

// getCover 获取数据集封面图
func getCover(filepath string) *canvas.Image {
	// 加载图片
	thumbnail := canvas.NewImageFromFile(filepath)
	thumbnail.FillMode = canvas.ImageFillContain // 或者 ImageFillStretch，确保缩放适应容器
	// 设置图片最小尺寸
	thumbnail.SetMinSize(fyne.NewSize(240, 300))

	return thumbnail
}
