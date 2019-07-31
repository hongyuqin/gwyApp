package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"regexp"
)

/**
判断推理	Judgment reasoning	JR
资料分析	date analyzing	DA
常识判断	Common sense judgment	CSJ
公共基础	Public foundation	PF

言语理解与表达	Speech understanding and expression	SUAE
数量关系	Quantity relationship	QR
*/
func testChineseCase(ch string) string {
	switch ch {
	case "判断推理":
		return "JR"
	case "资料分析":
		return "DA"
	case "常识判断":
		return "CSJ"
	case "公共基础":
		return "PF"
	case "言语理解与表达":
		return "SUAE"
	case "数量关系":
		return "QR"
	default:
		logrus.Panic("sth wrong :", ch)
		return ""
	}
}
func testRegex() {
	//正则表达式
	//[0-9]{1,3}.
	bool, err := regexp.MatchString("^([0-9]{1,3})．", "1．《国家“十二五”时期文化改革发展规划纲要》提出要加大政府对文化事")
	if err != nil {
		fmt.Println("error :", err)
	}
	fmt.Println("bool is :", bool)
}

//字符串截取

func main() {
	/*str := []rune("第一部分 常识判断")
	fmt.Println(string(str[5:]))*/
	testRegex()
}
