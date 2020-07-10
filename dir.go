package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
)

type ShortName struct {
	name [8]byte
	ext  [3]byte
}

func StringToShortName(origName string, fnMap map[string]uint8) *ShortName {
	name := bytes.ToUpper([]byte(origName))
	separatorIndex := -1
	modified := false

	for i := len(name) - 1; i > 0; i-- {
		if name[i]>>7 != 0 || name[i] == '+' {
			// We don't care about Unicode here
			name[i] = '_'
			modified = true
		} else if name[i] == '.' && i != 0 && separatorIndex == -1 {
			separatorIndex = i
		} else if name[i] == ' ' || name[i] == '.' {
			name = append(name[:i], name[i+1:]...)
			if separatorIndex != -1 {
				separatorIndex--
			}
			i++
			modified = true
		}
	}

	var ext []byte
	if separatorIndex == -1 {
		ext = nil
	} else {
		ext = make([]byte, len(name[separatorIndex+1:]))
		copy(ext, name[separatorIndex+1:])
		name = name[:separatorIndex]
	}

	if len(name) > 8 || len(ext) > 3 {
		modified = true
	}
	if modified {
		if len(name) > 6 {
			name = name[:6]
		}
		// Number to append to end of file (after '~')
		num := fnMap[string(name)]
		num++
		fnMap[string(name)] = num

		if num < 10 {
			name = append(
				name,
				'~',
				'0'+num,
			)
		} else {
			name = append(
				name[:len(name)-1],
				'~',
				'0'+num/10,
				'0'+num%10,
			)
		}
	}

	shortName := &ShortName{}

	for i := 0; i < len(shortName.name); i++ {
		if i < len(name) {
			shortName.name[i] = name[i]
		} else {
			shortName.name[i] = ' '
		}
	}
	for i := 0; i < len(shortName.ext); i++ {
		if i < len(ext) {
			shortName.ext[i] = ext[i]
		} else {
			shortName.ext[i] = ' '
		}
	}
	return shortName
}

func ShortNameToString(shortName *ShortName) string {
	return fmt.Sprintf("%s %s", shortName.name, shortName.ext)
}

func main() {
	var path string
	if len(os.Args) > 1 {
		path = os.Args[1]
	} else {
		path = "."
	}

	err, files := getFiles(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	for _, file := range files {
		fmt.Println(ShortNameToString(file))
	}
}

func getFiles(path string) (error, []*ShortName) {
	fileInfos, err := ioutil.ReadDir(path)
	if err != nil {
		return err, nil
	}

	files := make([]*ShortName, len(fileInfos))
	fnMap := make(map[string]uint8)

	for i, fileInfo := range fileInfos {
		files[i] = StringToShortName(fileInfo.Name(), fnMap)
	}

	return nil, files
}
