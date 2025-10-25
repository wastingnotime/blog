package main

import (
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s <command>\n", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), "\nCommands:\n")
		fmt.Fprintf(flag.CommandLine.Output(), "  check    Validate markdown frontmatter\n")
	}

	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		flag.Usage()
		os.Exit(2)
	}

	switch args[0] {
	case "check":
		if err := runCheck("content"); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	default:
		flag.Usage()
		os.Exit(2)
	}
}

func runCheck(root string) error {
	var issues []string

	walkErr := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			issues = append(issues, fmt.Sprintf("%s: %v", path, err))
			return nil
		}
		if entry.IsDir() {
			return nil
		}
		if filepath.Ext(path) != ".md" {
			return nil
		}

		rel, relErr := filepath.Rel(root, path)
		if relErr != nil {
			rel = path
		}

		fm, parseErr := parseFrontmatter(path)
		if parseErr != nil {
			issues = append(issues, fmt.Sprintf("%s: %v", rel, parseErr))
			return nil
		}

		kind := classifyContent(rel)
		fileIssues := validateFrontmatter(kind, fm)
		for _, msg := range fileIssues {
			issues = append(issues, fmt.Sprintf("%s: %s", rel, msg))
		}
		return nil
	})
	if walkErr != nil {
		return walkErr
	}

	if len(issues) > 0 {
		for _, msg := range issues {
			fmt.Fprintln(os.Stderr, msg)
		}
		return fmt.Errorf("validation failed (%d issue(s))", len(issues))
	}
	return nil
}

func parseFrontmatter(path string) (map[string]any, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	content := string(data)
	if !strings.HasPrefix(content, "---") {
		return nil, errors.New("missing opening frontmatter delimiter '---'")
	}

	rest := strings.TrimPrefix(content, "---")
	rest = strings.TrimPrefix(rest, "\r\n")
	rest = strings.TrimPrefix(rest, "\n")

	end := strings.Index(rest, "---")
	if end == -1 {
		return nil, errors.New("missing closing frontmatter delimiter '---'")
	}

	fmSection := rest[:end]
	var meta map[string]any
	if err := yaml.Unmarshal([]byte(fmSection), &meta); err != nil {
		return nil, fmt.Errorf("unable to parse frontmatter: %w", err)
	}
	return meta, nil
}

func classifyContent(rel string) string {
	parts := strings.Split(rel, string(os.PathSeparator))
	if len(parts) == 0 {
		return ""
	}

	switch parts[0] {
	case "posts":
		return "post"
	case "sagas":
		if len(parts) >= 3 && parts[2] == "index.md" {
			return "saga"
		}
		if len(parts) >= 4 && parts[3] == "index.md" {
			return "arc"
		}
		return "episode"
	case "library":
		return "library"
	default:
		return ""
	}
}

func validateFrontmatter(kind string, fm map[string]any) []string {
	switch kind {
	case "post":
		return validatePost(fm)
	case "saga":
		return validateSaga(fm)
	case "arc":
		return validateArc(fm)
	case "episode":
		return validateEpisode(fm)
	case "library":
		return validateLibrary(fm)
	default:
		return nil
	}
}

func validatePost(fm map[string]any) []string {
	var errs []string
	if !hasNonEmptyString(fm, "title") {
		errs = append(errs, "missing required field 'title'")
	}
	if !hasNonEmptyString(fm, "summary") {
		errs = append(errs, "missing required field 'summary'")
	}
	if v, ok := fm["date"]; !ok {
		errs = append(errs, "missing required field 'date'")
	} else if err := validateDate(v); err != nil {
		errs = append(errs, fmt.Sprintf("invalid 'date': %v", err))
	}
	if v, ok := fm["slug"]; ok {
		if !isString(v) {
			errs = append(errs, "'slug' must be a string")
		} else if strings.TrimSpace(v.(string)) == "" {
			errs = append(errs, "'slug' cannot be empty")
		}
	}
	if v, ok := fm["reading_time"]; ok {
		if !isString(v) {
			errs = append(errs, "'reading_time' must be a string")
		}
	}
	if v, ok := fm["tags"]; ok {
		if err := validateStringList(v); err != nil {
			errs = append(errs, fmt.Sprintf("'tags' %v", err))
		}
	}
	return errs
}

func validateSaga(fm map[string]any) []string {
	var errs []string
	if !hasNonEmptyString(fm, "title") {
		errs = append(errs, "missing required field 'title'")
	}
	if v, ok := fm["slug"]; ok {
		if !isString(v) {
			errs = append(errs, "'slug' must be a string")
		} else if strings.TrimSpace(v.(string)) == "" {
			errs = append(errs, "'slug' cannot be empty")
		}
	}
	if !hasNonEmptyString(fm, "summary") {
		errs = append(errs, "missing required field 'summary'")
	}
	if v, ok := fm["order"]; ok {
		if _, err := toInt(v); err != nil {
			errs = append(errs, "'order' must be an integer")
		}
	} else {
		errs = append(errs, "missing required field 'order'")
	}
	if v, ok := fm["links"]; ok {
		if err := validateLinks(v); err != nil {
			errs = append(errs, fmt.Sprintf("'links' %v", err))
		}
	}
	return errs
}

