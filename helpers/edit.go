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
	"io/ioutil"
	"regexp"
	"fmt"
	"github.com/ilyaran/Malina/app"
)


var pattern, substr string
// example pattern = `(app\.[a-zA-Z_]+)`
// substr = `$1()`
func GlobalReplacementByRegexp(p, s string){
	pattern, substr = p, s
	//start
	readAllProjectFiles("")
}

func readAllProjectFiles(dir string){
	files, err := ioutil.ReadDir(app.Root_path+dir)
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		if file.IsDir() {
			readAllProjectFiles(dir + "/" + file.Name())
		} else {
			// get only .go files
			match, _ := regexp.MatchString(`(\.go)$`, file.Name())
			if match {
				fmt.Println(dir + "/" + file.Name())

				// here example where all app.Base_url replaced to app.Base_url() by regular expression
				// you can use your own pattern and substring you need
				readReplaceAndWriteFile("./"+dir + "/" + file.Name(),pattern,substr)
			}
		}
	}
}

func readReplaceAndWriteFile(filename, pattern, substr string) {
	b, err := ioutil.ReadFile(filename) // just pass the file name
	if err != nil {
		panic(err)
	}
	str := string(b) // convert content to a 'string'

	// replace by pattern
	str = regexp.MustCompile(pattern).ReplaceAllString(str, substr)

	d1 := []byte(str)
	err = ioutil.WriteFile(filename, d1, 0644)
	if err != nil{
		panic(err)
	}
}
