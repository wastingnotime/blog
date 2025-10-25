package site

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"
)

type TaggedPost struct {
	Title     string
	Summary   string
	Date      time.Time
	Permalink string
	Saga      string
	Arc       string
	Type      string
}

type LibraryTag struct {
	Name string
	Slug string
}

type TagPageData struct {
	Tag     string
	Posts   []TaggedPost
	NowYear int
}

type LibraryIndexData struct {
	Tags    []LibraryTag
	NowYear int
}

var nonAlnum = regexp.MustCompile(`[^a-z0-9]+`)

func BuildTagIndex(sagas []*Saga, posts []Post) (map[string][]TaggedPost, []LibraryTag) {
	type bucket struct {
		name  string
		posts []TaggedPost
	}
	buckets := map[string]*bucket{}

	add := func(tag string, entry TaggedPost) {
		display := strings.TrimSpace(tag)
		if display == "" {
			return
		}
		key := strings.ToLower(display)
		b, ok := buckets[key]
		if !ok {
			b = &bucket{name: display}
			buckets[key] = b
		}
		b.posts = append(b.posts, entry)
	}

	for _, p := range posts {
		entry := TaggedPost{
			Title:     p.Title,
			Summary:   p.Summary,
			Date:      p.Date,
			Permalink: p.Permalink,
			Saga:      p.Saga,
			Arc:       p.Arc,
			Type:      strings.ToLower(strings.TrimSpace(p.Type)),
		}
		if entry.Type == "" {
			entry.Type = "post"
		}
		if len(p.Tags) == 0 {
			continue
		}
		for _, tag := range p.Tags {
			add(tag, entry)
		}
	}

	for _, s := range sagas {
		for _, a := range s.Arcs {
			for _, e := range a.Episodes {
				entry := TaggedPost{
					Title:     e.Title,
					Summary:   e.Summary,
					Date:      e.Date,
					Permalink: fmt.Sprintf("/sagas/%s/%s/%s/", s.Slug, a.Slug, e.Slug),
					Saga:      s.Title,
					Arc:       a.Title,
					Type:      "episode",
				}
				if len(e.Tags) == 0 {
					continue
				}
				for _, tag := range e.Tags {
					add(tag, entry)
				}
			}
		}
	}

	keys := make([]string, 0, len(buckets))
	for key := range buckets {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	usedSlugs := map[string]struct{}{}
	tagMap := make(map[string][]TaggedPost, len(buckets))
	tagList := make([]LibraryTag, 0, len(buckets))

	for _, key := range keys {
		bucket := buckets[key]
		posts := bucket.posts
		sort.SliceStable(posts, func(i, j int) bool {
			ti, tj := posts[i].Date, posts[j].Date
			if ti.Equal(tj) {
				return posts[i].Title < posts[j].Title
			}
			return ti.After(tj)
		})
		tagMap[bucket.name] = posts

		slug := uniqueTagSlug(bucket.name, usedSlugs)
		tagList = append(tagList, LibraryTag{Name: bucket.name, Slug: slug})
		usedSlugs[slug] = struct{}{}
	}

	sort.Slice(tagList, func(i, j int) bool {
		return strings.ToLower(tagList[i].Name) < strings.ToLower(tagList[j].Name)
	})

	return tagMap, tagList
}

func tagSlug(name string) string {
	slug := strings.ToLower(strings.TrimSpace(name))
	slug = nonAlnum.ReplaceAllString(slug, "-")
	slug = strings.Trim(slug, "-")
	if slug == "" {
		slug = "tag"
	}
	return slug
}

func uniqueTagSlug(name string, used map[string]struct{}) string {
	base := tagSlug(name)
	if _, exists := used[base]; !exists {
		return base
	}
	for i := 2; i < 1_000; i++ {
		candidate := fmt.Sprintf("%s-%d", base, i)
		if _, exists := used[candidate]; !exists {
			return candidate
		}
	}
	return fmt.Sprintf("%s-%d", base, time.Now().Unix())
}
