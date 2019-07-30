package api

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/unidoc/unioffice/document"
	"net/http"
	"gwyApp/models"
	"gwyApp/pkg/app"
	"gwyApp/pkg/e"
	"gwyApp/pkg/setting"
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

	i := 2
	str := fmt.Sprintf("%d.", i)
	for _, para := range doc.Paragraphs() {
		for _, run := range para.Runs() {
			//logrus.Info(run.Text())
			line := run.Text()
			if strings.HasPrefix(line, str) {
				//可以解析这个buffer了
				topic, err := analysisBuf(buf.Bytes())
				if err != nil {
					logrus.Error("analysis fail :", err)
					appG.Response(http.StatusOK, e.ERROR, nil)
					return
				}
				logrus.Info("topic is :", topic)
				buf.Reset()
				i++
				str = fmt.Sprintf("%d.", i)
			}
			buf.Write([]byte(run.Text()))
		}
	}

	appG.Response(http.StatusOK, e.SUCCESS, fmt.Sprintf("'%s' uploaded!", file.Filename))
}
func analysisBuf(byteBuf []byte) (topic *models.Topic, err error) {
	logrus.Info("analysisBuf :", string(byteBuf))
	bufStr := string(byteBuf)

	topicBegin := 2
	topicEnd := strings.Index(bufStr, "A.")
	optionABegin := strings.Index(bufStr, "A.") + 3
	optionAEnd := strings.Index(bufStr, "B.")
	optionBBegin := strings.Index(bufStr, "B.") + 3
	optionBEnd := strings.Index(bufStr, "C.")
	optionCBegin := strings.Index(bufStr, "C.") + 3
	optionCEnd := strings.Index(bufStr, "D.")
	optionDBegin := strings.Index(bufStr, "D.") + 3
	optionDEnd := strings.Index(bufStr, "【答案")
	answerBegin := strings.Index(bufStr, "【答案") + 12
	answerEnd := strings.Index(bufStr, "【解析】")
	analysisBegin := strings.Index(bufStr, "【解析】") + 12
	analysisEnd := len(bufStr)
	topicName := bufStr[topicBegin:topicEnd]
	optionA := bufStr[optionABegin:optionAEnd]
	optionB := bufStr[optionBBegin:optionBEnd]
	optionC := bufStr[optionCBegin:optionCEnd]
	optionD := bufStr[optionDBegin:optionDEnd]
	answer := strings.TrimSpace(bufStr[answerBegin:answerEnd])
	analysis := bufStr[analysisBegin:analysisEnd]
	topic = &models.Topic{
		TopicName:     topicName,
		OptionA:       optionA,
		OptionB:       optionB,
		OptionC:       optionC,
		OptionD:       optionD,
		Answer:        answer,
		TopicAnalysis: analysis,
	}
	err = models.AddTopic(topic)
	return
}
