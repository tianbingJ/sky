# Sky语言
Sky语言是一种脚本语言，支持的特性；

## 注释
注释包含两种，单行注释和多行注释

### 单行注释
单行注释以'//'打头，到行位结束
```
//this is a comment
```

### 多行注释
多行注释以'/*'打头，到最近的'*/'结束。
```azure
/*
this is multi line comment
another line
*/
```

## 数据类型

### nil
nil表示变量没有绑定值，变量未被初始化时是nil，也可以赋值为nil。

### 布尔类型
字面值包括：true 和 false.

### 整数类型
整数类型为了简单起见，全部都转为long型
```
var a = 1;
var b = 2;
var c = a + b;
```

## 表达式

|     操作      | 表头  |
|:-----------:|:-----:|
|     单元格     | 单元格 |
|     单元格     | 单元格 |


## 语句

### 变量声明
```
var x = 1
var y    //default nil
```


### if条件语句
个人比较喜欢Go的表达式不用圆括号括起来，借鉴过来。
```
if expr {
    //do something
} 

if expr1 {
    //do1
} elif expr 2 {
    //do2
} else {
    //do3
}
```

### 关于控制条件condition
if、while、for语句中会判断条件的真假来调整执行路径，python把0、''、False、None都判断成False；
Sky只把nil和false判断为false，其他全部为true。

### for 循环
```

for initializer; condition; incr {
}
```

initializer 可以是 变量声明语句，可以是任意表达式，可以是赋值语句;
initializer和incr可以不填。
目前initializer和incr都只支持单条语句。
```
for var i = 0; i < 10; i ++ {

}
```

### while 循环
```
while (condition) {

}
```

## 函数



## 闭包
