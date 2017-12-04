package libraries

import (
	"net/http"
	"runtime"
	"github.com/ilyaran/Malina/lang"
	"github.com/ilyaran/Malina/app"
	"encoding/base64"
	"fmt"
	"time"
	"strings"
	"image"
	_ "image/jpeg"
	_ "image/png"
	_ "image/gif"
	"github.com/ilyaran/Malina/caching"
	"os"
	"github.com/ilyaran/Malina/berry"
)

var UploadLib = &Upload{}
type Upload struct {
	
}

func(s *Upload) Img(malina *berry.Malina,keyName string, w http.ResponseWriter, r *http.Request)(string,bool){
	//maximize CPU usage for maximum performance
	runtime.GOMAXPROCS(runtime.NumCPU())
	var err error
	if r.Method != "POST" {
		malina.Status = 406
		malina.Result["error"] = lang.T("no post request")
		return "",false
	}
	//check request size to prevent malicious attack or accidental large request.
	if r.ContentLength > app.Image_max_size {
		malina.Status = 406
		malina.Result["error"] = lang.T(`upload_file_exceeds_limit`)
		return "",false
	}
	r.Body = http.MaxBytesReader(w, r.Body, app.Image_max_size)
	if err = r.ParseForm(); err != nil {
		malina.Status = 500
		malina.Result["error"] = lang.T(`parse form error`)
		return "",false
	}
	filename, prepareError, _ := s.PrepareToUploadPublicAjax(0, r.FormValue(keyName))
	if !prepareError {
		malina.Status = 500
		malina.Result["error"] = "error of prepare file: " + filename
		return filename,false
	}else if filename == "null"{
		malina.Status = -10
		return "null",true
	}else if filename == ""{
		malina.Status = -20
		return "",true
	}
	filenameResult, saveResult, _ := s.SaveUploadedFiles(r.FormValue(keyName), filename, 0, app.Root_path+app.Path_assets_uploads)
	if !saveResult {
		malina.Status = 500
		malina.Result["error"] = "save error with file: " + filenameResult
		return  "",false
	}
	malina.Status = 0
	return  filenameResult,true
}

func (s *Upload)PrepareToUploadPublicAjax(indexNumberOfImage int,img string) (string, bool, int) {
	if img == "null"{
		return "null", true, indexNumberOfImage // it means set image to null
	}
	if img == ""{
		return "", true, indexNumberOfImage // it means no image upload
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
	//fmt.Printf("w %v == h %v\n",bounds.Dx(),bounds.Dy())
	if bounds.Dx() > app.Image_max_width || bounds.Dy() > app.Image_max_height{
		return lang.T("upload invalid dimensions"),false, indexNumberOfImage
	}

	//Creating new filename
	var fileName = fmt.Sprintf("%d%d%d.%s", time.Now().UTC().UnixNano(), time.Now().Sub(caching.T0).Nanoseconds(),indexNumberOfImage, formatType)
	return fileName, true, indexNumberOfImage
}

func (s *Upload)SaveUploadedFiles(img,fileName string,indexNumberOfImage int, path string) (string, bool, int) {
	if strings.IndexByte(img,',') < 0{
		return lang.T("upload error base64 decoding"), false, indexNumberOfImage
	}
	d1 := strings.SplitN(img, ",", 2)

	d2, err := base64.StdEncoding.DecodeString(d1[1])
	if err != nil {
		//fmt.Println(err)
		//panic(err)
		return lang.T("upload error base64 decoding"), false, indexNumberOfImage
	}

	f, err := os.Create(path + fileName)
	if err != nil {
		/*fmt.Println(path + fileName)
		fmt.Println(err)*/
		panic(err)

		return lang.T("upload destination error"), false, indexNumberOfImage
	}
	defer f.Close()
	//***************

	_, err = f.Write(d2)
	if err != nil {
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