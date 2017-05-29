/**
 * Category entity class.  Malina eCommerce application
 *
 *
 * @author		John Aran (Ilyas Toxanbayev)
 * @version		1.0.0
 * @based on
 * @email      		il.aranov@gmail.com
 * @link
 * @github      	https://github.com/ilyaran/Malina
 * @license		MIT License Copyright (c) 2017 John Aran (Ilyas Toxanbayev)
 */
package entity

import (
	"database/sql"
	"fmt"
	"strconv"
	"encoding/json"
	"github.com/ilyaran/Malina/config"
	"github.com/ilyaran/Malina/helpers"
	"strings"
)

type Category struct {
	id                  int64
	parent              int64
	level               int
	sort                int64
	title               string
	description         string
	lang                string
	enable              bool
	img                 []string
	quantity            int64

	children            []*Category
	ancestors           []*Category
	descendants         []*Category

	path                string

	descendantIdsString string
}

func NewCategory(id, parent, sort int64, title, description string, enabled bool) *Category {
	category := &Category{
		id:id,
		parent:parent,
		sort:sort,
		title:title,
		description:description,
		enable:enabled,
	}

	return category
}
func (this *Category) Exec() []interface{} {
	return []interface{}{
		helper.NewNullInt64(this.parent),
		this.sort,
		this.title,

		this.description,
		this.lang,
		this.enable,
	}
}
func (this *Category) ExecWithId() []interface{} {
	return []interface{}{
		this.id,

		helper.NewNullInt64(this.parent),
		this.sort,
		this.title,

		this.description,
		this.lang,
		this.enable,
	}
}
func CategoryScan(row *sql.Row,rows *sql.Rows) *Category {
	var s = &Category{}
	var images_sting string
	//var imgArray []sql.NullString
	var err error
	if row!=nil {
		err = row.Scan(
			&s.id,
			&s.parent,
			&s.sort,

			&s.title,
			&s.description,
			&s.lang,

			&s.enable,
			&s.quantity,
			&images_sting)
	}
	if rows!=nil {
		err = rows.Scan(
			&s.id,
			&s.parent,
			&s.sort,

			&s.title,
			&s.description,
			&s.lang,

			&s.enable,
			&s.quantity,
			&images_sting)
	}
	if err == sql.ErrNoRows {
		return nil
	}
	if err != nil {
		panic(err)
		return nil
	}
	//fmt.Println(s.id," == ",images_sting)
	if images_sting != `` {
		s.img = strings.Split(images_sting,"|")
	}

	return s
}

func (this *Category) GetId() int64 {
	return this.id
}
func (this *Category) GetParent() int64 {
	return this.parent
}
func (this *Category) GetLevel() int {
	return this.level
}
func (this *Category) GetSort() int64 {
	return this.sort
}
func (this *Category) GetTitle() string {
	return this.title
}
func (this *Category)  Get_description() string {
	return this.description
}
func (this *Category)  Get_quantity() int64 {
	return this.quantity
}
func (this *Category)  Get_lang() string {
	return this.lang
}
func (this *Category)  Get_enable() bool {
	return this.enable
}
func (this *Category)  Get_path() string {
	return this.path
}
func (this *Category)  Get_descendantIdsString() string {
	return this.descendantIdsString
}

func (this *Category)  Get_img() []string {
	return this.img
}
func (this *Category)  Get_logo() string {
	if len(this.img) > 0 {
		return  this.img[0]
	}
	return app.No_image()
}

func (this *Category) Get_children() []*Category {
	return this.children
}
func (this *Category) Get_ancestors() []*Category {
	return this.ancestors
}
func (this *Category) Get_descendants() []*Category {
	return this.descendants
}

func (this *Category)  Set_level(level int) {
	this.level = level
}

func (this *Category) Set_children(treeList []*Category) {
	this.children = []*Category{}
	for _, v := range treeList {
		if v.parent == this.id {
			this.children = append(this.children, v)
		}
	}
}

func (this *Category) Set_ancestors(treeList []*Category) {
	this.path = ""
	this.ancestors = []*Category{}
	var anon func(int64)
	anon = func(parent int64) {
		for _, v := range treeList {
			if v.id == parent {
				this.ancestors = append([]*Category{v}, this.ancestors...)
				this.path = v.title + " / " + this.path
				anon(v.parent)
			}
		}
	}
	this.path += this.title
	anon(this.parent)
}

