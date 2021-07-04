package link

import "github.com/ppdraga/go-shortener/internal/shortener/link/datatype"

type Controller struct {
	rw LinkReadWriter
}

type LinkReadWriter interface {
	ReadLink(id int64) (*datatype.Link, error)
	WriteLink(link *datatype.Link) (error, int64)
	FindLink(shortLink string) (*datatype.Link, error)
}

func NewController(rw LinkReadWriter) *Controller {
	return &Controller{
		rw: rw,
	}
}

func (c *Controller) AddLink(link *datatype.Link) (error, int64) {
	err, linkID := c.rw.WriteLink(link)
	if err != nil {
		return err, -1
	}
	return nil, linkID
}

func (c *Controller) GetLink(id int64) (*datatype.Link, error) {
	return c.rw.ReadLink(id)
}

func (c *Controller) FindLink(shortLink string) (*datatype.Link, error) {
	return c.rw.FindLink(shortLink)
}
