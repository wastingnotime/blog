package site

import (
	"bytes"
	"fmt"
	"html/template"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/adrg/frontmatter"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"gopkg.in/yaml.v3"
)

type Saga struct {
	Title         string
	Slug          string
	Summary       string
	Emoji         string
	Tags          []string
	Arcs          []*Arc
	EpisodeCount  int
	LastRelease   *time.Time
	Status        string
	Repo          string
	RSS           string
	FirstEpisode  *EpisodeRef
	LatestEpisode *EpisodeRef
}

type Arc struct {
	Title        string
	Slug         string
	Summary      string
	Emoji        string
	Episodes     []*Episode
	EpisodeCount int
	LastRelease  *time.Time
}

type Episode struct {
	Title       string    `yaml:"title"`
	Slug        string    `yaml:"slug"`
	Number      int       `yaml:"number"`
	Date        time.Time `yaml:"date"`
	Summary     string    `yaml:"summary"`
	Tags        []string  `yaml:"tags"`
	ReadingTime string    `yaml:"reading_time"`

	SagaSlug string
	ArcSlug  string
	BodyHTML template.HTML
}

type EpisodeRef struct {
	Title     string
	Slug      string
	Number    int
	Date      time.Time
	Summary   string
	ArcSlug   string
	ArcTitle  string
	SagaSlug  string
	SagaTitle string
}

type Post struct {
	Title     string
	Saga      string
	Arc       string
	Type      string
	Studio    string
	Tags      []string
	Summary   string
	Date      time.Time
	Permalink string
}

// templateHTML is a minimal wrapper so html/template treats it as safe when you know it's safe.
type templateHTML string

func (h templateHTML) String() string { return string(h) }

