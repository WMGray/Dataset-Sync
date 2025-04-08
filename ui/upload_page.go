package ui

import (
	"fmt"
	"image/color"
	"os"
	"path/filepath"
	"sort"
	"syscall"
	"time"
	"unsafe"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// DropZone 自定义拖放区域组件
type DropZone struct {
	widget.BaseWidget
	background *canvas.Rectangle
	content    *fyne.Container
	onDropped  func([]fyne.URI)
}

// NewDropZone 创建新的拖放区域
func NewDropZone(content *fyne.Container, onDropped func([]fyne.URI)) *DropZone {
	d := &DropZone{
		background: canvas.NewRectangle(color.NRGBA{R: 240, G: 240, B: 240, A: 255}),
		content:    content,
		onDropped:  onDropped,
	}
	d.ExtendBaseWidget(d)
	return d
}

// CreateRenderer 实现自定义渲染
func (d *DropZone) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(container.NewStack(d.background, d.content))
}

// DragEnter 处理拖入事件
func (d *DropZone) DragEnter() {
	d.background.FillColor = color.NRGBA{R: 220, G: 220, B: 220, A: 255}
	d.background.Refresh()
}

// DragEnd 处理拖出事件
func (d *DropZone) DragEnd() {
	d.background.FillColor = color.NRGBA{R: 240, G: 240, B: 240, A: 255}
	d.background.Refresh()
}

// FileItem 文件项结构
type FileItem struct {
	Name    string
	Path    string
	IsDir   bool
	Size    int64
	ModTime string
}

// FileSelectDialog 自定义文件选择对话框
type FileSelectDialog struct {
	window     fyne.Window
	onSelected func(fyne.URI)
	currentDir string
	fileList   *widget.List
	items      []FileItem
	selected   int
}

// NewFileSelectDialog 创建新的文件选择对话框
func NewFileSelectDialog(window fyne.Window, onSelected func(fyne.URI)) *FileSelectDialog {
	d := &FileSelectDialog{
		window:     window,
		onSelected: onSelected,
		currentDir: ".",
	}
	return d
}

// loadDirectory 加载目录内容
func (d *FileSelectDialog) loadDirectory(path string) error {
	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	d.items = make([]FileItem, 0, len(entries))
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue
		}

		item := FileItem{
			Name:    entry.Name(),
			Path:    filepath.Join(path, entry.Name()),
			IsDir:   entry.IsDir(),
			Size:    info.Size(),
			ModTime: info.ModTime().Format("2006-01-02 15:04:05"),
		}
		d.items = append(d.items, item)
	}

	// 排序：文件夹在前，文件在后
	sort.Slice(d.items, func(i, j int) bool {
		if d.items[i].IsDir != d.items[j].IsDir {
			return d.items[i].IsDir
		}
		return d.items[i].Name < d.items[j].Name
	})

	d.currentDir = path
	if d.fileList != nil {
		d.fileList.Refresh()
	}
	return nil
}

