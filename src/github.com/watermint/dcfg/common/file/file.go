package file

import (
	"encoding/json"
	"os"
)

func FileExist(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func FileExistAndReadable(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		} else if os.IsPermission(err) {
			return false
		}
	}
	if info.IsDir() {
		return false
	}
	return true
}

func IsDirectory(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

func LoadJSON(path string, data interface{}) (result interface{}, err error) {
	j, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer j.Close()
	err = json.NewDecoder(j).Decode(data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func SaveJSON(path string, data interface{}) error {
	j, err := os.Create(path)
	if err != nil {
		return err
	}
	defer j.Close()

	err = json.NewEncoder(j).Encode(data)
	if err != nil {
		return err
	}
	return nil
}
