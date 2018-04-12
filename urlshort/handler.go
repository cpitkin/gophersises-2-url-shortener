package urlshort

import (
	"net/http"

	yamlv2 "gopkg.in/yaml.v2"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if redURL, ok := pathsToUrls[r.URL.Path]; ok {
			http.Redirect(w, r, redURL, 301)
		} else {
			fallback.ServeHTTP(w, r)
		}
	}
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {

	parsedYaml, err := parseYaml(yml)
	if err != nil {
		return nil, err
	}
	m := toMap(parsedYaml)
	return MapHandler(m, fallback), nil
}

func parseYaml(yml []byte) (dst []map[string]string, err error) {
	// var pathsToUrls pathData
	// err := yaml.Unmarshal([]byte(yml), &pathsToUrls)
	// if err != nil {
	// 	log.Fatalf("%v", err)
	// }
	// return structs.Map(pathsToUrls)
	err = yamlv2.Unmarshal(yml, &dst)
	return dst, err
}

func toMap(pyaml []map[string]string) map[string]string {
	m := make(map[string]string)
	for _, v := range pyaml {
		m[v["path"]] = v["url"]
	}
	return m
}
