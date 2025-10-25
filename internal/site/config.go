package site

import (
	"os"
	"path"
	"strings"
)

type Config struct {
	BasePath string
	BaseURL  string
}

func ConfigFromEnv() Config {
	base := os.Getenv("SITE_BASE_PATH")
	if base == "" {
		base = os.Getenv("BASE_PATH")
	}
	baseURL := os.Getenv("SITE_BASE_URL")
	if baseURL == "" {
		baseURL = os.Getenv("BASE_URL")
	}
	return Config{
		BasePath: normalizeBasePath(base),
		BaseURL:  normalizeBaseURL(baseURL),
	}
}

func (c Config) Href(p string) string {
	return buildHref(c.BasePath, p)
}

func (c Config) AbsoluteURL(p string) string {
	href := c.Href(p)
	if c.BaseURL == "" {
		return href
	}
	if strings.HasPrefix(href, "http://") || strings.HasPrefix(href, "https://") {
		return href
	}
	if href == "" {
		return c.BaseURL
	}
	if strings.HasPrefix(href, "/") {
		return c.BaseURL + href
	}
	return c.BaseURL + "/" + href
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

func normalizeBaseURL(in string) string {
	in = strings.TrimSpace(in)
	if in == "" {
		return ""
	}
	return strings.TrimRight(in, "/")
}
