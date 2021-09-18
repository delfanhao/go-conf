package conf

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"reflect"
)

type ymlContext struct {
	cfgCtx *configContext
	data   map[string]interface{}
}

func (ctx *ymlContext) scan(filename string) {
	fullName := fmt.Sprintf("%s/%s", ctx.cfgCtx.workDir, filename)
	if ctx.data == nil {
		ctx.data = make(map[string]interface{})
	}
	ymlFile, err := ioutil.ReadFile(fullName)
	if err != nil {
		ctx.data[filename] = nil
		return
	}

	data := make(map[interface{}]interface{})
	err = yaml.Unmarshal(ymlFile, data)
	if err != nil {
		println(err.Error())
	}

	valueMap := make(map[string]interface{})
	ctx.parse(data, valueMap, "")

	ctx.data = valueMap
}

func (ctx *ymlContext) parse(data map[interface{}]interface{}, valueMap map[string]interface{}, prefix string) {
	for k, v := range data {
		defType := reflect.TypeOf(v).Kind()
		if defType == reflect.Map {
			ctx.parse(v.(map[interface{}]interface{}), valueMap, k.(string))
		} else {
			dataKey := fmt.Sprintf("%s.%s", prefix, k)
			valueMap[dataKey] = v
		}
	}
}

// getVal 获取指定key的值
func (ctx *ymlContext) getVal(key string) (interface{}, bool) {

	v, ok := ctx.data[key]
	return v, ok
}

func (ctx *configContext) _getYmlVal(filename, key string) (interface{}, bool) {
	ymlCtx := ctx.fileMapping[filename]
	if ymlCtx == nil {
		return nil, false
	}

	return ymlCtx.(*ymlContext).getVal(key)
}

//
func (ctx *configContext) getConfDefaultYmlVal(filename, key string) (interface{}, bool) {
	return ctx._getYmlVal(filename, key)
}

func (ctx *configContext) getConfYmlVal(filename, key string) (interface{}, bool) {
	return ctx._getYmlVal(filename, key)
}

func (ctx *configContext) getDefaultYmlVal(filename, key string) (interface{}, bool) {
	return ctx._getYmlVal(filename, key)
}

func (ctx *configContext) getYmlVal(filename, key string) (interface{}, bool) {
	return ctx._getYmlVal(filename, key)
}
