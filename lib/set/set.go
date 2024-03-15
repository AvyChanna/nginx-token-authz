package set

import (
	"encoding/json"

	"gopkg.in/yaml.v3"
)

type StrSet = Set[string]

var (
	_ yaml.Unmarshaler = (*StrSet)(nil)
	_ yaml.Marshaler   = (*StrSet)(nil)
	_ yaml.Marshaler   = StrSet{}

	_ json.Unmarshaler = (*StrSet)(nil)
	_ json.Marshaler   = (*StrSet)(nil)
	_ json.Marshaler   = StrSet{}
)

type (
	empty             = struct{}
	Set[T comparable] map[T]empty
)

func (s Set[T]) Add(v T) {
	s[v] = empty{}
}

func (s Set[T]) Remove(v T) {
	delete(s, v)
}

func (s Set[T]) Values() []T {
	ret := make([]T, 0, len(s))
	for v := range s {
		ret = append(ret, v)
	}
	return ret
}

func (s Set[T]) Contains(v T) bool {
	_, ok := s[v]
	return ok
}

func (s Set[T]) MarshalJSON() ([]byte, error) {
	keys := make([]T, 0, len(s))
	for v := range s {
		keys = append(keys, v)
	}

	return json.Marshal(keys)
}

func (s *Set[T]) UnmarshalJSON(data []byte) error {
	var keys []T
	err := json.Unmarshal(data, &keys)
	if err != nil {
		return err
	}

	for _, key := range keys {
		s.Add(key)
	}
	return nil
}

func (s Set[T]) MarshalYAML() (any, error) {
	keys := make([]T, 0, len(s))
	for v := range s {
		keys = append(keys, v)
	}
	return keys, nil
}

func (s *Set[T]) UnmarshalYAML(data *yaml.Node) error {
	var keys []T
	err := data.Decode(&keys)
	if err != nil {
		return err
	}

	*s = Set[T]{}

	for _, key := range keys {
		s.Add(key)
	}
	return nil
}
