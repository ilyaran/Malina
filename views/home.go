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
 */ package views

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