package link

import (
	"github.com/opentracing/opentracing-go"
	"github.com/ppdraga/go-shortener/internal/shortener/link/datatype"
	"go.uber.org/zap"
)

type Controller struct {
	rw     LinkReadWriter
	Logger *zap.Logger
	Tracer opentracing.Tracer
}

type LinkReadWriter interface {
	ReadLink(id int64) (*datatype.Link, error)
	WriteLink(link *datatype.Link) (error, int64)
	FindLink(shortLink string) (*datatype.Link, error)
}

func NewController(rw LinkReadWriter, logger *zap.Logger, tracer opentracing.Tracer) *Controller {
	return &Controller{
		rw:     rw,
		Logger: logger,
		Tracer: tracer,
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
