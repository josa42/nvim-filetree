package files

import "os"

func childrenNames(path string) []string {
	names := []string{}

	file, err := os.Open(path)
	if err != nil {
		return names
	}
	defer file.Close()

	list, _ := file.Readdirnames(0) // 0 to read all files and folders
	for _, name := range list {
		if !dummyIsIgnored(name) {
			names = append(names, name)
		}
	}
	return names
}

func dummyIsIgnored(name string) bool {
	for _, i := range []string{".git", ".DS_Store"} {
		if i == name {
			return true
		}
	}
	return false
}

func isDir(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		return false
	}
	return fi.Mode().IsDir()
}

