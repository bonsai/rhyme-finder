package rhymfinder

import (
    "sort"
    "strings"
    "unicode/utf8"
)

type Mode string

const (
    ModeStrict     Mode = "strict"
    ModeAggressive Mode = "aggressive"
)

type Query struct {
    Text  string `json:"text"`
    Mode  Mode   `json:"mode,omitempty"`
    Limit int    `json:"limit,omitempty"`
}

type Candidate struct {
    Text       string  `json:"text"`
    Score      float64 `json:"score"`
    Phonetic   float64 `json:"phonetic"`
    Rhythm     float64 `json:"rhythm"`
    Position   float64 `json:"position"`
}

type Finder struct { Words []string }

func New(words []string) *Finder { return &Finder{Words: append([]string(nil), words...)} }

func (f *Finder) Find(q Query) []Candidate {
    if q.Limit <= 0 { q.Limit = 20 }
    target := normalize(q.Text)
    if target == "" { return nil }
    out := make([]Candidate, 0, len(f.Words))
    for _, raw := range f.Words {
        word := normalize(raw)
        if word == "" || word == target { continue }
        p := phoneticScore(target, word)
        r := rhythmScore(target, word)
        pos := positionScore(target, word)
        score := 0.60*p + 0.25*r + 0.15*pos
        if q.Mode == ModeAggressive { score = 0.50*p + 0.25*r + 0.25*pos }
        if q.Mode == ModeStrict && p < 0.80 { continue }
        out = append(out, Candidate{Text: raw, Score: score, Phonetic:p, Rhythm:r, Position:pos})
    }
    sort.SliceStable(out, func(i,j int) bool { if out[i].Score == out[j].Score { return out[i].Text < out[j].Text }; return out[i].Score > out[j].Score })
    if len(out) > q.Limit { out = out[:q.Limit] }
    return out
}

func normalize(s string) string { return strings.TrimSpace(strings.ToLower(s)) }

func phoneticScore(a,b string) float64 {
    aa, bb := []rune(a), []rune(b)
    if len(aa)==0 || len(bb)==0 { return 0 }
    // Japanese-friendly suffix comparison: mora-like trailing units dominate rhyme strength.
    n := min(len(aa), len(bb)); same := 0
    for i:=1; i<=n; i++ { if aa[len(aa)-i] == bb[len(bb)-i] { same++ } else { break } }
    suffix := float64(same)/float64(n)
    edit := 1.0 - float64(levenshtein(aa,bb))/float64(max(len(aa),len(bb)))
    return clamp(0.70*suffix + 0.30*edit)
}

func rhythmScore(a,b string) float64 {
    ar, br := utf8.RuneCountInString(a), utf8.RuneCountInString(b)
    d := ar-br; if d<0 { d=-d }
    return clamp(1.0-float64(d)/float64(max(ar,br)))
}

func positionScore(a,b string) float64 {
    if strings.HasSuffix(a,b) || strings.HasSuffix(b,a) { return 1 }
    if strings.Contains(a,b) || strings.Contains(b,a) { return 0.75 }
    return 0.25
}

func levenshtein(a,b []rune) int { prev:=make([]int,len(b)+1); for j:=range prev {prev[j]=j}; for i:=1;i<=len(a);i++ { cur:=make([]int,len(b)+1); cur[0]=i; for j:=1;j<=len(b);j++ { cost:=0; if a[i-1]!=b[j-1] {cost=1}; cur[j]=min(cur[j-1]+1,min(prev[j]+1,prev[j-1]+cost)) }; prev=cur }; return prev[len(b)] }
func min(a,b int) int {if a<b{return a};return b}
func max(a,b int) int {if a>b{return a};return b}
func clamp(v float64) float64 {if v<0{return 0};if v>1{return 1};return v}