// Show 显示文件选择对话框
func (d *FileSelectDialog) Show() {
	// 创建文件列表
	d.fileList = widget.NewList(
		func() int { return len(d.items) },
		func() fyne.CanvasObject { return widget.NewLabel("Template") },
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			label := obj.(*widget.Label)
			item := d.items[id]
			if item.IsDir {
				label.SetText(fmt.Sprintf("📁 %s", item.Name))
			} else {
				label.SetText(fmt.Sprintf("📄 %s (%d bytes)", item.Name, item.Size))
			}
		},
	)

	// 处理双击事件
	d.fileList.OnSelected = func(id widget.ListItemID) {
		d.selected = int(id)
		item := d.items[id]
		if item.IsDir {
			d.loadDirectory(item.Path)
		} else {
			uri := storage.NewFileURI(item.Path)
			d.onSelected(uri)
			d.window.Close()
		}
	}

	// 创建路径输入框
	pathEntry := widget.NewEntry()
	pathEntry.SetPlaceHolder("输入路径或选择位置")
	pathEntry.SetText(d.currentDir)
	pathEntry.OnSubmitted = func(path string) {
		if err := d.loadDirectory(path); err != nil {
			dialog.ShowError(err, d.window)
		}
	}

	// 创建导航按钮
	backButton := widget.NewButtonWithIcon("", theme.NavigateBackIcon(), func() {
		parentDir := filepath.Dir(d.currentDir)
		if err := d.loadDirectory(parentDir); err != nil {
			dialog.ShowError(err, d.window)
		}
		pathEntry.SetText(d.currentDir)
	})

	// 创建刷新按钮
	refreshButton := widget.NewButtonWithIcon("", theme.ViewRefreshIcon(), func() {
		if err := d.loadDirectory(d.currentDir); err != nil {
			dialog.ShowError(err, d.window)
		}
	})

	// 创建按钮容器
	buttonContainer := container.NewHBox(
		backButton,
		refreshButton,
		widget.NewSeparator(),
		widget.NewButton("选择", func() {
			if len(d.items) > 0 && d.selected >= 0 {
				item := d.items[d.selected]
				uri := storage.NewFileURI(item.Path)
				d.onSelected(uri)
				d.window.Close()
			}
		}),
		widget.NewButton("取消", func() {
			d.window.Close()
		}),
	)

	// 创建主容器
	content := container.NewBorder(
		container.NewVBox(pathEntry, buttonContainer),
		nil,
		nil,
		nil,
		d.fileList,
	)

	// 加载初始目录
	if err := d.loadDirectory(d.currentDir); err != nil {
		dialog.ShowError(err, d.window)
	}

	// 创建对话框窗口
	dialog.ShowCustom("选择文件或文件夹", "关闭", content, d.window)
}

// UploadPage 上传页面组件
type UploadPage struct {
	container *fyne.Container
	window    fyne.Window
}

// NewUploadPage 创建新的上传页面
func NewUploadPage(window fyne.Window) *UploadPage {
	page := &UploadPage{
		window: window,
	}
	page.init()
	return page
}

