package site

import (
	"os"
	"path"
	"strings"
)

type Config struct {
	BasePath string
}

func ConfigFromEnv() Config {
	base := os.Getenv("SITE_BASE_PATH")
	if base == "" {
		base = os.Getenv("BASE_PATH")
	}
	return Config{BasePath: normalizeBasePath(base)}
}

func (c Config) Href(p string) string {
	return buildHref(c.BasePath, p)
}

func normalizeBasePath(in string) string {
	in = strings.TrimSpace(in)
	if in == "" {
		return ""
	}
	if !strings.HasPrefix(in, "/") {
		in = "/" + in
	}
	cleaned := path.Clean(in)
	if cleaned == "." || cleaned == "/" {
		return ""
	}
	return strings.TrimSuffix(cleaned, "/")
}

func buildHref(basePath, target string) string {
	target = strings.TrimSpace(target)
	if target == "" {
		target = "/"
	}
	if !strings.HasPrefix(target, "/") {
		target = "/" + target
	}
	if basePath == "" {
		return target
	}
	if target == "/" {
		return basePath + "/"
	}
	return basePath + target
}
