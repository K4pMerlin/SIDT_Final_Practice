package process

type Teacher struct {
	Name  string `json:"name"`  // 陈亮
	Title string `json:"title"` // 教授
}

// TeachInfo 授课信息
type TeachInfo struct {
	CourseType string    `json:"courseType"` // 专业课 （专业课、通识课、英语课）
	CourseName string    `json:"courseName"` // 位置服务与实践
	Teachers   []Teacher `json:"teachers"`   // [{陈亮, 教授}]
	Faculty    string    `json:"faculty"`    // 测绘学院
	Address    []string  `json:"address"`    // [3区1-304, 3区1-304]
	CourseTime []string  `json:"courseTime"` // [星期一第6-7节{1-9周}, 星期四第6-7节{1-9周}]
	RegionName string    `json:"region"`     // 信息学部
	RegionId   string    `json:"regionId"`   // 3区 (更准确？)
}
