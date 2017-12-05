/**
 * @author		John Aran (Ilyas Aranzhanovich Toxanbayev)
 * @version		1.0.0
 * @based on
 * @email      		il.aranov@gmail.com
 * @link
 * @github      	https://github.com/ilyaran/github.com/ilyaran/Malina
 * @license		MIT License Copyright (c) 2017 John Aran (Ilyas Toxanbayev)
 */ 
 
 package controllers

import (
	"net/http"
	"github.com/ilyaran/Malina/views/publicView"
	"github.com/ilyaran/Malina/berry"
)

var CabinetController = &Cabinet{view:publicView.Cabinet{}}
type Cabinet struct {

	view publicView.Cabinet

}
func(s *Cabinet)Index(malina *berry.Malina,w http.ResponseWriter, r *http.Request){

	s.view.Index(malina,w)
}











