package views

import "net/http"

var HomeView * Home
func HomeViewInit(){
	HomeView = &Home{}

}


type Home struct {

}

func (s *Home)Index(w http.ResponseWriter){
	w.Write([]byte(HomeLayoutView.Layout))
}