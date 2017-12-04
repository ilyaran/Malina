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
 */ package controllers

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



