package tool

import (
	"errors"
	"os"
)

//判断name对应的文件是否存在，如果error不空则文件不存在
func ExistFile(name string) error {
	_, err := os.Stat(name)
	if err != nil {
		return err
	}

	if os.IsExist(err) {
		return errors.New("文件不存在")
	}

	return nil
}
