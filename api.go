package main

import "github.com/imroc/req/v3"

func CreateFile(client *req.Client, isDir bool, path string) {
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
		panic(err)
	}
	if !resp.IsSuccessState() {
		panic(resp)
	}
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

func SaveFile(client *req.Client, content string, path string) {
	var response Response[interface{}]
	resp, err := client.R().
		SetBody(&File{
			Content: content,
			Path:    path,
		}).
		SetSuccessResult(&response).
		Post("/api/v1/files/save")
	if err != nil {
		panic(err)
	}
	if !resp.IsSuccessState() {
		panic(resp)
	}
}
