package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const layout = "20060102"

type MonthRule struct {
	Days   []int
	Months []int
}

func nextDay(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	var err error

	nowStr := r.FormValue("now")
	if nowStr != "" {
		now, err = time.Parse(layout, nowStr)
		if err != nil {
			writeError(w, err, http.StatusInternalServerError)
			return
		}
	}

	dStart := r.FormValue("date")
	if dStart == "" {
		writeError(w, fmt.Errorf("empty date."), http.StatusBadRequest)
		return
	}

	repeat := r.FormValue("repeat")
	res, err := NextDate(now, dStart, repeat)
	if err != nil {
		writeError(w, err, http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(res))
}
func NextDate(now time.Time, dStart string, repeat string) (string, error) {
	parts := strings.Fields(repeat)
	if len(parts) == 0 {
		return "", fmt.Errorf("wrong rule.")
	}
	rule := parts[0]
	date, err := time.Parse(layout, dStart)

	if err != nil {
		return "", err
	}
	switch rule {

	case "y":
		if len(parts) > 1 {
			return "", fmt.Errorf("wrong rule.")
		}
		for {
			date = date.AddDate(1, 0, 0)
			if afterNow(date, now) {
				return date.Format(layout), nil
			}
		}
	case "d":
		if len(parts) != 2 {
			return "", fmt.Errorf("wrong rule.")
		}
		daysStr := parts[1]
		days, err := strconv.Atoi(daysStr)
		if err != nil {
			return "", err
		}
		if days > 400 {
			return "", fmt.Errorf("More than 400 days.")
		}
		if date.Equal(now) {
			return now.Format(layout), nil
		}
		for {
			date = date.AddDate(0, 0, days)
			if afterNow(date, now) {
				return date.Format(layout), nil
			}

		}
	case "w":
		if len(parts) != 2 {
			return "", fmt.Errorf("wrong rule.")
		}
		if date.Before(now) {
			date = now
		}
		daysStr := parts[1]
		weekdaysRule, err := parseWeekList(daysStr)
		if err != nil {
			return "", err
		}
		weekday := int(now.Weekday())
		if weekday == 0 {
			weekday = 7
		}
		daysCount := 7
		for _, w := range weekdaysRule {
			diff := w - weekday
			if diff <= 0 {
				diff = 7 - weekday + w
			}
			if daysCount >= diff {
				daysCount = diff
			}
		}
		for {
			date = date.AddDate(0, 0, daysCount)
			if afterNow(date, now) {
				return date.Format(layout), nil
			}
		}

	case "m":
		if len(parts) < 2 {
			return "", fmt.Errorf("wrong rule.")
		}
		if date.Before(now) {
			date = now
		}
		rule, err := parseMonthRule(repeat)
		if err != nil {
			return "", err
		}

		year := now.Year()

		for {
			next, ok := findNextInYear(date, year, rule)
			if ok {
				return next.Format(layout), nil
			}
			year++
			now = time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
		}
	default:
		return "", fmt.Errorf("unknown rule.")
	}

}
func afterNow(date, now time.Time) bool {
	date, now = date.Truncate(24*time.Hour), now.Truncate(24*time.Hour)

	return date.After(now)
}
func parseWeekList(str string) ([]int, error) {
	daysStr := strings.Split(str, ",")
	var res []int
	for _, ch := range daysStr {
		weekday, err := strconv.Atoi(ch)
		if err != nil {
			return nil, err
		}
		if weekday > 7 {
			return nil, fmt.Errorf("wrong rule.")
		}
		res = append(res, weekday)
	}
	return res, nil
}
func parseMonthRule(text string) (*MonthRule, error) {
	parts := strings.Fields(text)
	if len(parts) < 2 || parts[0] != "m" {
		return nil, fmt.Errorf("wrong repeat format")
	}
	days, err := parseDaysList(parts[1])
	if err != nil {
		return nil, err
	}
	var months []int
	if len(parts) >= 3 {
		months, err = parseMonthsList(parts[2])
		if err != nil {
			return nil, err
		}
	} else {
		for i := 1; i <= 12; i++ {
			months = append(months, i)
		}
	}

	return &MonthRule{
		Days:   days,
		Months: months,
	}, nil
}
func parseMonthsList(s string) ([]int, error) {
	parts := strings.Split(s, ",")
	res := make([]int, 0, len(parts))

	for _, p := range parts {
		if p == "" {
			continue
		}
		v, err := strconv.Atoi(strings.TrimSpace(p))
		if err != nil {
			return nil, err
		}
		if v < 1 || v > 12 {
			return nil, fmt.Errorf("Incorrect months.")
		}
		res = append(res, v)
	}
	return res, nil
}
func parseDaysList(s string) ([]int, error) {
	parts := strings.Split(s, ",")
	res := make([]int, 0, len(parts))

	for _, p := range parts {
		if p == "" {
			continue
		}
		v, err := strconv.Atoi(strings.TrimSpace(p))
		if err != nil {
			return nil, err
		}
		if v < -2 || v > 31 {
			return nil, fmt.Errorf("Incorrect days.")
		}
		res = append(res, v)
	}
	return res, nil
}
func findNextInYear(now time.Time, year int, r *MonthRule) (time.Time, bool) {
	best := time.Time{}

	for _, month := range r.Months {
		daysInMonth := daysInMonth(year, month)

		for _, d := range r.Days {
			day := d
			if d < 0 {
				day = daysInMonth + d + 1
			}

			if day <= 0 || day > daysInMonth {
				continue
			}

			candidate := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)

			if !candidate.After(now) {
				continue
			}

			if best.IsZero() || candidate.Before(best) {
				best = candidate
			}
		}
	}

	if best.IsZero() {
		return time.Time{}, false
	}
	return best, true
}

func daysInMonth(year int, month int) int {
	return time.Date(year, time.Month(month)+1, 0, 0, 0, 0, 0, time.UTC).Day()
}
