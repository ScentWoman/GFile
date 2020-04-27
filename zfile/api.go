package zfile

// IsInstalled responses /is-installed
type IsInstalled struct {
	Code int         `json:"code"` // -1
	Msg  string      `json:"msg"`  // ok
	Data interface{} `json:"data"` // null
}

// ListDriveResp responses /api/drive/list
type ListDriveResp struct {
	Code int     `json:"code"`
	Msg  string  `json:"msg"`
	Data []Drive `json:"data"`
}

// ConfigResp response /api/config
type ConfigResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data Config `json:"data"`
}

// ListResp responses /api/list
type ListResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		TotalPage int    `json:"totalPage"`
		FileList  []File `json:"fileList"`
	} `json:"data"`
}

// Config illustrates drive config.
type Config struct {
	Announcement     string  `json:"announcement"`
	Readme           *string `json:"readme"`
	SiteName         string  `json:"siteName"`
	Username         string  `json:"username"`
	Domain           string  `json:"domain"`
	Layout           string  `json:"layout"`
	TableSize        string  `json:"tableSize"`
	SearchEnable     bool    `json:"searchEnable"`
	ShowAnnouncement bool    `json:"showAnnouncement"`
	ShowDocument     bool    `json:"showDocument"`
	ShowOperator     bool    `json:"showOperator"`
	CustomCSS        string  `json:"customCss"`
	CustomJS         string  `json:"customJs"`
}

// Drive illustrates a drive.
type Drive struct {
	ID                         int    `json:"id"`
	Name                       string `json:"name"`
	EnableCache                bool   `json:"enableCache"`
	AutoRefreshCache           bool   `json:"autoRefreshCache"`
	SearchContainEncryptedFile bool   `json:"searchContainEncryptedFile"`
	SearchEnable               bool   `json:"searchEnable"`
	SearchIgnoreCase           bool   `json:"searchIgnoreCase"`
	Type                       struct {
		Key         string `json:"key"`
		Description string `json:"description"`
	} `json:"type"`
}

// File illustrates a file.
type File struct {
	Name string  `json:"name"`
	Path string  `json:"path"`
	Size int64   `json:"size"`
	Time string  `json:"time"`
	Type string  `json:"type"`
	URL  *string `json:"url"`
}
