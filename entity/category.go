/**
 *
 *
 *
 * @author		John Aran (Ilyas Aranzhanovich Toxanbayev)
 * @version		1.0.0
 * @based on
 * @email      		il.aranov@gmail.com
 * @link
 * @github      	https://github.com/ilyaran/github.com/ilyaran/Malina
 * @license		MIT License Copyright (c) 2017 John Aran (Ilyas Toxanbayev)
 */ package entity

import (
	"database/sql"
	"github.com/lib/pq"
	"github.com/ilyaran/Malina/app"
	"strconv"
	"strings"
	"html/template"

)

type Category struct {
	Id          int64       `json:"id"`
	Parent      int64       `json:"parent"`
	Level       int         `json:"level"`
	Sort        int64       `json:"sort"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	Lang        string      `json:"lang"`
	Enable      bool        `json:"enable"`
	Img         []string      `json:"img"`
	Parameter	[]int64		`json:"parameter"`
	Quantity    int64       `json:"quantity"`

	Children            []*Category `json:"-"`
	Ancestors           []*Category	`json:"-"`
	Descendants         []*Category	`json:"-"`

	ParametersListView       template.HTML `json:"-"`
	DescendantsString        string        `json:"descendants_string"`
	PublicBreadcrumbPathView string        `json:"-"`
}

func (this *Category) SetId(id int64) {this.Id = id}
func (this *Category) GetId() int64 {return this.Id}
func (this *Category) GetLogo() string {
	if this.Img!=nil && len(this.Img)>0{
		return this.Img[0]
	}
	return app.Url_no_image
}
func (this *Category) GetParent() int64 {
	return this.Parent
}
func (this *Category) GetLevel() int {
	return this.Level
}
func (this *Category) GetSort() int64 {
	return this.Sort
}
func (this *Category) GetTitle() string {
	return this.Title
}
func (this *Category)  GetDescription() string {
	return this.Description
}
func (this *Category)  GetQuantity() int64 {
	return this.Quantity
}
func (this *Category)  GetLang() string {
	return this.Lang
}
func (this *Category)  GetEnable() bool {
	return this.Enable
}

func (this *Category)  SetImg(img []string) {
	this.Img = img
}
func (this *Category)  GetImg() []string {
	return this.Img
}

func (this *Category)  SetLevel(level int) {
	this.Level = level
}

func (s *Category) Scanning(row *sql.Row,rows *sql.Rows)byte {
	var err error
	if row!=nil {
		err = row.Scan(
			&s.Id,
			&s.Parent,
			&s.Sort,
			&s.Title,
			&s.Description,
			&s.Lang,
			&s.Enable,
			pq.Array(&s.Img),
			pq.Array(&s.Parameter),
			&s.Quantity,
		)
	}
	if rows!=nil {
		err = rows.Scan(
			&s.Id,
			&s.Parent,
			&s.Sort,
			&s.Title,
			&s.Description,
			&s.Lang,
			&s.Enable,
			pq.Array(&s.Img),
			pq.Array(&s.Parameter),
			&s.Quantity,
		)
	}
	if err == sql.ErrNoRows {
		panic(err)
		return 'n'
	}
	if err != nil {
		panic(err)
		return 'e'
	}

	return 'o'
}
func (s *Category) Appending(rows *sql.Rows,size int64)*[]Scanable{
	list := make([]Scanable, 0, size)
	for rows.Next() {
		item := &Category{}
		if item.Scanning(nil, rows) == 'o' {
			list = append(list, item)
		}else {
			return nil
		}
	}
	return &list
}
func (s *Category) Item(row *sql.Row) Scanable {
	item := Category{}
	if item.Scanning(row, nil) == 'o' {
		return &item
	}
	return nil
}
func (s *Category) SetAncestorsDescendants(parent int64,getDescendants bool)int{
	if getDescendants {
		if s.Parent == parent {
			s.Descendants = append(s.Descendants, s)
			//s.closure(s.Id, getDescendants)
			return 1
		}
	}else {
		if s.Id == parent {
			s.Ancestors = append([]*Category{s}, s.Ancestors...)
			//s.closure(s.Id, getDescendants)
			return 2
		}
	}
	return 3
}

func (this *Category) SetChildren(item Hierarchical) {
	if this.Children == nil{
		this.Children = []*Category{}
	}
	this.Children = append(this.Children, item.(*Category))
}

func (this *Category) SetAncestors(item Hierarchical) {
	if this.Ancestors == nil{
		this.Ancestors = []*Category{}
	}
	this.Ancestors = append([]*Category{item.(*Category)}, this.Ancestors...)
	//this.Ancestors = append(this.Ancestors,item.(*Category))
}

func (this *Category) SetDescendants(item Hierarchical) {
	if this.Descendants == nil{
		this.Descendants = []*Category{}
	}
	this.Descendants = append(this.Descendants, item.(*Category))
	if this.DescendantsString==""{
		//this.DescendantsString+= strconv.FormatInt(item.(*Category).Id,10)
	}else {

	}
	this.DescendantsString+= ","+strconv.FormatInt(item.(*Category).Id,10)
}

func (s *Category) SetPublicParametersListView(treeListParam *[]Parameter) {
	var start bool
	var level int
	var out string
	for _,a:=range s.Parameter {
		level = 0
		for _,i:=range *treeListParam {
			if i.GetLevel() < level && start {
				out += `
					</div>
				</div>`
				start = false
			}

			if i.GetId() == a && start == false {
				out += `
				<div class=" chain-grid menu-chain">
					<div class="grid-chain-bottom chain-watch">
				`
				start=true
			}

			if start {
				if i.Children != nil {
					out += strings.Repeat("&nbsp;", i.GetLevel()) + `<span class="actual">` + i.Title + `</span><br/>`
				} else {
					out += strings.Repeat("&nbsp;", i.GetLevel()) + `<span class="parameter_list" data-parameter_id="">` + i.GetDescription() + `</span><span>` + i.GetTitle() + `</span><br/>
					`
					level=i.GetLevel()
				}
			}
		}
	}
	if s.Parameter!=nil{
		out += `
				</div>
			</div>`
		s.ParametersListView=template.HTML(out)
	}

}

func (s *Category) SetPublicBreadcrumbPathView() {
	for _,i:=range s.Ancestors {
		s.PublicBreadcrumbPathView +=`<li><a href="`+app.Url_public_product_list+`?category_id=`+strconv.FormatInt(i.GetId(),10)+`">`+i.Title+`</a></li>`
	}
}





















