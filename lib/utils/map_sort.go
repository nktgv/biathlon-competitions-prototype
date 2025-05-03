package utils

import (
	"sort"
)

type kv struct {
	Key   int
	Value *Result
}

func Sort(m map[int]*Result) []int {
	var ss []kv
	for k, v := range m {
		ss = append(ss, kv{k, v})
	}

	sort.Slice(ss, func(i, j int) bool {
		iNotReady := ss[i].Value.Status == "[NotStarted]" || ss[i].Value.Status == "[NotFinished]"
		jNotReady := ss[j].Value.Status == "[NotStarted]" || ss[j].Value.Status == "[NotFinished]"

		if iNotReady && jNotReady {
			return ss[i].Value.TotalTime < ss[j].Value.TotalTime
		}

		if iNotReady {
			return true
		}
		if jNotReady {
			return false
		}

		return ss[i].Value.TotalTime < ss[j].Value.TotalTime
	})

	keys := make([]int, len(ss))
	for i, kv := range ss {
		keys[i] = kv.Key
	}

	return keys
}