// Load walks a content root like:
// content/sagas/<saga>/index.md
// content/sagas/<saga>/<arc>/index.md
// content/sagas/<saga>/<arc>/ep-*.md
func Load(contentRoot string) ([]*Saga, []*EpisodeRef, error) {
	sagas := map[string]*Saga{}
	var latest []*EpisodeRef
	//md := goldmark.New()
	md := goldmark.New(
		goldmark.WithExtensions(
			highlighting.NewHighlighting(
				highlighting.WithStyle("dracula"),
			),
		),
	)

	err := filepath.WalkDir(contentRoot, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if !strings.HasSuffix(d.Name(), ".md") {
			return nil
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		var fm map[string]any
		body, err := frontmatter.Parse(bytes.NewReader(data), &fm)
		if err != nil {
			return err
		}

		// Determine where we are by path depth
		rel, _ := filepath.Rel(contentRoot, path)
		parts := strings.Split(rel, string(os.PathSeparator))
		// expect: sagas/<saga>/<arc>/file.md or sagas/<saga>/index.md
		if len(parts) < 2 || parts[0] != "sagas" {
			return nil
		}
		sagaSlug := parts[1]
		s := ensureSaga(sagas, sagaSlug, fm)

		if len(parts) == 3 && parts[2] == "index.md" {
			// saga index
			applySagaFrontmatter(s, fm)
			return nil
		}
		if len(parts) >= 4 {
			arcSlug := parts[2]
			a := ensureArc(s, arcSlug, fm)
			if parts[3] == "index.md" {
				applyArcFrontmatter(a, fm)
				return nil
			}
			// episode file
			ep := &Episode{SagaSlug: sagaSlug, ArcSlug: arcSlug, Slug: strings.TrimSuffix(filepath.Base(path), ".md")}
			applyEpisodeFrontmatter(ep, fm)

			var buf bytes.Buffer
			if err := md.Convert(body, &buf); err != nil {
				return err
			}
			ep.BodyHTML = template.HTML(buf.String())

			a.Episodes = append(a.Episodes, ep)
			s.EpisodeCount++

			ref := &EpisodeRef{Title: ep.Title, Slug: ep.Slug, Number: ep.Number, Date: ep.Date, Summary: ep.Summary, ArcSlug: arcSlug, ArcTitle: a.Title, SagaSlug: sagaSlug, SagaTitle: s.Title}
			latest = append(latest, ref)

			// last release tracking
			if s.LastRelease == nil || ep.Date.After(*s.LastRelease) {
				tmp := ep.Date
				s.LastRelease = &tmp
				s.LatestEpisode = ref
			}
			if a.LastRelease == nil || ep.Date.After(*a.LastRelease) {
				tmp := ep.Date
				a.LastRelease = &tmp
			}
		}
		return nil
	})
	if err != nil {
		return nil, nil, err
	}

	// ensure episode refs have final saga/arc titles
	for _, ref := range latest {
		if s, ok := sagas[ref.SagaSlug]; ok {
			if s.Title != "" {
				ref.SagaTitle = s.Title
			}
			for _, a := range s.Arcs {
				if a.Slug == ref.ArcSlug {
					if a.Title != "" {
						ref.ArcTitle = a.Title
					}
					break
				}
			}
		}
	}

	// sort arcs & episodes
	out := make([]*Saga, 0, len(sagas))
	for _, s := range sagas {
		for _, a := range s.Arcs {
			sort.SliceStable(a.Episodes, func(i, j int) bool {
				if a.Episodes[i].Number != 0 && a.Episodes[j].Number != 0 {
					return a.Episodes[i].Number < a.Episodes[j].Number
				}
				return a.Episodes[i].Date.Before(a.Episodes[j].Date)
			})
			a.EpisodeCount = len(a.Episodes)
			if a.EpisodeCount > 0 {
				s.FirstEpisode = &EpisodeRef{Title: a.Episodes[0].Title, Slug: a.Episodes[0].Slug, Number: a.Episodes[0].Number, Date: a.Episodes[0].Date, Summary: a.Episodes[0].Summary, ArcSlug: a.Slug, ArcTitle: a.Title, SagaSlug: s.Slug, SagaTitle: s.Title}
			}
		}
		s.Status = SagaStatus(s.LastRelease)
		out = append(out, s)
		sort.SliceStable(s.Arcs, func(i, j int) bool { return s.Arcs[i].Title < s.Arcs[j].Title })
	}

	// sort latest desc
	sort.SliceStable(latest, func(i, j int) bool { return latest[i].Date.After(latest[j].Date) })
	// pick current arc for each saga = most recent episode's arc
	for _, s := range out {
		var cur *EpisodeRef
		for _, ref := range latest {
			//hack: in case of a declared saga, with declared arcs, but without episodies, the last release is nil
			if s.LastRelease != nil && ref.ArcSlug != "" && ref.Date.Equal(*s.LastRelease) && s.Slug == ref.ArcSlug { // placeholder; we’ll re-derive per-saga if needed
				//if ref.ArcSlug != "" && ref.Date.Equal(*s.LastRelease) && s.Slug == ref.ArcSlug { // placeholder; we’ll re-derive per-saga if needed
				cur = ref
				break
			}
		}
		_ = cur // optional: compute .CurrentArc at render time by scanning arcs for latest date
	}

	// final sort sagas by last release desc
	sort.SliceStable(out, func(i, j int) bool {
		li, lj := out[i].LastRelease, out[j].LastRelease
		if li == nil {
			return false
		}
		if lj == nil {
			return true
		}
		return li.After(*lj)
	})

	return out, latest, nil
}

func collectPosts(contentDir string) ([]Post, error) {
	var posts []Post
	err := filepath.WalkDir(contentDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(d.Name(), ".md") {
			return nil
		}
		p, err := parseFrontmatter(path)
		if err != nil {
			// skip silently or log
			return nil
		}
		posts = append(posts, p)
		return nil
	})
	return posts, err
}

func BuildRecent(contentDir string, limit int) ([]Post, error) {
	if _, err := os.Stat(contentDir); err != nil {
		if os.IsNotExist(err) {
			return []Post{}, nil
		}
		return nil, err
	}
	all, err := collectPosts(contentDir)
	if err != nil {
		return nil, err
	}
	sort.Slice(all, func(i, j int) bool {
		return all[i].Date.After(all[j].Date)
	})
	if len(all) > limit {
		all = all[:limit]
	}
	return all, nil
}

func ensureSaga(m map[string]*Saga, slug string, fm map[string]any) *Saga {
	if s, ok := m[slug]; ok {
		return s
	}
	s := &Saga{Slug: slug}
	if v, _ := fm["sagaTitle"].(string); v != "" {
		s.Title = v
	}
	m[slug] = s
	return s
}

func ensureArc(s *Saga, slug string, fm map[string]any) *Arc {
	for _, a := range s.Arcs {
		if a.Slug == slug {
			return a
		}
	}
	a := &Arc{Slug: slug}
	s.Arcs = append(s.Arcs, a)
	return a
}

