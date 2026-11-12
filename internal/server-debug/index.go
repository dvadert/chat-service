package serverdebug

import (
	_ "embed"
	"html/template"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/dvadert/chat-service/internal/logger"
)

//go:embed index.html
var indexHTML string

var indexTmpl = template.Must(template.New("index").Parse(indexHTML))

type logLevelOption struct {
	Value string
	Label string
}

var logLevels = []logLevelOption{
	{Value: "debug", Label: "DEBUG"},
	{Value: "info", Label: "INFO"},
	{Value: "warn", Label: "WARN"},
	{Value: "error", Label: "ERROR"},
}

type page struct {
	Path        string
	Description string
}

type indexPage struct {
	pages []page
}

func newIndexPage() *indexPage {
	return &indexPage{}
}

func (i *indexPage) addPage(path string, description string) {
	i.pages = append(i.pages, page{
		Path:        path,
		Description: description,
	})
}

func (i indexPage) handler(eCtx echo.Context) error {
	eCtx.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTMLCharsetUTF8)
	eCtx.Response().WriteHeader(http.StatusOK)

	return indexTmpl.Execute(eCtx.Response(), struct {
		Pages     []page
		LogLevels []logLevelOption
		LogLevel  string
	}{
		Pages:     i.pages,
		LogLevels: logLevels,
		LogLevel:  logger.Level().Level().String(),
	})
}
