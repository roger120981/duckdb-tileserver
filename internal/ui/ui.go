package ui

import (
	"bytes"
	"html/template"

	"github.com/tobilg/duckdb-tileserver/internal/conf"
)

// HTMLDynamicLoad sets whether HTML templates are loaded every time they are used
// this allows rapid prototyping of changes
var HTMLDynamicLoad bool

// RenderHTML tbd
func RenderHTML(temp *template.Template, content interface{}, context interface{}) ([]byte, error) {
	bodyData := map[string]interface{}{
		"config":  conf.Configuration,
		"context": context,
		"data":    content}
	contentBytes, err := renderTemplate(temp, bodyData)
	if err != nil {
		return contentBytes, err
	}
	return contentBytes, err
}

func renderTemplate(temp *template.Template, data map[string]interface{}) ([]byte, error) {
	var buf bytes.Buffer

	if err := temp.ExecuteTemplate(&buf, "page", data); err != nil {
		return buf.Bytes(), err
	}
	return buf.Bytes(), nil
}

// LoadTemplate loads a simple template file (for standalone HTML pages)
func LoadTemplate(filename string) (*template.Template, error) {
	filePath := conf.Configuration.Server.AssetsPath + "/" + filename
	tmpl, err := template.ParseFiles(filePath)
	if err != nil {
		return nil, err
	}
	return tmpl, nil
}
