package main

import "time"

type Report struct {
	Location       string
	ICUUnits       int
	InpatientUnits int
	Time           time.Time
}

func (r *Report) Total() int {
	return r.ICUUnits + r.InpatientUnits
}

func (r *Report) HasChangedFrom(report *Report) bool {
	return r.ICUUnits != report.ICUUnits || r.InpatientUnits != report.InpatientUnits
}
