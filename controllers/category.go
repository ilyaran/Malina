package controllers




import (
	"net/http"
	"github.com/ilyaran/Malina/views"
	"github.com/ilyaran/Malina/models"
	"github.com/ilyaran/Malina/app"
	"github.com/ilyaran/Malina/entity"
	"github.com/ilyaran/Malina/lang"
	"strconv"
	"github.com/ilyaran/Malina/libraries"
	"encoding/json"

	"github.com/ilyaran/Malina/berry"
	"html/template"
	"github.com/ilyaran/Malina/filters"
)

var CategoryController *Category
func CategoryControllerInit(){
	CategoryController = &Category{base:&base{
		dbtable:"category",
		item:&entity.Category{},
		itemHierarchical:&entity.Category{},
		selectSqlFieldsDefault 	:
		`
			category_id,
			coalesce(category_parent,0),
			category_sort,
			category_title,
			coalesce(category_description,''),
			category_lang,
			category_enable,
			category_img,
			category_parameter,
			category_quantity `,
	}}

	//helpers.GenerateControllerFields(CategoryController.base.dbtable,CategoryController.base.item)

	// model init
	CategoryController.model = models.Category{}
	// end model init

	// orders init
	CategoryController.base.orderList =[][2]string{
		[2]string{"ORDER BY category_sort ASC",lang.T("sort")+`&uarr;`},
		[2]string{"ORDER BY category_sort DESC",lang.T("sort")+`&darr;`},
		[2]string{"ORDER BY category_title DESC",lang.T("title")+`&darr;`},
		[2]string{"ORDER BY category_title ASC",lang.T("title")+`&uarr;`},
		[2]string{"ORDER BY category_id ASC",lang.T("id")+`&uarr;`},
		[2]string{"ORDER BY category_id DESC",lang.T("id")+`&darr;`},
	}

	if CategoryController.base.orderList==nil || len(CategoryController.base.orderList) < 0{
		panic("order list is not init")
	}
	CategoryController.base.orderListLength = int64(len(CategoryController.base.orderList))
	// end orders init

	CategoryController.base.searchSqlTemplate=`category_title LIKE '%~%' OR category_description LIKE '%~%' `

	CategoryController.base.inlistSqlFields = map[string]byte{
		"category_title"     :'s',
		"category_description"     :'s',
		"category_sort"    :'i',
		"category_enable"    :'b',
	}

	// view init
	CategoryController.view = views.Category{}
	for k,_ := range CategoryController.base.inlistSqlFields{
		CategoryController.view.InlistFields += `<input class="inlist_fields" type="hidden" value="`+k+`">`
	}
	for k,v := range CategoryController.base.orderList{
		CategoryController.view.OrderSelectOptions += `<option value="`+strconv.Itoa(k)+`">`+v[1]+`</option>`
	}
	// end view init

	//finally set trees
	tempFlow := &berry.Malina{}
	tempFlow.TableSql = "category"

	CategoryController.base.setTree(tempFlow,libraries.CategoryLib)
}

type Category struct {
	base     *base
	model			models.Category
	view            views.Category
}
func(s *Category)Index(malina *berry.Malina, w http.ResponseWriter, r *http.Request){
	malina.Controller = s
	malina.TableSql = s.base.dbtable
	s.base.index(malina,w,r,"after delete")
}

func(s *Category)GetList(malina *berry.Malina,w http.ResponseWriter, r *http.Request, condition string){

	//if public request
	if malina.Department == "public"{
		// public browser
		if malina.Device == "browser" {

			// public device (mobile, thing etc.)
		}else{
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.Write(libraries.CategoryLib.PublicJsonTreeList)
		}
		return
	}
	//end if public request


	s.base.listingPrepare(malina, w, r)

	if malina.Status > 0 {
		return
	}

	//if request contains search arg
	if malina.Search != ""{

		s.base.getList(malina, w,r)
		// if via ajax
		if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {

			//helpers.SetAjaxHeader(w)

			malina.Paging = s.base.paging(malina,"","")
			w.Write([]byte(s.view.ListingNonTree(malina)))

			// if non ajax req
		}else {
			malina.Paging = s.base.paging(malina,"", app.Url_home_category_list+"?"+malina.Url[1:])
			s.view.Index(malina,false,w)
		}
		return
	}
	//end if request contains search arg

	malina.ChannelBool = make(chan bool)
	go models.CrudGeneral.Count(malina,"")

	//set tree
	malina.TableSql = "category"
	s.base.setTree(malina,libraries.CategoryLib)
	//end set tree

	// sync count items
	<- malina.ChannelBool

	if malina.Device == "browser" {
		// if via ajax
		if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
			malina.Paging = s.base.paging(malina,"","")
		}else {
			malina.Paging = s.base.paging(malina,"", app.Url_home_category_list+"?"+malina.Url[1:])
		}
	}

	// slice limitations count
	var treeLen = int64(len(libraries.CategoryLib.Trees.List))
	if malina.Page >= treeLen{
		malina.Page=0
	}

	if malina.Page + malina.Per_page < treeLen{
		malina.Per_page = malina.Page + malina.Per_page
	}else {
		malina.Per_page = treeLen
	}
	// end slice limitations count

	// if the user agent is a browser
	if malina.Device == "browser" {
		// if via ajax
		if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
			//helpers.SetAjaxHeader(w)
			w.Write([]byte(s.view.Listing(malina)))
			// if non ajax req
		}else {
			s.view.Index(malina,true,w)
		}

		// if the user agent is a mobile or another device
	}else {
		var out []byte
		if malina.Department == "public"{
			out = libraries.CategoryLib.PublicJsonTreeList
		}else {
			out, _ = json.Marshal(libraries.CategoryLib.Trees.List)
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.Write(out)
	}
}

