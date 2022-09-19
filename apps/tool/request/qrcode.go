package request

type QrCodeReq struct {
	Url  string `json:"url"`
	Size int    `json:"size"`
}
