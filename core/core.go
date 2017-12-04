package core

import (
	"runtime"
	"github.com/ilyaran/Malina/app"
	"github.com/ilyaran/Malina/dao"
	"github.com/ilyaran/Malina/views"
	"github.com/ilyaran/Malina/controllers"
	"fmt"
	"github.com/ilyaran/Malina/views/publicView"
	"github.com/ilyaran/Malina/models"
)

func Init(){
	malina:=&Core{}
	malina.Start()
}

type Core struct {}


func(s *Core)Start() error {
	var err error

	// recognize a platform
	// set Root_path, Port_addr, Base_url
	app.PlatformOs = runtime.GOOS

	app.SetConfig(app.PlatformOs)

	fmt.Println("Platform: "+app.PlatformOs)
	// end recognize a platform

	// database connection for general globe crud
	models.CrudGeneral = &models.CrudPostgres{}
	err=models.CrudGeneral.SetDBConnection()

	// init controllers
	sqlSelectFieldsAccount:=controllers.AccountControllerInit()
	controllers.HomeControllerInit()
	controllers.ProductControllerInit()
	controllers.ParameterControllerInit()
	controllers.CategoryControllerInit()
	controllers.SettingsControllerInit()


	// init DAOs
	dao.AuthDaoInit(sqlSelectFieldsAccount)
	// end init DAOs

	s.SetViews()

	return err
}

func(s *Core) SetViews(){

	views.HomeLayoutViewInit()
	views.HomeViewInit()

	publicView.PublicLayoutViewInit()

	publicView.NewsViewInit()

}







