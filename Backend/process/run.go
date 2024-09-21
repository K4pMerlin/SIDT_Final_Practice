package process

import (
    "strings"
)

var ChoseSpecialDate = new(CourseDate)

type CourseDate struct {
	Chose bool `json:"chose"`
	Month *int `json:"month,omitempty"`
	Day   *int `json:"day,omitempty"`
	Hour  *int `json:"hour,omitempty"`
	Min   *int `json:"min,omitempty"`
}

type RespTeachInfo struct {
	Room         string `json:"room"`
	Faculty      string `json:"faculty"`
	CourseName   string `json:"courseName"`
	TeacherName  string `json:"teacherName"`
	TeacherTitle string `json:"teacherTitle"`
	CourseTime   string `json:"courseTime"`
	CourseType   string `json:"courseType"`
}

// BuildingTeachInfos 每个学部各个教学楼的课程信息
type BuildingTeachInfos struct {
	Building string          `json:"building"`
	Infos    []RespTeachInfo `json:"infos"`
}

var buildingMapArray []map[string][]TeachInfo

// 5个学部
var respTeachInfos = make([][]BuildingTeachInfos, 5)

var curWeekNum = atomic.Int32{}
var curWeekday = atomic.Int32{}
var curLessonNum = atomic.Int32{}


func init() {
	num, weekday, lessonNum := CurCourseTime()
	curWeekNum.Store(int32(num))
	curWeekday.Store(int32(weekday))
	curLessonNum.Store(int32(lessonNum))

	for i := 1; i <= 5; i++ {
		buildingMap, _ := findBuildingMapWithKeys(i)
		buildingMapArray = append(buildingMapArray, buildingMap)
	}
}

func SearchCourses(query string) []RespTeachInfo {
    var results []RespTeachInfo
    query = strings.ToLower(query)

    for _, info := range respTeachInfos {
        if strings.Contains(strings.ToLower(info.CourseName), query) ||
           strings.Contains(strings.ToLower(info.TeacherName), query) ||
           strings.Contains(strings.ToLower(info.ClassRoom), query) {
            results = append(results, info)
        }
    }

    return results
}