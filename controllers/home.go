package controllers

import (
	"net/http"
	"github.com/ilyaran/Malina/views"
)

var HomeController *Home
func HomeControllerInit(){
	HomeController = &Home{}

}

type Home struct {

}

func(s *Home)Index(w http.ResponseWriter, r *http.Request){

	views.HomeView.Index(w)
}



