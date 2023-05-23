package URL_Shortener

import (
	"gopkg.in/yaml.v2"
	"io"
	"net/http"
	"os"
)

func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if dest, ok := pathsToUrls[path]; ok {
			http.Redirect(w, r, dest, http.StatusFound)
			return
		}
		fallback.ServeHTTP(w, r)
	}
}

func YAMLHandler(yamlFile *os.File, fallback http.Handler) (http.HandlerFunc, error) {
	decoder := yaml.NewDecoder(yamlFile)
	pathUrls, err := decode(decoder)
	if err != nil {
		return nil, err
	}
	pathsToUrls := buildMap(pathUrls)
	return MapHandler(pathsToUrls, fallback), nil
}

type decoder interface {
	Decode(v interface{}) error
}

func decode(d decoder) ([]pathUrl, error) {
	var pu []pathUrl
	for {
		err := d.Decode(&pu)
		if err == io.EOF {
			return pu, nil
		} else if err != nil {
			return nil, err
		}
	}
}

func buildMap(pathUrls []pathUrl) map[string]string {
	pathsToUrls := make(map[string]string)
	for _, pu := range pathUrls {
		pathsToUrls[pu.Path] = pu.URL
	}
	return pathsToUrls
}

func parseYaml(data []byte) ([]pathUrl, error) {
	var pathUrls []pathUrl
	err := yaml.Unmarshal(data, &pathUrls)
	if err != nil {
		return nil, err
	}
	return pathUrls, nil
}

type pathUrl struct {
	Path string `yaml:"path"`
	URL  string `yaml:"url"`
}
