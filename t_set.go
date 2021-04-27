package redis

import (
	mapset "github.com/deckarep/golang-set"
)

const SetType = uint64(2)
const SetTypeFancy = "set"

var _ Item = &Set{}

type Set struct {
	goSet mapset.Set
}

func NewSet() *Set {
	return &Set{goSet: mapset.NewSet()}
}

func (s *Set) Value() interface{} {
	return s.goSet
}

func (s *Set) Type() uint64 {
	return SetType
}

func (s *Set) TypeFancy() string {
	return SetTypeFancy
}

func (s *Set) OnDelete(key *string, db *RedisDb) {
	s.goSet.Clear()
}

func (s *Set) Set(key string) bool {
	return s.goSet.Add(key)
}

func (s *Set) Len() int {
	return s.goSet.Cardinality()
}

func (s *Set) Remove(value string) bool {
	if s.goSet.Contains(value) {
		s.goSet.Remove(value)
		return true
	}
	return false
}

func (s *Set) Equal(other *Set) bool {
	return s.goSet.Equal(other.goSet)
}

func (s *Set) Iter() <-chan interface{} {
	return s.goSet.Iter()
}

func (s *Set) Contains(i ...interface{}) bool {
	return s.goSet.Contains(i...)
}
