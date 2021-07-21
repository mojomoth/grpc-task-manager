package grpctaskmanager

import (
	"fmt"
	"testing"
)

func TestRemoveString(t *testing.T) {
	arr := []interface{}{"A", "B", "C", "D", "E", "F"}
	result := RemoveIndexFromSlice(arr, 2)
	fmt.Println(result, len((result)))
	fmt.Println(arr, len(arr))
}

func TestRemoveInterface(t *testing.T) {
	strs := []string{"A", "B", "C", "D", "E", "F"}
	interfaces := make([]interface{}, len(strs))
	for i := range strs {
		interfaces[i] = strs[i]
	}

	RemoveIndexFromSlice(interfaces, 2)
	fmt.Println(interfaces)
}

func TestRemoveIndex(t *testing.T) {
	arr := []interface{}{"A", "B", "C", "D", "E", "F"}
	result := FindIndexFromSlice(arr, "C")
	fmt.Println(result)
}
