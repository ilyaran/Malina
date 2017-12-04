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
 */ package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"log"
	"github.com/ilyaran/Malina/controllers"
	"github.com/ilyaran/Malina/app"
	"github.com/ilyaran/Malina/core"
)

func main() {
	core.Init()

	fmt.Println("Listening port "+app.Port_addr)

	//helpers.GlobalReplacementByRegexp()
	//helpers.ModuleGenerator("jhjhj",false)

	router := mux.NewRouter().StrictSlash(true)
	router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets", http.FileServer(http.Dir(app.Root_path+"assets/"))))
	router.HandleFunc("/", controllers.PublicDepartment).Methods("GET")
	router.HandleFunc("/home/{controller}/{action}/", controllers.HomeDepartment).Methods("GET")
	router.HandleFunc("/home/{controller}/", controllers.HomeDepartment).Methods("GET")
	router.HandleFunc("/home/{controller}/{action}/", controllers.HomeDepartment).Methods("POST")


	router.HandleFunc("/public/{controller}/{action}/", controllers.PublicDepartment).Methods("GET")
	router.HandleFunc("/public/{controller}/{action}/", controllers.PublicDepartment).Methods("POST")

	router.HandleFunc("/cabinet/{controller}/{action}/", controllers.CabinetDepartment).Methods("GET")
	router.HandleFunc("/cabinet/{controller}/{action}/", controllers.CabinetDepartment).Methods("POST")

	router.HandleFunc("/auth/{action}/", controllers.AuthDepartment).Methods("POST","GET")

	router.HandleFunc("/filemanager/{action:(?:dirtree|createdir|deletedir|movedir|copydir|renamedir|fileslist|upload|download|downloaddir|deletefile|movefile|copyfile|renamefile)}/", controllers.Filemanager.Index).Methods("POST")
	router.HandleFunc("/filemanager/{action:(?:thumb)}/", controllers.Filemanager.Index).Methods("GET")

	log.Fatal(http.ListenAndServe(app.Port_addr, router))
}




















