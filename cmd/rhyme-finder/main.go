package main

import (
    "flag"
    "log"
    "net/http"
    "os"
    "strings"

    rhymefinder "github.com/bonsai/rhyme-finder"
)

func main() {
    addr := flag.String("addr", ":8080", "HTTP listen address")
    dict := flag.String("dict", "", "newline-separated word dictionary")
    flag.Parse()
    words := []string{}
    if *dict != "" {
        b, err := os.ReadFile(*dict); if err != nil { log.Fatal(err) }
        for _, w := range strings.Split(string(b), "\n") { if strings.TrimSpace(w)!="" { words=append(words, strings.TrimSpace(w)) } }
    }
    if len(words)==0 { words=[]string{"強引","チラク","東京","旅行","交響","証拠","方向","成功","最高","妄想"} }
    server := rhymefinder.Server{Finder:rhymefinder.New(words)}
    log.Printf("rhyme-finder listening on %s", *addr)
    log.Fatal(http.ListenAndServe(*addr, server.Handler()))
}
