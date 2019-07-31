package common

func GetMapKeys(m map[int]struct{}) []int {
	keys := []int{}
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
