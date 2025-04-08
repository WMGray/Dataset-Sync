// 3. 右侧： 上半部分 -- 一些查找、搜索操作（一个横条）， 下面部分 -- 部分数据集展示
package ui

// 主界面UI布局
// 1. 上方横条 软件信息
// 2. 左侧： 上半部分 -- 功能模块， 下面部分 -- 账户、设置
// 3. 右侧： 上半部分 -- 一些查找、搜索操作（一个横条）， 下面部分 -- 部分数据集展示
import (
	"dataset-sync/conf"
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// customTheme implements fyne.Theme
type customTheme struct {
	variant fyne.ThemeVariant
}

func (t *customTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	return theme.DefaultTheme().Color(name, t.variant)
}

func (t *customTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (t *customTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (t *customTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}

// ShowMainUI 显示主界面
func ShowMainUI() {
	// 确保配置已加载
	if conf.Conf.AppConfig.AppName == "" || conf.Conf.AppConfig.AppVersion == "" {
		fmt.Println("警告: 配置未正确加载，使用默认值")
		conf.Conf.AppConfig.AppName = "Dataset-Sync"
		conf.Conf.AppConfig.AppVersion = "1.0.0"
	}

	// 创建应用
	myApp := app.New()

	// 创建自定义主题
	customTheme := &customTheme{
		variant: theme.VariantLight,
	}

	myApp.Settings().SetTheme(customTheme)

	window := myApp.NewWindow(conf.Conf.AppConfig.AppName + " v" + conf.Conf.AppConfig.AppVersion)
	window.Resize(fyne.NewSize(1200, 800))

	// ===== 左侧功能区 =====
	// 功能区折叠状态
	isCollapsed := false
	// 主题切换状态
	isDarkMode := false

	// 上半部分：主要功能区
	topFuncs := container.NewVBox(
		widget.NewButton("数据集", func() {}),
	)

	// 下半部分：系统功能区
	bottomFuncs := container.NewVBox(
		widget.NewButton("账户", func() {}),
		widget.NewButton("设置", func() {}),
	)

	// 创建左侧整体布局（使用两个部分，中间用弹性空间分隔）
	leftContent := container.NewVBox(
		topFuncs,
		layout.NewSpacer(), // 添加弹性空间，将上下部分分开
		bottomFuncs,
	)

	// ===== 右侧功能区 =====
	// 搜索栏
	searchEntry := widget.NewEntry()
	searchEntry.SetPlaceHolder("搜索...")
	searchButton := widget.NewButtonWithIcon("", theme.SearchIcon(), func() {})
	filterButton := widget.NewButtonWithIcon("", theme.SettingsIcon(), func() {})
	sortButton := widget.NewButtonWithIcon("", theme.ListIcon(), func() {})
	moreButton := widget.NewButtonWithIcon("", theme.MoreHorizontalIcon(), func() {})

	searchBar := container.NewHBox(
		searchEntry,
		searchButton,
		filterButton,
		sortButton,
		moreButton,
	)

	// 创建数据集展示网格
	grid := container.NewGridWithColumns(4) // 每行4个项目

	// 添加示例数据集卡片
	for i := 0; i < 20; i++ {
		// 创建卡片容器
		card := widget.NewCard(
			fmt.Sprintf("数据集 %d", i+1), // 标题
			"示例描述",                     // 副标题
			container.NewVBox( // 内容
				widget.NewLabel("数据量: 1000"),
				widget.NewLabel("更新时间: 2024-04-07"),
			),
		)
		// 设置卡片最小尺寸
		card.Resize(fyne.NewSize(250, 200))
		grid.Add(card)
	}

	// 创建滚动容器包装网格
	gridScroll := container.NewScroll(grid)

	// 右侧整体布局
	rightSide := container.NewBorder(
		searchBar, // 顶部搜索栏
		nil,       // 底部无内容
		nil, nil,  // 左右无内容
		gridScroll, // 中间是网格
	)

	// 创建顶部控制栏（先创建一个空的，后面再添加按钮）
	controlBar := container.NewHBox()

	// 创建左侧主容器（包含顶部控制栏）
	leftSide := container.NewBorder(
		controlBar, // 顶部放置控制栏
		nil,        // 底部为空
		nil, nil,   // 左右为空
		leftContent, // 中间放置主要内容
	)

	// ===== 整体布局 =====
	// 使用HSplit分割左右两侧
	splitContainer := container.NewHSplit(leftSide, rightSide)
	splitContainer.Offset = 0.2 // 设置分割比例，使左侧更窄

	// 创建折叠按钮（使用箭头图标）
	var collapseBtn *widget.Button
	collapseBtn = widget.NewButtonWithIcon("", theme.NavigateBackIcon(), func() {
		isCollapsed = !isCollapsed
		if isCollapsed {
			collapseBtn.SetIcon(theme.NavigateNextIcon())
			leftContent.Hide()              // 隐藏所有功能区
			splitContainer.SetOffset(0.001) // 几乎完全折叠
		} else {
			collapseBtn.SetIcon(theme.NavigateBackIcon())
			leftContent.Show()            // 显示所有功能区
			splitContainer.SetOffset(0.2) // 恢复原始宽度
		}
		splitContainer.Refresh()
	})
	collapseBtn.Resize(fyne.NewSize(24, 24))

	// 创建主题切换按钮（使用太阳/月亮图标）
	var themeBtn *widget.Button
	themeBtn = widget.NewButtonWithIcon("", theme.ColorPaletteIcon(), func() {
		isDarkMode = !isDarkMode
		if isDarkMode {
			themeBtn.SetIcon(theme.ColorPaletteIcon())
			customTheme.variant = theme.VariantDark
		} else {
			themeBtn.SetIcon(theme.ColorPaletteIcon())
			customTheme.variant = theme.VariantLight
		}
		myApp.Settings().SetTheme(customTheme)
	})
	themeBtn.Resize(fyne.NewSize(24, 24))

	// 更新顶部控制栏
	controlBar.Add(collapseBtn)
	controlBar.Add(widget.NewSeparator())
	controlBar.Add(themeBtn)

	// 设置主容器
	mainContainer := container.NewBorder(nil, nil, nil, nil, splitContainer)

	// 显示窗口
	window.SetContent(mainContainer)
	window.ShowAndRun()
}
