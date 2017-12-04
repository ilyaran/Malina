package libraries




import (
	"github.com/ilyaran/Malina/entity"
	"encoding/json"
	"strconv"
	"github.com/ilyaran/Malina/app"
)

var CategoryLib *Category = &Category{
	Trees: Tree{},
}


type Category struct {
	Trees              Tree
	PublicJsonTreeList []byte
	MapPublic          map[int64]*entity.Category
	ListPublic         []*entity.Category
	PublicListView string
}
func(s *Category)SetTrees(heapList *[]entity.Scanable){
	s.Trees.HeapList = []entity.Hierarchical{}
	for _,v:=range *heapList{
		s.Trees.HeapList = append(s.Trees.HeapList,v.(*entity.Category))
	}
	s.Trees.SetTrees()
}

func(s *Category)SetPublicTrees(){
	s.ListPublic = []*entity.Category{}
	s.MapPublic = map[int64]*entity.Category{}
	for _,v := range s.Trees.List {
		if v.(*entity.Category).GetEnable() {
			v.(*entity.Category).SetPublicParametersListView(&ParameterLib.ListPublic)
			s.ListPublic = append(s.ListPublic, v.(*entity.Category))
			s.MapPublic[v.(*entity.Category).GetId()] = v.(*entity.Category)
		}
	}
	s.PublicJsonTreeList, _ = json.Marshal(s.ListPublic)

	s.SetPublicListView()
	s.SetPublicBreadcrumbsView()
}


func (s *Category) SetPublicBreadcrumbsView() {
	for _, v := range s.ListPublic {
		v.SetPublicBreadcrumbPathView()
	}
}

func (this *Category) SetPublicListView() string {
	this.PublicListView = ``
	for _, v := range this.ListPublic {
		if v.GetEnable() && v.GetParent() == 0 {
			if len(v.Children) == 0 {
				this.PublicListView += `
				<li class="subitem1"><span><input class="categoryRadio" type="radio" data-item_id="`+strconv.FormatInt(v.GetId(),10) + `" name="group1"/><a href="` + app.Base_url+`public/product/list/?category_id=`+strconv.FormatInt(v.GetId(),10) + `">` + v.GetTitle() + `</a></span></li>
		`
			} else {
				this.PublicListView += `
				<li class="item3"><span><input class="categoryRadio" type="radio" data-item_id="`+strconv.FormatInt(v.GetId(),10) + `" name="group1"/><a href="` + app.Base_url+"public/product/list/?category_id="+strconv.FormatInt(v.GetId(),10) + `">` + v.GetTitle() + `</a><button class="arrow-btn btn btn-info">V</button></span>
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

func (this *Category) getChild(id int64) {

	for _, v := range this.ListPublic {
		if v.GetEnable() && v.GetParent() == id {
			if len(v.Children) == 0 {
				this.PublicListView += `
				<li class="subitem1"><span><input class="categoryRadio" type="radio" data-item_id="`+strconv.FormatInt(v.GetId(),10) + `" name="group1"/><a href="` + app.Base_url+"public/product/list/?category_id="+strconv.FormatInt(v.GetId(),10) + `">` + v.GetTitle() + `</a></span></li>`
			} else {
				this.PublicListView += `
				<li class="item3"><span><input class="categoryRadio" type="radio" data-item_id="`+strconv.FormatInt(v.GetId(),10) + `" name="group1"/><a href="` + app.Base_url+"public/product/list/?category_id="+strconv.FormatInt(v.GetId(),10) + `">` + v.GetTitle() + `</a><button class="arrow-btn btn btn-info">V</button></span>
					<ul>`
				this.getChild(v.GetId())
				this.PublicListView += `
					</ul>
				</li>`
			}
		}
	}
}



























