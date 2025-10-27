package site

import (
	"html/template"
	"time"
)

type ArcLink struct {
	Title string
	URL   string
}

type EpisodeLink struct {
	Title  string
	URL    string
	Number int
}

type ArcPageData struct {
	Section       string
	NowYear       int
	SagaTitle     string
	SagaURL       string
	ArcTitle      string
	ArcSummary    string
	Episodes      []*EpisodeRef
	FirstEpisode  *EpisodeRef
	LatestEpisode *EpisodeRef
	PrevArc       *ArcLink
	NextArc       *ArcLink
	LastRelease   *time.Time
	IsProd        bool
}

type EpisodePageData struct {
	Section            string
	NowYear            int
	SagaTitle          string
	SagaURL            string
	ArcTitle           string
	ArcURL             string
	EpisodeTitle       string
	EpisodeNumber      int
	EpisodeSummary     string
	EpisodeAuthor      string
	EpisodeDate        time.Time
	EpisodeReadingTime string
	EpisodeBody        template.HTML
	PrevEpisode        *EpisodeLink
	NextEpisode        *EpisodeLink
	IsProd             bool
}
