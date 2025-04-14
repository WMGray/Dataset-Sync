package ui

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// DatasetItem 表示数据集项
type DatasetItem struct {
	name       string
	count      int
	updateTime time.Time
}

var Datasets []DatasetItem                                 // 全局数据集列表
var grid *fyne.Container = container.NewGridWithColumns(3) // 全局网格容器
var curDatasets []DatasetItem                              // 当前数据集列表(网格中显示)

func CreateDatasetContainer() *fyne.Container {
	// 创建数据集列表
	for i := range 20 {
		Datasets = append(Datasets, DatasetItem{
			name:       fmt.Sprintf("数据集 %d", i+1),
			count:      1000 * (i + 1),
			updateTime: time.Now().AddDate(0, 0, -i),
		})
	}
	curDatasets = append([]DatasetItem{}, Datasets...) // 初始化当前数据集列表

	// 创建滚动容器
	gridScroll := container.NewScroll(grid)
	// 更新网格显示  -- 初始化
	updateGrid(Datasets)

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
		nil, nil, nil,
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

	return container.NewBorder(topNav, nil, nil, nil, gridScroll)
}

// searchDatasets 搜索数据集
func searchDatasets(keyword string) {
	if keyword == "" {
		curDatasets = append([]DatasetItem{}, Datasets...) // 重置当前数据集列表
		updateGrid(Datasets)
		return
	}
	// 搜索过滤
	filtered := []DatasetItem{}
	for _, item := range Datasets {
		if strings.Contains(strings.ToLower(item.name), strings.ToLower(keyword)) {
			filtered = append(filtered, item)
		}
	}
	// 更新当前数据集列表
	curDatasets = filtered
	updateGrid(curDatasets)
}

// sortDatasets 对数据集进行排序
func sortDatasets(option string) {
	sorted := make([]DatasetItem, len(curDatasets))
	copy(sorted, curDatasets)

	// 根据选择的排序类型和顺序进行排序
	switch option {
	case "名称":
		sort.Slice(sorted, func(i, j int) bool { return sorted[i].name < sorted[j].name })
	case "数量":
		sort.Slice(sorted, func(i, j int) bool { return sorted[i].count < sorted[j].count })
	case "日期":
		sort.Slice(sorted, func(i, j int) bool { return sorted[i].updateTime.Before(sorted[j].updateTime) })
	}
	// 更新全局数据集并刷新网格显示
	updateGrid(sorted)
}

// updateGrid 更新网格显示
func updateGrid(items []DatasetItem) {
	grid.Objects = nil // 清空网格
	if len(items) == 0 {
		grid.Add(widget.NewLabel("无该数据集"))
		return
	}
	for _, item := range items {
		card := widget.NewCard(
			item.name,
			"",
			container.NewVBox(
				widget.NewLabel(fmt.Sprintf("数据量: %d", item.count)),
				widget.NewLabel(fmt.Sprintf("更新时间: %s", item.updateTime.Format("2006-01-02"))),
			),
		)
		grid.Add(card)
	}
	grid.Refresh()
}
