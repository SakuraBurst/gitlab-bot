package helpers

import "os"

func FileStatMust(file *os.File) os.FileInfo {
	stats, err := file.Stat()
	if err != nil {
		panic(err)
	}
	return stats
}
