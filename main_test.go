package main

import (
	"io/ioutil"
	"reflect"
	"testing"
)

const testDataPath = "test_autojump.txt"
const testDataContent = `10.000,/Users/xiaotan/Toy/pdfminer/pdfminer
14.142,/Users/xiaotan/.vim/plugged/YouCompleteMe/python/ycm
14.142,/Users/xiaotan/Toy/tx/build/lib
10.000,/Users/xiaotan/.local/venvs
22.361,/Users/xiaotan/Work/slate
14.142,/Users/xiaotan/Toy/tx/elasticsearch-py
10.000,/Users/ycm/test
`

func setupTestCase(t *testing.T) {
	if err := ioutil.WriteFile(testDataPath, []byte(testDataContent), 0644); err != nil {
		panic(err)
	}
}

func Assert(v, e interface{}, t *testing.T) {
	if !reflect.DeepEqual(v, e) {
		t.Errorf("Expected %s, got %s", e, v)
	}
}

func TestDataAdd(t *testing.T) {
	setupTestCase(t)

	d := Data{value: make(map[string]float64)}
	d.Load(testDataPath)
	v := d.value["/Users/xiaotan/Toy/pdfminer/pdfminer"]
	Assert(v, float64(10), t)
	d.Add("add_new_dir")
	d.Add("add_new_dir")
	d.Add("中文目录")
	d.Add("中文目录")
	d.Close()

	d = Data{value: make(map[string]float64)}
	d.Load(testDataPath)
	Assert(d.value["add_new_dir"], float64(10), t)
	Assert(d.value["中文目录"], float64(10), t)
	d.Close()
}

func TestMatch(t *testing.T) {
	setupTestCase(t)
	d := Data{value: make(map[string]float64)}
	d.Load(testDataPath)
	Assert(Match("Work", &d), "/Users/xiaotan/Work/slate", t)
	Assert(Match("ycm", &d), "/Users/xiaotan/.vim/plugged/YouCompleteMe/python/ycm", t)
	d.Close()
}
