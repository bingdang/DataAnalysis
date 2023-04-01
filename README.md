# DataAnalysis
数据分析demo

一份泄漏的开房记录，需要进行清洗归类。

1.数据清洗
按照有效身份证进行清洗

2.数据整理
按照省份进行归档
```go
felix@MacBook-Pro 0319 % go run main.go
....

felix@MacBook-Pro 0319 % ls 省份
上海市.txt                      安徽省.txt                      河南省.txt                      贵州省.txt                      内蒙古自治区.txt
云南省.txt                      山东省.txt                      浙江省.txt                      辽宁省.txt                      宁夏回族自治区.txt
北京市.txt                      山西省.txt                      海南省.txt                      重庆市.txt                      广西壮族自治区.txt
台湾省.txt                      广东省.txt                      湖北省.txt                      陕西省.txt                      澳门特别行政区.txt
吉林省.txt                      江苏省.txt                      湖南省.txt                      青海省.txt                      香港特别行政区.txt
四川省.txt                      江西省.txt                      甘肃省.txt                      黑龙江省.txt                    新疆维吾尔自治区.txt
天津市.txt                      河北省.txt                      福建省.txt                      西藏自治区.txt
felix@MacBook-Pro 0319 % ls
go.mod                  go.sum                  kaifang-gbk.txt         kaifang-utf8_bad.txt    kaifang-utf8_good.txt   main.go                 省份
felix@MacBook-Pro 0319 % 

```
### 思路：
- 为每个省封装自己的对象，对象中写明名称，id（身份证前两位），对应的文件，管道
- 将省份放入map中通过省份id取对应省份对象 并给管道写入数据。main.go第166行
- 每个省份开一个协程处理数据从管道读出，并写入省份对应文件。main.go第182行