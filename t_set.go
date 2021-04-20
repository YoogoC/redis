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
	return &Set{goSet: mapset.NewThreadUnsafeSet()}
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
	panic("implement me")
}

func (s *Set) Set(key string) bool {
	return s.goSet.Add(key)
}
