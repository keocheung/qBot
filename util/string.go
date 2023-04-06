package util

func StringsContain(s []string, str string) bool {
	for _, s := range s {
		if s == str {
			return true
		}
	}
	return false
}

func StringArraysHasCommon(s1, s2 []string) bool {
	for _, s1 := range s1 {
		for _, s2 := range s2 {
			if s1 == s2 {
				return true
			}
		}
	}
	return false
}
