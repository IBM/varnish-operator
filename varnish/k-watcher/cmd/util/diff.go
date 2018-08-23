package util

// DiffBackends returns
// backends that only exist in a1
// backends that only exist in a2
func DiffBackends(a1, a2 []string) (onlyA1, onlyA2 []string) {
	diffM := make(map[string]int)

	for _, v := range a1 {
		diffM[v] = 1
	}
	for _, v := range a2 {
		diffM[v] = diffM[v] - 1
	}
	for k, v := range diffM {
		if v == 1 {
			onlyA1 = append(onlyA1, k)
		}
		if v == -1 {
			onlyA2 = append(onlyA2, k)
		}
	}
	return
}
