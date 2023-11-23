# alpaca
alpaca 是一个简单的游戏配置生成工具，能过读取 excel 文件，生成 json 和各种语言的配置类

## 使用方法
```bash
alpaca -dir=/path/to/excels -json_out=/path/to/gen/json -go_out=/path/to/gen/go
```
> 目前仅支持生成 json 文件和 go 配置类
> 其它语言配置后序补充
> // TODO: java c/c++ c# js/ts python erlang etc.

## 配置方式
![image](https://github.com/lgynico/alpaca/assets/2893568/0dc005af-3958-449f-8188-1a7a60362ac9)
首列为行类型配置，支持 5 种行类型配置（两端为双下划线）：
- __type__ 列的类型
- __name__ 列的名称
- __side__ 生成端：c 表示客户端, s 表示服务端（// TODO 待实现）
- __desc__ 列的描述
- __rule__ 可对列进行一些规则限制

## 支持的 __type__
- int int8 int16 int32 int64
- uint uint8 uint16 uint32 uint64
- float double
- bool
- string
- array:type array2:type
- map:ktype,vtype

## 支持的 __rule__
- key 主键，必须，唯一，目前仅支持一个主键
- unique 唯一
- require 必填
- range[min,max] 对数值类型有效，限定数值的范围，区间可开可闭，如 range(1,10] range[1,10)
- length[min,max] 对 string 和集合类有效，限定长度

## Todo List
- 支持多配置数据来源
  - Excel ✔️
  - etc......
- 支持多语言配置类生成
  - Golang ✔️
  - Java
  - C/C++
  - C#
  - Js/Ts
  - Python
  - Rust
  - etc......
- 支持多格式导出
  - json ✔️
  - xml
  - lua
  - yaml
  - etc......
