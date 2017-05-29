package library

import (
	"github.com/ilyaran/Malina/entity"
	"strings"
	"fmt"
	"github.com/ilyaran/Malina/config"
	"strconv"
)

var CATEGORY *CategoryLibrary

type CategoryLibrary struct {
	TreeList          []*entity.Category
	TreeMap           map[int64]*entity.Category
	SelectOptionsList string
	PublicListView    string
	Id, Parent        int64
}

func (this *CategoryLibrary) SetTrees(categoryList []*entity.Category) {
	if len(categoryList) > 0 {
		this.TreeList = []*entity.Category{}
		this.TreeMap = map[int64]*entity.Category{}

		var level int = 0

		var closure func(int64)
		closure = func(parent int64) {
			for i := 0; i < len(categoryList); i++ {
				if categoryList[i].GetParent() == parent {
					categoryList[i].Set_level(level)
					this.TreeMap[categoryList[i].GetId()] = categoryList[i]
					this.TreeList = append(this.TreeList, categoryList[i])
					level ++
					closure(categoryList[i].GetId())
				}
			}
			level --
		}
		closure(0)

		if len(this.TreeList) > 0 {
			for _, v := range this.TreeList {
				v.Set_children(this.TreeList)
				v.Set_ancestors(this.TreeList)
				v.Set_descendants(this.TreeList)
			}
		}
		this.SetSelectOptionsView()
		this.SetPublicListView()
	}
}
func (this *CategoryLibrary) SetSelectOptionsView() {
	this.SelectOptionsList = ""
	const dis string = ` disabled="disabled"`
	const sel string = ` selected="selected"`
	for _, v := range this.TreeList {
		this.SelectOptionsList += `<option`
		if v.GetId() == this.Id {
			this.SelectOptionsList += dis
		}
		if v.GetId() == this.Parent {
			this.SelectOptionsList += sel
		}
		this.SelectOptionsList += ` value="` + fmt.Sprintf("%d", v.GetId()) + `">` + strings.Repeat("&rarr;", v.GetLevel()) + v.GetTitle() + `</option>`
	}
}

func (this *CategoryLibrary) SetPublicListView() string {
	this.PublicListView = ``
	for _, v := range this.TreeList {
		if v.Get_enable() && v.GetParent() == 0 {
			if len(v.Get_children()) == 0 {
				this.PublicListView += `
				<li><a href="` + app.Base_url() + app.Uri_public_product_list() + `?category_id=` + strconv.FormatInt(v.GetId(), 10) + `">
				<img src="` + v.Get_logo() + `" width="30" height="25">&nbsp;` + v.GetTitle() + `</a></li>`
			} else {
				this.PublicListView += `
				<li class="item11"><a href="` + app.Base_url() + app.Uri_public_product_list() + `?category_id=` + strconv.FormatInt(v.GetId(), 10) + `">
				<img src="` + v.Get_logo() + `" width="30" height="25">&nbsp;` + v.GetTitle() + `<img class="arrow-img" src="` + app.Assets_public_path() + `images/arrow1.png" alt=""/> </a>
					<ul>`

					this.getChild(v.GetId())

					this.PublicListView += `

					</ul>
				</li>`
			}
		}
	}
	return this.PublicListView
}

func (this *CategoryLibrary) getChild(id int64) {

	for _, v := range this.TreeList {
		if v.Get_enable() && v.GetParent() == id {
			if len(v.Get_children()) == 0 {
				this.PublicListView += `
				<li><a href="` + app.Base_url() + app.Uri_public_product_list() + strconv.FormatInt(v.GetId(), 10) + `">
				<img src="` + v.Get_logo() + `" width="30" height="25">&nbsp;` + v.GetTitle() + `</a></li>`
			} else {
				this.PublicListView += `
				<li class="item11"><a href="` + app.Base_url() + app.Uri_public_product_list() + `?category_id=` + strconv.FormatInt(v.GetId(), 10) + `">
				<img src="` + v.Get_logo() + `" width="30" height="25">&nbsp;` + v.GetTitle() + `<img class="arrow-img" src="` + app.Assets_public_path() + `images/arrow1.png" alt=""/> </a>
					<ul>`

				this.getChild(v.GetId())

				this.PublicListView += `
					</ul>
				</li>`
			}
		}
	}
}
func (this *CategoryLibrary) BuildSelectOptionsView(disable map[int64]bool, enable map[int64]bool, selected int64) string {
	var out string
	const dis string = ` disabled="disabled"`
	const sel string = ` selected="selected"`
	for _, v := range this.TreeList {
		out += `<option`
		if disable != nil {
			if _, ok := disable[v.GetId()]; ok {
				out += dis
			}}
		if enable != nil {
			if _, ok := enable[v.GetId()]; !ok {
				out += dis
			}}
		if v.GetId() == selected {
			out += sel
		}
		out += ` value="` + fmt.Sprintf("%d", v.GetId()) + `">` + strings.Repeat("&rarr;", v.GetLevel()) + v.GetTitle() + `</option>`
	}
	return out
}
