package conf

import (
	"os"
	"path/filepath"
	"strings"
)

// getAppInfo 获取应用名称
func getAppInfo() (string, string) {
	return filepath.Dir(os.Args[0]), filepath.Base(os.Args[0])
}

func getFormatKey(key string) string {
	newKey := strings.ToUpper(key)
	return strings.ReplaceAll(newKey, ".", "_")
}
