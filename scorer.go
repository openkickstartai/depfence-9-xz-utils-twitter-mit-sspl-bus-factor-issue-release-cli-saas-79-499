package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"strings"
	"time"
)

type Result struct {
	Name      string  `json:"name"`
	Score     float64 `json:"score"`
	Risk      string  `json:"risk"`
	Activity  float64 `json:"maintainer_activity"`
	BusFactor float64 `json:"bus_factor"`
	License   string  `json:"license"`
}

type Scorer struct {
	Token  string
	client http.Client
}

func NewScorer(token string) *Scorer {
	return &Scorer{Token: token, client: http.Client{Timeout: 10 * time.Second}}
}

func (s *Scorer) ScoreAll(deps []Dep) []Result {
	results := make([]Result, len(deps))
	for i, d := range deps { results[i] = s.scoreDep(d) }
	return results
}

func (s *Scorer) scoreDep(d Dep) Result {
	r := Result{Name: d.Name, License: "unknown", Activity: 50, BusFactor: 50, Score: 50, Risk: "unknown"}
	repo := resolveGitHubRepo(d)
	if repo == "" || s.Token == "" { return r }
	var repoData struct {
		PushedAt string `json:"pushed_at"`
		License  *struct{ SpdxID string `json:"spdx_id"` } `json:"license"`
	}
	if body := s.ghGet("https://api.github.com/repos/" + repo); body != nil {
		json.Unmarshal(body, &repoData)
	}
	if t, err := time.Parse(time.RFC3339, repoData.PushedAt); err == nil {
		days := time.Since(t).Hours() / 24
		r.Activity = math.Max(0, math.Min(100, 100-days/3.65))
	}
	if repoData.License != nil && repoData.License.SpdxID != "" {
		r.License = repoData.License.SpdxID
	}
	var contribs []struct{}
	url := fmt.Sprintf("https://api.github.com/repos/%s/contributors?per_page=100&anon=true", repo)
	if body := s.ghGet(url); body != nil {
		json.Unmarshal(body, &contribs)
		r.BusFactor = math.Min(100, float64(len(contribs))*10)
	}
	licScore := 100.0
	for _, risky := range []string{"SSPL", "BSL", "BUSL", "NOASSERTION"} {
		if strings.EqualFold(r.License, risky) { licScore = 30 }
	}
	if r.License == "unknown" { licScore = 50 }
	r.Score = math.Round(r.Activity*0.4 + r.BusFactor*0.3 + licScore*0.3)
	r.Risk = riskLevel(r.Score)
	return r
}

func (s *Scorer) ghGet(url string) []byte {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil { return nil }
	req.Header.Set("Authorization", "Bearer "+s.Token)
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	resp, err := s.client.Do(req)
	if err != nil { return nil }
	defer resp.Body.Close()
	if resp.StatusCode != 200 { return nil }
	body, _ := io.ReadAll(resp.Body)
	return body
}

func resolveGitHubRepo(d Dep) string {
	for _, src := range []string{d.Repo, d.Name} {
		if strings.HasPrefix(src, "github.com/") {
			parts := strings.SplitN(strings.TrimPrefix(src, "github.com/"), "/", 3)
			if len(parts) >= 2 { return parts[0] + "/" + parts[1] }
		}
	}
	return ""
}

func riskLevel(score float64) string {
	switch {
	case score >= 70: return "low"
	case score >= 40: return "medium"
	default: return "critical"
	}
}
