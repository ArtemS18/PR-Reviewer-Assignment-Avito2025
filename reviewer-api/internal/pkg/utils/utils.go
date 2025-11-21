package utils

import (
	"math/rand"
	"time"
)

func Remove[T comparable](s []T, i int) []T {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func GetRandomSlice(baseSlice []string) []string {
	if len(baseSlice) == 0 {
		return []string{}
	}
	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)
	tmp := make([]string, len(baseSlice))
	copy(tmp, baseSlice)
	randomSlice := make([]string, 0)
	for i := 0; i < 2; i++ {
		var rand_id int
		if len(tmp) > 0 {
			rand_id = r.Intn(len(tmp))
			randomSlice = append(randomSlice, tmp[rand_id])
			tmp = Remove(tmp, rand_id)
		} else {
			break
		}

	}
	return randomSlice
}
func GetRandomNumber(baseSlice []string) string {
	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)
	var rand_id int
	if len(baseSlice) > 0 {
		rand_id = r.Intn(len(baseSlice))
	} else {
		return ""
	}
	return baseSlice[rand_id]
}
