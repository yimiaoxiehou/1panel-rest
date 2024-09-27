package main

import "github.com/imroc/req/v3"

func GetToken(username string, password string, server string) string {
	var authResp Response[AuthResp]
	resp, err := req.C().DevMode().R().
		SetBody(&LoginPayload{
			Name:       username,
			Password:   password,
			AuthMethod: "jwt",
			Language:   "zh",
		}).
		SetHeaders(map[string]string{
			"entrancecode": "eWltaWFv",
			"content-Type": "application/json",
		}).
		SetSuccessResult(&authResp).
		Post(server + "/api/v1/auth/login")
	if err != nil {
		panic(err)
	}
	if !resp.IsSuccessState() {
		panic(resp)
	}
	if authResp.Code != 200 {
		panic(authResp.Message)
	}
	return authResp.Data.Token
}

func GetAuthClient(username string, password string, server string) req.Client {
	client := req.C()
	client.OnBeforeRequest(func(client *req.Client, req *req.Request) error {
		req = req.SetURL(server + req.RawURL)
		return nil
	})

	token := GetToken(username, password, server)
	authHeader := map[string]string{
		"accept":             "application/json, text/plain, */*",
		"content-type":       "application/json",
		"PanelAuthorization": token,
	}
	client = client.SetCommonHeaders(authHeader)
	return *client
}
