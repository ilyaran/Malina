package libraries

import (
	"github.com/ilyaran/Malina/entity"
	"encoding/json"
)

var ParameterLib *Parameter = &Parameter{
	Trees: Tree{},

}


type Parameter struct {
	Trees              Tree
	PublicJsonTreeList []byte
	MapPublic          map[int64]entity.Parameter
	ListPublic         []entity.Parameter
}
func(s *Parameter)SetTrees(heapList *[]entity.Scanable){
	s.Trees.HeapList = []entity.Hierarchical{}
	for _,v:=range *heapList{
		s.Trees.HeapList = append(s.Trees.HeapList,v.(*entity.Parameter))
	}
	s.Trees.SetTrees()
}

func(s *Parameter)SetPublicTrees(){
	s.ListPublic = []entity.Parameter{}
	s.MapPublic = map[int64]entity.Parameter{}
	for _,v := range s.Trees.List {
		s.ListPublic = append(s.ListPublic, *v.(*entity.Parameter))
		s.MapPublic[v.(*entity.Parameter).GetId()] = *v.(*entity.Parameter)
	}
	s.PublicJsonTreeList, _ = json.Marshal(s.ListPublic)
}

