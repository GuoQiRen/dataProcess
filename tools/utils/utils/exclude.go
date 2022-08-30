package utils

func ExcludeSuffixFileName(srcFileName string, suffix string) (fileName string) {
	postIndex := len(srcFileName) - len(suffix)
	if postIndex < 0 {
		return
	}

	fileName = srcFileName[:postIndex]
	return
}
