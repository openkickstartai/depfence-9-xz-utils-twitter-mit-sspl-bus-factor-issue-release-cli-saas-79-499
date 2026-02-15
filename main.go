package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type Dep struct {
	Name, Version, Repo string
}

var goModRe = regexp.MustCompile(`^\s+(\S+)\s+v(\S+)`)

func main() {
	f := flag.String("f", "", "dependency file (auto-detect if empty)")
	format := flag.String("format", "table", "output: table, json, csv")
	token := flag.String("token", os.Getenv("GITHUB_TOKEN"), "GitHub API token")
	minScore := flag.Float64("min-score", 40, "exit 1 if any dep below this")
	flag.Parse()
	if *f == "" {
		for _, n := range []string{"go.mod", "package.json", "requirements.txt"} {
			if _, e := os.Stat(n); e == nil { *f = n; break }
		}
	}
	if *f == "" { fmt.Fprintln(os.Stderr, "no dependency file found"); os.Exit(2) }
	deps, err := parseDeps(*f)
	if err != nil { fmt.Fprintln(os.Stderr, err); os.Exit(2) }
	results := NewScorer(*token).ScoreAll(deps)
	switch *format {
	case "json":
		e := json.NewEncoder(os.Stdout); e.SetIndent("", "  "); e.Encode(results)
	case "csv":
		fmt.Println("name,score,risk,activity,bus_factor,license")
		for _, r := range results {
			fmt.Printf("%s,%.0f,%s,%.0f,%.0f,%s\n", r.Name, r.Score, r.Risk, r.Activity, r.BusFactor, r.License)
		}
	default:
		fmt.Printf("\n  %-40s %5s %-8s %4s %3s %s\n", "DEPENDENCY", "SCORE", "RISK", "ACT", "BUS", "LICENSE")
		fmt.Printf("  %s\n", strings.Repeat("─", 78))
		for _, r := range results {
			ic := "🟢"; if r.Score < 60 { ic = "🟡" }; if r.Score < 40 { ic = "🔴" }
			fmt.Printf("  %-40s %s%3.0f %-8s %3.0f %3.0f %s\n", r.Name, ic, r.Score, r.Risk, r.Activity, r.BusFactor, r.License)
		}
		fmt.Println()
	}
	for _, r := range results { if r.Score < *minScore { os.Exit(1) } }
}

func parseDeps(path string) ([]Dep, error) {
	data, err := os.ReadFile(path)
	if err != nil { return nil, err }
	b := filepath.Base(path)
	switch {
	case strings.Contains(b, "go.mod"): return parseGoMod(data), nil
	case strings.Contains(b, "package.json"): return parsePkgJSON(data), nil
	case strings.Contains(b, "requirements"): return parseRequirements(data), nil
	}
	return nil, fmt.Errorf("unsupported file: %s", b)
}

func parseGoMod(data []byte) []Dep {
	var deps []Dep
	sc := bufio.NewScanner(strings.NewReader(string(data)))
	for sc.Scan() {
		if m := goModRe.FindStringSubmatch(sc.Text()); m != nil {
			deps = append(deps, Dep{m[1], m[2], m[1]})
		}
	}
	return deps
}

func parsePkgJSON(data []byte) []Dep {
	var pkg map[string]json.RawMessage
	json.Unmarshal(data, &pkg)
	var deps []Dep
	for _, key := range []string{"dependencies", "devDependencies"} {
		var d map[string]string
		if raw, ok := pkg[key]; ok {
			json.Unmarshal(raw, &d)
			for n, v := range d { deps = append(deps, Dep{n, strings.TrimLeft(v, "^~>=<"), ""}) }
		}
	}
	return deps
}

func parseRequirements(data []byte) []Dep {
	var deps []Dep
	sc := bufio.NewScanner(strings.NewReader(string(data)))
	for sc.Scan() {
		l := strings.TrimSpace(sc.Text())
		if l == "" || l[0] == '#' { continue }
		p := strings.SplitN(l, "==", 2)
		v := ""; if len(p) > 1 { v = p[1] }
		deps = append(deps, Dep{strings.TrimSpace(p[0]), v, ""})
	}
	return deps
}
