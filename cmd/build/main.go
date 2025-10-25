package main

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/wastingnotime/blog/internal/site"
)

// ---- Data fed to templates ----

type HomeSaga struct {
	Title        string
	Slug         string
	Summary      string
	Emoji        string
	EpisodeCount int
	LastRelease  *time.Time
	Status       string
}

type HomeEpisode struct {
	Title     string
	Slug      string
	Number    int
	Date      time.Time
	Summary   string
	SagaTitle string
	SagaSlug  string
	ArcSlug   string
	ArcTitle  string
	ArcEmoji  string
}

type HomeRecent struct {
	Title     string
	Summary   string
	Date      time.Time
	Type      string
	Saga      string
	Arc       string
	Permalink string
	Tags      []string
}

type SagaSummary struct {
	Title       string
	Description template.HTML
	SagaURL     string
	StartURL    string
}

const homeRecentLimit = 10
const feedItemLimit = 20

// episode header payload used by templates/episode.tmpl
// (we pass it as a map[string]any for flexibility)

// loadBase loads the base layout with the site's helper funcs applied.
func loadBase(cfg site.Config) *template.Template {
	return template.Must(template.New("base").Funcs(site.Funcs(cfg)).ParseFiles("templates/base.gohtml"))
}

// renderView executes a named view into the target output file.
func renderView(base *template.Template, viewName string, outPath string, data any, viewFiles ...string) {
	t := template.Must(base.Clone())
	t = template.Must(t.ParseFiles(viewFiles...))
	if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
		panic(err)
	}
	f, err := os.Create(outPath)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if err := t.ExecuteTemplate(f, viewName, data); err != nil {
		panic(err)
	}
}

func buildSagaSummaries(contentRoot string, sagas []*site.Saga) ([]SagaSummary, error) {
	summaries := make([]SagaSummary, 0, len(sagas))
	for _, saga := range sagas {
		pagePath := filepath.Join(contentRoot, "sagas", saga.Slug, "index.md")
		page, err := site.LoadPage(pagePath)
		if err != nil {
			if !errors.Is(err, os.ErrNotExist) {
				return nil, fmt.Errorf("load saga index %s: %w", saga.Slug, err)
			}
			page = nil
		}

		summary := SagaSummary{
			SagaURL: fmt.Sprintf("/sagas/%s/", saga.Slug),
		}

		if page != nil {
			summary.Title = page.Title
			summary.Description = page.Body
		}

		if summary.Title == "" {
			summary.Title = saga.Title
		}

		if len(summary.Description) == 0 && saga.Summary != "" {
			escaped := template.HTMLEscapeString(saga.Summary)
			summary.Description = template.HTML("<p>" + escaped + "</p>")
		}

		startURL, err := firstArcPermalink(contentRoot, saga.Slug)
		if err != nil {
			return nil, fmt.Errorf("find first arc for %s: %w", saga.Slug, err)
		}
		summary.StartURL = startURL

		summaries = append(summaries, summary)
	}
	return summaries, nil
}

func firstArcPermalink(contentRoot, sagaSlug string) (string, error) {
	sagaDir := filepath.Join(contentRoot, "sagas", sagaSlug)
	entries, err := os.ReadDir(sagaDir)
	if err != nil {
		return "", err
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		arcIndex := filepath.Join(sagaDir, entry.Name(), "index.md")
		if _, err := os.Stat(arcIndex); err == nil {
			return fmt.Sprintf("/sagas/%s/%s/", sagaSlug, entry.Name()), nil
		} else if !os.IsNotExist(err) {
			return "", err
		}
	}
	return "", nil
}

// ---- main ----

