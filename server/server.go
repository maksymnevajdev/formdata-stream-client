package server

import (
	"crypto/md5"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"os"
	"os/signal"

	"github.com/gofiber/fiber/v2"
)

func Start(addr string) error {
	app := fiber.New(fiber.Config{DisablePreParseMultipartForm: true, StreamRequestBody: true})

	app.Post("/", func(c *fiber.Ctx) error {
		v := c.Get("Content-Type")
		if v == "" {
			return nil
		}
		d, params, err := mime.ParseMediaType(v)
		if err != nil || !(d == "multipart/form-data" || d == "multipart/mixed") {
			return nil
		}
		boundary, ok := params["boundary"]
		if !ok {
			return nil
		}
		reader := multipart.NewReader(c.Context().RequestBodyStream(), boundary)
		hash := md5.New()
		for {
			part, err := reader.NextPart()

			if err != nil {
				if err == io.EOF {
					fmt.Println("EOF")
					break
				} else {
					fmt.Println("Other type of error", err)
					return nil
				}
			}
			fileSave, _ := os.OpenFile("./"+part.FileName(), os.O_RDWR|os.O_APPEND|os.O_CREATE, 0660)
			w := io.MultiWriter(fileSave, hash)
			io.Copy(w, part)
			fileSave.Close()
		}

		c.WriteString(fmt.Sprintf("%x", hash.Sum(nil)))
		return nil
	})

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	go func() {
		_ = <-c
		fmt.Println("Завершение работы сервера")
		_ = app.Shutdown()
	}()

	if err := app.Listen(addr); err != nil {
		return fmt.Errorf("%s", err)
	}

	return nil
}