func applySagaFrontmatter(s *Saga, fm map[string]any) {
	if v, _ := fm["title"].(string); v != "" {
		s.Title = v
	}
	if v, _ := fm["summary"].(string); v != "" {
		s.Summary = v
	}
	if v, _ := fm["emoji"].(string); v != "" {
		s.Emoji = v
	}
	if v, _ := fm["tags"].([]any); len(v) > 0 {
		s.Tags = toStrings(v)
	}
	if v, _ := fm["repo"].(string); v != "" {
		s.Repo = v
	}
	if v, _ := fm["rss"].(string); v != "" {
		s.RSS = v
	}
}

func applyArcFrontmatter(a *Arc, fm map[string]any) {
	if v, _ := fm["title"].(string); v != "" {
		a.Title = v
	}
	if v, _ := fm["summary"].(string); v != "" {
		a.Summary = v
	}
	if v, _ := fm["emoji"].(string); v != "" {
		a.Emoji = v
	}
}

func applyEpisodeFrontmatter(e *Episode, fm map[string]any) {
	if v, _ := fm["title"].(string); v != "" {
		e.Title = v
	}
	if v, _ := fm["slug"].(string); v != "" {
		e.Slug = v
	}
	if v, ok := fm["number"].(int); ok {
		e.Number = v
	}
	if v, ok := fm["date"].(string); ok {
		theTime, err := time.Parse("2006-01-02", v)
		if err != nil {
			fmt.Println("Could not parse time:", err)
		}
		e.Date = theTime
	}
	if v, _ := fm["summary"].(string); v != "" {
		e.Summary = v
	}
	if v, _ := fm["tags"].([]any); len(v) > 0 {
		e.Tags = toStrings(v)
	}
	if v, _ := fm["reading_time"].(string); v != "" {
		e.ReadingTime = v
	}
}

func toStrings(xs []any) []string {
	out := make([]string, 0, len(xs))
	for _, x := range xs {
		if s, ok := x.(string); ok {
			out = append(out, s)
		}
	}
	return out
}

func parseFrontmatter(path string) (Post, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Post{}, err
	}

	content := string(data)
	start := strings.Index(content, "---")
	if start != 0 {
		return Post{}, fmt.Errorf("no frontmatter in %s", path)
	}
	end := strings.Index(content[3:], "---")
	if end == -1 {
		return Post{}, fmt.Errorf("no closing --- in %s", path)
	}
	fm := content[3 : end+3]
	var meta map[string]any
	if err := yaml.Unmarshal([]byte(fm), &meta); err != nil {
		return Post{}, err
	}

	var p Post
	if v, _ := meta["title"].(string); v != "" {
		p.Title = v
	}
	if v, _ := meta["saga"].(string); v != "" {
		p.Saga = v
	}
	if v, _ := meta["arc"].(string); v != "" {
		p.Arc = v
	}
	if v, _ := meta["type"].(string); v != "" {
		p.Type = v
	}
	if v, _ := meta["studio"].(string); v != "" {
		p.Studio = v
	}
	if v, _ := meta["summary"].(string); v != "" {
		p.Summary = v
	}
	if v, ok := meta["tags"].([]any); ok && len(v) > 0 {
		p.Tags = toStrings(v)
	}
	switch v := meta["date"].(type) {
	case string:
		if t, err := time.Parse(time.RFC3339, v); err == nil {
			p.Date = t
		} else if t, err := time.Parse("2006-01-02", v); err == nil {
			p.Date = t
		}
	case time.Time:
		p.Date = v
	}

	// Compute permalink
	if v, _ := meta["permalink"].(string); v != "" {
		p.Permalink = v
	} else {
		if rel, err := filepath.Rel("content", path); err == nil {
			dir := filepath.Dir(rel)
			if dir == "." {
				dir = ""
			}
			name := strings.TrimSuffix(filepath.Base(rel), ".md")
			joined := filepath.Join(dir, name)
			p.Permalink = "/" + filepath.ToSlash(joined)
		} else {
			p.Permalink = "/"
		}
	}
	if p.Type == "" {
		p.Type = "Post"
	}
	if p.Date.IsZero() {
		if info, err := os.Stat(path); err == nil {
			p.Date = info.ModTime()
		}
	}
	return p, nil
}
