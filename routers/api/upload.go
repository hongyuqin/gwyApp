package api

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/unidoc/unioffice/document"
	"gwyApp/common"
	"gwyApp/models"
	"gwyApp/pkg/app"
	"gwyApp/pkg/e"
	"gwyApp/pkg/setting"
	"net/http"
	"regexp"
	"strings"
)

func UploadFile(c *gin.Context) {
	appG := app.Gin{C: c}
	// single file
	file, _ := c.FormFile("file")
	if file == nil {
		logrus.Error("file empty...")
		appG.Response(http.StatusOK, e.ERROR, nil)
		return
	}
	logrus.Info("file name is :", file.Filename)
	dst := setting.AppSetting.FileSavePath + file.Filename
	// Upload the file to specific dst.
	err := c.SaveUploadedFile(file, dst)
	//1.保存文件
	if err != nil {
		logrus.Error("save uploaded file fail :", err)
		appG.Response(http.StatusOK, e.ERROR, nil)
		return
	}
	//2.读取文件
	doc, err := document.Open(dst)
	if err != nil {
		logrus.Error("error opening document: ", err)
		appG.Response(http.StatusOK, e.ERROR, nil)
		return
	}
	//全部输出到一个buffer里面
	buf := bytes.NewBuffer([]byte{})

	eleOne := ""
	begin := true
	answer := false
	answerBegin := true
	for _, para := range doc.Paragraphs() {
		for _, run := range para.Runs() {
			//logrus.Info(run.Text())
			line := run.Text()
			bool, err := regexp.MatchString("^([0-9]{1,3})．", line)
			if err != nil {
				fmt.Println("error :", err)
				panic("not match")
			}
			if strings.HasPrefix(line, "答案与解析") {
				answer = true
				topic, _ := analysisBuf(2013, file.Filename, eleOne, buf.Bytes())
				buf.Write([]byte(run.Text()))
				buf.Reset()
				logrus.Info("topic is :", topic)
				continue
			}
			if answer {
				//开始解析答案 第二个开始解析答案
				if bool {
					if answerBegin {
						answerBegin = false
						buf.Write([]byte(run.Text()))
						continue
					}
					//可以解析这个buffer了
					topic, err := analysisAnswerBuf(file.Filename, buf.Bytes())
					if err != nil {
						logrus.Error("analysisAnswer fail :", err)
						appG.Response(http.StatusOK, e.ERROR, nil)
						return
					}
					logrus.Info("topic is :", topic)
					buf.Reset()
				}
			} else if bool {
				if begin {
					begin = false
					buf.Write([]byte(run.Text()))
					continue
				}
				//可以解析这个buffer了
				topic, err := analysisBuf(2013, file.Filename, eleOne, buf.Bytes())
				if err != nil {
					logrus.Error("analysis fail :", err)
					appG.Response(http.StatusOK, e.ERROR, nil)
					return
				}
				logrus.Info("topic is :", topic)
				buf.Reset()
			} else if strings.HasPrefix(line, "第") {
				eleOne = getCode(line)
				continue
			} else {
				logrus.Info("nothing to do .continue")
			}
			buf.Write([]byte(run.Text()))
		}
	}
	topic, err := analysisAnswerBuf(file.Filename, buf.Bytes())
	if err != nil {
		logrus.Error("analysisAnswer fail :", err)
		appG.Response(http.StatusOK, e.ERROR, nil)
		return
	}
	logrus.Info("topic is :", topic)

	appG.Response(http.StatusOK, e.SUCCESS, fmt.Sprintf("'%s' uploaded!", file.Filename))
}
func analysisBuf(year int, fileName string, elementOne string, byteBuf []byte) (topic *models.Topic, err error) {
	logrus.Info("analysisBuf :", string(byteBuf))
	bufStr := []rune(string(byteBuf))

	topicBegin := search(bufStr, "．") + 1
	topicEnd := search(bufStr, "A．")
	optionABegin := search(bufStr, "A．") + 2
	optionAEnd := search(bufStr, "B．")
	optionBBegin := search(bufStr, "B．") + 2
	optionBEnd := search(bufStr, "C．")
	optionCBegin := search(bufStr, "C．") + 2
	optionCEnd := search(bufStr, "D．")
	optionDBegin := search(bufStr, "D．") + 2
	photoBegin := search(bufStr, "http")
	photoEnd := search(bufStr, "png") + 3
	analysisEnd := len(bufStr)
	topicName := strings.TrimSpace(string(bufStr[topicBegin:topicEnd]))
	optionA := strings.TrimSpace(string(bufStr[optionABegin:optionAEnd]))
	optionB := strings.TrimSpace(string(bufStr[optionBBegin:optionBEnd]))
	optionC := strings.TrimSpace(string(bufStr[optionCBegin:optionCEnd]))
	optionD := strings.TrimSpace(string(bufStr[optionDBegin:analysisEnd]))

	topicNo := fmt.Sprintf("%s_%s", fileName, string(bufStr[0:topicBegin-1]))
	topic = &models.Topic{
		TopicNo:        topicNo,
		TopicName:      topicName,
		OptionA:        optionA,
		OptionB:        optionB,
		OptionC:        optionC,
		OptionD:        optionD,
		Year:           year,
		ExamType:       common.CIVIL_SERVICE_EXAM,
		ElementTypeOne: elementOne,
		Region:         common.CN,
	}
	if photoBegin > 0 {
		photo := strings.TrimSpace(string(bufStr[photoBegin:photoEnd]))
		topic.Photo = photo
	}
	err = models.AddTopic(topic)
	return
}
func analysisAnswerBuf(fileName string, byteBuf []byte) (topic *models.Topic, err error) {
	logrus.Info("analysisAnswerBuf :", string(byteBuf))
	bufStr := []rune(string(byteBuf))

	begin := search(bufStr, "【")
	answerBegin := search(bufStr, "】")
	answerEnd := search(bufStr, "。")
	analysisBegin := search(bufStr, "。") + 2
	analysisEnd := len(bufStr)

	answer := strings.TrimSpace(string(bufStr[answerBegin+1 : answerEnd]))
	analysis := strings.TrimSpace(string(bufStr[analysisBegin:analysisEnd]))

	topicNo := fmt.Sprintf("%s_%s", fileName, string(bufStr[0:begin-1]))
	topic = &models.Topic{
		TopicNo:       topicNo,
		Answer:        answer,
		TopicAnalysis: analysis,
	}

	err = models.UpdateTopicByTopicNo(topic)
	return
}
func getCode(ch string) string {
	begin := strings.Index(ch, " ")
	//str := []rune(ch)
	ch = strings.TrimSpace(ch[begin+1:])
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
		logrus.Error("sth wrong :", ch)
		return ""
	}
	return ""
}
func search(text []rune, what string) int {
	whatRunes := []rune(what)

	for i := range text {
		found := true
		for j := range whatRunes {
			if text[i+j] != whatRunes[j] {
				found = false
				break
			}
		}
		if found {
			return i
		}
	}
	return -1
}
