package main

import (
	"fmt"
	"strings"

	"github.com/imroc/req/v3"
)

func CreateFile(client *req.Client, isDir bool, path string) error {
	var response Response[interface{}]
	resp, err := client.R().
		SetBody(map[string]interface{}{
			"isDir":     isDir,
			"isLink":    false,
			"isSymlink": true,
			"linkPath":  "",
			"path":      path,
		}).
		SetSuccessResult(&response).
		Post("/api/v1/files")
	if err != nil {
		return fmt.Errorf("failed to create file %s: %v", path, err)
	}
	if !resp.IsSuccessState() {
		return fmt.Errorf("failed to create file %s: %s", path, resp.String())
	}
	return nil
}

func GetRedirects(client *req.Client, id int) []Redirect {
	f := func(id int) []Redirect {
		var redirectResp Response[[]Redirect]
		resp, err := client.R().
			SetBody(map[string]interface{}{
				"websiteID": id,
			}).
			SetSuccessResult(&redirectResp).
			Post("/api/v1/websites/redirect")
		if err != nil {
			panic(err)
		}
		if !resp.IsSuccessState() {
			panic(resp)
		}
		return redirectResp.Data
	}
	return f(id)
}

func GetProxies(client *req.Client, id int) []Proxy {
	f := func(id int) []Proxy {
		var proxyResp Response[[]Proxy]
		resp, err := client.R().
			SetBody(map[string]interface{}{
				"id": id,
			}).
			SetSuccessResult(&proxyResp).
			Post("/api/v1/websites/proxies")
		if err != nil {
			panic(err)
		}
		if !resp.IsSuccessState() {
			panic(resp)
		}
		return proxyResp.Data
	}
	return f(id)
}

func UpdateRedirects(client *req.Client, redirect Redirect) {
	resp, err := client.R().
		SetBody(redirect).
		Post("/api/v1/websites/redirect/update")
	if err != nil {
		panic(err)
	}
	if !resp.IsSuccessState() {
		panic(resp)
	}
}

func UpdateProxy(client *req.Client, proxy Proxy) error {
	proxy.Operate = "edit"
	content := strings.ReplaceAll(proxyContent, "location ^~ / {", fmt.Sprintf("location %s %s {", proxy.Modifier, proxy.Match))
	content = strings.ReplaceAll(content, "proxy_pass https://127.0.0.1;", fmt.Sprintf("proxy_pass %s%s; ", proxy.ProxyProtocol, proxy.ProxyAddress))
	content = strings.ReplaceAll(content, "proxy_set_header Host $host;", fmt.Sprintf("pproxy_set_header Host %s; ", proxy.ProxyHost))
	proxy.Content = content
	resp, err := client.R().
		SetBody(proxy).
		Post("/api/v1/websites/proxies/update")
	if err != nil {
		return fmt.Errorf("failed to update proxy: %v", err)
	}
	if !resp.IsSuccessState() {
		return fmt.Errorf("failed to update proxy: %s", resp.String())
	}
	return nil
}

func GetWebsites(client *req.Client, name string) []Website {
	f := func(name string, page int, pageSize int) DataList[Website] {
		var websitesResp Response[DataList[Website]]
		resp, err := client.R().
			SetBody(map[string]interface{}{
				"name":           name,
				"page":           page,
				"pageSize":       pageSize,
				"orderBy":        "created_at",
				"order":          "null",
				"websiteGroupId": 0,
			}).
			SetSuccessResult(&websitesResp).
			Post("/api/v1/websites/search")
		if err != nil {
			panic(err)
		}
		if !resp.IsSuccessState() {
			panic(resp)
		}
		return websitesResp.Data
	}

	total := 0
	length := -1
	var websites []Website
	page := 1
	size := 100

	for total != length {
		resp := f(name, page, size)
		total = resp.Total
		websites = append(websites, resp.Items...)
		length = len(websites)
		page++
	}
	return websites
}

func SaveFile(client *req.Client, content string, path string) error {
	var response Response[interface{}]
	resp, err := client.R().
		SetBody(&File{
			Content: content,
			Path:    path,
		}).
		SetSuccessResult(&response).
		Post("/api/v1/files/save")
	if err != nil {
		return fmt.Errorf("failed to save file: %v", err)
	}
	if !resp.IsSuccessState() {
		return fmt.Errorf("failed to save file: %s", resp.String())
	}
	return nil
}
