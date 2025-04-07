# 组织架构
- 模块	主要功能	主要文件
- ui/	构建桌面 UI，文件上传和统计展示	main_ui.go, upload.go, stats.go
- upload/	处理图片上传、去重和存储	upload.go, process.go, storage.go
- database/	数据库操作和缓存	mysql.go, redis.go, sqlite.go
- utils/	通用工具函数，如哈希计算	hash.go, file.go
- exporter/	上传结果到 GitHub 和 Hugging Face	github.go, huggingface.go
- models/	定义数据结构	image.go, upload_result.go
- config/	配置文件和常量	config.go