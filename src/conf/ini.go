package conf

import (
	"fmt"
	"io/ioutil"
	"strings"
)

type iniContext struct {
	cfgCtx *configContext
	data   map[string]interface{}
}

func (ctx *iniContext) scan(filename string) {
	fullName := fmt.Sprintf("%s/%s", ctx.cfgCtx.workDir, filename)

	if ctx.data == nil {
		ctx.data = make(map[string]interface{})
	}
	iniFile, err := ioutil.ReadFile(fullName)
	if err != nil {
		ctx.data[filename] = nil
		return
	}

	str := string(iniFile)
	lines := strings.Split(str, "\n")

	valueMap := make(map[string]interface{})
	for idx := range lines {
		line := strings.Trim(lines[idx], " ")
		line = strings.Trim(line, "\r")
		if len(line) > 0 {
			pos := strings.Index(line, "=")
			k := strings.Trim(line[:pos], " ")
			v := strings.Trim(line[pos+1:], " ")
			valueMap[k] = v
		}
	}

	ctx.data = valueMap
}

func (ctx *iniContext) getVal(key string) (interface{}, bool) {
	key = strings.ToLower(key)
	v, ok := ctx.data[key]
	return v, ok
}

func (ctx *configContext) _getIniVal(filename, key string) (interface{}, bool) {
	iniCtx := ctx.fileMapping[filename]
	if iniCtx == nil {
		return nil, false
	}

	return iniCtx.(*iniContext).getVal(key)
}

func (ctx *configContext) getConfDefaultIniVal(filename, key string) (interface{}, bool) {
	return ctx._getIniVal(filename, key)
}

func (ctx *configContext) getConfIniVal(filename, key string) (interface{}, bool) {
	return ctx._getIniVal(filename, key)
}

func (ctx *configContext) getDefaultIniVal(filename, key string) (interface{}, bool) {
	return ctx._getIniVal(filename, key)
}

func (ctx *configContext) getIniVal(filename, key string) (interface{}, bool) {
	return ctx._getIniVal(filename, key)
}
