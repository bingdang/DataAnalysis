package main

import (
	"bufio"
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"io"
	"os"
	"strings"
	"sync"
)

/*
//数据清洗：
现有一份泄漏的酒店开房记录，需要进行数据清洗和分析
1.数据清洗
获取有效身份证数据

2.数据分类
按身份归类

./kaifang-gbk.txt 为GBK字符集

解决中文乱码
go get golang.org/x/text
simplifiedchinese.GBK.NewEncoder().Bytes()   //utf-8 转 gbk
simplifiedchinese.GBK.NewDecoder().Bytes()  //gbk 转 utf-8
*/
var ps []string

type Province struct {
	Id    string
	Name  string
	Queue chan string
	File  *os.File
}

var wr = sync.WaitGroup{}

/* UTF8转GBK
iconv -f UTF-8 -t GBK ./kaifang.txt > ./kaifang-gbk.txt
*/

func HandleError(err error, where string) {
	if err != nil {
		fmt.Println(where, err)
	}
}

// 字符集转换
// 需要处理的数据的字节流
func CharacterSet(LineB []byte) (dstStr string) {
	// 创建GBK解码处理器
	encoder := simplifiedchinese.GBK.NewDecoder()

	// 转换为UTF8的字节流
	UTF8Line, _ := encoder.Bytes(LineB)

	// 转为UTF8的字符串
	dstStr = string(UTF8Line)
	return
}

// 数据初步清洗
func initialData(LineStr string) {
	//创建好数据文件
	goodfile, _ := os.OpenFile("./kaifang-utf8_good.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	defer goodfile.Close()

	//创建坏数据文件
	badfile, _ := os.OpenFile("./kaifang-utf8_bad.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	defer badfile.Close()

	//根据行数据取身份证
	//按逗号切割行数据
	dataSplit := strings.Split(LineStr, ",")
	if len(dataSplit) >= 2 && len(dataSplit[1]) == 18 {
		_, err := goodfile.WriteString(LineStr + "\n")
		HandleError(err, "WriteGoodLine")
		fmt.Println("GoodData", LineStr)
	} else {
		_, err := badfile.WriteString(LineStr + "\n")
		HandleError(err, "WriteBadLine")
		fmt.Println("BadData", LineStr)
	}
}

func ReadFile(filename string) {
	//打开源文件
	open, err := os.Open(filename)
	HandleError(err, "OpenFile")
	defer open.Close()

	//缓冲读取对象
	reader := bufio.NewReader(open)
	for {
		LineB, _, err := reader.ReadLine()
		if err == io.EOF {
			fmt.Println("读取完成")
			break
		}
		HandleError(err, "ReadLine")
		LineStr := CharacterSet(LineB)
		initialData(LineStr)

	}
}

// 省份划分
func ReadGoodFile(filename string) {
	//创建省份字典
	ps = []string{"北京市11", "天津市12", "河北省13",
		"山西省14", "内蒙古自治区15", "辽宁省21", "吉林省22",
		"黑龙江省23", "上海市31", "江苏省32", "浙江省33", "安徽省34",
		"福建省35", "江西省36", "山东省37", "河南省41", "湖北省42",
		"湖南省43", "广东省44", "广西壮族自治区45", "海南省46",
		"重庆市50", "四川省51", "贵州省52", "云南省53", "西藏自治区54",
		"陕西省61", "甘肃省62", "青海省63", "宁夏回族自治区64", "新疆维吾尔自治区65",
		"香港特别行政区81", "澳门特别行政区82", "台湾省83"}

	var psMap = make(map[string]*Province)

	// 封装省份对象，创建省份文件、省份管道
	CapitalConstruction(ps, psMap)

	// 读取管道数据写入对应的省份
	ReadChenToWriteFile(psMap)

	// 读取数据写入对应省份的管道
	WriteGoodDataToChen(filename, psMap)

}

// 封装省份对象，创建省份文件、省份管道
func CapitalConstruction(provinces []string, psMap map[string]*Province) {
	for _, p := range provinces {
		province := Province{Id: p[len(p)-2:], Name: p[:len(p)-2]}
		province.Queue = make(chan string, 256)
		province.File, _ = os.OpenFile("./省份/"+province.Name+".txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		fmt.Println(province.Name, "管道已创建", province.File, "文件已创建")

		//构建map
		psMap[province.Id] = &province
	}
}

// 读取数据写入对应省份的管道
func WriteGoodDataToChen(filename string, psMap map[string]*Province) {
	//打开源文件
	open, err := os.Open(filename)
	HandleError(err, "OpenFile")
	defer open.Close()

	//缓冲读取对象
	reader := bufio.NewReader(open)
	for {
		LineB, _, err := reader.ReadLine()
		if err == io.EOF {
			fmt.Println("读取完成")
			break
		}
		HandleError(err, "ReadLine")
		LineStr := string(LineB)
		Linesplit := strings.Split(LineStr, ",")
		id := Linesplit[1][:2]
		if province, ok := psMap[id]; ok {
			province.Queue <- LineStr
			fmt.Println(LineStr, "已写入", province.Name+"管道")
		} else {
			fmt.Println("未知的省份")
		}
	}

	//关闭管道
	for _, ps := range psMap {
		close(ps.Queue)
	}
}

// 读取管道数据写入对应的省份文件
func ReadChenToWriteFile(psMap map[string]*Province) {
	for _, province := range psMap {
		wr.Add(1)
		go writeProvinceFile(province)
	}
}

// 写入各省文件
func writeProvinceFile(province *Province) {
	for LineStr := range province.Queue {
		province.File.WriteString(LineStr + "\n")
		defer province.File.Close()
		fmt.Println(LineStr, "写入", province.Name)
	}
	wr.Done()
}

func main() {
	//筛选有效数据
	ReadFile("./kaifang-gbk.txt")

	//有效数据再次归类
	ReadGoodFile("./kaifang-utf8_good.txt")
	wr.Wait()
}
