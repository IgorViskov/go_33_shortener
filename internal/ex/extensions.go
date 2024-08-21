package ex

import (
	"github.com/emirpasic/gods/sets/hashset"
)

func AnyString(set *hashset.Set, values []string) bool {
	for _, v := range values {
		if set.Contains(v) {
			return true
		}
	}
	return false
}
