package helpers

import (
	"net/http"
	"time"
)

func SetAjaxHeader(w http.ResponseWriter){

	w.Header().Set("Cache-Control", "no-cache, private, max-age=0")
	w.Header().Set("Expires", time.Unix(0, 0).Format(http.TimeFormat))
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("X-Access-Expires", "0")
	w.Header().Set("Content-Type", "text/html, charset=utf-8")
}
