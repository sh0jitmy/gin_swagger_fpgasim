package util 

import (
	"testing"
)

func TestLoad(t *testing.T) {
	timestr := TimeNowString()
	t.Log(timestr)

	tm ,err := TimeStringParse(timestr)
	if err != nil {
		t.Errorf("error:%v\n",err)
	} else {
		t.Log(tm)
	}
}
