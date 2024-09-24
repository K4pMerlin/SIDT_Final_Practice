package process

import (
	"CengkeHelper/logger"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

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

func (receiver *TeachInfo) getWeekday() int {
	pattern := `(星期\S+)第`
	regex := regexp.MustCompile(pattern)
	match := regex.FindStringSubmatch(receiver.CourseTime[0])
	if match == nil {
		logger.Error(receiver.CourseTime[0])
	}
	switch match[1] {
	case "星期日":
		return 0
	case "星期一":
		return 1
	case "星期二":
		return 2
	case "星期三":
		return 3
	case "星期四":
		return 4
	case "星期五":
		return 5
	case "星期六":
		return 6
	default:
		logger.Error("不合法的周格式: ", match[1], receiver)
		return -1
	}
}

func (receiver *TeachInfo) getLessonDuration() string {
	_, _, lessonNum := CurCourseTime()

	durationMatch := regexp.MustCompile(`第(\S+)节`).
		FindStringSubmatch(receiver.CourseTime[0])
	if durationMatch == nil {
		logger.Error(receiver.CourseTime[0])
	}
	durationStr := durationMatch[1]

	// 单独匹配剩下的部分
	regex := regexp.MustCompile(`(\d+)-(\d+)`)
	matches := regex.FindAllStringSubmatch(durationStr, -1)

	if matches == nil {
		logger.Error(receiver.CourseTime[0])
	}
	for _, match := range matches {
		begin, err := strconv.Atoi(match[1])
		if err != nil {
			logger.ErrorF("getNum from %v error: %v", receiver.CourseTime[0], err)
		}

		end, err := strconv.Atoi(match[2])
		if err != nil {
			logger.ErrorF("getNum from %v error: %v", receiver.CourseTime[0], err)
		}
		// 只要有一个时间段被匹配到就行
		if lessonNum >= begin && lessonNum <= end {
			return "第" + match[0] + "节"
		}
	}
	logger.Error("未匹配到任何时间段", receiver)

	return "第 -1 节"
}

func (receiver *TeachInfo) getRoom() string {
	addr := receiver.Address[0]
	//if strings.Contains(addr, "A318") {
	//	println("tst")
	//}

	if addr == "--" ||
		strings.Contains(addr, "虚拟教室") ||
		strings.Contains(addr, "体育馆") {
		return addr
	}

	areas := []string{
		"1区", "2区", "3区", "国软", "4区", "新珈楼",
	}

	for _, area := range areas {
		if !strings.Contains(addr, area) {
			continue
		}

		res := strings.ReplaceAll(addr, area, "")
		return res
		//if i > 3 {
		//	return res
		//}
		//
		//splits := strings.Split(res, "-")
		//if len(splits) == 2 {
		//	return splits[1]
		//}
		//if strings.Contains(splits[0], "A") {
		//	return strings.TrimSuffix(res, "A")
		//}
		//
		//logger.Error(addr)
	}

	logger.Error(receiver)
	return ""
}

func (receiver *TeachInfo) MatchType(courseType string) bool {
	if !slices.Contains([]string{"英语课", "通识课", "专业课"}, courseType) {
		logger.Error("不合法的类型: ", courseType)
	}

	if courseType == receiver.CourseType {
		return true
	}
	return false
}
func (receiver *TeachInfo) MatchTime(weekNum int, weekday int, lessonNum int) bool {
	return receiver.matchWeekNum(weekNum) &&
		receiver.matchWeekday(weekday) &&
		receiver.matchLessonNum(lessonNum)
}

// 匹配第几周
func (receiver *TeachInfo) matchWeekNum(weekNum int) bool {

	if weekNum == -1 {
		return true
	}

	pattern := `(?:(\d+)[-](\d+)周)|(?:(\d+)周)`
	regex := regexp.MustCompile(pattern)
	originMatch := regex.FindStringSubmatch(receiver.CourseTime[0])
	if originMatch == nil {
		logger.Error(receiver, receiver.CourseTime[0])
		return false
	}

	// 一共使用了3个捕获组,所以会分三组
	// 但是使用了或, 必然会有无效的捕获组(或者说未捕获到数据的捕获组)
	// 这里可知,捕获组是从左到右按顺序输出到切片的
	match := make([]string, 0)
	for _, s := range originMatch {
		//logger.Debug(s)
		if len(s) > 0 {
			match = append(match, s)
		}
	}

	if len(match) > 3 || len(match) <= 1 {
		for _, s := range match {
			logger.Warning(s)
		}
		logger.ErrorF("不合法的 %v 分组: %v", len(match), receiver.CourseTime[0])
	}

	if len(match) == 3 {
		// begin
		begin, err := strconv.Atoi(match[1])
		if err != nil {
			logger.ErrorF("getNum from %v error: %v", receiver.CourseTime[0], err)
		}

		// end
		end, err := strconv.Atoi(match[2]) // 获取第二个数字！
		if err != nil {
			logger.ErrorF("getNum from %v error: %v", receiver.CourseTime[0], err)
		}

		if weekNum >= begin && weekNum <= end {
			return true
		}
		return false
	}

	// 长度为2
	// match == 2
	num, err := strconv.Atoi(match[1])
	if err != nil {
		logger.Warning(err)
	}
	return weekNum == num

}

// 匹配哪一天(周几)
func (receiver *TeachInfo) matchWeekday(weekday int) bool {
	if weekday == -1 {
		return true
	}

	if weekday < 0 || weekday > 6 {
		logger.Error("使用了不合法的 weekday: ", weekday)
		return false
	}
	return weekday == receiver.getWeekday()
}

// 匹配第几节课
func (receiver *TeachInfo) matchLessonNum(lessonNum int) bool {

	if lessonNum == -1 {
		return true
	}

	if lessonNum == -2 || lessonNum == -3 || lessonNum == -4 {
		// 没课时间
		return false
	}
	durationMatch := regexp.MustCompile(`第(\S+)节`).
		FindStringSubmatch(receiver.CourseTime[0])
	if durationMatch == nil {
		logger.Error(receiver.CourseTime[0])
	}
	durationStr := durationMatch[1]

	// 单独匹配剩下的部分
	pattern := `(\d+)-(\d+)`
	// 第1-2,11-13节
	regex := regexp.MustCompile(pattern)
	matches := regex.FindAllStringSubmatch(durationStr, -1)
	//logger.Warning(matches)
	//if receiver.CourseTime[0] == "星期二第1-2,11-13节{1-11周}" {
	//	logger.Error(durationStr)
	//}

	if matches == nil {
		logger.Error(receiver.CourseTime[0])
	}
	for _, match := range matches {
		begin, err := strconv.Atoi(match[1])
		if err != nil {
			logger.ErrorF("getNum from %v error: %v", receiver.CourseTime[0], err)
		}

		end, err := strconv.Atoi(match[2])
		if err != nil {
			logger.ErrorF("getNum from %v error: %v", receiver.CourseTime[0], err)
		}
		// 只要有一个时间段被匹配到就行
		if lessonNum >= begin && lessonNum <= end {
			return true
		}
	}

	return false

}

// 初步整理课程信息并写入文件
