package main

type Response[T any] struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    T      `json:"data"`
}

type AuthResp struct {
	Name      string `json:"name"`
	Token     string `json:"token"`
	MfaStatus string `json:"mfaStatus"`
}

type DataList[T any] struct {
	Total int `json:"total"`
	Items []T `json:"items"`
}

type File struct {
	Content string
	Path    string
}

type Website struct {
	ID            int    `json:"id"`
	CreatedAt     string `json:"createdAt"`
	Protocol      string `json:"protocol"`
	PrimaryDomain string `json:"primaryDomain"`
	Type          string `json:"type"`
	Alias         string `json:"alias"`
	Remark        string `json:"remark"`
	Status        string `json:"status"`
	ExpireDate    string `json:"expireDate"`
	SitePath      string `json:"sitePath"`
	AppName       string `json:"appName"`
	RuntimeName   string `json:"runtimeName"`
	SslExpireDate string `json:"sslExpireDate"`
	SslStatus     string `json:"sslStatus"`
	AppInstallId  int    `json:"appInstallId"`
	RuntimeType   string `json:"runtimeType"`
}

type Redirect struct {
	WebsiteID    int      `json:"websiteID"`
	Name         string   `json:"name"`
	Domains      []string `json:"domains,omitempty"`
	KeepPath     bool     `json:"keepPath"`
	Enable       bool     `json:"enable"`
	Type         string   `json:"type"`
	Redirect     string   `json:"redirect"`
	Path         string   `json:"path"`
	Target       string   `json:"target"`
	FilePath     string   `json:"filePath"`
	Content      string   `json:"content"`
	RedirectRoot bool     `json:"redirectRoot"`
	Operate      string   `json:"operate,omitempty"`
}

type Proxy struct {
	ID            int                    `json:"id"`
	Operate       string                 `json:"operate" default:"edit"`
	Enable        bool                   `json:"enable"`
	Cache         bool                   `json:"cache"`
	CacheTime     int                    `json:"cacheTime"`
	CacheUnit     string                 `json:"cacheUnit"`
	Name          string                 `json:"name"`
	Modifier      string                 `json:"modifier"`
	Match         string                 `json:"match"`
	ProxyPass     string                 `json:"proxyPass"`
	ProxyHost     string                 `json:"proxyHost"`
	Replaces      map[string]interface{} `json:"replaces"`
	ProxyProtocol string                 `json:"proxyProtocol"`
	ProxyAddress  string                 `json:"proxyAddress"`
	Content       string                 `json:"content"`
	FilePath      string                 `json:"filePath"`
	SNI           bool                   `json:"sni"`
	ProxySSLName  string                 `json:"proxySSLName"`
}

type LoginPayload struct {
	Name       string `json:"name"`
	Password   string `json:"password"`
	AuthMethod string `json:"authMethod"`
	Language   string `json:"language"`
}

type RedirectConfigPayload struct {
	WebsiteID int    `json:"websiteID"`
	Name      string `json:"Name"`
	Redirect  string `json:"redirect"`
	Content   string `json:"content"`
}
