package main

import (
	"fmt"
	"strings"
)

func savePaste(key, content string) {
	if conf.Save == false {
		return
	}

	if len(content) > conf.MaxSize {
		return
	}

	writeDB(conf.db, "pastes", key, []byte(content))
}

func processRegexes(key, content string) {
	save := false
	for i, _ := range conf.Regexes {
		r := conf.Regexes[i]

		switch r.Match {
		case "all":
			items := r.compiled.FindAllString(content, -1)

			if items != nil {
				save = true
			}

			for k := range items {
				rKey := fmt.Sprintf("%s-%d", key, k)
				writeDB(conf.db, r.Bucket, rKey, []byte(items[k]))
			}
		case "one":
			match := r.compiled.FindString(content)

			if match != "" {
				save = true
				writeDB(conf.db, r.Bucket, key, []byte(match))
			}
		default:
		}
	}

	if save {
		savePaste(key, content)
	}
}

func processKeywords(key, content string) {
	save := false
	for i, _ := range conf.Keywords {
		kwd := conf.Keywords[i]

		if strings.Contains(strings.ToLower(content), strings.ToLower(kwd.Keyword)) {
			save = true
			writeDB(conf.db, kwd.Bucket, key, nil)
		}
	}

	if save {
		savePaste(key, content)
	}
}

func processContent(key, content string) {
	conf.db = getDBConn()
	defer conf.db.Close()

	processRegexes(key, content)
	processKeywords(key, content)
}
