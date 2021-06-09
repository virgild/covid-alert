package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const pageURL = `https://unityhealth.to/patients-and-visitors/covid-19/`

func getPageReport() (*Report, error) {
	res, err := http.Get(pageURL)
	if err != nil {
		return nil, fmt.Errorf("http get: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status not ok: %d", res.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, fmt.Errorf("new document: %w", err)
	}

	nodes := doc.Find("table").First().
		Find("tr").
		Find("td").Map(func(i int, selection *goquery.Selection) string {
		return selection.Text()
	})

	// for i, n := range nodes {
	// 	log.Info().Msgf("%d. %s", i, n)
	// }

	icuCount, err := strconv.Atoi(nodes[1])
	if err != nil {
		return nil, fmt.Errorf("convert icu count: %w", err)
	}

	inpatientCount, err := strconv.Atoi(nodes[2])
	if err != nil {
		return nil, fmt.Errorf("convert inpatients count: %w", err)
	}

	report := &Report{
		Location:       nodes[0],
		ICUUnits:       icuCount,
		InpatientUnits: inpatientCount,
		Time:           time.Now(),
	}

	return report, nil
}
