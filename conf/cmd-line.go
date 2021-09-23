package conf

import (
	"os"
	"strings"
)

// getCmdVal 从命令行获取指定key的参数
func (ctx *configContext) getCmdVal(_, key string) (interface{}, bool) {
	key = getFormatKey(key)
	result, ok := "", false
	for i := 1; i < len(os.Args); i++ {
		flag := os.Args[i]
		if flag[0] == '-' && strings.Index(flag, key+"=") == 1 {
			result, ok = flag[len(key)+2:], true
		}
	}
	if ok {
		trace("CMD - found: %s = %v", key, result)
	}
	return result, ok
}
