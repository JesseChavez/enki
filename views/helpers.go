package views

import (
	"fmt"
	"html/template"
)

var funcMap = template.FuncMap{
	"stylesheetPath": stylesheetPath,
	"javascriptPath": javascriptPath,
}

func stylesheetPath(name string) string {
	if envRender != "development" {
		return fmt.Sprintf("%v-123456789.css", name)
	}

	return fmt.Sprintf("%v.css", name)
}

func javascriptPath(name string) string {
	if envRender != "development" {
		return fmt.Sprintf("%v-1234567890.js", name)
	}

	return fmt.Sprintf("%v.js", name)
}

func assetPath(name string) string {
	if envRender != "development" {
		return fmt.Sprintf("%v", name)
	}

	return fmt.Sprintf("%v", name)
}
