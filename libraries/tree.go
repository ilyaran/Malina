package libraries

import (
	"github.com/ilyaran/Malina/entity"
)


type TreeSetter interface {
	SetTrees(heapList *[]entity.Scanable)
	SetPublicTrees()
}

type Tree struct {
	HeapList []entity.Hierarchical
	List     []entity.Hierarchical
	Map      map[int64]entity.Hierarchical

	ListPublic     []entity.Hierarchical
}

func (s *Tree) SetTrees() {
	if len(s.HeapList) > 0 {
		s.List = []entity.Hierarchical{}
		s.Map = map[int64]entity.Hierarchical{}

		var level int = 0

		var closure func(int64)
		closure = func(parent int64) {
			for i := 0; i < len(s.HeapList); i++ {
				if s.HeapList[i].GetParent() == parent {
					s.HeapList[i].SetLevel(level)
					s.Map[s.HeapList[i].GetId()] = s.HeapList[i]
					s.List = append(s.List, s.HeapList[i])
					level ++
					closure(s.HeapList[i].GetId())
				}
			}
			level --
		}
		closure(0)
		if len(s.List) > 0 {
			for _, a := range s.List {
				for _, v := range s.List {
					if v.GetParent() == a.GetId() {
						a.SetChildren(v)
					}
				}

				s.closure(a,a,"descendants")
				s.closure(a,a,"ancestors")
			}
		}
	}
}

func (this *Tree) closure(root,item entity.Hierarchical,kind string) {
	for _, v := range this.List {
		if kind == "descendants" {
			if v.GetParent() == item.GetId() {
				root.SetDescendants(v)
				this.closure(root,v, kind)
			}
		}
		if kind == "ancestors" {
			if v.GetId() == item.GetParent() {
				root.SetAncestors(v)
				this.closure(root,v, kind)
			}
		}
	}
}















