package helper

import (
	"net/http"
	"time"
	"database/sql"
)

func SetAjaxHeaders(w http.ResponseWriter) {
	w.Header().Set("Cache-Control", "no-cache, private, max-age=0")
	w.Header().Set("Expires", time.Unix(0, 0).Format(http.TimeFormat))
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("X-Accel-Expires", "0")
	//w.Header().Set("Content-Type", "text/html, charset=utf-8")
}
func NewNullInt64(s int64) *sql.NullInt64{
	if s > 0 {
		return &sql.NullInt64{
			Int64: s,
			Valid: true,
		}
	}
	return &sql.NullInt64{}
}