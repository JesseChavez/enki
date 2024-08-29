package views

import (
	"fmt"
	"html/template"
)

var funcMap = template.FuncMap{
	"stylesheetPath": stylesheetPath,
	"javascriptPath": javascriptPath,
}

func stylesheetPath(assetName string) string {
	key := fmt.Sprintf("%v.css", assetName)

	filePath := Manifest[key]

	return fmt.Sprintf("%s/assets/%v", prefixPath, filePath)
}

func javascriptPath(assetName string) string {
	key := fmt.Sprintf("%v.js", assetName)

	filePath := Manifest[key]

	return fmt.Sprintf("%s/assets/%v", prefixPath, filePath)
}

func assetPath(assetName string) string {
	key := fmt.Sprintf("%v", assetName)

	filePath := Manifest[key]

	return fmt.Sprintf("%s/assets/%v", prefixPath, filePath)
}