// showFileDialog 显示文件选择对话框
func (p *UploadPage) showFileDialog() {
	// 使用 Windows API 创建文件选择对话框
	var (
		modcomdlg32             = syscall.NewLazyDLL("comdlg32.dll")
		modshell32              = syscall.NewLazyDLL("shell32.dll")
		modole32                = syscall.NewLazyDLL("ole32.dll")
		procGetOpenFileName     = modcomdlg32.NewProc("GetOpenFileNameW")
		procSHBrowseForFolder   = modshell32.NewProc("SHBrowseForFolderW")
		procSHGetPathFromIDList = modshell32.NewProc("SHGetPathFromIDListW")
		procCoTaskMemFree       = modole32.NewProc("CoTaskMemFree")
	)

	// 创建按钮容器
	buttonContainer := container.NewHBox(
		widget.NewButton("选择文件", func() {
			// 创建缓冲区
			buf := make([]uint16, 260)

			// 创建 OPENFILENAME 结构
			var ofn struct {
				StructSize    uint32
				Owner         uintptr
				Instance      uintptr
				Filter        *uint16
				CustomFilter  *uint16
				MaxCustFilter uint32
				FilterIndex   uint32
				File          *uint16
				MaxFile       uint32
				FileTitle     *uint16
				MaxFileTitle  uint32
				InitialDir    *uint16
				Title         *uint16
				Flags         uint32
				FileOffset    uint16
				FileExtension uint16
				DefExt        *uint16
				CustData      uintptr
				FnHook        uintptr
				TemplateName  *uint16
				PvReserved    uintptr
				DwReserved    uint32
				FlagsEx       uint32
			}

			// 设置标题和过滤器
			title := "选择文件"
			filter := "所有文件\000*.*\000"

			// 转换为 UTF-16
			titlePtr, _ := syscall.UTF16PtrFromString(title)
			filterPtr, _ := syscall.UTF16PtrFromString(filter)

			// 初始化结构
			ofn.StructSize = uint32(unsafe.Sizeof(ofn))
			ofn.Owner = 0
			ofn.Filter = filterPtr
			ofn.File = &buf[0]
			ofn.MaxFile = uint32(len(buf))
			ofn.Title = titlePtr
			ofn.Flags = 0x00080000 | 0x00001000 | 0x00000800 // OFN_FILEMUSTEXIST | OFN_PATHMUSTEXIST | OFN_EXPLORER

			// 调用 Windows API
			ret, _, _ := procGetOpenFileName.Call(uintptr(unsafe.Pointer(&ofn)))

			if ret != 0 {
				// 获取选择的文件路径
				path := syscall.UTF16ToString(buf[:])
				uri := storage.NewFileURI(path)
				fmt.Printf("Selected file: %s\n", uri.Path())
				// 模拟上传文件
				p.simulateUpload(uri)
			}
		}),
		widget.NewButton("选择文件夹", func() {
			// 创建 BROWSEINFO 结构
			var bi struct {
				Owner        uintptr
				Root         uintptr
				DisplayName  *uint16
				Title        *uint16
				Flags        uint32
				CallbackFunc uintptr
				LParam       uintptr
				Image        int32
			}

			// 设置标题
			title := "选择文件夹"
			titlePtr, _ := syscall.UTF16PtrFromString(title)

			// 初始化结构
			bi.Title = titlePtr
			bi.Flags = 0x00000001 | 0x00000040 // BIF_RETURNONLYFSDIRS | BIF_NEWDIALOGSTYLE

			// 调用 Windows API
			pidl, _, _ := procSHBrowseForFolder.Call(uintptr(unsafe.Pointer(&bi)))
			if pidl != 0 {
				// 获取选择的文件夹路径
				pathBuf := make([]uint16, 260)
				procSHGetPathFromIDList.Call(pidl, uintptr(unsafe.Pointer(&pathBuf[0])))
				path := syscall.UTF16ToString(pathBuf[:])

				// 释放 PIDL
				procCoTaskMemFree.Call(pidl)

				uri := storage.NewFileURI(path)
				fmt.Printf("Selected folder: %s\n", uri.Path())
				// 模拟上传文件夹
				p.simulateUpload(uri)
			}
		}),
	)

	// 显示按钮容器
	dialog.ShowCustom("选择文件或文件夹", "关闭", buttonContainer, p.window)
}

// simulateUpload 模拟上传文件或文件夹
func (p *UploadPage) simulateUpload(uri fyne.URI) {
	// 显示上传进度对话框
	progress := dialog.NewProgress("上传中", "正在上传...", p.window)
	progress.Show()

	// 模拟上传进度
	go func() {
		for i := 0.0; i <= 1.0; i += 0.1 {
			progress.SetValue(i)
			time.Sleep(200 * time.Millisecond)
		}
		progress.Hide()
		dialog.ShowInformation("上传完成", fmt.Sprintf("已上传: %s", uri.Path()), p.window)
	}()
}

// init 初始化上传页面
func (p *UploadPage) init() {
	// 创建提示文本
	dropText := widget.NewLabel("拖放文件到这里\n或点击选择文件")
	dropText.Alignment = fyne.TextAlignCenter

	// 创建选择按钮
	selectButton := widget.NewButtonWithIcon("选择文件或文件夹", theme.FolderOpenIcon(), func() {
		p.showFileDialog()
	})

	// 创建垂直布局容器
	content := container.NewVBox(
		dropText,
		selectButton,
	)

	// 创建拖放区域
	dropZone := NewDropZone(container.NewCenter(content), func(uris []fyne.URI) {
		// 处理拖放的文件
		for _, uri := range uris {
			fmt.Printf("Dropped file: %s\n", uri.Path())
			// 模拟上传拖放的文件
			p.simulateUpload(uri)
		}
	})

	// 创建主容器
	p.container = container.NewMax(dropZone)
}

// Container 返回页面容器
func (p *UploadPage) Container() fyne.CanvasObject {
	return p.container
}
