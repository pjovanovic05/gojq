package gojq

import (
	"testing"
)

func TestFromBytes(t *testing.T) {
	testCases := []struct {
		tag       string
		jsonStr   string
		expectErr bool
	}{
		{
			tag:       "valid json",
			jsonStr:   `{"name":"John Doe", "age": 30}`,
			expectErr: false,
		},
		{
			tag:       "invalid json",
			jsonStr:   `{"name":"John Doe", "age": 30, "oops"}`,
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		jn := FromBytes([]byte(tc.jsonStr))
		if jn.Err != nil && !tc.expectErr {
			t.Errorf("Failed %s", tc.tag)
		}
	}
}

func TestString(t *testing.T) {
	testStr := `{"age":"150","name":"abbot"}`
	jn := FromBytes([]byte(testStr))
	rep := jn.String()
	if rep != testStr {
		t.Errorf("[%v]Not equal string reps. Got <<%s>>", jn.Err, rep)
	}
}

func TestSelect(t *testing.T) {
	testStr := `{
		"_source": {
			"docid": 123,
			"content": "hello",
			"arr": [
				{"obj": 1, "name": "abbot"},
				{"obj": 2, "name": "costello"}
			]
		},
		"hits": 1
	}`
	jn := FromBytes([]byte(testStr))
	var res float64
	err := jn.Select("_source", "docid").As(&res)
	if err != nil {
		t.Error(err)
	}
	if res != 123.0 {
		t.Errorf("Expected 123.0. got %v", res)
	}

	var gotName string
	err = jn.Select("_source", "arr", "0", "name").As(&gotName)
	if err != nil {
		t.Error(err)
	}
	if gotName != "abbot" {
		t.Errorf("Expected \"abbot\", got: \"%s\"", gotName)
	}

	var hits float64
	err = jn.Select("hits").As(&hits)
	if err != nil {
		t.Error(err)
	}
	if hits != 1 {
		t.Errorf("Expected 1, got %v", hits)
	}
}

func TestSet(t *testing.T) {
	testStr := `{
		"_source": {
			"docid": 123,
			"content": "hello",
			"arr": [
				{"obj": 1, "name": "abbot"},
				{"obj": 2, "name": "costello"}
			]
		},
		"hits": 1
	}`
	jn := FromBytes([]byte(testStr))

	tmp := jn.Select("_source", "docid").Set(1234.0)
	if tmp.Err != nil {
		t.Error(tmp.Err)
	}
	var docid float64
	err := jn.Select("_source", "docid").As(&docid)
	if err != nil {
		t.Error(err)
	}
	if docid != 1234.0 {
		t.Errorf("Expecter 1234.0 got %v", docid)
	}
	// test appending to array
	var arr []interface{}
	err = jn.Select("_source", "arr").As(&arr)
	if err != nil {
		t.Error(err)
	}
	newElement := make(map[string]interface{})
	newElement["obj"] = 3.0
	newElement["name"] = "stanlio"
	arr = append(arr, newElement)
	tmp = jn.Select("_source", "arr").Set(arr)
	if tmp.Err != nil {
		t.Error(tmp.Err)
	}
	var arr2 []interface{}
	err = jn.Select("_source", "arr").As(&arr2)
	if err != nil {
		t.Error(err)
	}
	if len(arr2) != 3 {
		t.Errorf("Wrong number of elements: %d", len(arr2))
	}
	var name string
	err = jn.Select("_source", "arr", "2", "name").As(&name)
	if err != nil {
		t.Error(err)
	}
	if name != "stanlio" {
		t.Errorf("Expected name 'stanlio', got '%s'", name)
	}

	olio := make(map[string]interface{})
	olio["obj"] = 4.0
	olio["name"] = "olio"
	tmp = jn.Select("_source", "arr", "2").Set(olio)
	if tmp.Err != nil {
		t.Error(tmp.Err)
	}
	err = jn.Select("_source", "arr", "2", "name").As(&name)
	if err != nil {
		t.Error(err)
	}
	if name != "olio" {
		t.Errorf("Expected name 'olio', got '%s'", name)
	}
}

func TestDelete(t *testing.T) {
	testStr := `{
		"_source": {
			"docid": 123,
			"content": "hello",
			"arr": [
				{"obj": 1, "name": "abbot"},
				{"obj": 2, "name": "costello"},
				{"obj": 3.0, "name": "stanlio"}
			]
		},
		"hits": 1
	}`
	jn := FromBytes([]byte(testStr))

	tmp := jn.Select("_source", "docid").Delete()
	if tmp.Err != nil {
		t.Error(tmp.Err)
	}
	exists := jn.Select("_source", "docid").Exists()
	if exists {
		t.Error("expected docid not exist")
	}
	// test appending to array
	tmp = jn.Select("_source", "arr", "1").Delete()
	if tmp.Err != nil {
		t.Error(tmp.Err)
	}
	var arr []interface{}
	err := jn.Select("_source", "arr").As(&arr)
	if err != nil {
		t.Error(err)
	}
	if len(arr) != 2 {
		t.Errorf("Wrong number of elements: %d", len(arr))
	}
	var name string
	err = jn.Select("_source", "arr", "1", "name").As(&name)
	if err != nil {
		t.Error(err)
	}
	if name != "stanlio" {
		t.Errorf("Expected name 'stanlio', got '%s'", name)
	}
}

func TestExists(t *testing.T) {
	testStr := `{
		"_source": {
			"docid": 123,
			"content": "hello",
			"arr": [
				{"obj": 1, "name": "abbot"},
				{"obj": 2, "name": "costello"}
			]
		},
		"hits": 1
	}`
	jn := FromBytes([]byte(testStr))
	exists := jn.Exists("_source", "docid")
	if !exists {
		t.Error("Should exist")
	}
	exists = jn.Exists("loch", "ness", "monster")
	if exists {
		t.Error("should not exist")
	}
	jn = jn.Select("loch", "ness", "monster")
	exists = jn.Exists()
	if exists {
		t.Error("should not exist")
	}
}

func TestIterator(t *testing.T) {
	testStr := `{
		"_source": {
			"docid": 123,
			"content": "hello",
			"arr": [
				{"obj": 1, "name": "abbot"},
				{"obj": 2, "name": "costello"}
			]
		},
		"hits": 1
	}`
	jn := FromBytes([]byte(testStr))
	it := jn.Select("_source", "arr").Iterator()
	count := 0
	for it.Next() {
		count++
	}
	if count != 2 {
		t.Errorf("Wrong number of elements: %d", count)
	}
}
