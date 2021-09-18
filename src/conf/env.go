package conf

import (
	"os"
)

// getValFromEnv 从环境变量中获取对应的值
func (ctx *configContext) getEnvVal(_, key string) (interface{}, bool) {
	realKey := getFormatKey(key)
	result := os.Getenv(realKey)
	ok := len(result) > 0
	return result, ok
}
