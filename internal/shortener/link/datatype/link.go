package datatype

type Link struct {
	ID           *int64  `json:"id,omitempty"`
	Resource     *string `json:"resource,omitempty"`
	ShortLink    *string `json:"short_link,omitempty"`
	ShortLinkNum *string `json:"short_link_num,omitempty"`
	CustomName   *string `json:"custom_name,omitempty"`
}
