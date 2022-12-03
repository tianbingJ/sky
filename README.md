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

```
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

整数类型为了简单起见，全部都转为int64类型

### 浮点类型

```
var a = 1;
var b = 2.2;
var c = a + b;
```

### string
```
var s = "hello " + "world";
```

## 表达式

### 字面值

### 算法表达式

优先级从低到高：

|        操作        |      含义      | 结合性 |
|:----------------:|:------------:|:---:|
|        =         |      赋值      | 右结合 |
|      and or      |    逻辑 与或    | 右结合 |
| & &#124; ^ << >> | 与、或、异或,左移，右移 | 左结合 |
|      == !=       |      比较      | 左结合 |
|    < <= > >=     |      比较      | 左结合 |
|       + -        |      加       | 左结合 |
|      * / %       |      乘       | 左结合 |
|       - !        |    一元运算符     | 右结合 |


## 语句


### 变量声明语句

```
var x = 1;
var y;    //default nil
var a, b = x, x * y;
```

### 赋值语句
```
a, b = x, y;
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

initializer和incr可以不填。
initializer支持声明语句和复制语句，最多支持一个'=',多个变量可以使用多变量声明的形式。

for var i, j = 0; i < 10; i, j = i + 1, j + 1 {
    //do something
}
```

### while 循环

```
while (condition) {
    //do something
}
```

## 函数

```
func add(a, b) {
    return a + b;
}

函数可以有返回值,没有return默认返回nil。
```

## 闭包
```
    func makeTimes(factor) {
        func times(x ) {
            return  x * factor;
        }
        return times;
    }
    var f10 = makeTimes(10);
    print(f10(3)); //30
```