package core

import (
	"github.com/ilyaran/Malina/libraries"
	"github.com/ilyaran/Malina/models"
	"github.com/ilyaran/Malina/config"
	"github.com/ilyaran/Malina/language"
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/ilyaran/Malina/entity"
	"github.com/ilyaran/Malina/views/publicView"
	"github.com/ilyaran/Malina/views/auth"
	"github.com/ilyaran/Malina/views/account"
)

var MALINA *Malina

type IMalina interface {

}

type Malina struct {

}
//var Repeat int
func (s *Malina) DB_init()  bool {
	db, err := sql.Open("postgres", "postgres://" + app.DB_USER() + ":" + app.DB_PASSWORD() + "@" + app.HOST + "/" + app.DB_NAME() + "?sslmode=disable")
	//db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))

	if err != nil {
		model.Crud = nil
		return false
	}
	if db.Ping() != nil {
		model.Crud = nil
		return false
	}
	model.Crud = &model.CrudPostgres{db}
	if model.Crud != nil {
		return true
	}
	return false
}

func (s *Malina) DB_close(){
	model.Crud.DB.Close()
}


func (s *Malina) LibrariesInit(){
	library.VALIDATION = &library.Validation{}
	library.SESSION = &library.Session{
		SessionObj:&entity.Session{},
		StdChars:[]byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789!@#$%^&*()-_=+,.?/:;{}[]`~"),
	}

	library.POSITION = &library.PositionLibrary{}
	library.CATEGORY = &library.CategoryLibrary{Id:-1,Parent:-1}

	// set position globals
	s.SetPositionGlobals()

	// set category globals
	s.SetCategoryGlobals()

	//set global language list
	for _,v := range app.Language_list(){
		if i,ok := lang.FullLanguageList[v];ok{
			lang.LangSelectOptions += `<option value="` + v + `" class="in-of">` + i + `</option>`
		}
	}

}

func (s *Malina) PublicViewObjectsInit(){

	authView.AuthViewObjInit()

	accountView.AccountViewObjInit()

	publicView.WelcomeViewObjInit()

	publicView.ProductViewObjInit()

	publicView.CartViewObjInit()


}

func (s *Malina) SetPositionGlobals(){
	//set globals for category
	///model.PositionModel.Id = int64(-1)
	var positionList = model.PositionModel.GetList("","","","position_sort ASC",library.POSITION.TreeMap)
	library.POSITION.SetTrees(positionList)
	/*fmt.Println("================")
	for k,v:=range library.POSITION.TreeMap{
		fmt.Println(k,v)
	}*/

	var tempMap = map[int64]*entity.Position{}
	for k,v := range library.POSITION.TreeMap {
		tempMap[k] = v
	}
	positionList = model.PositionModel.GetList("","","","position_sort ASC",tempMap)
	library.POSITION.SetTrees(positionList)
	//fmt.Println("~~~~~~~~~~~~~~~")
	/*for k,v := range library.POSITION.TreeMap {
		fmt.Println(k,v)
	}*/

	//fmt.Println(library.POSITION.TreeMap)


}
func (s *Malina) SetCategoryGlobals(){
	//set globals for category
	var categoryList = model.CategoryModel.GetList("","","","category_sort ASC")
	library.CATEGORY.SetTrees(categoryList)

	library.CATEGORY.BuildSelectOptionsView(nil,nil,0)
}