func main() {
	cfg := site.ConfigFromEnv()
	// 1) Load content (sagas, arcs, episodes as parsed structs)
	sagas, latest, err := site.Load("content")
	if err != nil {
		log.Fatalf("load: %v", err)
	}

	sagaSummaries, err := buildSagaSummaries("content", sagas)
	if err != nil {
		log.Fatalf("build saga summaries: %v", err)
	}

	posts, err := site.LoadPosts(filepath.Join("content", "posts"))
	if err != nil {
		log.Fatalf("load posts: %v", err)
	}
	recentPosts := posts
	if len(recentPosts) > homeRecentLimit {
		recentPosts = recentPosts[:homeRecentLimit]
	}

	base := loadBase(cfg)
	nowYear := time.Now().Year()

	aboutPage, err := site.LoadPage(filepath.Join("content", "about", "index.md"))
	if err != nil {
		log.Fatalf("load about page: %v", err)
	}

	tagIndex, libraryTags := site.BuildTagIndex(sagas, posts)

	// 2) Render Home
	homeData := map[string]any{
		"Section":        "home",
		"Sagas":          toHomeSagas(sagas),
		"LatestEpisodes": toHomeLatest(latest),
		"RecentPosts":    toHomeRecent(latest, recentPosts, homeRecentLimit),
		"NowYear":        nowYear,
	}

	//write(t, "home", "public/index.html", homeData)
	renderView(base,
		"home",
		"public/index.html",
		homeData,
		"templates/home.gohtml",
	)

	// 3) Render Library
	libraryData := map[string]any{
		"Section": "library",
		"Tags":    libraryTags,
		"NowYear": nowYear,
	}

	renderView(base,
		"library_index",
		"public/library/index.html",
		libraryData,
		"templates/library_index.gohtml",
	)

	for _, tag := range libraryTags {
		pageData := map[string]any{
			"Section": "library",
			"Tag":     tag.Name,
			"Posts":   tagIndex[tag.Name],
			"NowYear": nowYear,
		}
		renderView(base,
			"tag",
			filepath.Join("public", "library", tag.Slug, "index.html"),
			pageData,
			"templates/tag.gohtml",
		)
	}

	aboutData := map[string]any{
		"Section": "about",
		"Page":    aboutPage,
		"NowYear": nowYear,
	}

	renderView(base,
		"about",
		"public/about/index.html",
		aboutData,
		"templates/about.gohtml",
	)

	sagasData := map[string]any{
		"Section":       "sagas",
		"SagaSummaries": sagaSummaries,
		"NowYear":       nowYear,
	}

	renderView(base,
		"sagas",
		"public/sagas/index.html",
		sagasData,
		"templates/sagas.gohtml",
	)

	// 5) Render all Saga pages, Arc pages, and Episode pages
	for _, s := range sagas {
		// Sort arcs by last release desc for the saga timeline derivation & spotlight
		arcs := append([]*site.Arc{}, s.Arcs...)
		sort.SliceStable(arcs, func(i, j int) bool {
			li, lj := arcs[i].LastRelease, arcs[j].LastRelease
			if li == nil { // empty goes last
				return false
			}
			if lj == nil {
				return true
			}
			return li.After(*lj)
		})

		// Current arc = arc with most recent episode
		var currentArc *site.Arc
		if len(arcs) > 0 {
			currentArc = arcs[0]
			if currentArc.LastRelease == nil && len(arcs) > 1 {
				// all empty? then just pick first
				currentArc = arcs[0]
			}
		}

		// SagaTimeline = all episodes across arcs sorted newest → oldest
		var timeline []*site.EpisodeRef
		for _, a := range s.Arcs {
			for _, e := range a.Episodes {
				ref := &site.EpisodeRef{Title: e.Title, Slug: e.Slug, Number: e.Number, Date: e.Date, Summary: e.Summary, ArcSlug: a.Slug, ArcTitle: a.Title, SagaSlug: s.Slug, SagaTitle: s.Title}
				timeline = append(timeline, ref)
			}
		}
		sort.SliceStable(timeline, func(i, j int) bool { return timeline[i].Date.After(timeline[j].Date) })

		// Build Arcs lite for listing
		arcsLite := make([]map[string]any, 0, len(s.Arcs))
		for _, a := range s.Arcs {
			m := map[string]any{
				"Title":        a.Title,
				"Slug":         a.Slug,
				"Summary":      a.Summary,
				"Emoji":        a.Emoji,
				"EpisodeCount": len(a.Episodes),
			}
			if a.LastRelease != nil {
				m["LastRelease"] = *a.LastRelease
			}
			arcsLite = append(arcsLite, m)
		}

		// Current Arc compact data
		var currentArcData map[string]any
		if currentArc != nil {
			latestFromArc := lastNEpisodes(currentArc, 3)
			refs := make([]*site.EpisodeRef, 0, len(latestFromArc))
			for _, e := range latestFromArc {
				refs = append(refs, &site.EpisodeRef{Title: e.Title, Slug: e.Slug, Number: e.Number, Date: e.Date, Summary: e.Summary, ArcSlug: currentArc.Slug, ArcTitle: currentArc.Title, SagaSlug: s.Slug, SagaTitle: s.Title})
			}
			currentArcData = map[string]any{
				"Title":          currentArc.Title,
				"Slug":           currentArc.Slug,
				"Summary":        currentArc.Summary,
				"Emoji":          currentArc.Emoji,
				"LatestEpisodes": refs,
			}
		}

		// Saga hero aggregation
		var firstEp, latestEp *site.EpisodeRef
		if s.FirstEpisode != nil {
			firstEp = s.FirstEpisode
		}
		if s.LatestEpisode != nil {
			latestEp = s.LatestEpisode
		}

		sagaHero := map[string]any{
			"Title":        s.Title,
			"Slug":         s.Slug,
			"Summary":      s.Summary,
			"Emoji":        s.Emoji,
			"EpisodeCount": s.EpisodeCount,
			"Status":       s.Status,
			"Repo":         s.Repo,
			"RSS":          s.RSS,
		}
		if s.LastRelease != nil {
			sagaHero["LastRelease"] = *s.LastRelease
		}
		if firstEp != nil {
			sagaHero["FirstEpisode"] = firstEp
		}
		if latestEp != nil {
			sagaHero["LatestEpisode"] = latestEp
		}

		// Render Saga page
		sagaData := map[string]any{
			"Section":      "sagas",
			"Saga":         sagaHero,
			"CurrentArc":   currentArcData,
			"Arcs":         arcsLite,
			"SagaTimeline": timeline,
			"NowYear":      time.Now().Year(),
		}
		//write(t, "saga", filepath.Join("public", "sagas", s.Slug, "index.html"), sagaData)
		renderView(base,
			"saga",
			"public/sagas/"+s.Slug+"/index.html",
			sagaData,
			"templates/saga.gohtml",
		)

		// Render Arc pages and Episode pages
		for ai, a := range s.Arcs {
			// Prev/Next arc references
			var prevArc, nextArc *site.Arc
			if ai > 0 {
				prevArc = s.Arcs[ai-1]
			}
			if ai+1 < len(s.Arcs) {
				nextArc = s.Arcs[ai+1]
			}

			arcLite := func(x *site.Arc) map[string]any {
				if x == nil {
					return nil
				}
				return map[string]any{"Title": x.Title, "Slug": x.Slug}
			}

			// Episodes list for arc (refs)
			epRefs := make([]*site.EpisodeRef, 0, len(a.Episodes))
			for _, e := range a.Episodes {
				epRefs = append(epRefs, &site.EpisodeRef{Title: e.Title, Slug: e.Slug, Number: e.Number, Date: e.Date, Summary: e.Summary, ArcSlug: a.Slug, ArcTitle: a.Title, SagaSlug: s.Slug, SagaTitle: s.Title})
			}

			var firstArcEp, latestArcEp *site.EpisodeRef
			if len(epRefs) > 0 {
				firstArcEp = epRefs[0]
				latestArcEp = epRefs[len(epRefs)-1]
			}

			arcHero := map[string]any{"Title": a.Title, "Slug": a.Slug, "Summary": a.Summary, "Emoji": a.Emoji}
			var lastRel time.Time
			if a.LastRelease != nil {
				lastRel = *a.LastRelease
			}

			arcData := map[string]any{
				"Section":       "sagas",
				"Saga":          map[string]any{"Title": s.Title, "Slug": s.Slug},
				"Arc":           arcHero,
				"Episodes":      epRefs,
				"FirstEpisode":  firstArcEp,
				"LatestEpisode": latestArcEp,
				"PrevArc":       arcLite(prevArc),
				"NextArc":       arcLite(nextArc),
				"LastRelease":   lastRel,
				"NowYear":       time.Now().Year(),
			}
			//write(t, "arc", filepath.Join("public", "sagas", s.Slug, a.Slug, "index.html"), arcData)
			renderView(base,
				"arc",
				"public/sagas/"+s.Slug+"/"+a.Slug+"/index.html",
				arcData,
				"templates/arc.gohtml",
			)

			// Episode pages with prev/next
			for ei, e := range a.Episodes {
				var prevEp, nextEp *site.Episode
				if ei > 0 {
					prevEp = a.Episodes[ei-1]
				}
				if ei+1 < len(a.Episodes) {
					nextEp = a.Episodes[ei+1]
				}

				// Episode header payload
				eh := map[string]any{
					"SagaTitle":   s.Title,
					"ArcTitle":    a.Title,
					"Number":      e.Number,
					"Title":       e.Title,
					"Subtitle":    e.Summary,
					"Author":      "Henrique Riccio",
					"Date":        e.Date,
					"ReadingTime": e.ReadingTime,
					"SagaSlug":    s.Slug,
					"ArcSlug":     a.Slug,
					"Recap":       "", // optional: pull from previous episode summary if you like
				}

				var prevRef, nextRef map[string]any
				if prevEp != nil {
					prevRef = map[string]any{"Title": prevEp.Title, "Slug": prevEp.Slug, "Number": prevEp.Number}
				}
				if nextEp != nil {
					nextRef = map[string]any{"Title": nextEp.Title, "Slug": nextEp.Slug, "Number": nextEp.Number}
				}

				episodeData := map[string]any{
					"Section":       "sagas",
					"EpisodeHeader": eh,
					"EpisodeHTML":   e.BodyHTML, // already safe via goldmark
					"Prev":          prevRef,
					"Next":          nextRef,
					"Saga":          map[string]any{"Title": s.Title, "Slug": s.Slug},
					"Arc":           map[string]any{"Title": a.Title, "Slug": a.Slug},
					"NowYear":       time.Now().Year(),
				}
				//out := filepath.Join("public", "sagas", s.Slug, a.Slug, e.Slug, "index.html")
				//write(t, "episode", out, episodeData)
				renderView(base,
					"episode",
					"public/sagas/"+s.Slug+"/"+a.Slug+"/"+e.Slug+"/index.html",
					episodeData,
					"templates/episode.gohtml",
				)
			}
		}
	}

	for _, p := range posts {
		out := postOutputPath(p.Permalink)
		postData := map[string]any{
			"Section": "library",
			"Post":    p,
			"NowYear": time.Now().Year(),
		}
		renderView(base,
			"post",
			out,
			postData,
			"templates/post.gohtml",
		)
	}

	feedItems := collectRecent(latest, posts, feedItemLimit)
	if len(feedItems) > feedItemLimit {
		feedItems = feedItems[:feedItemLimit]
	}
	if err := writeFeed(cfg, feedItems); err != nil {
		log.Fatalf("write feed: %v", err)
	}

	if err := writeSitemap(cfg, sagas, posts); err != nil {
		log.Fatalf("write sitemap: %v", err)
	}

	if err := writeSearchIndex(cfg, sagas); err != nil {
		log.Fatalf("write search index: %v", err)
	}
}

