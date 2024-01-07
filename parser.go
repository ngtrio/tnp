package tnp

import (
	"regexp"
	"slices"
	"sort"
	"strconv"
	"strings"

	"github.com/samber/lo"
)

type MediaType string

const (
	TypeTVSeries = MediaType("tv_series")
	TypeTVSeason = MediaType("tv_season")
	TypeMovie    = MediaType("movie")
)

type Match struct {
	Type    string
	Start   int
	End     int
	Content string
}

type Parsed struct {
	Name       string
	Year       string
	Producer   string
	Codec      string
	BitDepth   string
	Audio      string
	Resoluton  string
	Quality    string
	Season     []int
	Episode    []int
	Network    string
	Excess     []string
	MediaType  MediaType
	AudioTrack string
}

func Parse(title string, standard bool) *Parsed {
	var unmatch []string
	matches := make(map[string]Match)

	for t, ps := range Patterns {
		for _, p := range ps {
			r := regexp.MustCompile("(?i)" + p.Regex)
			ms := r.FindAllStringSubmatchIndex(title, -1)
			if ms == nil {
				continue
			}
			m := ms[0]
			if t == "year" {
				m = ms[len(ms)-1]
			}
			start, end := m[0], m[1]
			if len(m) > 2 && m[2] != -1 {
				start, end = m[2], m[3]
			}
			if _, has := matches[t]; !has {
				matches[t] = Match{
					t,
					m[0],
					m[1],
					lo.If(standard && p.Replace != "", p.Replace).Else(title[start:end]),
				}
			}
		}
	}

	mts := lo.Values(matches)
	sort.Slice(mts, func(i, j int) bool {
		return mts[i].Start < mts[j].Start
	})

	i, mIdx := 0, 0
	for i < len(title) {
		var unmatched string
		if mIdx < len(matches) {
			m := mts[mIdx]
			if m.Start > i {
				unmatched = title[i:m.Start]
			}
			if i < m.End {
				i = m.End
			}
			mIdx++
		} else {
			unmatched = title[i:]
			i = len(title)
		}

		if len(unmatched) > 0 {
			unmatch = append(unmatch, unmatched)
		}
	}

	var (
		name     string
		producer string
		excess   []string
	)
	if len(unmatch) > 0 {
		name = strings.TrimSpace(unmatch[0])
		unmatch = slices.Delete(unmatch, 0, 1)
		trimRe := regexp.MustCompile(`(^[-_.\s(),]+)|([-.\\s,]+$)`)
		splitRe := regexp.MustCompile(`\.\.+|\s+`)
		for _, u := range unmatch {
			trimed := trimRe.ReplaceAllString(u, "")
			for _, e := range splitRe.Split(trimed, -1) {
				if len(e) > 0 {
					excess = append(excess, e)
				}
			}
		}
	}

	if len(excess) > 0 {
		producer = excess[len(excess)-1]
		excess = slices.Delete(excess, len(excess)-1, len(excess))
	}

	parsed := new(Parsed)
	for _, m := range matches {
		if m.Content != "" {
			switch m.Type {
			case "year":
				parsed.Year = m.Content
			case "codec":
				parsed.Codec = m.Content
			case "audio":
				parsed.Audio = m.Content
			case "resolution":
				parsed.Resoluton = m.Content
			case "quality":
				parsed.Quality = m.Content
			case "season":
				parsed.Season = getSeasonEpisode(m.Content)
			case "episode":
				parsed.Episode = getSeasonEpisode(m.Content)
			case "network":
				parsed.Network = m.Content
			case "bit_depth":
				parsed.BitDepth = m.Content
			case "audio_track":
				parsed.AudioTrack = m.Content
			}
		}
	}
	parsed.Name = name
	parsed.Producer = producer
	parsed.Excess = excess

	if len(parsed.Season) > 1 {
		parsed.MediaType = TypeTVSeries
	} else if len(parsed.Season) > 0 {
		parsed.MediaType = TypeTVSeason
	} else {
		parsed.MediaType = TypeMovie
	}

	return parsed
}

func getSeasonEpisode(match string) []int {
	re := regexp.MustCompile(`[0-9]+`)
	matches := re.FindAllString(match, -1)
	var res []int
	if len(matches) == 2 {
		left, _ := strconv.Atoi(matches[0])
		right, _ := strconv.Atoi(matches[1])
		for i := left; i <= right; i++ {
			res = append(res, i)
		}
	} else if len(matches) > 0 {
		n, _ := strconv.Atoi(matches[0])
		res = append(res, n)
	}

	return res
}
