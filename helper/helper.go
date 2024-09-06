package helper

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Helper struct {
	env        string
	prefixPath string
}

func New(env string, contextPath string) *Helper {
	prefixPath := ""

	if contextPath != "/" {
		prefixPath = contextPath
	}

	helper := Helper{
		env:        env,
		prefixPath: prefixPath,
	}

	return &helper
}

func (hp *Helper) RoutePath(path string) string {
	return fmt.Sprintf("%s/%v", hp.prefixPath, path)
}

func (hp *Helper) AssetPath(fileName string) string {
	return fmt.Sprintf("%s/assets/%v", hp.prefixPath, fileName)
}

// URLParam returns the url parameter from a http.Request object.
func (hp *Helper) URLParam(r *http.Request, key string) string {
	value := chi.URLParam(r, key)

	return value
}
