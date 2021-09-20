# GO-CONF

配置文件解析模块，类似Spring Boot，带优先级，从指定或缺省的配置类别中进行配置项的解析。


## 使用``go-conf``
``go mod github.com/delfanhao/go-conf``

代码如下：
```
package main

import (
	"go-conf/src/conf"
)

type Config {
    Flag string
}

func main(){
    cfg := Config{}
	conf.Load(&cfg)
	printlnf(cfg.Flag)
}
```
如果您想了解解析过程，可使用conf包中的一个全局布尔变量``TRACE``来跟踪解析过程，这个变量的默认值是关闭的。

## 配置项规则 
1. 配置项定义
   配置项使用全路径方式配置，以全小写为主，根据位置不同，略有不同：
- 命令行/环境变量中，使用全大写规则，主从结构关系间使用下划线间隔
- ini 文件使用全小写，主从结构关系用点符号``.``间隔
- yml,json文件使用各自的标准结构，配置项均为小写
如下配置结构例子如下：
```
type SampleConfig struct {
    MainItem struct {
        SubItem string
    }
}
```
SubItem为其中一个配置项，全路径(全小写)表示为： ``mainitem.subitem``
不同配置位置如下：
命令行配置为: ``-MAINITEM_SUBITEM=val``
环境变量中配置为 :  ``MAINITEM_SUBITEM=val``
ini文件中为 ``mainitem.subimte=val``
yml文件中为
```
mainitem:
  subitem: val
```
json文件中为:
```
{
  "mainitem": {
    "subitem":"val"
  }
}
```

2. 配置项查找顺序

类似于Spring Boot的配置项加载顺序， go-conf的配置项加载顺序如下(优先级从高到低排列):  
- 命令行  
- 环境变量  
- conf/<可执行文件名>.yml  
- conf/<可执行文件名>.json  
- conf/<可执行文件名>.ini  
- ./<可执行文件名>.yml  
- ./<可执行文件名>.json  
- ./<可执行文件名>.ini  
- conf/config.yml  
- conf/config.json  
- conf/config.ini  
- ./config.yml  
- ./config.json  
- ./config.yml  
- 配置项初始时指定值  
- 配置项tag中定义的缺省值  

指定配置项根据优先级从高到底进行对应值的查找，找到对应值后，就设置为当前找到的值。

> ⚠️注意！在golang中，处于代码安全考虑，任何变量都会有一个默认的空值，如int为0，string为"", interface{}为nil
> 等，因此以上判断顺序的最后两个类型，如果程序代码中出现显式赋值，但却赋为类型对应的零值，如果定义时配置了缺省值，则依然会被设置为缺省值。

举个例子：
```
type Cfg struct {
    item string `default:"default value"`
}

func main() {
    cfg := Cfg{
        item: ""
    }
    conf.Load(&cfg)
    println(cfg.item)
}
```
此时根据配置加载顺序，cfg.item应该为 "", 但由于 ""为字符串类型的零值， 
所以 cfg.item 被设置成了 tag 中的 "default value"

