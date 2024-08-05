package client

import (
	"crypto/md5"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"os"
	"testing"

	"testformdata/server"
)

var FileSize = flag.Int("filesize", 100000, "Размер файла в байтах")
var FileName = flag.String("filename", "default", "Имя файла на сервере")
var RemoveFile = flag.Bool("remove", false, "Удалять ли файл, сервером?")

func TestHttpClient_SendFile(t *testing.T) {
	os.Remove("./" + *FileName)
	go server.Start(":3000")
	h := NewClient()
	hashCheck := md5.New()

	f, wr := io.Pipe()
	w := io.MultiWriter(wr, hashCheck)

	go func() {
		GenRandomBytes(4096, *FileSize, w)
		wr.Close()
	}()

	hash, _, err := h.SendFile("http://localhost:3000/", *FileName, f)
	if err != nil {
		t.Fatal(err)
		return
	}

	mdHash1 := fmt.Sprintf("%x", hashCheck.Sum(nil))
	mdHash2 := fmt.Sprintf("%s", hash)

	if mdHash1 == mdHash2 {
		t.Logf("md5 хеши файлов совпадают %x %x", mdHash1, mdHash2)
	} else {
		t.Fatalf("md5 хеши файлов не совпадают %x %x", mdHash1, mdHash2)
	}

	if *RemoveFile {
		os.Remove("./" + *FileName)
	}

}

func GenRandomBytes(size int, fileSizeInBytes int, w io.Writer) {

	bytesCount := 0
	for {

		if bytesCount >= fileSizeInBytes {
			fmt.Println(bytesCount)
			return
		}

		if bytesCount+size > fileSizeInBytes {
			size = fileSizeInBytes - bytesCount

		}

		blk := make([]byte, size)
		rand.Read(blk)
		n, err := w.Write(blk)
		if err != nil {
			return
		}
		bytesCount += n

	}
}
