package ffuf

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	domain "github.com/sudankdk/offense/internal/domain/ffuf"
)

// ffufRawResult matches ffuf's -of json line output
type ffufRawResult struct {
	Input struct {
		FUZZ string `json:"FUZZ"`
	} `json:"input"`
	Position         int    `json:"position"`
	Status           int    `json:"status"`
	Length           int    `json:"length"`
	Words            int    `json:"words"`
	Lines            int    `json:"lines"`
	ContentType      string `json:"content-type"`
	RedirectLocation string `json:"redirectlocation"`
	URL              string `json:"url"`
	Duration         int64  `json:"duration"`
	Host             string `json:"host"`
}

// ffufRawConfig is embedded in ffuf's JSON summary output
type ffufRawConfig struct {
	CommandLine string          `json:"commandline"`
	Time        string          `json:"time"`
	Results     []ffufRawResult `json:"results"`
	Config      struct {
		Method          string   `json:"method"`
		URL             string   `json:"url"`
		Wordlists       []string `json:"wordlists"`
		Timeout         int      `json:"requesttimeout"`
		Threads         int      `json:"threads"`
		FollowRedirects bool     `json:"followredirects"`
		Calibration     bool     `json:"autocalibration"`
		Matchers        struct {
			Status string `json:"status"`
		} `json:"matchers"`
	} `json:"config"`
	Stats struct {
		Errors int `json:"errors"`
		Total  int `json:"total"`
	} `json:"stats"`
}

func parseOutput(raw string) (domain.FfufData, error) {
	// ffuf -of json writes a single JSON object at the end
	// but also streams per-result JSON lines — find the summary object
	summary, err := extractSummary([]byte(raw))
	if err != nil {
		return domain.FfufData{}, err
	}
	return mapToDomain(summary), nil
}

// extractSummary finds the top-level ffuf JSON summary block.
// ffuf streams progress lines + per-result lines, then writes one
// big summary JSON object. We scan for the object that has a "config" key.
func extractSummary(raw []byte) (ffufRawConfig, error) {
	scanner := bufio.NewScanner(bytes.NewReader(raw))
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if !strings.HasPrefix(line, "{") {
			continue
		}
		var candidate ffufRawConfig
		if err := json.Unmarshal([]byte(line), &candidate); err != nil {
			continue
		}
		// The summary block always has config.url populated
		if candidate.Config.URL != "" {
			return candidate, nil
		}
	}
	return ffufRawConfig{}, errors.New("ffuf: could not find summary JSON in output")
}

func mapToDomain(raw ffufRawConfig) domain.FfufData {
	results := make([]domain.FfufResult, 0, len(raw.Results))
	for _, r := range raw.Results {
		results = append(results, domain.FfufResult{
			Position:    r.Position,
			Fuzz:        decodeFuzz(r.Input.FUZZ),
			URL:         r.URL,
			RedirectURL: r.RedirectLocation,
			Status:      r.Status,
			Length:      r.Length,
			Words:       r.Words,
			Lines:       r.Lines,
			ContentType: r.ContentType,
			DurationMs:  float64(r.Duration) / 1e6, // nanoseconds → ms
		})
	}

	wordlist := ""
	if len(raw.Config.Wordlists) > 0 {
		wordlist = raw.Config.Wordlists[0]
	}

	return domain.FfufData{
		Config: domain.FfufConfig{
			Method:          raw.Config.Method,
			URL:             raw.Config.URL,
			Wordlist:        wordlist,
			FollowRedirects: raw.Config.FollowRedirects,
			Calibration:     raw.Config.Calibration,
			Timeout:         raw.Config.Timeout,
			Threads:         raw.Config.Threads,
			MatcherCodes:    parseStatusCodes(raw.Config.Matchers.Status),
		},
		Progress: domain.FfufProgress{
			Total:     raw.Stats.Total,
			Completed: raw.Stats.Total, // summary only written on completion
			JobsTotal: 1,
			JobsDone:  1,
			Errors:    raw.Stats.Errors,
		},
		Results: results,
	}
}

// decodeFuzz handles ffuf base64-encoding its FUZZ values in JSON mode
func decodeFuzz(s string) string {
	decoded, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return s // already plaintext
	}
	return string(decoded)
}

// parseStatusCodes turns "200-299,301,302,307" into []int
func parseStatusCodes(s string) []int {
	codes := []int{}
	for _, part := range strings.Split(s, ",") {
		part = strings.TrimSpace(part)
		if strings.Contains(part, "-") {
			// range like "200-299"
			var lo, hi int
			if _, err := fmt.Sscanf(part, "%d-%d", &lo, &hi); err == nil {
				for c := lo; c <= hi; c++ {
					codes = append(codes, c)
				}
			}
		} else {
			var code int
			if _, err := fmt.Sscanf(part, "%d", &code); err == nil {
				codes = append(codes, code)
			}
		}
	}
	return codes
}
