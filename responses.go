package veritrans

type VTWebChargeResponse struct {
	StatusCode    string `json:"status_code"`
	StatusMessage string `json:"status_message"`
	RedirectUrl   string `json:"redirect_url"`
}