func (this *Category) Set_descendants(treeList []*Category) {
	this.descendants = []*Category{}
	this.descendantIdsString = ""
	var anon func(int64)
	anon = func(parent int64) {
		for _, v := range treeList {
			if v.parent == parent {
				this.descendants = append(this.descendants, v)
				this.descendantIdsString += "," + fmt.Sprintf("%v", v.id)
				anon(v.id)
			}
		}
	}
	anon(this.id)
	this.descendantIdsString = fmt.Sprintf("%v", this.id) + this.descendantIdsString
}



//**********************************

func (this *Category)BreadCrumb(home, href, selfTitle string) string {
	var out string = `<ol class="breadcrumb"><li>` + home + `</li>`
	for _, i := range this.ancestors {
		out += `<li><a href="` + href + strconv.FormatInt(i.GetId(), 64) + `">` + i.GetTitle() + `</a></li>`
	}
	return out + selfTitle + `</ol>`
}

func (t *Category) JsonEncode() string {
	b, err := json.Marshal(t.GetJsonAdaptedType())
	if err != nil {
		return ""
		//log.Fatal(err)
	}
	//log.Printf("%s", b)
	return string(b)
}
func (t *Category) JsonDecode(jsonString []byte) *Category {
	o := new(Category)
	err := json.Unmarshal(jsonString, &o)
	if err != nil {
		//log.Fatal(err)
		return nil
	}
	//log.Printf("%#v", t)
	return o
}
func (t *Category) GetJsonAdaptedType() *JsonCategory {
	return &JsonCategory{
		t.id,
		t.parent,
		t.level,

		t.sort,
		t.title,
		t.description,

		t.quantity,
		t.lang,
		t.enable,

		t.img,

		t.children,
		t.ancestors,
		t.descendants,
	}
}

type JsonCategory struct {
	Id          int64              `json:"id"`
	Parent      int64              `json:"parent"`
	Level       int                `json:"level"`

	Sort        int64              `json:"sort"`
	Title       string             `json:"title"`
	Description string             `json:"description"`

	Quantity    int64              `json:"quantity"`
	Lang        string             `json:"lang"`
	Enable      bool               `json:"enable"`

	Img         []string           `json:"img"`

	Children    []*Category        `json:"children"`
	Ancestors   []*Category        `json:"ancestors"`
	Descendants []*Category        `json:"descendants"`
}

func (t *Category) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.GetJsonAdaptedType())
}

func (t *Category) UnmarshalJSON(b []byte) error {
	temp := &JsonCategory{}

	if err := json.Unmarshal(b, &temp); err != nil {
		return err
	}

	t.id = temp.Id
	t.parent = temp.Parent
	t.level = temp.Level
	t.sort = temp.Sort
	t.title = temp.Title
	t.description = temp.Description
	t.quantity = temp.Quantity
	t.lang = temp.Lang
	t.enable = temp.Enable
	t.img = temp.Img
	t.children = temp.Children
	t.ancestors = temp.Ancestors
	t.descendants = temp.Descendants

	return nil
}


//*********** collection
func GetCategoryCollectionFromJsonString(jsonString []byte) []*Category {

	ic := new(CategoryCollection)
	err := ic.FromJson(jsonString)
	if err != nil {
		panic(err)
		return nil
	}
	var categoryList = []*Category{}
	for _, v := range ic.Pool {
		categoryList = append(categoryList, &Category{
			id:v.Id, parent:v.Parent, level:v.Level, sort:v.Sort,
			title:v.Title, description:v.Description, quantity:v.Quantity,
			lang:v.Lang, enable:v.Enable, img:v.Img,
			children:v.Children, ancestors:v.Ancestors, descendants:v.Descendants})
	}

	return categoryList
}

type CategoryCollection struct {
	Pool []*JsonCategory
}

func (mc *CategoryCollection) FromJson(jsonStr []byte) error {
	var data = &mc.Pool
	b := jsonStr
	return json.Unmarshal(b, data)
}



