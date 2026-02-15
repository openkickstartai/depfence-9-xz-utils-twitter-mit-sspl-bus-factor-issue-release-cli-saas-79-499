package main

import (
	"os"
	"testing"
)

func TestParseGoMod(t *testing.T) {
	data := []byte("module test\n\ngo 1.21\n\nrequire (\n\tgithub.com/gin-gonic/gin v1.9.1\n\tgithub.com/lib/pq v1.10.9\n\tgolang.org/x/crypto v0.14.0\n)\n")
	deps := parseGoMod(data)
	if len(deps) != 3 { t.Fatalf("expected 3 deps, got %d", len(deps)) }
	if deps[0].Name != "github.com/gin-gonic/gin" { t.Errorf("got %s", deps[0].Name) }
	if deps[0].Version != "1.9.1" { t.Errorf("got version %s", deps[0].Version) }
}

func TestParsePkgJSON(t *testing.T) {
	data := []byte(`{"dependencies":{"express":"^4.18.0","lodash":"^4.17.21"},"devDependencies":{"jest":"^29.0.0"}}`)
	deps := parsePkgJSON(data)
	if len(deps) != 3 { t.Fatalf("expected 3, got %d", len(deps)) }
	found := map[string]bool{}
	for _, d := range deps { found[d.Name] = true }
	for _, n := range []string{"express", "lodash", "jest"} {
		if !found[n] { t.Errorf("missing: %s", n) }
	}
}

func TestParseRequirements(t *testing.T) {
	data := []byte("# comment\nflask==2.3.0\nrequests==2.31.0\n\nnumpy==1.24.0\n")
	deps := parseRequirements(data)
	if len(deps) != 3 { t.Fatalf("expected 3, got %d", len(deps)) }
	if deps[0].Name != "flask" || deps[0].Version != "2.3.0" { t.Errorf("got %+v", deps[0]) }
}

func TestRiskLevel(t *testing.T) {
	cases := []struct{ s float64; w string }{{85, "low"}, {55, "medium"}, {20, "critical"}}
	for _, c := range cases {
		if g := riskLevel(c.s); g != c.w { t.Errorf("riskLevel(%.0f)=%s, want %s", c.s, g, c.w) }
	}
}

func TestResolveGitHubRepo(t *testing.T) {
	if r := resolveGitHubRepo(Dep{Name: "github.com/gin-gonic/gin"}); r != "gin-gonic/gin" {
		t.Errorf("got %s", r)
	}
	if r := resolveGitHubRepo(Dep{Name: "lodash"}); r != "" {
		t.Errorf("expected empty, got %s", r)
	}
}

func TestStaticScoring(t *testing.T) {
	s := NewScorer("")
	r := s.scoreDep(Dep{Name: "github.com/test/repo", Version: "1.0.0"})
	if r.Score != 50 { t.Errorf("expected 50, got %.0f", r.Score) }
	if r.Risk != "medium" { t.Errorf("expected medium, got %s", r.Risk) }
}

func TestScoreAll(t *testing.T) {
	s := NewScorer("")
	deps := []Dep{{"github.com/a/b", "1.0", ""}, {"github.com/c/d", "2.0", ""}}
	results := s.ScoreAll(deps)
	if len(results) != 2 { t.Fatalf("expected 2 results, got %d", len(results)) }
	for _, r := range results {
		if r.Score < 0 || r.Score > 100 { t.Errorf("score out of range: %.0f", r.Score) }
	}
}

func TestParseDepsGoMod(t *testing.T) {
	f, _ := os.CreateTemp("", "go.mod.*")
	f.WriteString("module test\n\ngo 1.21\n\nrequire (\n\tgithub.com/pkg/errors v0.9.1\n)\n")
	f.Close()
	defer os.Remove(f.Name())
	deps, err := parseDeps(f.Name())
	if err != nil { t.Fatal(err) }
	if len(deps) != 1 { t.Fatalf("expected 1, got %d", len(deps)) }
	if deps[0].Name != "github.com/pkg/errors" { t.Errorf("got %s", deps[0].Name) }
}

func TestParseDepsRequirements(t *testing.T) {
	f, _ := os.CreateTemp("", "requirements*.txt")
	f.WriteString("flask==2.3.0\ndjango==4.2.0\n")
	f.Close()
	defer os.Remove(f.Name())
	deps, err := parseDeps(f.Name())
	if err != nil { t.Fatal(err) }
	if len(deps) != 2 { t.Fatalf("expected 2, got %d", len(deps)) }
}