func (s *Category) ResultHandler(malina *berry.Malina, args ... interface{}){
	//set tree
	s.base.setTree(malina,libraries.CategoryLib)
	malina.Result["select_options"] = views.HomeLayoutView.Form_category_select(lang.T("root"),"parent")
}

func (s *Category) FormHandler(malina *berry.Malina,action byte,w http.ResponseWriter,r *http.Request) {
	if action=='e'{
		malina.IdInt64, malina.IdStr = filters.IsUint(malina,true, "id", 20, r)
		if malina.IdInt64 > 0 && malina.Status == 0{
			if _,ok:=libraries.CategoryLib.Trees.Map[malina.IdInt64];!ok{
				malina.Status = http.StatusNotAcceptable
				malina.Result["id"]=lang.T("not found")
			}
		}
	}
	imgUrls, _, _ := filters.ImgUrls(malina,false, "img", r)

	sort,_ := filters.IsUint(malina,false, "sort", 20, r)

	parent, _ := filters.IsUint(malina,true, "parent", 20, r)
	if parent > 0 && malina.Status == 0{
		if _,ok:=libraries.CategoryLib.Trees.Map[parent];!ok{
			malina.Status = http.StatusNotAcceptable
			malina.Result["parent"]=lang.T("not found")
		}
	}

	enable := filters.CheckBox(malina,"enable",r)

	title := template.HTMLEscapeString(r.FormValue("title"))
	if title==""{
		malina.Status = http.StatusNotAcceptable
		malina.Result["title"]=lang.T("required")
	}else{
		if len(title) > app.Title_max_len{
			malina.Status = http.StatusNotAcceptable
			malina.Result["title"]=lang.T("length must have no more then:")+strconv.Itoa(app.Title_max_len)+" symbols"
		}
	}

	description := r.FormValue("description")
	short_description := template.HTMLEscapeString(r.FormValue("short_description"))
	if len(short_description) > app.Short_description_max_len{
		malina.Status = http.StatusNotAcceptable
		malina.Result["short_description"]=lang.T("length must have no more then:")+strconv.Itoa(app.Short_description_max_len)+" symbols"
	}

	if malina.Status > 0 {
		return
	}
	var res int64
	if action == 'a' {
		res = s.model.Upsert(action,imgUrls,parent,title,description,sort,enable)
	} else {
		res = s.model.Upsert(action,imgUrls,parent,title,description,sort,enable,malina.IdInt64)
	}
	if res > 0 {
		models.CrudGeneral.WhereAnd(malina,s.base.dbtable+"_id =","$1")
		malina.SelectSql=s.base.selectSqlFieldsDefault
		if action == 'a'{
			models.CrudGeneral.GetItem(malina,s.base.item,"",res)
		} else {
			models.CrudGeneral.GetItem(malina,s.base.item,"",malina.IdInt64)
		}

		malina.Status = http.StatusOK

		//set tree
		s.base.setTree(malina,libraries.CategoryLib)

		malina.Result["item"] = malina.Item
		malina.Result["select_options"] = views.HomeLayoutView.Form_category_select(lang.T("root"),"parent")

	} else {
		malina.Status = http.StatusInternalServerError
		malina.Result["error"] = lang.T(`server error`)
	}

}


func(s *Category)Default(malina *berry.Malina, w http.ResponseWriter, r *http.Request){
	if malina.Action==""{
		s.GetList(malina,w,r,"")
	}
}




