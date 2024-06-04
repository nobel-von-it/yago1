package files

import (
	"errors"
	"os"
	"strings"
)

func DirExists(path string) error {
	if !strings.Contains(path, "/") || len(strings.Split(path, "/")) != 2 {
		return errors.New("invalid path: " + path)
	}
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		err = os.MkdirAll(strings.Split(path, "/")[0], os.ModePerm)
		if err != nil {
			return err
		}
		file, err := os.Create(path)
		if err != nil {
			return err
		}
		defer func(file *os.File) {
			_ = file.Close()
		}(file)
	} else if err != nil {
		return err
	}
	return nil
}
