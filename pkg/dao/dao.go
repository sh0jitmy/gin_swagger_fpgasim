package dao

import (
        "log"
        "math"
        "os"
        "syscall"
        "unsafe"	
        "strconv"	
	"gopkg.in/yaml.v2"
	"github.com/sh0jitmy/gin_swagger_fpgasim/pkg/model"
)

const ADDRBITMASK = 0xFFF //64KByte
const REGSIZE = 4 
const DECMODE = 10 

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
	InitValue uint32  `yaml:initvalue`
}

type PLRegDao struct {
	conf string
	reglist[] RegEntry //key propname, value regentry
	propmap map[string] int
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
	//log.Println(cd)
	d.propmap = make(map[string]int)
	d.regbase,err = d.IORemapReg32(cd)
	if err != nil {
		return err
	}	
	for i,v := range cd.RegList {
		d.propmap[v.PropName] = i
		//log.Printf("addr:%v,value:%v\n",v.AddrOffset,v.InitValue)
		d.regbase[v.AddrOffset] = v.InitValue	
		d.reglist = append(d.reglist,v)
	}
	return nil
}

func (d *PLRegDao) GetAll()([]model.Property,error) {
	var tmpprop model.Property
	var props []model.Property		
	
	for _,v := range d.reglist {
		tmpprop.ID = v.PropName
		tmpprop.Value = strconv.FormatUint(uint64(d.regbase[v.AddrOffset]),DECMODE)
		props = append(props,tmpprop)
	}
	return props,nil
}

/*
func (d *PLRegDao) Get(id string)(model.Property,error) {

}

func (d *PLRegDao) Set(prop model.Property)(error) {

}
*/
func (d *PLRegDao) IORemapReg32(c Config) ([]uint32,error) {
	var nildata []uint32
	f, err := os.OpenFile(c.DevNode,os.O_RDWR | os.O_CREATE,0777)
	if err != nil {
		log.Fatal(err)
		return nildata,err
	}
	initdata := make([]byte,c.DevSize)
	_,werr := f.Write(initdata)
	if werr != nil {
		log.Fatal(werr)
		return nildata,werr
	}
	
	offset := int64(c.BaseAddr) &^ ADDRBITMASK
	data, ferr := syscall.Mmap(int(f.Fd()), offset, (c.DevSize+ADDRBITMASK)&^ADDRBITMASK, 
		syscall.PROT_READ | syscall.PROT_WRITE, syscall.MAP_SHARED)
	if err != nil {
		log.Fatalf("mmap %s: %v", f.Name(), err)
	}
	f.Close()
	map_array := (*[math.MaxInt32 / REGSIZE]uint32)(unsafe.Pointer(&data[0]))
	//log.Printf("start:%v,stop:%v\n",(c.BaseAddr&ADDRBITMASK)/REGSIZE,((int(c.BaseAddr)&ADDRBITMASK)+c.DevSize))
	return map_array[(c.BaseAddr&ADDRBITMASK)/REGSIZE : ((int(c.BaseAddr)&ADDRBITMASK)+c.DevSize)/REGSIZE],ferr
}

