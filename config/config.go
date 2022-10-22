package config

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

func ReadConfig(path string) ([]Config, error) {
	f, err := os.Open(path)
	if err != nil {
		// 打开文件失败
		log.Fatal(err)
		return nil, err
	}
	defer f.Close()
	var data []byte
	buf := make([]byte, 1024)
	for {
		n, err := f.Read(buf)
		if err != nil && err != io.EOF {
			log.Fatal(err)
		}
		if n == 0 {
			break
		}
		data = append(data, buf[:n]...)
	}

	var configs []Config

	err = json.Unmarshal(data, &configs)
	if err != nil {
		// 打开文件失败
		log.Fatal(err)
		return nil, err
	}
	return configs, nil
}
