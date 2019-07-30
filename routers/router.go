package routers

import (
	"github.com/gin-gonic/gin"
	"gwyApp/routers/api"
)

func InitRouter() *gin.Engine {
	r := gin.New()
	//测试接口
	/*r.GET("/article", api.GetArticle)
	r.POST("/upload", api.UploadFile)
	r.GET("/getUser", api.GetUser)*/
	r.GET("/getAccessToken", api.GetAccessToken)
	//1.首页
	r.GET("/home", api.HomePage)
	//2.开始答题/下一题
	r.GET("/nextTopic", api.NextTopic)
	//3.下一道错题
	r.GET("/nextWrongTopic", api.NextWrongTopic)
	//4.下一道收藏的题
	r.GET("/nextCollect", api.NextCollect)
	//5.提交答案
	r.GET("/answer", api.Answer)
	//6.修改计划
	r.GET("/plan", api.Plan)
	//7.获取计划
	r.GET("/getPlan", api.GetPlan)
	//8.收藏题目
	r.GET("/collect", api.Collect)
	//9.取消收藏
	r.GET("/cancelCollect", api.CancelCollect)
	//10.查看解析
	r.GET("/getAnalysis", api.GetAnalysis)

	//11.设置要素轴
	r.GET("/setting", api.Plan)

	return r
}
