package process

import (
	"CengkeHelper/logger"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"
)

func findCurTimeTeachInfosMap(buildingMap map[string][]TeachInfo) map[string][]TeachInfo {
	weekNum, weekday, lessonNum := CurCourseTime()
	//fmt.Printf("%d %d %d", weekNum, weekday, lessonNum)
	courses := filterMapByCondition(buildingMap, func(info TeachInfo) bool {
		if info.MatchTime(weekNum, weekday, lessonNum) {
			return true
		}
		return false
	})
	return courses
}

func CurCourseTime() (weekNum int, weekday int, lessonNum int) {
	now := time.Now()
	if ChoseSpecialDate.Chose {
		month := now.Month()
		if ChoseSpecialDate.Month != nil {
			month = time.Month(*ChoseSpecialDate.Month)
		}
		day := now.Day()
		if ChoseSpecialDate.Day != nil {
			day = *ChoseSpecialDate.Day
		}
		hour := now.Hour()
		if ChoseSpecialDate.Hour != nil {
			hour = *ChoseSpecialDate.Hour
		}
		minute := now.Minute()
		if ChoseSpecialDate.Min != nil {
			minute = *ChoseSpecialDate.Min
		}
		//now = time.Date(2024, time.March, 1,
		//	10, 14, 0, 0, time.Local)
		now = time.Date(now.Year(), month, day,
			hour, minute, 0, 0, time.Local)

		//logger.Debug("self specialize time: ", now)
	}

	// 计算第几周
	beginDate := time.Date(2024, time.September, 9,
		0, 0, 0, 0, time.Local)

	sub := now.Sub(beginDate)
	durationDay := int(sub.Hours()) / 24
	weekNum = durationDay/7 + 1

	// 计算周几
	weekday = int(now.Weekday())

	// 计算第几节课
	if isTimeBeforeHourAndMin(now, 7, 50) { // 8点前，早上
		lessonNum = -2
	} else if isTimeBeforeHourAndMin(now, 8, 45) { // 8点45前，第1节
		lessonNum = 1
	} else if isTimeBeforeHourAndMin(now, 9, 35) { // 9点35前，第2节
		lessonNum = 2
	} else if isTimeBeforeHourAndMin(now, 10, 35) { // 10点35前，第3节
		lessonNum = 3
	} else if isTimeBeforeHourAndMin(now, 11, 25) { // 11点25前，第4节
		lessonNum = 4
	} else if isTimeBeforeHourAndMin(now, 12, 15) { // 12点15前，第5节
		lessonNum = 5
	} else if isTimeBeforeHourAndMin(now, 13, 55) { // 中午，13点55之前, 没课
		lessonNum = -3
		// 中文没课捏
	} else if isTimeBeforeHourAndMin(now, 14, 50) { // 14点50前，第6节
		lessonNum = 6
	} else if isTimeBeforeHourAndMin(now, 15, 40) { // 15点40前，第7节课
		lessonNum = 7
	} else if isTimeBeforeHourAndMin(now, 16, 35) { // 16点35前，第8节课
		lessonNum = 8
	} else if isTimeBeforeHourAndMin(now, 17, 25) { // 17点25前，第9节课
		lessonNum = 9
	} else if isTimeBeforeHourAndMin(now, 18, 15) { // 18点15前，第10节课（一般不排课吧）
		lessonNum = 10
	} else if isTimeBeforeHourAndMin(now, 18, 20) { // 18点20前，晚饭时间，没课
		lessonNum = -4
		// 晚饭没课捏
	} else if isTimeBeforeHourAndMin(now, 19, 15) { // 19点15前，第11节课
		lessonNum = 11
	} else if isTimeBeforeHourAndMin(now, 20, 05) { // 20点05前，第12节课
		lessonNum = 12
	} else if isTimeBeforeHourAndMin(now, 20, 55) { // 20点55前，第13节课
		lessonNum = 13
	} else { // 晚上了
		// 今天的课上完了
		lessonNum = -5
	}
	return weekNum, weekday, lessonNum
}

