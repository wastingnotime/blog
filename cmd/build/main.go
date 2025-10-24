package main

import (
	"html/template"
	"log"
	"os"
	"path/filepath"
	"sort"
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
	LastRelease  time.Time
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

// ---- main ----

func main() {
	cfg := site.ConfigFromEnv()
	// 1) Load content (sagas, arcs, episodes as parsed structs)
	sagas, latest, err := site.Load("content")
	if err != nil {
		log.Fatalf("load: %v", err)
	}

	base := loadBase(cfg)

	// 2) Render Home
	homeData := map[string]any{
		"Sagas":          toHomeSagas(sagas),
		"LatestEpisodes": toHomeLatest(latest),
		"NowYear":        time.Now().Year(),
	}

	//write(t, "home", "public/index.html", homeData)
	renderView(base,
		"home",
		"public/index.html",
		homeData,
		"templates/home.gohtml",
	)

	// 3) Render Library
	libData := map[string]any{
		"Sagas":   sagas,
		"Query":   "",
		"NowYear": time.Now().Year(),
	}
	//write(t, "library", "public/library/index.html", libData)
	renderView(base,
		"library",
		"public/library/index.html",
		libData,
		"templates/library.gohtml",
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
				ref := &site.EpisodeRef{Title: e.Title, Slug: e.Slug, Number: e.Number, Date: e.Date, Summary: e.Summary, ArcSlug: a.Slug, ArcTitle: a.Title}
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
				refs = append(refs, &site.EpisodeRef{Title: e.Title, Slug: e.Slug, Number: e.Number, Date: e.Date, Summary: e.Summary, ArcSlug: currentArc.Slug, ArcTitle: currentArc.Title})
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
				epRefs = append(epRefs, &site.EpisodeRef{Title: e.Title, Slug: e.Slug, Number: e.Number, Date: e.Date, Summary: e.Summary})
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
		var last time.Time
		if s.LastRelease != nil {
			last = *s.LastRelease
		}
		out = append(out, HomeSaga{
			Title:        s.Title,
			Slug:         s.Slug,
			Summary:      s.Summary,
			Emoji:        s.Emoji,
			EpisodeCount: s.EpisodeCount,
			LastRelease:  last,
			Status:       s.Status,
		})
	}
	// sort by last release desc
	sort.SliceStable(out, func(i, j int) bool { return out[i].LastRelease.After(out[j].LastRelease) })
	return out
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
			SagaTitle: findSagaTitle(latest, ref), // optional; not strictly available here—override below if needed
			SagaSlug:  findSagaSlug(latest, ref),  // optional; not strictly available here—override below if needed
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

func findSagaSlug(_ []*site.EpisodeRef, _ *site.EpisodeRef) string {
	// If you need this mapping, pass SagaSlug through EpisodeRef in loader.
	return ""
}

func findSagaTitle(_ []*site.EpisodeRef, _ *site.EpisodeRef) string {
	// If you need this mapping, pass SagaSlug through EpisodeRef in loader.
	return ""
}
