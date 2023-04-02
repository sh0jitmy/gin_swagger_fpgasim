package dao 

import (
	"time"
	"testing"
	"github.com/sh0jitmy/gin_swagger_fpgasim/pkg/model"
)

func TestLoad(t *testing.T) {
	path := "./config.yaml"
	d ,err := NewPLRegDAO(path)	
	if err != nil {
		t.Errorf("new error:%v\n",err)
	}
	t.Log("GetAll")
	ps,_ := d.GetAll()
	t.Log(ps)
	
	t.Log("Sleep 3sec")
	time.Sleep(3*time.Second)	

	t.Log("Set")
	p := model.Property{ID: "frequency",Value: "310000000"}
	err = d.Set(p)
	if err != nil {
		t.Errorf("error:%v\n",err)
	}
	t.Log("GetAll")
	ps,_ = d.GetAll()
	t.Log(ps)
}
