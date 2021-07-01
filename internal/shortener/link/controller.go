package link

import "github.com/ppdraga/go-shortener/internal/shortener/link/datatype"

type Controller struct {
	rw LinkReadWriter
}

type LinkReadWriter interface {
	ReadLink(id int64) (*datatype.Link, error)
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

func (c *Controller) GetLink(id int64) (*datatype.Link, error) {
	return c.rw.ReadLink(id)
}
