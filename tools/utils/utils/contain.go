package utils

func StringContainKeyOrNot(key string, keys []string) (result bool) {
	for _, compareKey := range keys {
		if key == compareKey {
			result = true
			return
		}
	}
	return
}

func Int32ContainKeyOrNot(key int32, keys []int32) (result bool) {
	for _, compareKey := range keys {
		if key == compareKey {
			result = true
			return
		}
	}
	return
}
