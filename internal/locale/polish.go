package locale

import (
	"fmt"
	"time"
)

var polishMonths = map[time.Month]string{
	time.January:   "Styczeń",
	time.February:  "Luty",
	time.March:     "Marzec",
	time.April:     "Kwiecień",
	time.May:       "Maj",
	time.June:      "Czerwiec",
	time.July:      "Lipiec",
	time.August:    "Sierpień",
	time.September: "Wrzesień",
	time.October:   "Październik",
	time.November:  "Listopad",
	time.December:  "Grudzień",
}

var polishMonthsGen = map[time.Month]string{
	time.January:   "Stycznia",
	time.February:  "Lutego",
	time.March:     "Marca",
	time.April:     "Kwietnia",
	time.May:       "Maja",
	time.June:      "Czerwca",
	time.July:      "Lipca",
	time.August:    "Sierpnia",
	time.September: "Września",
	time.October:   "Października",
	time.November:  "Listopada",
	time.December:  "Grudnia",
}

var polishWeekdays = map[time.Weekday]string{
	time.Monday:    "Pn",
	time.Tuesday:   "Wt",
	time.Wednesday: "Śr",
	time.Thursday:  "Cz",
	time.Friday:    "Pt",
	time.Saturday:  "Sb",
	time.Sunday:    "Nd",
}

func PolishMonth(m time.Month) string {
	return polishMonths[m]
}

func PolishWeekday(d time.Weekday) string {
	return polishWeekdays[d]
}

func PolishMonthGen(m time.Month) string {
	return polishMonthsGen[m]
}

func FormatEventDate(s string) string {
	if len(s) < 10 {
		return s
	}
	t, err := time.Parse("2006-01-02", s[:10])
	if err != nil {
		return s // fall back to raw value rather than breaking the page
	}
	out := fmt.Sprintf("%d %s", t.Day(), PolishMonthGen(t.Month()))
	if t.Year() != time.Now().Year() {
		out = fmt.Sprintf("%s %d", out, t.Year())
	}
	return out
}
