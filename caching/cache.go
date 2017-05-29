package caching

import "time"

var T0 time.Time
var T1 time.Time


var PublicPages = map[string][]byte{}

type Page struct {
	content []byte
	timestamp time.Time
}
func NewPage(content []byte)*Page{
	return &Page{content:content,timestamp:time.Now()}
}

func Get(keyName string)[]byte{

	if v,ok := PublicPages[keyName]; ok{
		return v
	}
	return nil
}

func Set(keyName string, cacheData []byte){
	PublicPages[keyName]=cacheData
}