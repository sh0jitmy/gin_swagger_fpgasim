package dao

import (
        "log"
        "errors"
        "math"
        "os"
        "syscall"
        "unsafe"	
        "strconv"	
	"gopkg.in/yaml.v2"
	"github.com/sh0jitmy/gin_swagger_fpgasim/pkg/model"
	"github.com/sh0jitmy/gin_swagger_fpgasim/pkg/util"
)

const ADDRBITMASK = 0xFFF //64KByte
const REGSIZE = 4 
const DECMODE = 10 

type Config struct {
	DevNode string `yaml:"devnode"`
	BaseAddr int64 `yaml:"baseaddr"`
	Initial bool `yaml:"initial"`
	DevSize int `yaml:"devsize"`
	RegList []RegEntry `yaml:inline` 
}

type RegEntry struct {
	PropName string `yaml:propname`
	RegName string  `yaml:regname`
	AddrOffset int32  `yaml:addroffset`
	RegValue uint32  `yaml:regvalue`
	UpdatedAt string 
}

type PLRegDao struct {
	conf string
	regList[] RegEntry 
	propmap map[string] int //key propname ,value regentry index
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
		//log.Printf("addr:%v,value:%v\n",v.AddrOffset,v.RegValue)
		d.regbase[v.AddrOffset] = v.RegValue	
		v.UpdatedAt = util.TimeNowString()	
		d.regList = append(d.regList,v)
	}
	return nil
}

func (d *PLRegDao) GetAll()([]model.Property,error) {
	var tmpprop model.Property
	var props []model.Property		
	
	for _,v := range d.regList {
		tmpprop.ID = v.PropName
		tmpprop.Value = strconv.FormatUint(uint64(d.regbase[v.AddrOffset]),DECMODE)
		if (d.regbase[v.AddrOffset] != v.RegValue) {
			v.RegValue = d.regbase[v.AddrOffset]
			v.UpdatedAt = util.TimeNowString()
		}
		tmpprop.UpdatedAt = v.UpdatedAt
		props = append(props,tmpprop)
	}
	return props,nil
}

func (d *PLRegDao) Get(id string)(model.Property,error) {
	var tmpprop model.Property
	if val, ok := d.propmap[id]; ok {
		tmpprop.ID = id
		tmpprop.Value = strconv.FormatUint(uint64(d.regbase[d.regList[val].AddrOffset]),DECMODE)
		if (d.regbase[d.regList[val].AddrOffset] != d.regList[val].RegValue) {
			d.regList[val].RegValue = d.regbase[d.regList[val].AddrOffset]
			d.regList[val].UpdatedAt = util.TimeNowString()
		}
		tmpprop.UpdatedAt = d.regList[val].UpdatedAt
	} 
	return tmpprop,nil	
}

func (d *PLRegDao) Set(setp model.Property)(error) {
	err := errors.New("id is notfound")	
	if val, ok := d.propmap[setp.ID]; ok {
		value, cerr  := strconv.ParseUint(setp.Value,DECMODE,32)
		if cerr  != nil {
			err = cerr
		} else {
			d.regbase[d.regList[val].AddrOffset] = uint32(value) 
			d.regList[val].RegValue = d.regbase[d.regList[val].AddrOffset]
			log.Printf("updatedAt(before):%v\n",d.regList[val].UpdatedAt)
			d.regList[val].UpdatedAt = util.TimeNowString()
			log.Printf("updatedAt(after):%v\n",d.regList[val].UpdatedAt)
			err = nil
		}
	} 
	return err
}

func (d *PLRegDao) IORemapReg32(c Config) ([]uint32,error) {
	var nildata []uint32
	f, err := os.OpenFile(c.DevNode,os.O_RDWR | os.O_CREATE,0777)
	if err != nil {
		log.Fatal(err)
		return nildata,err
	}
	if c.Initial {
		log.Println("initial on")
		initdata := make([]byte,c.DevSize)
		_,werr := f.Write(initdata)
		if werr != nil {
			log.Fatal(werr)
			return nildata,werr
		}
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

