package grpctaskmanager

func FindIndexFromSlice(s []interface{}, d interface{}) int {
	for i, v := range s {
		if d == v {
			return i
		}
	}
	return -1
}

func RemoveIndexFromSlice(s []interface{}, index int) []interface{} {
	return append(s[:index], s[index+1:]...)
}
