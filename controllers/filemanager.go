package controllers

import (
	"net/http"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"image"
	"io"
	"bytes"
	"image/jpeg"
	"github.com/disintegration/imaging"
	"strconv"
	"image/gif"
	"image/png"
	"path"
	"github.com/ilyaran/Malina/app"
	"github.com/gorilla/mux"
)

var Filemanager = &filemanager{crud:&base{dbtable:"product"}}
type filemanager  struct {
	crud *base
	out string
}

func (this *filemanager) Index(w http.ResponseWriter, r *http.Request) {

	//dao.AuthDao.Authentication(w,r)

	this.out = this.resultError("")

	//if dao.AuthDao.Session.Account.Role == app.Admin_role_id {
		switch mux.Vars(r)["action"] {
		case "dirtree"      : this.dirtree(w,r)
		case "createdir"    : this.createdir(w,r)
		case "deletedir"    : this.deletedir(w,r)
		case "movedir" 	    : this.movedir(w,r)
		case "copydir" 	    : this.copydir(w,r)
		case "renamedir"    : this.renamedir(w,r)
		case "fileslist"    : this.fileslist(w,r)
		case "upload" 	    : this.upload(w,r)
		case "download"     : this.download(w,r)
		case "downloaddir"  : this.downloaddir(w,r)
		case "deletefile"   : this.deletefile(w,r)
		case "movefile"     : this.movefile(w,r)
		case "copyfile"     : this.copyfile(w,r)
		case "renamefile"   : this.renamefile(w,r)
		case "thumb" 	    : this.thumb(w,r);return
		}
	//}else {
		//this.out = this.resultError(lang.T("unauth access"))
	//}
	fmt.Fprintf(w, this.out)
}

func (this *filemanager) DirRead(dir string)string{
	f,d := 0,0
	files, err := ioutil.ReadDir("."+dir)
	if err != nil {
		f,d = 0,0
	}

	for _, file := range files {
		if !file.IsDir() {
			f++
		}else {
			d++
			this.DirRead(dir+"/"+file.Name())
		}
	}
	this.out = fmt.Sprintf(`,{"p":"`+dir+`","f":"%d","d":"%d"}`, f,d) + this.out
	return this.out
}
func (this *filemanager) dirtree(w http.ResponseWriter, r *http.Request) {
	//"/assets/uploads" ///r.FormValue("d")
	dir := "/"+app.Path_assets_uploads[:len(app.Path_assets_uploads)-1]
	this.out = ``
	this.DirRead(dir)
	this.out = "["+this.out[1:]+"]"
}
func (this *filemanager) createdir(w http.ResponseWriter, r *http.Request) {
	/*
	d:/assets/Uploads/Documents/gghh/jjhhgg
	n:ytytryt
	 */
	d := r.FormValue("d")
	n := r.FormValue("n")
	path_dir := "."+d+"/"+n
	this.out = `{res: "error", msg: ""}`
	_, err := os.Stat(path_dir)
	if os.IsNotExist(err) {
		os.Mkdir(path_dir, 0777)
		this.out = this.resultSuccess("")
	}else {
		this.out = this.resultError("allready exists")
	}
}
func (this *filemanager) deletedir(w http.ResponseWriter, r *http.Request) {
	/*
	d:/assets/Uploads/Documents/new_Name/temp_hfhf
	 */
	d := r.FormValue("d");
	if path.IsAbs(d) {
		dir1 := "." + d
		_, err := os.Stat(dir1)
		if os.IsNotExist(err) {
			this.out = this.resultError("not found")
		}else {
			os.Remove(dir1)
			this.out = this.resultSuccess("")
		}
	}else {
		this.out = this.resultError("invalid")
	}
}
func (this *filemanager) movedir(w http.ResponseWriter, r *http.Request) {
	d,n := r.FormValue("d"),r.FormValue("n")
	dir1,dir2 := "."+d , "."+n
	if path.IsAbs(d) && path.IsAbs(n) {
		_, err := os.Stat(dir1)
		if os.IsNotExist(err) {
			this.out = this.resultError(d + " is not exists")
		} else {
			_, err = os.Stat(dir2)
			if os.IsNotExist(err) {
				this.out = this.resultError(n + " is not exists")
			} else {
				//os.Mkdir(dir2 + "/" + path.Base(dir1), 0777)
				os.Rename(dir1,dir2 + "/" + path.Base(dir1))
				this.out = this.resultSuccess("")
			}
		}
	}
}
func (this *filemanager) copydir(w http.ResponseWriter, r *http.Request) {
	/*
	d:/assets/Uploads/Images
	n:/assets/Uploads/oiiiiii
	 */
	d,n := r.FormValue("d"),r.FormValue("n")
	dir1,dir2 := "."+d , "."+n
	if path.IsAbs(d) && path.IsAbs(n) {
		_, err := os.Stat(dir1)
		if os.IsNotExist(err) {
			this.out = this.resultError(d + " is not exists")
		} else {
			_, err = os.Stat(dir2)
			if os.IsNotExist(err) {
				this.out = this.resultError(n + " is not exists")
			} else {
				var newpath = dir2 + "/" + path.Base(dir1)
				var findName func(int)
				findName = func(n int) {
					var pref = ""
					if n > 0 {
						pref = "_"+strconv.Itoa(n)
					}
					_, err1 := os.Stat(newpath)
					if os.IsNotExist(err1) {
						newpath = newpath + pref
					} else {
						n++
						findName(n)
					}
				}
				findName(0)
				os.Mkdir(newpath, 0777)
				this.out = this.resultSuccess("")
			}
		}
	}
}
func (this *filemanager) renamedir(w http.ResponseWriter, r *http.Request) {
	d,n := r.FormValue("d"),r.FormValue("n")
	if path.IsAbs(d) {
		dir1 := "."+d;
		_, err := os.Stat(dir1)
		if os.IsNotExist(err) {
			this.out = `{res: "error", msg: "`+d+` is not exists"}`
		}else {
			if path.IsAbs(path.Dir(d)+"/"+n){
				dir2 := "." + path.Dir(d) + "/" + n;
				_, err = os.Stat(dir2)
				if os.IsNotExist(err) {
					os.Rename(dir1,dir2)
					this.out = this.resultSuccess("")
				}else {
					this.out = this.resultError("allready exists")
				}
			}
		}
	}
}
func (this *filemanager) fileslist(w http.ResponseWriter, r *http.Request) {
	d := r.FormValue("d")
	if path.IsAbs(d){
		dir := "."+d
		_, err := os.Stat(dir)
		if os.IsNotExist(err) {
			this.out = this.resultError(dir+" is not exists directory")
			return
		}else {
			files, err := ioutil.ReadDir(dir)
			if err != nil {
				log.Fatal(err)
				return
			}
			if len(files)>0{
				this.out = ``
				for _, file := range files {
					if !file.IsDir() {
						if reader, err := os.Open(filepath.Join(dir, file.Name())); err == nil {
							defer reader.Close()
							im, _, err := image.DecodeConfig(reader)
							if err != nil {
								this.out += fmt.Sprintf(`,{"p":"%v","s":"%v","t":"%v","w":"%v","h":"%v"}`,d+"/"+file.Name(),file.Size(),file.ModTime().Unix(),0, 0)
							}else {
								this.out += fmt.Sprintf(`,{"p":"%v","s":"%v","t":"%v","w":"%v","h":"%v"}`,d+"/"+file.Name(),file.Size(),file.ModTime().Unix(),im.Width, im.Height)
							}
						}
					}
				}
				if len(this.out) > 0{
					this.out =  "[" + this.out[1:] + "]"
					return
				}
			}
		}
	}
	this.out = `[]`
}

