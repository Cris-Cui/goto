# URLShortener -- 短URL web项目

## 功能
1. **添加**: 给定一个长URL, 返回一个短的版本
2. **重定向**: 当请求短URL的时候, 将用户重定向到原始的长的URL

## 实现
### Version1 -- 实现添加和重定向两项业务功能

### Version2 -- 短链数据持久化
### Version3 -- goroutine 和 channel 重构
### Version4 -- 持久化层改为json格式
### Version5 -- 支持`rpc`协议的分布式版本

# Tutorial -- 项目中用到的语法
## 基本结构和基本数据类型
1. 可见性原则: 
   - 如果常量、变量、类型、函数名、结构字段等等以一个大写字母开头, 那么该对象就可以被外部包的代码所使用 -- public
   - 标识符如果以小写字母开头, 对包外不可见 -- private
2. Go程序一般结构体
   - 在完成包的 import 之后，开始对常量、变量和类型的定义或声明。
   - 如果存在 init 函数的话，则对该函数进行定义
   - 如果当前包是 main 包，则定义 main 函数
   - 然后定义其余的函数，首先是类型的方法，接着是按照 main 函数中先后调用的顺序来定义相关函数，如果有很多函数，则可以按照字母顺序来进行排序
3. 常量 -- `const`
   [**itoa**](https://github.com/Unknwon/go-fundamental-programming/blob/master/lectures/lecture4.md)每遇到一次 const 关键字，iota 就重置为 0 
4. 变量 -- `var`