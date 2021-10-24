package main

import (
	"fmt"
	"log"

	"github.com/0penTools/panel/system"
	"github.com/0penTools/panel/web"
	"github.com/digitalocean/go-libvirt"
	"github.com/digitalocean/go-libvirt/socket/dialers"
	"github.com/gin-gonic/gin"
)

func main() {
	system.Libvirt = libvirt.NewWithDialer(dialers.NewLocal())
	l := system.Libvirt
	if err := l.Connect(); err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	fmt.Println("Connected to libvirt")

	fmt.Println("Starting gin...")
	app := gin.Default()
	app.GET("/domains", web.ListDomain)
	app.POST("/domains/create", web.CreateDomain)
	app.POST("/domains/start", web.StartDomain)
	app.Run(":8080")

	if err := l.Disconnect(); err != nil {
		log.Fatalf("failed to disconnect: %v", err)
	}
}
