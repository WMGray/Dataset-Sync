package utils

import (
	"dataset-sync/conf"
	"errors"
	"reflect"
)

// ChangeSettings 修改设置中的内容
func ChangeSettings(settings interface{}, key string, value interface{}) error {
	// 不开放API给外部使用，不做额外的错误处理之类的
	val := reflect.ValueOf(settings).Elem()
	field := val.FieldByName(key)

	// Update field
	if field.IsValid() && field.CanSet() && reflect.TypeOf(value) == field.Type() {
		field.Set(reflect.ValueOf(value))
		err := conf.SaveConfig()
		if err != nil {
			return errors.New("failed to save configuration: " + err.Error())
		}
		return nil
	}
	return errors.New("failed to update field: " + key)
}