func isTimeBeforeHourAndMin(inputTime time.Time, hour int, min int) bool {
	return inputTime.Before(
		time.Date(inputTime.Year(), inputTime.Month(), inputTime.Day(),
			hour, min, 0, 0,
			inputTime.Location()))
}

func filterMapByCondition(buildingMap map[string][]TeachInfo,
	condition func(info TeachInfo) bool) map[string][]TeachInfo {
	res := make(map[string][]TeachInfo)
	for building, teachInfos := range buildingMap {
		tempArray := make([]TeachInfo, 0)
		for _, info := range teachInfos {
			if len(info.CourseTime) != 1 {
				logger.Error("assert error!", info)
			}

			if condition(info) {
				//logger.Warning(building)
				tempArray = append(tempArray, info)
			}
		}
		if len(tempArray) == 0 {
			continue
		}

		// 为tempArray排序
		sort.Slice(tempArray, func(i, j int) bool {
			return tempArray[i].Address[0] < tempArray[j].Address[0]
		})
		res[building] = tempArray
	}

	return res
}

func findBuildingMapWithKeys(partNum int) (map[string][]TeachInfo, []string) {

	infos := readTeachInfosByPart(partNum)
	buildingMap := make(map[string][]TeachInfo)

	switch partNum {
	case 1, 2, 3:
		// 1-3区这里处理
		partName := fmt.Sprintf("%d区", partNum)
		for _, info := range infos {
			if len(info.Address) != 1 {
				logger.Error("处理后的授课地点应当只有一个: ", info)
			}

			if info.Address[0] == "--" {
				// 跳过
				continue
			}

			splits := strings.Split(info.Address[0], "-")
			if len(splits) > 2 {
				logger.Error("授课地点切分后大于2,有误: ", info)
			}

			curBuilding := ""
			if len(splits) == 2 {
				curBuilding = strings.ReplaceAll(splits[0], partName, "")
				curBuilding += "-教学楼"
				buildingMap[curBuilding] = append(buildingMap[curBuilding], info)
				continue
			}

			if strings.Contains(info.Address[0], partName) {
				// 包含 x区 且不能切分,可能是字母
				curBuilding = string([]rune(info.Address[0])[2:3])
				curBuilding += "-教学楼"
				buildingMap[curBuilding] = append(buildingMap[curBuilding], info)
				continue
			}

			// 否则默认
			curBuilding = info.Address[0]
			buildingMap[curBuilding] = append(buildingMap[curBuilding], info)

			// 否则不合法
			//logger.ErrorF("不合法的上课教室: %v \n%v", info.Address[0], info)
		}
		break

	case 4, 5:
		// 4-5区 这里处理
		// 教室少, 只显示有课的教室
		for _, info := range infos {
			if len(info.Address) != 1 {
				logger.Error("处理后的授课地点应当只有一个: ", info)
			}
			curBuilding := info.Address[0]

			if strings.Contains(curBuilding, "4区") {
				curBuilding = strings.ReplaceAll(curBuilding, "4区", "")
			}

			buildingMap[curBuilding] = append(buildingMap[curBuilding], info)
		}
		break

	default:
		logger.Error("不合法的 partNum: ", partNum)
		break
	}
	keys := make([]string, 0)
	for key := range buildingMap {
		keys = append(keys, key)
	}

	//sort.Slice(keys, func(i, j int) bool {
	//
	//})
	sort.Strings(keys)

	return buildingMap, keys
}

func readTeachInfosByPart(partNum int) []TeachInfo {
	if partNum < 1 || partNum > 5 {
		return nil
	}

	fileName := fmt.Sprintf("data/part%d/teach_info.json", partNum)
	data, err := os.ReadFile(fileName)
	if err != nil {
		logger.Error(err)
		return nil
	}
	teachInfos := make([]TeachInfo, 0)
	err = json.Unmarshal(data, &teachInfos)
	if err != nil {
		logger.ErrorF("解析文件 %v 失败: %v", fileName, err)
		return nil
	}

	//logger.Debug(len(teachInfos))
	return teachInfos
}

// 分区
