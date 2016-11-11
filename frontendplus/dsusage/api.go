package dsusage

type ResponseGetNewUsageObj struct {
	Installs int `json:"new_install"`
	Opens    int `json:"new_open"`
}
