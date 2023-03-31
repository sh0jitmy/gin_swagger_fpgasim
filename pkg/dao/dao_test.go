package dao 

import "testing"


func TestLoad(t *testing.T) {
	path := "./config.yaml"
	_ ,err := NewPLRegDAO(path)	
	if err != nil {
		t.Errorf("new error:%v",err)
	}
}
