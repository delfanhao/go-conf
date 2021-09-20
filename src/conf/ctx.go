package conf

import (
	"fmt"
	"reflect"
	"strconv"
)

// funcMapping 字段类型处理函数映射, 参数为： key/字段反射/值反射
type funcMapping map[reflect.Kind]func(string, *reflect.StructField, *reflect.Value)

// getterMapping 获取值属性函数的映射, 参数：
type getterMapping map[int]func(string, string) (interface{}, bool)

// configContext 配置过程中的上下文
type configContext struct {
	appName     string // 应用名称
	workDir     string // 应用所在路径
	entry       interface{}
	typeIndex   map[int]string //
	caller      funcMapping
	parserFunc  getterMapping
	fileMapping map[string]interface{} // 文件名 -> 对应类型实体 的映射，为nil 说明文件不存在
	// 反射相关
	defTypes  reflect.Type
	defValues reflect.Value
}

const (
	CmdLine = iota
	Env
	ConfAppYml
	ConfAppJson
	ConfAppIni
	AppYml
	AppJson
	AppIni
	ConfDefaultYml
	ConfDefaultJson
	ConfDefaultIni
	DefaultYml
	DefaultJson
	DefaultIni
	DefaultDefine
)

// Load 加载配置项目，按照指定的加载顺序进行处理
func Load(confStruct interface{}) {
	ctx := &configContext{
		entry:     confStruct,
		defTypes:  reflect.TypeOf(confStruct),
		defValues: reflect.ValueOf(confStruct),
		typeIndex: make(map[int]string),
	}
	if ctx.defTypes.Kind() != reflect.Ptr && ctx.defTypes.Kind() != reflect.Struct {
		panic("Not struct type, Check your source code please.")
	}

	ctx.workDir, ctx.appName = getAppInfo()
	// 解析配置结构
	ctx.initCallerMapping()
	// 预扫描文件
	ctx.scanAllFile()
	// 解析配置结构
	ctx.parseConfigStruct()
}

// setFileState
func (ctx *configContext) setFileState(index int, filename, fileType string) {
	ctx.typeIndex[index] = filename
	ctx.fileMapping[filename] = ctx.prepareFile(filename, fileType)
	if ctx.fileMapping[filename] == nil {
		trace("Prepare file [%s], but not found or parse error.", filename)
	} else {
		trace("Prepare file [%s], found and parse ok.", filename)
	}
}

// scanAllFile
func (ctx *configContext) scanAllFile() {
	ctx.fileMapping = make(map[string]interface{})
	ctx.setFileState(ConfAppYml, "conf/"+ctx.appName+".yml", "yml")
	ctx.setFileState(ConfAppJson, "conf/"+ctx.appName+"config.json", "json")
	ctx.setFileState(ConfAppIni, "conf/"+ctx.appName+"config.ini", "ini")
	ctx.setFileState(AppYml, ctx.appName+".yml", "yml")
	ctx.setFileState(AppJson, ctx.appName+".json", "json")
	ctx.setFileState(AppIni, ctx.appName+".ini", "ini")
	ctx.setFileState(ConfDefaultYml, "conf/config.yml", "yml")
	ctx.setFileState(ConfDefaultJson, "conf/config.json", "json")
	ctx.setFileState(ConfDefaultIni, "conf/config.ini", "ini")
	ctx.setFileState(DefaultYml, "config.yml", "yml")
	ctx.setFileState(DefaultJson, "config.json", "json")
	ctx.setFileState(DefaultIni, "config.ini", "ini")
}

// initCallerMapping 设置每种类型的对应解析函数
func (ctx *configContext) initCallerMapping() {
	ctx.caller = make(funcMapping)
	ctx.caller[reflect.Struct] = ctx.parseStruct
	ctx.caller[reflect.String] = ctx.setValue
	ctx.caller[reflect.Int] = ctx.setValue
	ctx.caller[reflect.Float32] = ctx.setValue
	ctx.caller[reflect.Float64] = ctx.setValue
	ctx.caller[reflect.Array] = ctx.parseArray
	ctx.caller[reflect.Bool] = ctx.setValue

	// 按照优先级进行解析函数的映射
	ctx.parserFunc = make(getterMapping)
	ctx.parserFunc[CmdLine] = ctx.getCmdVal
	ctx.parserFunc[Env] = ctx.getEnvVal
	ctx.parserFunc[ConfDefaultYml] = ctx.getYmlValFromFile
	ctx.parserFunc[ConfDefaultJson] = ctx.getJsonValFromFile
	ctx.parserFunc[ConfDefaultIni] = ctx.getIniValFromFile
	ctx.parserFunc[ConfAppYml] = ctx.getYmlValFromFile
	ctx.parserFunc[ConfAppJson] = ctx.getJsonValFromFile
	ctx.parserFunc[ConfAppIni] = ctx.getIniValFromFile
	ctx.parserFunc[DefaultYml] = ctx.getYmlValFromFile
	ctx.parserFunc[DefaultJson] = ctx.getJsonValFromFile
	ctx.parserFunc[DefaultIni] = ctx.getIniValFromFile
	ctx.parserFunc[AppYml] = ctx.getYmlValFromFile
	ctx.parserFunc[AppJson] = ctx.getJsonValFromFile
	ctx.parserFunc[AppIni] = ctx.getIniValFromFile
	ctx.parserFunc[DefaultDefine] = ctx.getTagDefault

	ctx.fileMapping = make(map[string]interface{})
}

