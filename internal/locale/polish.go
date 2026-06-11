package locale

import "time"

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

func PolishMonth(m time.Month) string {
	return polishMonths[m]
}
