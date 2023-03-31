package dao 

import "testing"


func TestLoad(t *testing.T) {
	path := "./config.yaml"
	d ,err := NewPLRegDAO(path)	
	if err != nil {
		t.Errorf("new error:%v\n",err)
	}
	t.Log("GetAll")
	ps,_ := d.GetAll()
	t.Log(ps)
}