func validateArc(fm map[string]any) []string {
	var errs []string
	if !hasNonEmptyString(fm, "title") {
		errs = append(errs, "missing required field 'title'")
	}
	if v, ok := fm["slug"]; ok {
		if !isString(v) {
			errs = append(errs, "'slug' must be a string")
		} else if strings.TrimSpace(v.(string)) == "" {
			errs = append(errs, "'slug' cannot be empty")
		}
	}
	if !hasNonEmptyString(fm, "summary") {
		errs = append(errs, "missing required field 'summary'")
	}
	if v, ok := fm["order"]; ok {
		if _, err := toInt(v); err != nil {
			errs = append(errs, "'order' must be an integer")
		}
	} else {
		errs = append(errs, "missing required field 'order'")
	}
	return errs
}

func validateEpisode(fm map[string]any) []string {
	var errs []string
	if !hasNonEmptyString(fm, "title") {
		errs = append(errs, "missing required field 'title'")
	}
	if v, ok := fm["number"]; !ok {
		errs = append(errs, "missing required field 'number'")
	} else if n, err := toInt(v); err != nil {
		errs = append(errs, "'number' must be an integer")
	} else if n <= 0 {
		errs = append(errs, "'number' must be greater than 0")
	}
	if v, ok := fm["date"]; !ok {
		errs = append(errs, "missing required field 'date'")
	} else if err := validateDate(v); err != nil {
		errs = append(errs, fmt.Sprintf("invalid 'date': %v", err))
	}
	if !hasNonEmptyString(fm, "summary") {
		errs = append(errs, "missing required field 'summary'")
	}
	if v, ok := fm["reading_time"]; ok {
		if !isString(v) {
			errs = append(errs, "'reading_time' must be a string")
		}
	}
	if v, ok := fm["tags"]; ok {
		if err := validateStringList(v); err != nil {
			errs = append(errs, fmt.Sprintf("'tags' %v", err))
		}
	}
	if v, ok := fm["slug"]; ok {
		if !isString(v) {
			errs = append(errs, "'slug' must be a string")
		} else if strings.TrimSpace(v.(string)) == "" {
			errs = append(errs, "'slug' cannot be empty")
		}
	}
	return errs
}

func validateLibrary(fm map[string]any) []string {
	var errs []string
	if !hasNonEmptyString(fm, "title") {
		errs = append(errs, "missing required field 'title'")
	}
	if !hasNonEmptyString(fm, "summary") {
		errs = append(errs, "missing required field 'summary'")
	}
	if v, ok := fm["tags"]; ok {
		if err := validateStringList(v); err != nil {
			errs = append(errs, fmt.Sprintf("'tags' %v", err))
		}
	}
	return errs
}

func hasNonEmptyString(fm map[string]any, key string) bool {
	v, ok := fm[key]
	if !ok {
		return false
	}
	s, ok := v.(string)
	if !ok {
		return false
	}
	return strings.TrimSpace(s) != ""
}

func isString(v any) bool {
	_, ok := v.(string)
	return ok
}

func toInt(v any) (int, error) {
	switch n := v.(type) {
	case int:
		return n, nil
	case int64:
		return int(n), nil
	case float64:
		if n != float64(int(n)) {
			return 0, errors.New("not an integer")
		}
		return int(n), nil
	case uint64:
		return int(n), nil
	default:
		return 0, errors.New("not an integer")
	}
}

func validateDate(v any) error {
	s, ok := v.(string)
	if !ok {
		return errors.New("must be a string")
	}
	s = strings.TrimSpace(s)
	if s == "" {
		return errors.New("cannot be empty")
	}
	if _, err := time.Parse(time.RFC3339, s); err == nil {
		return nil
	}
	if _, err := time.Parse("2006-01-02", s); err == nil {
		return nil
	}
	return fmt.Errorf("must be RFC3339 or YYYY-MM-DD formatted (got %q)", s)
}

func validateStringList(v any) error {
	list, ok := v.([]any)
	if !ok {
		return errors.New("must be a list of strings")
	}
	for i, item := range list {
		if _, ok := item.(string); !ok {
			return fmt.Errorf("must contain only strings (item %d)", i)
		}
	}
	return nil
}

func validateLinks(v any) error {
	m, ok := v.(map[string]any)
	if !ok {
		return errors.New("must be a map of strings")
	}
	for key, value := range m {
		if !isString(value) {
			return fmt.Errorf("key %q must be a string", key)
		}
	}
	return nil
}
