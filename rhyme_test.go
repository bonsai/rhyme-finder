package rhymefinder

import "testing"

func TestAggressiveRhyme(t *testing.T) {
    f := New([]string{"強引", "チラク", "東京", "方向"})
    got := f.Find(Query{Text:"強引", Mode:ModeAggressive, Limit:10})
    if len(got) == 0 { t.Fatal("expected candidates") }
}

func TestStrictFilters(t *testing.T) {
    f := New([]string{"強引", "東京", "旅行"})
    got := f.Find(Query{Text:"強引", Mode:ModeStrict, Limit:10})
    for _, c := range got { if c.Phonetic < .80 { t.Fatalf("weak candidate: %+v", c) } }
}