func (this *filemanager) upload(w http.ResponseWriter, r *http.Request) {
	d := r.FormValue("d");
	if path.IsAbs(d){
		d = "." + d
		_, err := os.Stat(d)
		if os.IsNotExist(err) {
			this.out = this.resultError(d+" is not exists directory")
			return
		}
	}else {
		this.out = this.resultError("invalid directory parameter")
		return
	}
	err := r.ParseMultipartForm(100000)
	if err != nil {
		this.out = this.resultError("")
		return
	}
	//get a ref to the parsed multipart form
	m := r.MultipartForm

	//get the *fileheaders
	files := m.File["files[]"]
	for i, _ := range files {
		//for each fileheader, get a handle to the actual file
		file, err := files[i].Open()
		defer file.Close()
		if err != nil {
			this.out = this.resultError("file open error")
			return
		}
		//create destination file making sure the path is writeable.
		var filePath string
		var findName func(int)
		findName = func(n int) {
			var pref = ""
			if n > 0 {
				pref = "("+strconv.Itoa(n)+")"
			}
			_, err1 := os.Stat(d + "/"+ pref + files[i].Filename)
			if os.IsNotExist(err1) {
				filePath = d + "/" + pref + files[i].Filename
			} else {
				n++
				findName(n)
			}
		}
		findName(0)

		dst, err := os.Create(filePath)
		defer dst.Close()
		if err != nil {
			//http.Error(w, err.Error(), http.StatusInternalServerError)
			this.out = this.resultError("")
			return
		}
		//copy the uploaded file to the destination file
		if _, err := io.Copy(dst, file); err != nil {
			//http.Error(w, err.Error(), http.StatusInternalServerError)
			this.out = this.resultError("")
			return
		}
	}
	this.out = this.resultSuccess("")
}

