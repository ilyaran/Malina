package publicView


import (
	"fmt"
	"github.com/ilyaran/Malina/entity"

)

var NewsView * News
func NewsViewInit(){
	NewsView = &News{}

}


type News struct {
	List   []entity.Scanable
	Paging string
}

func (s *News)Index(){

	s.Listing()

}
func (s *News)Listing()string{

	for i :=0; s.List[i].GetId()>0;i++{
		fmt.Println("id:",s.List[i].GetId())
	}
	fmt.Println("paging:",s.Paging)

	return s.Paging
}