// ---- utilities ----

func write(t *template.Template, name, out string, data any) {
	if err := os.MkdirAll(filepath.Dir(out), 0o755); err != nil {
		log.Fatalf("mkdir %s: %v", out, err)
	}
	f, err := os.Create(out)
	if err != nil {
		log.Fatalf("create %s: %v", out, err)
	}
	defer f.Close()
	if err := t.ExecuteTemplate(f, name, data); err != nil {
		log.Fatalf("render %s (%s): %v", name, out, err)
	}
	log.Printf("wrote %s", out)
}

func toHomeSagas(sagas []*site.Saga) []HomeSaga {
	out := make([]HomeSaga, 0, len(sagas))
	for _, s := range sagas {
		out = append(out, HomeSaga{
			Title:        s.Title,
			Slug:         s.Slug,
			Summary:      s.Summary,
			Emoji:        s.Emoji,
			EpisodeCount: s.EpisodeCount,
			LastRelease:  s.LastRelease,
			Status:       s.Status,
		})
	}
	// sort by last release desc
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
	return out
}

func postOutputPath(permalink string) string {
	perm := strings.TrimSpace(permalink)
	if perm == "" || perm == "/" {
		return filepath.Join("public", "index.html")
	}
	perm = strings.TrimPrefix(perm, "/")
	perm = strings.TrimSuffix(perm, "/")
	if perm == "" {
		return filepath.Join("public", "index.html")
	}
	if strings.HasSuffix(perm, ".html") {
		return filepath.Join("public", filepath.FromSlash(perm))
	}
	return filepath.Join("public", filepath.FromSlash(perm), "index.html")
}