func (this *filemanager) download(w http.ResponseWriter, r *http.Request) {

}
func (this *filemanager) downloaddir(w http.ResponseWriter, r *http.Request) {

}
func (this *filemanager) deletefile(w http.ResponseWriter, r *http.Request) {
	f := r.FormValue("f");
	if path.IsAbs(f) {
		dir1 := "." + f
		_, err := os.Stat(dir1)
		if os.IsNotExist(err) {
			this.out = this.resultError("not found")
		}else {
			os.Remove(dir1)
			this.out = this.resultSuccess("")
			var del_q = `
			WITH t AS (
				UPDATE category SET category_img = array_remove(category_img,'`+f+`')
				WHERE '`+f+`' = ANY(category_img)
			)
			UPDATE product SET product_img = array_remove(product_img,'`+f+`')
			WHERE '`+f+`' = ANY(product_img)`
			ProductController.model.Upsert('e',"",del_q)
		}
	}else {
		this.out = this.resultError("invalid")
	}
}

func (this *filemanager) movefile(w http.ResponseWriter, r *http.Request) {
	f,n := r.FormValue("f"), r.FormValue("n")
	if path.IsAbs(f) && path.IsAbs(n){
		dir1,dir2 := "."+f, "."+n
		_, err := os.Stat(dir1)
		if os.IsNotExist(err) {
			this.out = this.resultError(f+" is not exists directory")
		}else {
			_, err = os.Stat(dir2)
			if os.IsNotExist(err) {
				err =  os.Rename(dir1, dir2)
				if err != nil {
					panic(err)
					this.out = this.resultError("")
				}else {
					this.out = this.resultSuccess("")
					return
				}
			}else {
				this.out = this.resultError(n+" is allready exists file name")
			}
		}
	}
	this.out = this.resultError("")
}
func (this *filemanager) copyfile(w http.ResponseWriter, r *http.Request) {
	/*
	f:/assets/Uploads/oiiiiii/Images/DSC_2987.jpg
	n:/assets/Uploads/Images
	 */



}
func (this *filemanager) renamefile(w http.ResponseWriter, r *http.Request) {
	f,n := r.FormValue("f"),r.FormValue("n")
	if path.IsAbs(f) {
		dir1 := "."+ f;
		_, err := os.Stat(dir1)
		if os.IsNotExist(err) {
			this.out = `{res: "error", msg: "`+ f +` is not exists"}`
		}else {
			if path.IsAbs(path.Dir(f)+"/"+n){
				dir2 := "." + path.Dir(f) + "/" + n;
				_, err = os.Stat(dir2)
				if os.IsNotExist(err) {
					os.Rename(dir1,dir2)
					this.out = this.resultSuccess("")
				}else {
					this.out = this.resultError("allready exists")
				}
			}
		}
	}
}
func (this *filemanager) thumb(w http.ResponseWriter, r *http.Request) {
	f := "."+r.FormValue("f")
	_, err := os.Stat(f)
	if os.IsNotExist(err){
		fmt.Fprintf(w, this.resultError(f + " is not exists file"))
		return
	}
	widthStr,heightStr := r.FormValue("width"),r.FormValue("height")
	if (widthStr == "") { widthStr = "100";}
	if (heightStr == "") {heightStr = "0";}
	width,errW := strconv.Atoi(widthStr)
	if (errW!=nil){width=100}
	height, errH := strconv.Atoi(heightStr)
	if (errH!=nil){height=0}

	// get content type
	contentType := "application/octet-stream";
	file, err := os.Open(f)
	//fmt.Printf("%T",file)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Only the first 512 bytes are used to sniff the content type.
	buffer := make([]byte, 512)
	// function still work, if size of buffer < 512
	n, err1 := file.Read(buffer)
	if err1 != nil && err1 != io.EOF {
		panic(err1)
	}
	contentType = http.DetectContentType(buffer[:n])
	// end get content type

	img, err2 := imaging.Open(f)
	if err2 != nil {
		panic(err)
	}
	var thumb image.Image = imaging.Thumbnail(img, width, height, imaging.CatmullRom)

	buffer1 := new(bytes.Buffer)
	if contentType == "image/jpeg" {
		if err := jpeg.Encode(buffer1, thumb, nil); err != nil {
			log.Println("unable to encode image.")
		}
	}
	if contentType == "image/gif" {
		if err := gif.Encode(buffer1, thumb, nil); err != nil {
			log.Println("unable to encode image.")
		}
	}
	if contentType == "image/png"{
		if err := png.Encode(buffer1, thumb); err != nil {
			log.Println("unable to encode image.")
		}
	}

	w.Header().Set("Pragma", "cache")
	w.Header().Set("Cache-Control", "max-age=3600")
	w.Header().Set("Content-type", contentType)
	w.Header().Set("Content-Length", strconv.Itoa(len(buffer1.Bytes())))

	if _, err := w.Write(buffer1.Bytes()); err != nil {
		log.Println("unable to write image.")
	}
}

func (this *filemanager) resultSuccess(msg string)string{
	return `{"res" : "ok", "msg" : "`+msg+`"}`
}
func (this *filemanager) resultError(msg string)string{
	return `{"res" : "error", "msg" : "`+msg+`"}`
}











