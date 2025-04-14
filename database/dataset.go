package database

import (
	"dataset-sync/models"
	"fmt"
	"math/rand"
	"time"
)

func GetDatasets() []*models.Dataset {
	// 从数据库中获取数据集列表
	// 返回一个包含所有数据集的切片
	// 模拟数据 -- 数据集切片GetDatasets
	var datasets []*models.Dataset
	for i := range 20 {
		datasets = append(datasets, &models.Dataset{
			ID:          i,
			Name:        fmt.Sprintf("数据集%d", i),
			Description: fmt.Sprintf("描述%d", i),
			// 随机数量 -- 1000 + 随机数 使用round
			ImageCount: rand.Intn(i*1000 + 100),
			CreatedAt:  time.Now(),
			// 随机更新时间 -- 当前时间 + 随机时间
			UpdatedAt: time.Now().Add(time.Duration(i*20) * time.Hour),
			Status:    1,
			Cover:     "C:\\Users\\WMGray\\OneDrive\\Dev\\Workspaces\\Dataset-Sync\\ui\\luna.jpg",
		})
	}
	return datasets
}