func toHomeRecent(latest []*site.EpisodeRef, posts []site.Post, limit int) []HomeRecent {
	items := collectRecent(latest, posts, limit)
	if len(items) > limit {
		items = items[:limit]
	}
	return items
}

func collectRecent(latest []*site.EpisodeRef, posts []site.Post, maxEpisodes int) []HomeRecent {
	items := make([]HomeRecent, 0, len(posts)+maxEpisodes)

	for _, p := range posts {
		items = append(items, HomeRecent{
			Title:     p.Title,
			Summary:   p.Summary,
			Date:      p.Date,
			Type:      p.Type,
			Saga:      p.Saga,
			Arc:       p.Arc,
			Permalink: p.Permalink,
			Tags:      p.Tags,
		})
	}

	if maxEpisodes > len(latest) {
		maxEpisodes = len(latest)
	}
	for i := 0; i < maxEpisodes; i++ {
		ref := latest[i]
		perm := fmt.Sprintf("/sagas/%s/%s/%s/", ref.SagaSlug, ref.ArcSlug, ref.Slug)
		items = append(items, HomeRecent{
			Title:     ref.Title,
			Summary:   ref.Summary,
			Date:      ref.Date,
			Type:      "Episode",
			Saga:      ref.SagaTitle,
			Arc:       ref.ArcTitle,
			Permalink: perm,
		})
	}

	sort.SliceStable(items, func(i, j int) bool {
		if items[i].Date.Equal(items[j].Date) {
			return items[i].Title < items[j].Title
		}
		return items[i].Date.After(items[j].Date)
	})

	return items
}

