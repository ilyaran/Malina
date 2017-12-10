/**
 * @author		John Aran (Ilyas Aranzhanovich Toxanbayev)
 * @version		1.0.0
 * @based on
 * @email      		il.aranov@gmail.com
 * @link
 * @github      	https://github.com/ilyaran/github.com/ilyaran/Malina
 * @license		MIT License Copyright (c) 2017 John Aran (Ilyas Toxanbayev)
 */package berry

import (
	"html/template"
	"net/http"
	"time"

	"github.com/ilyaran/Malina/entity"
)

type IController interface {
	GetList(malina *Malina, w http.ResponseWriter, r *http.Request, condition string)
	FormHandler(malina *Malina, action byte, w http.ResponseWriter, r *http.Request)
	ResultHandler(malina *Malina, args ...interface{})
	Default(malina *Malina, w http.ResponseWriter, r *http.Request)
}

type Malina struct {
	ChannelBool    chan bool
	Lang           string
	CurrentAccount *entity.Account
	SessionId      string
	Status         int64
	Result         map[string]interface{}

	Controller IController
	List       *[]entity.Scanable
	Item       entity.Scanable

	IdInt64     int64
	IdStr       string
	Url         string
	Order_index int64

	Per_page int64

	Page    int64
	PageStr string

	Search string

	TableSql   string
	JoinSql    string
	SelectSql  string
	WhereSql   string
	LimitSql   string
	OrderBySql string

	All int64

	Department     string
	ControllerName string
	Action         string
	Paging         string
	Device         string

	NavAuth            template.HTML
	Content            template.HTML
	CategoryParameters template.HTML
}

//**** Context
// Done returns a channel that is closed when this Context is canceled
// or times out.
func (s *Malina) Done() <-chan struct{} {

	t := make(chan struct{})

	return t
}

// Err indicates why this context was canceled, after the Done channel
// is closed.
func (s *Malina) Err() error {
	var err error

	return err
}

// Deadline returns the time when this Context will be canceled, if any.
func (s *Malina) Deadline() (deadline time.Time, ok bool) {

	return time.Now(), false
}

// Value returns the value associated with key or nil if none.
func (s *Malina) Value(key interface{}) interface{} {

	return key
}

//**** end Context