// prepareFile 根据文件类型，获取对应类型的上下文以及预加载的数据
func (ctx *configContext) prepareFile(filename, fileType string) interface{} {

	switch fileType {
	case "yml":
		ymlCtx := &ymlContext{cfgCtx: ctx}
		ymlCtx.scan(filename)
		if ymlCtx.data == nil {
			ctx.fileMapping[filename] = nil
			return nil
		}
		return ymlCtx
	case "json":
		jsonCtx := &jsonContext{cfgCtx: ctx}
		jsonCtx.scan(filename)
		if jsonCtx.data == nil {
			ctx.fileMapping[filename] = nil
			return nil
		}
		return jsonCtx
	case "ini":
		iniCtx := &iniContext{cfgCtx: ctx}
		iniCtx.scan(filename)
		if iniCtx.data == nil {
			ctx.fileMapping[filename] = nil
			return nil
		}
		return iniCtx
	}

	return nil
}

// getTagDefault
func (ctx *configContext) getTagDefault(_, key string) (interface{}, bool) {
	// todo 从Tag中获取default
	return "", false
}

// parseConfigStruct 解析配置结构的映射， 执行完成后， 根据用户定义的结构，调用不同的解析函数
func (ctx *configContext) parseConfigStruct() {
	numEle := ctx.defTypes.Elem().NumField()
	for idx := 0; idx < numEle; idx++ {
		prefix := ""
		field := ctx.defTypes.Elem().Field(idx)
		value := ctx.defValues.Elem().Field(idx)
		if caller, ok := ctx.caller[field.Type.Kind()]; ok {
			caller(prefix, &field, &value)
		}
	}
}

// parseStruct 解析Struct
func (ctx *configContext) parseStruct(prefix string, field *reflect.StructField, value *reflect.Value) {
	prefix = prefix + field.Name + "."
	numEle := field.Type.NumField()
	for idx := 0; idx < numEle; idx++ {
		fld := field.Type.Field(idx)
		val := value.Field(idx)
		if caller, ok := ctx.caller[fld.Type.Kind()]; ok {
			caller(prefix, &fld, &val)
		}
	}
}

func (ctx *configContext) setValueByType(field *reflect.StructField, value *reflect.Value, val interface{}) {
	switch field.Type.Kind() {
	case reflect.String:
		value.SetString(val.(string))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if v := convToInt64(val); v != nil {
			value.SetInt(*convToInt64(val))
		}
	case reflect.Bool:
		i, err := strconv.ParseBool(val.(string))
		if err == nil {
			value.SetBool(i)
		}
	case reflect.Float32, reflect.Float64:
		i, err := strconv.ParseFloat(val.(string), 64)
		if err == nil {
			value.SetFloat(i)
		}
	}

}

// setValue 根据类型设置配置文件中对应的值
func (ctx *configContext) setValue(prefix string, field *reflect.StructField, value *reflect.Value) {
	key := fmt.Sprintf("%s%s", prefix, field.Name)
	for parseIdx := CmdLine; parseIdx < DefaultDefine; parseIdx++ {
		filename := ctx.typeIndex[parseIdx]
		if val, ok := ctx.parserFunc[parseIdx](filename, key); ok {
			ctx.setValueByType(field, value, val)
			return
		}
	}
	// 如果外部配置位置均无法找到可以设置的值，并且程序中也未设置值，则需要获取定义的缺省值
	if value.IsZero() {
		if v, ok := field.Tag.Lookup("default"); ok {
			ctx.setValueByType(field, value, v)
		}
	}

}

func (ctx *configContext) parseArray(prefix string, field *reflect.StructField, value *reflect.Value) {
	// todo 处理数组
}
