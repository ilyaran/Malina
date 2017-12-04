package controllers

import (
	"net/http"
	"github.com/ilyaran/Malina/views/publicView"
	"github.com/ilyaran/Malina/berry"
)


var WelcomeController = &Welcome{view:publicView.Welcome{}}
type Welcome struct {

	view publicView.Welcome

}

func(s *Welcome)Index(malina *berry.Malina,w http.ResponseWriter, r *http.Request){

	s.view.Index(malina,w)
}
