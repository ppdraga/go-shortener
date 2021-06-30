package link

import "github.com/ppdraga/go-shortener/internal/shortener/link/datatype"

type Controller struct {
	rw LinkReadWriter
}

type LinkReadWriter interface {
	ReadLink(params map[string]interface{}) (*datatype.Link, error)
	WriteLink(link *datatype.Link) error
}

func NewController(rw LinkReadWriter) *Controller {
	return &Controller{
		rw: rw,
	}
}

func (c *Controller) AddLink(link *datatype.Link) error {
	return c.rw.WriteLink(link)
}

func (c *Controller) GetLink(params map[string]interface{}) (*datatype.Link, error) {
	return c.rw.ReadLink(params)
}
