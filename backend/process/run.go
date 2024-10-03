package process

import (
	"sort"
	"strings"
	"sync/atomic"
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

	GetTeachInfos(false)
}

func ValidCache() bool {
	weekNum, weekday, lessonNum := CurCourseTime()
	if int32(weekNum) == curWeekNum.Load() &&
		int32(weekday) == curWeekday.Load() &&
		int32(lessonNum) == curLessonNum.Load() {
		return true
	}
	return false

}

func FreshCacheFlag() {
	weekNum, weekday, lessonNum := CurCourseTime()
	curWeekNum.Store(int32(weekNum))
	curWeekday.Store(int32(weekday))
	curLessonNum.Store(int32(lessonNum))

}

func GetTeachInfos(cacheable bool) [][]BuildingTeachInfos {
	if cacheable {
		return respTeachInfos
	}

	res := make([][]BuildingTeachInfos, 0)

	for _, buildingMap := range buildingMapArray {
		// 每个学部
		departmentRes := make([]BuildingTeachInfos, 0)

		teachInfosMap := findCurTimeTeachInfosMap(buildingMap)
		for building, teachInfos := range teachInfosMap {
			infosRes := make([]RespTeachInfo, 0)
			for _, info := range teachInfos {
				infosRes = append(infosRes, RespTeachInfo{
					Room:         info.getRoom(),
					Faculty:      info.Faculty,
					CourseName:   info.CourseName,
					TeacherName:  info.Teachers[0].Name,
					TeacherTitle: info.Teachers[0].Title,
					CourseTime:   info.getLessonDuration(),
					CourseType:   info.CourseType,
				})
			}
			departmentRes = append(departmentRes, BuildingTeachInfos{Building: building, Infos: infosRes})

		}

		// 教学楼排序
		sort.Slice(departmentRes, func(i, j int) bool {
			return departmentRes[i].Building < departmentRes[j].Building
		})

		res = append(res, departmentRes)
	}

	// 更新缓存
	respTeachInfos = res

	return res
}

func SearchCourses(query string) []RespTeachInfo {
	var results []RespTeachInfo
	query = strings.ToLower(query)

	for _, buildingInfos := range respTeachInfos {
		for _, building := range buildingInfos {
			for _, teachInfo := range building.Infos {
				if strings.Contains(strings.ToLower(teachInfo.CourseName), query) ||
					strings.Contains(strings.ToLower(teachInfo.TeacherName), query) ||
					strings.Contains(strings.ToLower(teachInfo.Room), query) {
					results = append(results, teachInfo)
				}
			}
		}
	}

	return results
}
