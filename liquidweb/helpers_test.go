package liquidweb

import (
	"reflect"
	"testing"
)

func TestExpandSetToStrings(t *testing.T) {
	names := []string{"My", "name", "is", "blars"}
	nameList := make([]interface{}, len(names))
	for i, s := range names {
		nameList[i] = s
	}

	stringNameList := expandSetToStrings(nameList)
	k := reflect.ValueOf(stringNameList)
	if k.Kind() != reflect.Slice {
		t.Errorf("list was not made up of strings, got %v", stringNameList)
	}
}
