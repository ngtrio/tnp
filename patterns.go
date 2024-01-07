package tnp

import (
	"fmt"

	"github.com/samber/lo"
)

type Pattern struct {
	Regex   string
	Replace string
}

const (
	delimiters  = "[.\\s\\-+_\\/(),]"
	seasonRange = "(?:Complete" +
		delimiters +
		"*)?" +
		"(?:s(?:easons?)?)" +
		delimiters +
		`*(?:s?[0-9]{1,2}[\s]*(?:(?:\-|(?:\s*to\s*))[\s]*s?[0-9]{1,2})+)(?:` +
		delimiters +
		"*Complete)?"
)

func init() {
	episode := []Pattern{
		{"(?:e|ep)(?:[0-9]{1,2}(?:-?(?:e|ep)?(?:[0-9]{1,2}))?)", ""},
		{"\\s-\\s\\d{1,3}\\s", ""},
		{`\b[0-9]{1,2}x([0-9]{2})\b`, ""},
		{fmt.Sprintf(`\bepisod(?:e|io)%s\d{1,2}\b`, delimiters), ""},
	}

	season := []Pattern{
		{`s?(\d{1,2})\s\-\s\d{1,2}\s`, ""},
		{fmt.Sprintf(`\b%s\b`, seasonRange), ""},
		{`(?:s\d{1,2}[.+\s]*){2,}\b`, ""},
		{fmt.Sprintf(`\b(?:Complete%s)?s([0-9]{1,2})%s?\b`, delimiters, link(episode).Regex), ""},
		{fmt.Sprintf(`Series%s\d{1,2}`, delimiters), ""},
		{fmt.Sprintf(`\b(?:Complete%s)?Season[\. -][0-9]{1,2}\b`, delimiters), ""},
	}

	Patterns["episode"] = episode
	Patterns["season"] = season
}

var Patterns = map[string][]Pattern{
	"network": {
		{`\bMyTVS\b`, "MyTVS"},
		{`\bATVP\b`, "Apple TV+"},
		{`\bAMZN|Amazon\b`, "Amazon Studios"},
		{`\bNF|Netflix\b`, "Netflix"},
		{`\bNICK\b`, "Nickelodeon"},
		{`\bRED\b`, "YouTube Premium"},
		{`\bDSNY?P\b`, "Disney Plus"},
		{`\bDSNY\b`, "DisneyNOW"},
		{`\bHMAX\b`, "HBO Max"},
		{`\bHBO\b`, "HBO"},
		{`\bHULU\b`, "Hulu Networks"},
		{`\bMS?NBC\b`, "MSNBC"},
		{`\bDCU\b`, "DC Universe"},
		{`\bID\b`, "Investigation Discovery"},
		{`\biT\b`, "iTunes"},
		{`\bAS\b`, "Adult Swim"},
		{`\bCRAV\b`, "Crave"},
		{`\bCC\b`, "Comedy Central"},
		{`\bSESO\b`, "Seeso"},
		{`\bVRV\b`, "VRV"},
		{`\bPCOK\b`, "Peacock"},
		{`\bCBS\b`, "CBS"},
		{`\biP\b`, "BBC iPlayer"},
		{`\bNBC\b`, "NBC"},
		{`\bAMC\b`, "AMC"},
		{`\bPBS\b`, "PBS"},
		{`\bSTAN\b`, "Stan."},
		{`\bRTE\b`, "RTE Player"},
		{`\bCR\b`, "Crunchyroll"},
		{`\bANPL\b`, "Animal Planet Live"},
		{`\bDTV\b`, "DirecTV Stream"},
		{`\bVICE\b`, "VICE"},
	},
	"quality": {
		{`\bWEB[ -.]?DL(?:Rip|Mux)?|HDRip`, "WEB-DL"},
		{`\bWEB[ -]?Cap\b`, "WEBCap"},
		{`\bW[EB]B[ -]?(?:Rip)|WEB\b`, "WEBRip"},
		{`\b(?:HD)?CAM(?:-?Rip)?\b`, "Cam"},
		{`\b(?:HD)?TS|TELESYNC|PDVD|PreDVDRip\b`, "Telesync"},
		{`\bWP|WORKPRINT\b`, "Workprint"},
		{`\b(?:HD)?TC|TELECINE\b`, "Telecine"},
		{`\b(?:DVD)?SCR(?:EENER)?|BDSCR\b`, "Screener"},
		{`\bDDC\b`, "Digital Distribution Copy"},
		{`\bDVD-?(?:Rip|Mux)\b`, "DVD-Rip"},
		{`\bDVDR|DVD-Full|Full-rip\b`, "DVD-R"},
		{`\bPDTV|DVBRip\b`, "PDTV"},
		{`\bDSR(?:ip)?|SATRip|DTHRip\b`, "DSRip"},
		{`\bAHDTV(?:Mux)?\b`, "AHDTV"},
		{`\bHDTV(?:Rip)?\b`, "HDTV"},
		{`\bD?TVRip|DVBRip\b`, "TVRip"},
		{`\bVODR(?:ip)?\b`, "VODRip"},
		{`\bHD-Rip\b`, "HD-Rip"},
		{fmt.Sprintf(`\bBlu-?Ray%sRip|BDR(?:ip)?\b`, delimiters), "BDRip"},
		{`\bBlu-?Ray|(?:US|JP)?BD(?:remux)?\b`, "Blu-ray"},
		{`\bBR-?Rip\b`, "BRRip"},
		{`\bHDDVD\b`, "HD DVD"},
		{`\bPPV(?:Rip)?\b`, "Pay-Per-View Rip"},
	},
	"resolution": {
		{`\b([0-9]{3,4}(?:p|i))\b`, ""},
		{`\b(1280x720p?)\b`, "720p"},
		{`\bFHD|1920x1080p?\b`, "1080p"},
		{`\bUHD\b`, "UHD"},
		{`\bHD\b`, "HD"},
		{`\b4K\b`, "4K"},
	},
	"year": {
		{`\b(?:19|20)[0-9]{2}-(?:19|20)[0-9]{2}\b`, ""},
		{`\b(?:19|20)[0-9]{2}\b`, ""},
	},
	"codec": {
		{"xvid", "Xvid"},
		{"av1", "AV1"},
		{fmt.Sprintf("[hx]%s?264", delimiters), "H.264"},
		{"AVC", "H.264"},
		{fmt.Sprintf("HEVC(?:{d}Main%s?10P?)", delimiters), "H.265 Main 10"},
		{fmt.Sprintf("[hx]%s?265", delimiters), "H.265"},
		{"HEVC", "H.265"},
		{fmt.Sprintf("[h]%s?263", delimiters), "H.263"},
		{"VC-1", "VC-1"},
	},
	"audio": genAudioPatterns(),
	"bit_depth": {
		{`(8|10)-?bits?`, ""},
	},
	"audio_track": {
		{`\d+Audio`, ""},
	},
}

