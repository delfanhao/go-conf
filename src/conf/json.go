package conf

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"
	"strings"
)

type jsonContext struct {
	cfgCtx *configContext
	data   map[string]interface{}
}

// scan 扫描json文件，构造平面kv
func (ctx *jsonContext) scan(filename string) {
	fullName := fmt.Sprintf("%s/%s", ctx.cfgCtx.workDir, filename)

	if ctx.data == nil {
		ctx.data = make(map[string]interface{})
	}
	jsonFile, err := ioutil.ReadFile(fullName)
	if err != nil {
		ctx.data[filename] = nil
		return
	}

	var data interface{}

	err = json.Unmarshal(jsonFile, &data)
	if err != nil {
		println(err.Error())
	}

	valueMap := make(map[string]interface{})
	ctx.parse(data.(map[string]interface{}), valueMap, "")

	ctx.data = valueMap
}

// parse 把json的树状结构转换成平面结构
func (ctx *jsonContext) parse(data map[string]interface{}, valueMap map[string]interface{}, prefix string) {
	for k, v := range data {
		defType := reflect.TypeOf(v).Kind()
		dataKey := fmt.Sprintf("%s.%s", prefix, k)
		if len(prefix) == 0 {
			dataKey = dataKey[1:]
		}
		if defType == reflect.Map {
			ctx.parse(v.(map[string]interface{}), valueMap, dataKey)
		} else {
			valueMap[dataKey] = v
		}
	}
}

// getVal 获取指定key的值
func (ctx *jsonContext) getVal(key string) (interface{}, bool) {
	key = strings.ToLower(key)
	v, ok := ctx.data[key]
	return v, ok
}

// _getJsonVal
func (ctx *configContext) _getJsonVal(filename, key string) (interface{}, bool) {
	jsonCtx := ctx.fileMapping[filename]
	if jsonCtx == nil {
		return nil, false
	}

	return jsonCtx.(*jsonContext).getVal(key)
}

// getConfDefaultJsonVal
func (ctx *configContext) getConfDefaultJsonVal(filename, key string) (interface{}, bool) {
	return ctx._getJsonVal(filename, key)
}

// getConfJsonVal
func (ctx *configContext) getConfJsonVal(filename, key string) (interface{}, bool) {
	return ctx._getJsonVal(filename, key)
}

// getDefaultJsonVal
func (ctx *configContext) getDefaultJsonVal(filename, key string) (interface{}, bool) {
	return ctx._getJsonVal(filename, key)
}

// getJsonVal
func (ctx *configContext) getJsonVal(filename, key string) (interface{}, bool) {
	return ctx._getJsonVal(filename, key)
}
