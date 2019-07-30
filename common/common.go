package common

//地域选择 国家:REGION_COUNTRY 深圳:REGION_SZ
const (
	REGION_COUNTRY = "REGION_COUNTRY"
	REGION_SZ      = "REGION_SZ"
)

//考试类型 公务员考：EXAM_TYPE_CIVIL 事业单位考：EXAM_TYPE_INSTITUTION
const (
	EXAM_TYPE_CIVIL       = "EXAM_TYPE_CIVIL"
	EXAM_TYPE_INSTITUTION = "EXAM_TYPE_INSTITUTION"
)

//还有两层分类

const OPENID_SET = "OPENID_SET"

const TODAY_FINISH_PREFIX = "TODAY_FINISH_"

const TODAY_PREFIX = "today_prefix_"

const TOPIC_LIST = "topic_list_"
const COLLECT_LIST = "collect_list_"
const WRONG_TOPIC_LIST = "wrong_topic_list_"

//上一题 下一题
const (
	OPERATE_LAST = "last"
	OPERATE_NEXT = "next"
)
