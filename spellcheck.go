package main

import (
    "github.com/trustmaster/go-aspell"
    "fmt"
    "strings"
    "net/http"
    "encoding/json"
    "os"
)

type Correction struct {
    Word string `json:"word"`
    Ud string `json:"ud"`
    Suggestions []string `json:"suggestions"`
}

func check(text string, lang string) string {
    speller, err := aspell.NewSpeller(map[string]string{
        "lang": lang,
    })

    if err != nil {
        fmt.Println(err.Error())
        return "{error: true}"
    }
    defer speller.Delete()

    words := strings.Split(text, ", ")
    corrections := []Correction{}
    maxSuggestions := 12;

    for _, word := range words {
        if !speller.Check(word) {
            c := Correction{Word: word, Ud: "false"}
            c.Suggestions = speller.Suggest(word)
            if len(c.Suggestions) > maxSuggestions {
                c.Suggestions = c.Suggestions[:maxSuggestions]
            }

//            fmt.Println(c)
//            j, _ := json.Marshal(c)
//            fmt.Println(string(j))

            corrections = append(corrections, c)
        }
    }

    data, _ := json.Marshal(corrections)
    fmt.Println(string(data))
    return string(data)
}

func handler(w http.ResponseWriter, r *http.Request) {
    query := r.URL.Query()
    command := query.Get("cmd")
    callback := query.Get("callback")
    data := ""

    switch command {
    case "check_spelling":
        lang := query.Get("slang")
        text := query.Get("text")
        data = check(text, lang)
    case "get_lang_list":
        data = "{langList: {ltr: {\"en_US\": \"American English\"}, rtl: {}}, verLang : 6}"
    case "getbanner":
        data = "{banner: false}"
    }

//    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Content-Type", "text/javascript; charset=UTF-8")
    w.Header().Set("X-XSS-Protection", "0")

    fmt.Fprintf(w, "%s(%s)", callback, data)
}

func main() {
    http.HandleFunc("/", handler)
    http.ListenAndServe(":" + os.Getenv("PORT"), nil)
}