type rss struct {
	XMLName xml.Name   `xml:"rss"`
	Version string     `xml:"version,attr"`
	Channel rssChannel `xml:"channel"`
}

type rssChannel struct {
	Title         string    `xml:"title"`
	Link          string    `xml:"link"`
	Description   string    `xml:"description,omitempty"`
	LastBuildDate string    `xml:"lastBuildDate,omitempty"`
	Items         []rssItem `xml:"item"`
}

type rssItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	GUID        string `xml:"guid"`
	PubDate     string `xml:"pubDate"`
	Description string `xml:"description,omitempty"`
}

func writeFeed(cfg site.Config, items []HomeRecent) error {
	feed := rss{Version: "2.0"}
	feed.Channel = rssChannel{
		Title:         "wasting no time",
		Link:          cfg.AbsoluteURL("/"),
		Description:   "Latest posts and episodes from wasting no time",
		LastBuildDate: time.Now().UTC().Format(time.RFC1123Z),
	}

	for _, item := range items {
		link := cfg.AbsoluteURL(item.Permalink)
		feed.Channel.Items = append(feed.Channel.Items, rssItem{
			Title:       item.Title,
			Link:        link,
			GUID:        link,
			PubDate:     item.Date.UTC().Format(time.RFC1123Z),
			Description: item.Summary,
		})
	}

	return writeXMLFile("public/feed.xml", feed)
}

type sitemapEntry struct {
	Loc     string
	LastMod time.Time
}

type sitemapURLSet struct {
	XMLName xml.Name          `xml:"urlset"`
	XMLNS   string            `xml:"xmlns,attr"`
	URLs    []sitemapURLEntry `xml:"url"`
}

type sitemapURLEntry struct {
	Loc     string `xml:"loc"`
	LastMod string `xml:"lastmod,omitempty"`
}

func writeSitemap(cfg site.Config, sagas []*site.Saga, posts []site.Post) error {
	if cfg.BaseURL == "" {
		return fmt.Errorf("site base URL must be configured for sitemap generation")
	}
	entries := map[string]*sitemapEntry{}

	add := func(path string, lastMod *time.Time) {
		loc := cfg.AbsoluteURL(path)
		if existing, ok := entries[loc]; ok {
			if lastMod != nil && (existing.LastMod.IsZero() || lastMod.After(existing.LastMod)) {
				existing.LastMod = *lastMod
			}
			return
		}
		entry := &sitemapEntry{Loc: loc}
		if lastMod != nil {
			entry.LastMod = *lastMod
		}
		entries[loc] = entry
	}

	add("/", nil)
	add("/library/", nil)

	for _, s := range sagas {
		add(fmt.Sprintf("/sagas/%s/", s.Slug), s.LastRelease)
		for _, a := range s.Arcs {
			add(fmt.Sprintf("/sagas/%s/%s/", s.Slug, a.Slug), a.LastRelease)
			for _, e := range a.Episodes {
				t := e.Date
				add(fmt.Sprintf("/sagas/%s/%s/%s/", s.Slug, a.Slug, e.Slug), &t)
			}
		}
	}

	for _, p := range posts {
		t := p.Date
		add(p.Permalink, &t)
	}

	list := make([]sitemapEntry, 0, len(entries))
	for _, e := range entries {
		list = append(list, *e)
	}
	sort.Slice(list, func(i, j int) bool { return list[i].Loc < list[j].Loc })

	urlset := sitemapURLSet{XMLNS: "http://www.sitemaps.org/schemas/sitemap/0.9"}
	for _, e := range list {
		entry := sitemapURLEntry{Loc: e.Loc}
		if !e.LastMod.IsZero() {
			entry.LastMod = e.LastMod.UTC().Format("2006-01-02")
		}
		urlset.URLs = append(urlset.URLs, entry)
	}

	return writeXMLFile("public/sitemap.xml", urlset)
}

