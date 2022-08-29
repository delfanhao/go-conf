package conf

import (
	"fmt"
	"reflect"
	"strconv"
)

// getterMapping 获取值属性函数的映射, 参数：
type getterMapping map[int]func(string, string) (interface{}, bool)

// configContext 配置过程中的上下文
type configContext struct {
	appName     string // 应用名称
	workDir     string // 应用所在路径
	entry       interface{}
	typeIndex   map[int]string //
	parserFunc  getterMapping
	fileMapping map[string]interface{} // 文件名 -> 对应类型实体 的映射，为nil 说明文件不存在
	// 反射相关
	defTypes  reflect.Type
	defValues reflect.Value
}

const (
	cmdLine = iota
	env
	cfgAppYml
	cfgAppJson
	cfgAppIni
	appYml
	appJson
	appIni
	cfgDefaultYml
	cfgDefaultJson
	cfgDefaultIni
	defaultYml
	defaultJson
	defaultIni
	defaultDefine
)

// Load 加载配置项目，按照指定的加载顺序进行处理
func Load(confStruct interface{}) {
	ctx := &configContext{
		entry:     confStruct,
		defTypes:  reflect.TypeOf(confStruct),
		defValues: reflect.ValueOf(confStruct),
		typeIndex: make(map[int]string),
	}
	if confStruct == nil || ctx.defTypes.Kind() != reflect.Ptr && ctx.defTypes.Kind() != reflect.Struct {
		println("Whoops! param is not a struct, check your code please.")
		return
	}

	ctx.workDir, ctx.appName = getAppInfo()
	// 解析配置结构
	ctx.initCallerMapping()
	// 预扫描文件
	ctx.scanAllFile()
	// 解析配置结构
	ctx.parseConfigStruct()
}

// setFileState 设置配置索引与配置文件名的关联
func (ctx *configContext) setFileState(index int, filename, fileType string) {
	ctx.typeIndex[index] = filename
	ctx.fileMapping[filename] = ctx.prepareFile(filename, fileType)
	if ctx.fileMapping[filename] == nil {
		trace("Prepare file [%s], but not found or parse error.", filename)
	} else {
		trace("Prepare file [%s], found and parse ok.", filename)
	}
}

// scanAllFile 扫描所有允许出现配置文件的位置
func (ctx *configContext) scanAllFile() {
	ctx.fileMapping = make(map[string]interface{})
	ctx.setFileState(cfgAppYml, "conf/"+ctx.appName+".yml", "yml")
	ctx.setFileState(cfgAppJson, "conf/"+ctx.appName+".json", "json")
	ctx.setFileState(cfgAppIni, "conf/"+ctx.appName+".ini", "ini")
	ctx.setFileState(appYml, ctx.appName+".yml", "yml")
	ctx.setFileState(appJson, ctx.appName+".json", "json")
	ctx.setFileState(appIni, ctx.appName+".ini", "ini")
	ctx.setFileState(cfgDefaultYml, "conf/config.yml", "yml")
	ctx.setFileState(cfgDefaultJson, "conf/config.json", "json")
	ctx.setFileState(cfgDefaultIni, "conf/config.ini", "ini")
	ctx.setFileState(defaultYml, "config.yml", "yml")
	ctx.setFileState(defaultJson, "config.json", "json")
	ctx.setFileState(defaultIni, "config.ini", "ini")
}

// initCallerMapping 设置每种类型的对应解析函数
func (ctx *configContext) initCallerMapping() {
	// 按照优先级进行解析函数的映射
	ctx.parserFunc = make(getterMapping)
	ctx.parserFunc[cmdLine] = ctx.getCmdVal
	ctx.parserFunc[env] = ctx.getEnvVal
	ctx.parserFunc[cfgDefaultYml] = ctx.getYmlValFromFile
	ctx.parserFunc[cfgDefaultJson] = ctx.getJsonValFromFile
	ctx.parserFunc[cfgDefaultIni] = ctx.getIniValFromFile
	ctx.parserFunc[cfgAppYml] = ctx.getYmlValFromFile
	ctx.parserFunc[cfgAppJson] = ctx.getJsonValFromFile
	ctx.parserFunc[cfgAppIni] = ctx.getIniValFromFile
	ctx.parserFunc[defaultYml] = ctx.getYmlValFromFile
	ctx.parserFunc[defaultJson] = ctx.getJsonValFromFile
	ctx.parserFunc[defaultIni] = ctx.getIniValFromFile
	ctx.parserFunc[appYml] = ctx.getYmlValFromFile
	ctx.parserFunc[appJson] = ctx.getJsonValFromFile
	ctx.parserFunc[appIni] = ctx.getIniValFromFile

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

// parseConfigStruct 解析配置结构的映射， 执行完成后， 根据用户定义的结构，调用不同的解析函数
func (ctx *configContext) parseConfigStruct() {
	numEle := ctx.defTypes.Elem().NumField()
	for idx := 0; idx < numEle; idx++ {
		prefix := ""
		field := ctx.defTypes.Elem().Field(idx)
		value := ctx.defValues.Elem().Field(idx)
		switch field.Type.Kind() {
		case reflect.Struct:
			ctx.parseStruct(prefix, &field, &value)
		default:
			ctx.setValue(prefix, &field, &value)
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
		switch fld.Type.Kind() {
		case reflect.Struct:
			ctx.parseStruct(prefix, &fld, &val)
		default:
			ctx.setValue(prefix, &fld, &val)
		}
	}
}

// setValueByType 根据不同类型，使用不同方法设置对应的值
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
	case reflect.Slice:
		ctx.setSliceValue(val, field.Type, value)
	}
}

func (ctx *configContext) setSliceValue(data interface{}, fieldType reflect.Type, value *reflect.Value) {
	typeName := reflect.SliceOf(fieldType).Elem().String()
	switch typeName {
	case "[]string":
		arr := make([]string, 0)
		vArr := data.([]string)
		for i := range vArr {
			arr = append(arr, vArr[i])
		}
		value.Set(reflect.ValueOf(arr))
	case "[]int":
		arr := arrStr2int(data.([]string))
		value.Set(reflect.ValueOf(arr))
	case "[]int32":
		arr := arrStr2int32(data.([]string))
		value.Set(reflect.ValueOf(arr))
	case "[]int64":
		arr := arrStr2int64(data.([]string))
		value.Set(reflect.ValueOf(arr))
	case "[]float32", "[]float64":
		arr := arrStr2Float(data.([]string))
		value.Set(reflect.ValueOf(arr))
	case "[]bool":
		arr := make([]bool, 0)
		vArr := data.([]bool)
		for i := range vArr {
			arr = append(arr, vArr[i])
		}
		value.Set(reflect.ValueOf(arr))
	default:
		//slice = val.([]interface{})
	}
}

// setValue 根据类型设置配置文件中对应的值
func (ctx *configContext) setValue(prefix string, field *reflect.StructField, value *reflect.Value) {
	key := fmt.Sprintf("%s%s", prefix, field.Name)
	for parseIdx := cmdLine; parseIdx < defaultDefine; parseIdx++ {
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
