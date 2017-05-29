/**
 * Upload library.  Malina eCommerce application
 *
 *
 * @author		John Aran (Ilyas Toxanbayev)
 * @version		1.0.0
 * @based on
 * @email      		il.aranov@gmail.com
 * @link
 * @github      	https://github.com/ilyaran/Malina
 * @license		MIT License Copyright (c) 2017 John Aran (Ilyas Toxanbayev)
 */
package library

import (
	"encoding/base64"
	"os"
	"strings"
	"github.com/ilyaran/Malina/language"
	"fmt"
	"time"
	"image"
	"github.com/ilyaran/Malina/config"
	_ "image/jpeg"
	_ "image/png"
	_ "image/gif"

)

var UPLOAD *Upload
type Upload struct {}

func (s *Upload)PrepareToUploadPublicAjax(indexNumberOfImage int,img string) (string, bool, int) {

	if img == ""{
		return "no_post_data", false, indexNumberOfImage // it means no image upload
	}

	d1 := strings.SplitN(img, ",", 2)
	if !strings.Contains(d1[0], ";base64") {
		return lang.T("upload_no_base64"), false, indexNumberOfImage
	}

	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(d1[1]))
	m, formatType, err := image.Decode(reader)
	if err != nil {
		return lang.T("upload no image"),false, indexNumberOfImage
	}
	//check file type
	if !(formatType == "jpg" || formatType == "jpeg" || formatType == "png" || formatType == "gif"){
		return lang.T("upload invalid filetype"), false, indexNumberOfImage
	}

	//check image dimensions
	bounds := m.Bounds()
	fmt.Printf("w %v == h %v\n",bounds.Dx(),bounds.Dy())
	if bounds.Dx() > app.Image_max_width() || bounds.Dy() > app.Image_max_height(){
		return lang.T("upload invalid dimensions"),false, indexNumberOfImage
	}

	//Creating new filename
	var fileName = fmt.Sprintf("%d%d.%s", time.Now().UTC().UnixNano(), indexNumberOfImage, formatType)

	return fileName, true, indexNumberOfImage
}


func (s *Upload)SaveUploadedFiles(img,fileName string,indexNumberOfImage int, path string) (string, bool, int) {

	d1 := strings.SplitN(img, ",", 2)

	d2, err := base64.StdEncoding.DecodeString(d1[1])
	if err != nil {
		fmt.Println(err)
		panic(err)
		return lang.T("upload error base64 decoding"), false, indexNumberOfImage
	}

	f, err := os.Create(path + fileName)
	if err != nil {
		fmt.Println(err)
		panic(err)

		return lang.T("upload destination error"), false, indexNumberOfImage
	}
	defer f.Close()
	//***************

	_, err = f.Write(d2)
	if err != nil {
		fmt.Println(err)
		panic(err)
		return lang.T("upload_unable_to_write_file"), false, indexNumberOfImage
	}

	//done <- fileName
	//close(done)
	return fileName,true,indexNumberOfImage
}

func (s *Upload)DelFile(path,filename string)bool{
	var pathToFile string = path + filename
	if _, err := os.Stat(pathToFile); err == nil {
		err = os.Remove(pathToFile)
		if err != nil{
			return false
		}
		return true
	}
	return false
}