func writeXMLFile(outPath string, data any) error {
	var buf bytes.Buffer
	buf.WriteString(xml.Header)
	enc := xml.NewEncoder(&buf)
	enc.Indent("", "  ")
	if err := enc.Encode(data); err != nil {
		return err
	}
	buf.WriteByte('\n')
	if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
		return err
	}
	return os.WriteFile(outPath, buf.Bytes(), 0o644)
}

func toHomeLatest(latest []*site.EpisodeRef) []HomeEpisode {
	max1 := 6
	if len(latest) < max1 {
		max1 = len(latest)
	}
	out := make([]HomeEpisode, 0, max1)
	for i := 0; i < max1; i++ {
		ref := latest[i]
		out = append(out, HomeEpisode{
			Title:     ref.Title,
			Slug:      ref.Slug,
			Number:    ref.Number,
			Date:      ref.Date,
			Summary:   ref.Summary,
			SagaTitle: ref.SagaTitle,
			SagaSlug:  ref.SagaSlug,
			ArcSlug:   ref.ArcSlug,
			ArcTitle:  ref.ArcTitle,
			ArcEmoji:  "⚙️",
		})
	}
	return out
}

func lastNEpisodes(a *site.Arc, n int) []*site.Episode {
	if n > len(a.Episodes) {
		n = len(a.Episodes)
	}
	// episodes already sorted oldest→newest; take tail
	start := len(a.Episodes) - n
	return a.Episodes[start:]
}

type searchEntry struct {
	Title   string   `json:"title"`
	URL     string   `json:"url"`
	Type    string   `json:"type"`
	Summary string   `json:"summary,omitempty"`
	Tags    []string `json:"tags,omitempty"`
	Context string   `json:"context,omitempty"`
	Date    string   `json:"date,omitempty"`
}

func writeSearchIndex(cfg site.Config, sagas []*site.Saga) error {
	entries := make([]searchEntry, 0)

	add := func(entry searchEntry) {
		entries = append(entries, entry)
	}

	for _, s := range sagas {
		add(searchEntry{
			Title:   s.Title,
			URL:     cfg.Href(fmt.Sprintf("/sagas/%s/", s.Slug)),
			Type:    "saga",
			Summary: s.Summary,
			Tags:    uniqueStrings(append([]string{}, s.Tags...)),
		})

		for _, a := range s.Arcs {
			add(searchEntry{
				Title:   a.Title,
				URL:     cfg.Href(fmt.Sprintf("/sagas/%s/%s/", s.Slug, a.Slug)),
				Type:    "arc",
				Summary: a.Summary,
				Context: s.Title,
				Tags:    uniqueStrings(append([]string{}, s.Tags...)),
			})

			for _, e := range a.Episodes {
				entry := searchEntry{
					Title:   e.Title,
					URL:     cfg.Href(fmt.Sprintf("/sagas/%s/%s/%s/", s.Slug, a.Slug, e.Slug)),
					Type:    "episode",
					Tags:    uniqueStrings(append(append([]string{}, e.Tags...), s.Tags...)),
					Context: fmt.Sprintf("%s · %s", s.Title, a.Title),
				}
				if !e.Date.IsZero() {
					entry.Date = e.Date.Format("2006-01-02")
				}
				if e.Summary != "" {
					entry.Summary = e.Summary
				}
				add(entry)
			}
		}
	}

	data, err := json.Marshal(entries)
	if err != nil {
		return err
	}
	if err := os.MkdirAll("public", 0o755); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join("public", "search.json"), data, 0o644)
}

func uniqueStrings(values []string) []string {
	seen := make(map[string]struct{}, len(values))
	out := make([]string, 0, len(values))
	for _, v := range values {
		if v == "" {
			continue
		}
		key := strings.ToLower(v)
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		out = append(out, v)
	}
	return out
}
