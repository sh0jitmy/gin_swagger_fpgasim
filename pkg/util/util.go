package util

import (
    "time"  
)

const TIMEFORMAT = "2006-01-02 15:04:05"

//JST Time Format Gen Function
func TimeNowString()(string) {
 	jst,_ := time.LoadLocation("Asia/Tokyo")
	now := time.Now().In(jst)
	return now.Format(TIMEFORMAT)
}

func TimeStringParse(timestr string)(time.Time,error) {
	jst, _ := time.LoadLocation("Asia/Tokyo")
	return time.ParseInLocation(TIMEFORMAT,timestr,jst)
} 
