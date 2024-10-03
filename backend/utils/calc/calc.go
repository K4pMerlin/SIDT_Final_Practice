package calc

import "sort"

func IsTargetInArray(target string, array []string) bool {
	// 切片必须升序
	sort.Strings(array)
	index := sort.SearchStrings(array, target)
	//index的取值：0 ~ (len(str_array)-1)
	return index < len(array) && array[index] == target
}
