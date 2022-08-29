package conf

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// getAppInfo 获取应用名称
func getAppInfo() (string, string) {
	name := filepath.Base(os.Args[0])
	// for windows
	name = strings.TrimRight(name, ".exe")
	name = strings.TrimRight(name, ".com")
	return filepath.Dir(os.Args[0]), filepath.Base(os.Args[0])
}

func getFormatKey(key string) string {
	newKey := strings.ToUpper(key)
	return strings.ReplaceAll(newKey, ".", "_")
}

func convToInt64(val interface{}) *int64 {
	switch val.(type) {
	case string:
		if v, ok := val.(string); ok {
			result, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				return nil
			}
			return &result
		}
	case int:
		re := int64(val.(int))
		return &re
	case float32, float64:
		re := int64(val.(float64))
		return &re
	}

	return nil
}

var TRACE = false

// trace 调试信息
func trace(msg string, params ...interface{}) {
	if !TRACE {
		return
	}
	t := time.Now()
	ts := t.Format("2006-01-02 15:04:05")
	prefix := fmt.Sprintf("%s - [T] ", ts)
	out := prefix + msg
	fmt.Println(fmt.Sprintf(out, params...))
}

func arrStr2int(sArr []string) []int {
	iArr := make([]int, len(sArr))
	for idx, s := range sArr {
		if i, e := strconv.Atoi(s); e == nil {
			iArr[idx] = i
		}
	}

	return iArr
}
func arrStr2int32(sArr []string) []int32 {
	iArr := make([]int32, len(sArr))
	for idx, s := range sArr {
		if i, e := strconv.Atoi(s); e == nil {
			iArr[idx] = int32(i)
		}
	}

	return iArr
}
func arrStr2int64(sArr []string) []int64 {
	iArr := make([]int64, len(sArr))
	for idx, s := range sArr {
		fmt.Printf("%d-[%s]\n", idx, s)
		if i, e := strconv.Atoi(s); e == nil {
			iArr[idx] = int64(i)
		}
	}

	return iArr
}

func arrStr2Float(sArr []string) []float64 {
	iArr := make([]float64, len(sArr))
	for idx, s := range sArr {
		if i, e := strconv.ParseFloat(s, 64); e == nil {
			iArr[idx] = i
		}
	}

	return iArr
}

func arrStr2Bool(sArr []string) []bool {
	iArr := make([]bool, len(sArr))
	for idx, s := range sArr {
		if i, e := strconv.ParseBool(s); e == nil {
			iArr[idx] = i
		}
	}

	return iArr
}
