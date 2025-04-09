package main

import (
	"api-gateway/config"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func main() {
	config := config.New()

	r := gin.Default()

	r.Any("/*proxyPath", func(c *gin.Context) {
		targetURL := ""
		path := c.Param("proxyPath")
		rawQuery := c.Request.URL.RawQuery

		// Выбор сервиса в зависимости от пути
		if strings.HasPrefix(path, "/orders") {
			targetURL = config.OrderService.Addr + path
		} else if strings.HasPrefix(path, "/products") {
			targetURL = config.InventoryService.Addr + path
		} else {
			c.JSON(http.StatusNotFound, gin.H{"error": "unknown route"})
			return
		}

		// Добавляем query параметры, если они есть
		if rawQuery != "" {
			targetURL += "?" + rawQuery
		}

		// Создаём новый запрос
		req, err := http.NewRequest(c.Request.Method, targetURL, c.Request.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot create request"})
			return
		}

		// Копируем заголовки
		for key, values := range c.Request.Header {
			for _, value := range values {
				req.Header.Add(key, value)
			}
		}

		// Отправляем запрос
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "target service unavailable"})
			return
		}
		defer resp.Body.Close()

		// Копируем заголовки ответа
		for k, vv := range resp.Header {
			for _, v := range vv {
				c.Writer.Header().Add(k, v)
			}
		}

		// Устанавливаем статус и возвращаем тело ответа
		c.Status(resp.StatusCode)
		io.Copy(c.Writer, resp.Body)
	})

	r.Run(":8080") // API Gateway будет слушать здесь
}
