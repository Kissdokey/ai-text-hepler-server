package utils

import (
	"encoding/json"
	"fmt"
	"os"
)

func WriteJSONToFile(data interface{}, filePath string) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(jsonData)
	if err != nil {
		return err
	}

	fmt.Println("JSON data has been written to", filePath)
	return nil
}

func WriteJson() {
	// 示例数据
	jsonData := map[string]interface{}{
		"name":  "John Doe",
		"age":   25,
		"email": "john.doe@example.com",
	}

	// 指定保存路径和文件名,上数据库前先这样
	filePath := "E:/ai-text-helper-server/userInfo/example.json"

	// 调用函数写入JSON数据到文件
	err := WriteJSONToFile(jsonData, filePath)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
}
