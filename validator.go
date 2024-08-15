package validator

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"reflect"
	"strings"
)

// CustomValidator 是一个包装了 validator 实例的结构体
type CustomValidator struct {
	Validator *validator.Validate
}

// NewCustomValidator 创建一个新的自定义验证器
func NewCustomValidator() *CustomValidator {
	return &CustomValidator{Validator: validator.New()}
}

// Struct 使用自定义错误消息进行验证
func (cv *CustomValidator) Struct(i interface{}) error {
	err := cv.Validator.Struct(i)
	if err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			var errMsgs []string
			for _, fieldError := range validationErrors {
				customMsg := getCustomErrorMessage(i, fieldError)
				if customMsg != "" {
					errMsgs = append(errMsgs, customMsg)
				} else {
					errMsgs = append(errMsgs, fieldError.Error())
				}
			}
			return fmt.Errorf(strings.Join(errMsgs, "; "))
		}
	}
	return nil
}

// getCustomErrorMessage 获取字段的自定义错误消息
func getCustomErrorMessage(i interface{}, err validator.FieldError) string {
	field, ok := reflect.TypeOf(i).Elem().FieldByName(err.Field())
	if !ok {
		return ""
	}

	customMsgs := field.Tag.Get("err_msg")
	ruleMsgs := parseErrorMessages(customMsgs)

	return ruleMsgs[err.Tag()]
}

// parseErrorMessages 解析 err_msg 标签的内容，返回一个规则-错误消息映射
func parseErrorMessages(tag string) map[string]string {
	ruleMsgs := make(map[string]string)
	parts := strings.Split(tag, ";")
	for _, part := range parts {
		ruleAndMsg := strings.SplitN(part, ":", 2)
		if len(ruleAndMsg) == 2 {
			ruleMsgs[ruleAndMsg[0]] = ruleAndMsg[1]
		}
	}
	return ruleMsgs
}
