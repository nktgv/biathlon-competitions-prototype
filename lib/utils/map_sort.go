package utils

import (
	"biathlon-competitions-prototype/lib"
	"sort"
)

type kv struct {
	Key   int
	Value *lib.Result
}

func Sort(m map[int]*lib.Result) map[int]*lib.Result {
	var ss []kv
	for k, v := range m {
		ss = append(ss, kv{k, v})
	}

	sort.Slice(ss, func(i, j int) bool {
		if ss[i].Value.Status == "[NotStarted]" || ss[i].Value.Status == "[NotFinished]" {
			return false
		}
		if ss[j].Value.Status == "[NotStarted]" || ss[j].Value.Status == "[NotFinished]" {
			return true
		}
		return ss[i].Value.TotalTime > ss[j].Value.TotalTime
	})

	for _, kv := range ss {
		m[kv.Key] = kv.Value
	}

	return m
}
