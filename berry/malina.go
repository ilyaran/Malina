package berry

import (
	"github.com/ilyaran/Malina/entity"
	"net/http"
	"html/template"
)



type IController interface {
	GetList(malina *Malina,w http.ResponseWriter, r *http.Request, condition string)
	FormHandler(malina *Malina,action byte,w http.ResponseWriter, r *http.Request)
	ResultHandler(malina *Malina,args ... interface{})
	Default(malina *Malina, w http.ResponseWriter, r *http.Request)
}

type Malina struct {
	ChannelBool				chan bool

	CurrentAccount 		*entity.Account
	SessionId			string
	Status 				int64
	Result 				map[string]interface{}

	Controller       	IController
	List             	*[]entity.Scanable
	Item             	entity.Scanable


	IdInt64            	int64
	IdStr              	string
	Url 				string
	Order_index     	int64


	Per_page           	int64

	Page               	int64
	PageStr            	string

	Search            	string

	Err 				error


	TableSql    			string
	JoinSql    				string
	SelectSql   			string
	WhereSql    			string
	LimitSql    			string
	OrderBySql  			string

	All         			int64

	Department 			string
	ControllerName 		string
	Action             	string
	Paging  string
	Device     string

	NavAuth 				template.HTML
	Content 				template.HTML
	CategoryParameters 		template.HTML


}