func link(patterns []Pattern) Pattern {
	return Pattern{
		Regex: lo.Reduce(patterns, func(agg string, pattern Pattern, idx int) string {
			end := lo.If(idx == len(patterns)-1, ")").Else("|")
			return agg + pattern.Regex + end
		}, "(?:"),
	}
}

var channels = [][]int{
	{1, 0},
	{2, 0},
	{5, 1},
	{6, 1},
	{7, 1},
}

func genAudioPatterns() []Pattern {
	var res []Pattern

	var raw = []Pattern{
		{"LPCM", "LPCM"},
		{"TrueHD", "Dolby TrueHD"},
		{"Atmos", "Dolby Atmos"},
		{"DD-EX", "Dolby Digital EX"},
		{"DDP|E-?AC-?3|EC-3", "Dolby Digital Plus"},
		{"DD|AC-?3|DolbyD", "Dolby Digital"},
		{fmt.Sprintf("DTS%s?HD(?:%s?(?:MA|Masters?(?:%sAudio)?))", delimiters, delimiters, delimiters), "DTS-HD MA"},
		{"DTSMA", "DTS-HD MA"},
		{fmt.Sprintf("DTS%s?HD", delimiters), "DTS-HD"},
		{"DTS", "DTS"},
		{"AAC[ \\.\\-]LC", "AAC-LC"},
		{"AAC", "AAC"},
		{fmt.Sprintf("Dual%sAudios?", delimiters), "Dual"},
		{"FLAC", "FLAC"},
		{"OGG", "OGG"},
	}

	for _, r := range raw {
		for _, c := range channels {
			speakers := c[0]
			subwoofers := c[1]
			res = append(res, Pattern{
				Regex:   fmt.Sprintf(`((?:%s)%s*%v[. \-]?%v(?:ch)?)`, r.Regex, delimiters, speakers, subwoofers),
				Replace: fmt.Sprintf(`%s %v.%v`, r.Replace, speakers, subwoofers),
			})
		}
	}

	for _, r := range raw {
		res = append(res, Pattern{
			Regex:   fmt.Sprintf("(%s)", r.Regex),
			Replace: r.Replace,
		})
	}

	return res
}
