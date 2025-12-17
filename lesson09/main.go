package main

import "github.com/J1407B-K/buff/buff"

func main() {
	b := buff.NewEngine()

	b.GET("/hello", func(c *buff.Context) {
		c.JSON(200, "我喜欢你")
	})

	b.Run(":8080")
}
