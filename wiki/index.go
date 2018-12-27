package main

import (
	"bufio"
	"fmt"
	"github.com/cybercongress/cyberd/client"
	"os"
	"regexp"
	"strings"
	"unicode"
)

func ContinueIndex(cbdClient client.CyberdClient) {

	startArticleId := int64(1)

	f, err := os.OpenFile("enwiki-latest-all-titles", 0, 0)
	if err != nil {
		panic(err)
	}
	br := bufio.NewReader(f)
	defer f.Close()

	reg, err := regexp.Compile("[^a-zA-Z0-9]+")

	counter := int64(0)
	links := make([]client.Link, 0, 1000)
	for {

		line, err := br.ReadString('\n')

		if err != nil {
			break
		}

		if counter < startArticleId {
			counter++
			continue
		}

		split := strings.Split(strings.TrimSuffix(line, "\n"), "\t")
		ids := strings.Split(split[1], "_")

		idsSet := make(map[string]struct{})
		for _, id := range ids {
			idsSet[id] = struct{}{}
		}

		for id := range idsSet {

			id = reg.ReplaceAllString(id, "")
			id = strings.ToLower(id)

			if len(id) == 0 || id == "" {
				continue
			}

			if len(id) == 1 && unicode.IsSymbol(rune(id[0])) {
				continue
			}

			page := split[1] + ".wiki"
			links = append(links, client.Link{From: Cid(id), To: Cid(page)})
			counter++

			if len(links) == 1000 {
				fmt.Printf("%d %s\n", counter, split[1])

				accBw, err := cbdClient.GetAccountBandwidth()
				if err == nil {
					per := int64(100 * float64(accBw.RemainedValue) / float64(accBw.MaxValue))
					fmt.Printf("Remaining acc bw: %d %v%%\n", accBw.RemainedValue, per)
				}

				err = cbdClient.SubmitLinksSync(links)
				if err != nil {
					panic(err.Error())
				}
				links = make([]client.Link, 0, 1000)
			}
		}
	}
}
