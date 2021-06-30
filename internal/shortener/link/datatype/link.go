package datatype

type Link struct {
	ID           *int64  `json:"id"`
	Resource     *string `json:"resource"`
	ShortLink    *string `json:"short_link"`
	ShortLinkNum *string `json:"short_link_num"`
	CustomName   *string `json:"custom_name"`
}
