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
	r.POST("/proxy", postProxy)
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

func postProxy(c *gin.Context) {
	var body PostProxyBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}

	client := GetAuthClient(body.Username, body.Password, body.Server)
	if err := UpsertProxy(&client, body.Domain, body.Name, body.Modifier, body.Match, body.ProxyProtocol, body.ProxyAddress); err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "success",
	})
}

type PostProxyBody struct {
	Server        string `json:"server"`
	Username      string `json:"username"`
	Password      string `json:"password"`
	Domain        string `json:"domain"`
	Name          string `json:"name"`
	Modifier      string `json:"modifier"`
	Match         string `json:"match"`
	ProxyHost     string `json:"proxyHost"`
	Replaces      string `json:"replaces"`
	ProxyProtocol string `json:"proxyProtocol"`
	ProxyAddress  string `json:"proxyAddress"`
}

var proxyContent = `
location ^~ / {
proxy_pass https://127.0.0.1; 
proxy_set_header Host $host; 
proxy_set_header X-Real-IP $remote_addr; 
proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for; 
proxy_set_header REMOTE-HOST $remote_addr; 
proxy_set_header Upgrade $http_upgrade; 
proxy_set_header Connection $http_connection; 
proxy_set_header X-Forwarded-Proto $scheme; 
proxy_http_version 1.1; 
add_header X-Cache $upstream_cache_status; 
add_header Cache-Control no-cache; 
proxy_ssl_server_name off; 
proxy_ssl_name $proxy_host; 
add_header Strict-Transport-Security "max-age=31536000"; 
}`

func UpsertProxy(client *req.Client, domain string, name string, modifier string, match string, proxyProtocol string, proxyAddress string) error {
	websites := GetWebsites(client, domain)
	if len(websites) == 0 {
		return fmt.Errorf("website not found: %s", domain)
	}

	var web *Website
	for _, _web := range websites {
		if _web.PrimaryDomain == domain {
			web = &_web
		}
	}

	if web == nil {
		return fmt.Errorf("website not found: %s", domain)
	}

	filePath := fmt.Sprintf("%s/proxy/%s.conf", web.SitePath, name)
	if err := CreateFile(client, false, filePath); err != nil {
		return fmt.Errorf("create file failed: %v", err)
	}

	if err := SaveFile(client, proxyContent, filePath); err != nil {
		return fmt.Errorf("save file failed: %v", err)
	}

	proxies := GetProxies(client, web.ID)
	for _, p := range proxies {
		if p.Name == name {
			p.Modifier = modifier
			p.Match = match
			p.ProxyProtocol = proxyProtocol
			p.ProxyAddress = proxyAddress
			p.ProxyPass = proxyProtocol + proxyAddress
			if err := UpdateProxy(client, p); err != nil {
				return fmt.Errorf("update proxy failed: %v", err)
			}
		}
	}
	return nil
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
