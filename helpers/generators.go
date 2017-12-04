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
 */ package helpers

import (
	"reflect"
	"strings"
	"github.com/ilyaran/Malina/entity"
	"fmt"
	"io/ioutil"
	"regexp"
	"github.com/ilyaran/Malina/app"
	"os"
)

func GenerateControllerFields(dbtable string, item entity.Scanable){

	// get fields
	var sqlFieldNames string
	var fieldNames string
	var val reflect.Value
	val = reflect.ValueOf(item).Elem()

	for i := 0; i < val.NumField(); i++ {
		//valueField := val.Field(i)
		typeField := val.Type().Field(i)
		//tag := typeField.Tag
		fieldNames += "&s."+typeField.Name+",\n"
		sqlFieldNames += strings.ToLower(dbtable+"_"+typeField.Name)+",\n"
		//fmt.Printf("Name: %s,\t Value: %v,\t Tag: %s\n", typeField.Name, valueField.Interface(), tag.Get("tag_name"))
	}
	fmt.Println(sqlFieldNames)
	fmt.Println(fieldNames)
	// end get fields

}

func ModuleGenerator(newModuleName string, isHierarchical bool){
	var listOfPaths = []string{
		"/entity/",
		"/controllers/",
		"/views/",
		"/models/",
	}

	if isHierarchical{

	}else{
		for _,v:=range listOfPaths{
			var filename = app.Root_path+v+"product.go"
			b, err := ioutil.ReadFile(filename) // just pass the file name
			if err != nil {
				panic(err)
				return
			}
			str := string(b) // convert content to a 'string'

			// replace by pattern
			str = regexp.MustCompile(`(Product)`).ReplaceAllString(str, newModuleName)
			str = regexp.MustCompile(`(product)`).ReplaceAllString(str, strings.ToLower(newModuleName))

			if _, err := os.Stat(app.Root_path+v+strings.ToLower(newModuleName)+".go"); err == nil {
				fmt.Println("file exists allready: ",app.Root_path+v+strings.ToLower(newModuleName)+".go")
				return
			}

			d1 := []byte(str)
			err = ioutil.WriteFile(app.Root_path+v+strings.ToLower(newModuleName)+".go", d1, 0644)
			if err != nil{
				panic(err)
				return
			}
		}

	}


}