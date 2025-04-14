// models/dataset.go
package models

import "time"

// Dataset 表示一个数据集
type Dataset struct {
	ID          int       `json:"id"`          // 数据集 ID，主键
	Name        string    `json:"name"`        // 数据集名称
	Description string    `json:"description"` // 数据集描述
	ImageCount  int       `json:"image_count"` // 图片数量
	CreatedAt   time.Time `json:"created_at"`  // 创建时间
	UpdatedAt   time.Time `json:"updated_at"`  // 更新时间
	Status      int       `json:"status"`      // 数据集状态 0: 更新后未同步 1: 更新后已同步
	Cover       string    `json:"cover"`       // 添加封面字段
}
