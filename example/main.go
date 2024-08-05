package main

import (
	"fmt"
	"os"

	"testformdata/client"
	"testformdata/server"
)

func main() {
	go server.Start(":3000")

	h := client.NewClient()
	file, _ := os.Open("./testdata/test")
	hashFile, status, err := h.SendFile("http://localhost:3000", "testfilename", file)

	fmt.Printf("MD5 Hash файла: %s, \r\n Статус ответа сервера: %d \r\n Ошибка: %v \r\n", hashFile, status, err)

}
