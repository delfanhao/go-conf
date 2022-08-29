package conf

import (
	"os"
	"strings"
)

// getValFromEnv 从环境变量中获取对应的值
func (ctx *configContext) getEnvVal(_, key string) (interface{}, bool) {
	realKey := getFormatKey(key)
	result := os.Getenv(realKey)
	ok := len(result) > 0

	if ok {
		trace("ENV - found: %s = %v", key, result)
	}
	// 如果是数组定义
	if strings.HasPrefix(result, "[") && strings.HasSuffix(result, "]") {
		s := strings.TrimPrefix(result, "[")
		s = strings.TrimSuffix(s, "]")
		return str2array(s)
	}
	return result, ok
}

// str2array 字符串转数组
func str2array(s string) ([]string, bool) {
	sArr := strings.Split(s, ",")
	result := make([]string, 0)
	for _, stmp := range sArr {
		var ss, sp string

		if strings.HasPrefix(stmp, "'") && strings.HasSuffix(stmp, "'") {
			sp = "'"
		} else if strings.HasPrefix(stmp, "\"") && strings.HasSuffix(stmp, "\"") {
			sp = "\""
		} else {
			sp = ""
		}

		if sp != "" {
			ss = strings.TrimPrefix(stmp, sp)
			ss = strings.TrimSuffix(ss, sp)
		} else {
			ss = stmp
		}

		if strings.Trim(ss, " ") != "" {
			result = append(result, ss)
		}
	}

	return result, len(result) > 0
}
