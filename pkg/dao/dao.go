package dao

import (
        "log"
        "math"
        "os"
        "syscall"
        "unsafe"	
	"gopkg.in/yaml.v2"
	"github.com/sh0jitmy/gin_swagger_fpgasim/pkg/model"
)

const ADDRBITMASK = 0xFFFF //64KByte
const REGSIZE = 4 

type Config struct {
	DevNode string `yaml:"devnode"`
	BaseAddr int64 `yaml:"baseaddr"`
	DevSize int `yaml:"devsize"`
	RegList []RegEntry `yaml:inline` 
}

type RegEntry struct {
	PropName string `yaml:propname`
	RegName string  `yaml:regname`
	AddrOffset int32  `yaml:addroffset`
	Value int32  `yaml:initvalue`
}

type PLRegDao struct {
	conf string
	propmap map[string] RegEntry //key propname, value regentry
	regbase []uint32
}

func NewPLRegDAO(confpath string) (*PLRegDao,error) {
	dao := PLRegDao{conf:confpath}
	err := dao.Initialize()
	return &dao,err
}
func (d *PLRegDao) Initialize() (error){
	cd := Config{}
	b,err := os.ReadFile(d.conf) 
	yaml.Unmarshal(b,&cd)
	log.Println(cd)
	d.propmap = make(map[string]RegEntry)
	d.regbase,err = d.IORemapReg32(cd)
	if err != nil {
		return err
	}	
	for _,v := range cd.RegList {
		d.propmap[v.PropName] = v
		d.regbase[v.AddrOffset] = v.InitValue	
	}
	return nil
}

func (d *PLRegDao) GetAll()([]model.Property,error) {
		
}

func (d *PLRegDao) Get(id string)(model.Property,error) {

}

func (d *PLRegDao) Set(prop model.Property)(error) {



}

func (d *PLRegDao) IORemapReg32(c Config) ([]uint32,error) {
	var nildata []uint32
	f, err := os.OpenFile(c.DevNode,os.O_RDWR | os.O_CREATE,0777)
	if err != nil {
		log.Fatal(err)
		return nildata,err
	}
	offset := int64(c.BaseAddr) &^ ADDRBITMASK
	data, ferr := syscall.Mmap(int(f.Fd()), offset, (c.DevSize+ADDRBITMASK)&^ADDRBITMASK, 
		syscall.PROT_READ | syscall.PROT_WRITE, syscall.MAP_SHARED)
	if err != nil {
		log.Fatalf("mmap %s: %v", f.Name(), err)
	}
	f.Close()
	map_array := (*[math.MaxInt32 / REGSIZE]uint32)(unsafe.Pointer(&data[0]))
	return map_array[(c.BaseAddr&ADDRBITMASK)/REGSIZE : ((int(c.BaseAddr)&ADDRBITMASK)+c.DevSize)/REGSIZE],ferr
}

