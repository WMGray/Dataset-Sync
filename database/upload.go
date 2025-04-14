package database

import (
	"dataset-sync/models"
	"strconv"
)

// 获取上传历史记录
func GetUploadHistory() []*models.UploadDetails {
	// 模拟数据
	var uploadHistory []*models.UploadDetails
	for i := 0; i < 200; i++ {
		// 使用rand包来随机生成
		uploadHistory = append(uploadHistory, &models.UploadDetails{
			ImageName:   "image" + strconv.Itoa(i),
			DatasetName: strconv.Itoa(i),
			ImagePath:   "path/to/image" + strconv.Itoa(i),
			ImageSize:   strconv.Itoa(i) + "MB",
			UploadTime:  "2023-10-01 12:00:00",
			// 状态根据随机值确定
			UploadStatus: func() string {
				if i%2 == 0 {
					return "成功"
				}
				return "失败"
			}(),
		})
	}
	return uploadHistory
}
