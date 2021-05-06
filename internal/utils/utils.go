package utils

import (
	"bufio"
	"os"
	"path/filepath"
)

func WriteStringToFile(file string, str string) error {
	if err := os.MkdirAll(filepath.Dir(file), 0755); err != nil {
		return err
	}

	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	_, err = w.WriteString(str)
	if err != nil {
		return err
	}
	w.Flush()
	return nil
}

/*
func AppendStringToFile(file string, str string) error {
	if err := os.MkdirAll(filepath.Dir(file), 0755); err != nil {
		return err
	}

	f, err := os.OpenFile(file, os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(str)
	if err != nil {
		return err
	}
	return nil
}
*/
