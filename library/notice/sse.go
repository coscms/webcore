package notice

import (
	"io"

	"github.com/admpub/sse"
	"github.com/webx-top/echo"
	sseRenderUtils "github.com/webx-top/echo/middleware/render/sse"
)

const SSEventName = `notice`

var SSERender = &sseRender{
	ServerSentEvents: sseRenderUtils.New(),
}

type sseRender struct {
	*sseRenderUtils.ServerSentEvents
}

type releaser interface {
	Release()
}

func (s *sseRender) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	err := s.ServerSentEvents.Render(w, name, data, c)
	data.(sse.Event).Data.(releaser).Release()
	return err
}

func (s *sseRender) RenderBy(w io.Writer, name string, f func(string) ([]byte, error), data interface{}, c echo.Context) error {
	err := s.ServerSentEvents.RenderBy(w, name, f, data, c)
	data.(sse.Event).Data.(releaser).Release()
	return err
}
