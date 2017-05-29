package entity

import (
	"encoding/json"
	"github.com/ilyaran/Malina/config"
)

type Img struct {
	id int64
	name string
}

func (s *Img) GetId()int64{return s.id}
func (s *Img) GetName()string{return s.name}
func (s *Img) GetUrl()string{return app.Base_url()+app.Upload_path()+s.name}

func (t *Img) JsonEncode() string {
	b, err := json.Marshal(t.GetJsonAdaptedType())
	if err != nil {
		return ""
	}
	return string(b)
}

func (t *Img) JsonDecode(jsonString []byte) *Img {
	o := new(Img)
	err := json.Unmarshal(jsonString, &o)
	if err != nil {
		return nil
	}
	return o
}

func (t *Img) GetJsonAdaptedType() *JsonImg {
	return &JsonImg{
		t.id,
		t.name,
		}
}

type JsonImg struct {
	Id int64 		`json:"img_id"`
	Name string		`json:"img_name"`
}

func (t *Img) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.GetJsonAdaptedType())
}

func (t *Img) UnmarshalJSON(b []byte) error {
	temp := &JsonImg{}
	if err := json.Unmarshal(b, &temp); err != nil {
		return err
	}
	t.id = temp.Id
	t.name = temp.Name
	return nil
}

//*********** collection
func GetImgCollectionFromJsonString(jsonString []byte) []*Img {
	ic := new(ImgCollection)
	err := ic.FromJson(jsonString)
	if err != nil{
		panic(err)
		return nil
	}
	var imgList = []*Img{}
	for _,v := range ic.Pool{
		imgList=append(imgList,&Img{v.Id, v.Name})
	}

	return imgList
}

type ImgCollection struct {
	Pool []*JsonImg
}

func (mc *ImgCollection) FromJson(jsonStr []byte) error {
	var data = &mc.Pool
	b := jsonStr
	return json.Unmarshal(b, data)
}

