# TODO LIST

1. 解析用户定义的 config struct 结构，生成平面结构，多层使用 xxx.xxx.xxx 的形式
> 例
```
   type Config struct {
       Ftp struct{
           Host string
           Port int
       }
       
       Root string
   }
```
解析后的数据为(k,v):, 其中 v 为指针，指向struct中对应的变量？（需测试直接操作field进行）
> 测试操纵自定义struct 中的值
```
ftp.host = xxx
ftp.port = xxx
root = xxx
```

2. 根据优先级从不同来源填充 k,v 结构

3. 根据(k,v)结构，回填 struct