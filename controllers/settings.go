package controllers

import (
	"net/http"
	"github.com/ilyaran/Malina/views"
	"go/importer"
	"fmt"

)

var SettingsController *Settings
func SettingsControllerInit(){
	SettingsController = &Settings{}

	// view init
	SettingsController.view = views.Settings{}
	// end view init
}

type Settings struct {

	view            views.Settings
}

func(s *Settings)Index(w http.ResponseWriter, r *http.Request){

	pkg, err := importer.Default().Import("github.com/ilyaran/Malina/app")
	if err != nil {
		fmt.Printf("error: %s\n", err.Error())
		//return
	}
	for _, declName := range pkg.Scope().Names() {
		fmt.Println(declName)
	}

	s.view.Index(w)
}

















