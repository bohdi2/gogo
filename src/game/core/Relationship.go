package core

import (
//. "container/vector"
//"fmt"
)

type Relationships struct {
	Matrix [][]CardId
}

func NewRelationships(count int) *Relationships {
	mat := make([][]CardId, count)
	for i := range mat {
		mat[i] = make([]CardId, 0)
	}
	return &Relationships{mat}
}

func (self *Relationships) Row(n CardId) []CardId {
	return self.Matrix[n]
}

func (self *Relationships) Append(rowId CardId, id CardId) {
	self.Matrix[rowId] = append(self.Matrix[rowId], id)
}

func (self Relationships) Visit(id CardId, f func(CardId)) {
	visited := make([]bool, len(self.Matrix))
	self.visit2(id, f, &visited)
}

func (self Relationships) visit2(id CardId, f func(CardId), visited *[]bool) {
	if false == (*visited)[id] {
		for _, child := range self.Row(id) {
			f(child)
			self.visit2(child, f, visited)
		}
		(*visited)[id] = true
	}
}
