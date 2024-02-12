package utils

func Contains(s []string, el string) bool {
	for _, x := range s {
		if x == el {
			return true
		}
	}
	return false
}
