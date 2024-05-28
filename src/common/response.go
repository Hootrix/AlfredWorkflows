package common

type AlfredResponse struct {
	Items []AlfredItem `json:"items"`
}

type AlfredItem struct {
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
	Arg      string `json:"arg,omitempty"`
}
