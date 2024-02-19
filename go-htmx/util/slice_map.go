package util

func SliceMap[E any, O any](slice []E, fn func(E) O) []O {
	var ret = make([]O, 0)
	if len(slice) == 0 {
		return nil
	}
	for _, e := range slice {
		ret = append(ret, fn(e))
	}
	return ret
}
