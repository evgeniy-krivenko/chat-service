package serverdebug

import (
	"html/template"
	"strings"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

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
	funcMap := template.FuncMap{
		"ToUpper": strings.ToUpper,
	}

	return template.Must(template.New("index").Funcs(funcMap).Parse(`<html>
	<title>Chat Service Debug</title>
<body>
	<h2>Chat Service Debug</h2>
	<ul>
		{{range .Pages}}
			<li>
				<a href="{{ .Path }}">{{ .Path}} {{ .Description }}</a>
			</li>
		{{end}}
	</ul>

	<h2>Log Level</h2>
	<form onSubmit="putLogLevel()">
		<select id="log-level-select">
			{{range $level := .LogLevels}}
				<option value="{{$level}}" {{if eq $level $.LogLevel}} selected {{end}}>
					{{- $level | ToUpper -}}
				</option>
			{{end}}
    	</select>
		<input type="submit" value="Change"></input>
	</form>
	
	<script>
		function putLogLevel() {
			const req = new XMLHttpRequest();
			req.open('PUT', '/log/level', false);
			req.setRequestHeader('Content-Type', 'application/x-www-form-urlencoded')
			req.onload = function() { window.location.reload(); };
			req.send('level='+document.getElementById('log-level-select').value);
		};
	</script>
</body>
</html>
`)).Execute(eCtx.Response(), struct {
		Pages     []page
		LogLevel  string
		LogLevels []string
	}{
		Pages: i.pages,
		LogLevels: []string{
			zap.DebugLevel.String(),
			zap.InfoLevel.String(),
			zap.WarnLevel.String(),
			zap.ErrorLevel.String(),
		},
		LogLevel: zap.L().Level().String(),
	})
}
