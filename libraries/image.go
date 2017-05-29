/**
 * Image library class.  Malina eCommerce application
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
	"github.com/disintegration/imaging"
	"os"
	"github.com/ilyaran/Malina/language"

)


var ImageLib *Image = &Image{}
type Image struct {}
func (this *Image)Imgresize(fromPath, fromFileName,destinationPath, destinationFileName string,width int,height int) (string,bool) {
	// maximize CPU usage for maximum performance
	//runtime.GOMAXPROCS(runtime.NumCPU())
	// load original image
	img, err := imaging.Open(fromPath + fromFileName)  //("./upload/advert/" + fileName)
	if err != nil {
		//fmt.Println(err)
		//os.Exit(1)
		return lang.T(`image_resize_error_open`),false
	}
	// resize image from  to 400 while preserving the aspect ration
	// Supported resize filters: NearestNeighbor, Box, Linear, Hermite, MitchellNetravali,
	// CatmullRom, BSpline, Gaussian, Lanczos, Hann, Hamming, Blackman, Bartlett, Welch, Cosine.
	dstimg := imaging.Resize(img, width, height, imaging.Box)
	// save resized image
	err = imaging.Save(dstimg, destinationPath + destinationFileName) //(dstimg, "./upload/advert/"+fileName)
	if err != nil {
		//fmt.Println(err)
		//os.Exit(1)
		return lang.T(`image_resize_error_save`),false
	}
	// everything ok
	return "ok",true
}

func(this *Image) GetFormat(file *os.File) string {
	bytes := make([]byte, 4)
	n, _ := file.ReadAt(bytes, 0)
	if n < 4 {
		return ""
	}
	if bytes[0] == 0x89 && bytes[1] == 0x50 && bytes[2] == 0x4E && bytes[3] == 0x47 {
		return "png"
	}
	if bytes[0] == 0xFF && bytes[1] == 0xD8 {
		return "jpg"
	}
	if bytes[0] == 0x47 && bytes[1] == 0x49 && bytes[2] == 0x46 && bytes[3] == 0x38 {
		return "gif"
	}
	if bytes[0] == 0x42 && bytes[1] == 0x4D {
		return "bmp"
	}
	return ""
}

func(this *Image) GetGifDimensions(file *os.File) (width int, height int) {
	bytes := make([]byte, 4)
	file.ReadAt(bytes, 6)
	width = int(bytes[0]) + int(bytes[1])*256
	height = int(bytes[2]) + int(bytes[3])*256
	return
}

func (this *Image)GetBmpDimensions(file *os.File) (width int, height int) {
	bytes := make([]byte, 8)
	file.ReadAt(bytes, 18)
	width = int(bytes[3])<<24 | int(bytes[2])<<16 | int(bytes[1])<<8 | int(bytes[0])
	height = int(bytes[7])<<24 | int(bytes[6])<<16 | int(bytes[5])<<8 | int(bytes[4])
	return
}

func (this *Image)GetPngDimensions(file *os.File) (width int, height int) {
	bytes := make([]byte, 8)
	file.ReadAt(bytes, 16)
	width = int(bytes[0])<<24 | int(bytes[1])<<16 | int(bytes[2])<<8 | int(bytes[3])
	height = int(bytes[4])<<24 | int(bytes[5])<<16 | int(bytes[6])<<8 | int(bytes[7])
	return
}

func (this *Image)GetJpgDimensions(file *os.File) (width int, height int) {
	fi, _ := file.Stat()
	fileSize := fi.Size()

	position := int64(4)
	bytes := make([]byte, 4)
	file.ReadAt(bytes[:2], position)
	length := int(bytes[0]<<8) + int(bytes[1])
	for position < fileSize {
		position += int64(length)
		file.ReadAt(bytes, position)
		length = int(bytes[2])<<8 + int(bytes[3])
		if (bytes[1] == 0xC0 || bytes[1] == 0xC2) && bytes[0] == 0xFF && length > 7 {
			file.ReadAt(bytes, position+5)
			width = int(bytes[2])<<8 + int(bytes[3])
			height = int(bytes[0])<<8 + int(bytes[1])
			return
		}
		position += 2
	}
	return 0, 0
}


func(this *Image) GetFormatByBytes(bytes []byte) string {

	/*if !bytes[3] {
		return ""
	}*/
	if bytes[0] == 0x89 && bytes[1] == 0x50 && bytes[2] == 0x4E && bytes[3] == 0x47 {
		return "png"
	}
	if bytes[0] == 0xFF && bytes[1] == 0xD8 {
		return "jpg"
	}
	if bytes[0] == 0x47 && bytes[1] == 0x49 && bytes[2] == 0x46 && bytes[3] == 0x38 {
		return "gif"
	}
	if bytes[0] == 0x42 && bytes[1] == 0x4D {
		return "bmp"
	}
	return ""
}

func(this *Image) GetGifDimensionsByBytes(bytes []byte) (width int, height int) {
	width = int(bytes[0]) + int(bytes[1])*256
	height = int(bytes[2]) + int(bytes[3])*256
	return
}

func (this *Image)GetBmpDimensionsByBytes(bytes []byte) (width int, height int) {

	width = int(bytes[3])<<24 | int(bytes[2])<<16 | int(bytes[1])<<8 | int(bytes[0])
	height = int(bytes[7])<<24 | int(bytes[6])<<16 | int(bytes[5])<<8 | int(bytes[4])
	return
}

func (this *Image)GetPngDimensionsByBytes(bytes []byte) (width int, height int) {

	width = int(bytes[0])<<24 | int(bytes[1])<<16 | int(bytes[2])<<8 | int(bytes[3])
	height = int(bytes[4])<<24 | int(bytes[5])<<16 | int(bytes[6])<<8 | int(bytes[7])
	return
}

func (this *Image)GetJpgDimensionsByBytes(bytes []byte) (width int, height int) {

	fileSize := len(bytes)
	position := 4
	length := int(bytes[0]<<8) + int(bytes[1])
	for position < fileSize {
		position += length
		length = int(bytes[2])<<8 + int(bytes[3])
		if (bytes[1] == 0xC0 || bytes[1] == 0xC2) && bytes[0] == 0xFF && length > 7 {
			width = int(bytes[2])<<8 + int(bytes[3])
			height = int(bytes[0])<<8 + int(bytes[1])
			return
		}
		position += 2
	}
	return 0, 0
}
