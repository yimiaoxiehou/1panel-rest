package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/imroc/req/v3"
)

func main() {
	r := gin.Default()
	r.Use(Handler)
	r.POST("/redirect", postRedirect)
	r.POST("/save", saveFile)
	r.Run(":8080") // 监听并在 0.0.0.0:8080 上启动服务
}

type PostRedirectBody struct {
	Server       string `json:"server"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	Domain       string `json:"domain"`
	Name         string `json:"name"`
	Path         string `json:"path"`
	Target       string `json:"target"`
	RedirectType int    `json:"redirectType"`
}

type SaveFileBody struct {
	Server   string `json:"server"`
	Username string `json:"username"`
	Password string `json:"password"`
	Content  string `json:"content"`
	Path     string `json:"path"`
}

func saveFile(c *gin.Context) {
	var body SaveFileBody
	if c.ShouldBind(&body) == nil {
		client := GetAuthClient(body.Username, body.Password, body.Server)
		SaveFile(&client, body.Content, body.Path)
	}
}

func postRedirect(c *gin.Context) {
	var body PostRedirectBody
	if c.ShouldBind(&body) == nil {
		client := GetAuthClient(body.Username, body.Password, body.Server)
		UpsertRedirect(&client, body.Domain, body.Name, body.Path, body.Target, body.RedirectType)
	}
}

func UpsertRedirect(client *req.Client, domain string, name string, path string, target string, redirectype int) {
	permanent := "redirect"
	if redirectype == 302 {
		permanent = "permanent"
	}
	website := GetWebsites(client, domain)[0]
	filePath := fmt.Sprintf("%s/redirect/%s.conf", website.SitePath, name)
	CreateFile(client, false, filePath)
	SaveFile(client, fmt.Sprintf("rewrite %s %s %s;", path, target, permanent), filePath)
	redirects := GetRedirects(client, website.ID)
	for _, r := range redirects {
		if r.Name == name {
			r.Operate = "edit"
			UpdateRedirects(client, r)
		}
	}
}
