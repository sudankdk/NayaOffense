package domain

// FfufConfig captures scan settings
type FfufConfig struct {
	Method          string `json:"method"`
	URL             string `json:"url"`
	Wordlist        string `json:"wordlist"`
	FollowRedirects bool   `json:"follow_redirects"`
	Calibration     bool   `json:"calibration"`
	Timeout         int    `json:"timeout"`
	Threads         int    `json:"threads"`
	MatcherCodes    []int  `json:"matcher_codes"`
}

// FfufProgress captures scan progress
type FfufProgress struct {
	Total     int `json:"total"`
	Completed int `json:"completed"`
	JobsTotal int `json:"jobs_total"`
	JobsDone  int `json:"jobs_done"`
	Errors    int `json:"errors"`
	ReqPerSec int `json:"req_per_sec"`
}

// FfufResult is a single matched endpoint
type FfufResult struct {
	Position    int     `json:"position"`
	Fuzz        string  `json:"fuzz"`
	URL         string  `json:"url"`
	RedirectURL string  `json:"redirect_url,omitempty"`
	Status      int     `json:"status"`
	Length      int     `json:"length"`
	Words       int     `json:"words"`
	Lines       int     `json:"lines"`
	ContentType string  `json:"content_type"`
	DurationMs  float64 `json:"duration_ms"`
}

// FfufData is what goes into Response.Data for the ffuf tool
type FfufData struct {
	Config   FfufConfig   `json:"config"`
	Progress FfufProgress `json:"progress"`
	Results  []FfufResult `json:"results"`
}
