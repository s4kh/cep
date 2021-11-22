package main

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"text/tabwriter"
)

type cep struct {
	minute     string
	hour       string
	dayOfMonth string
	month      string
	dayOfWeek  string
	command    string
}

type limit struct {
	min int
	max int
}

var limits = map[string]limit{
	"minute":     {0, 59},
	"hour":       {0, 23},
	"dayOfMonth": {1, 31},
	"month":      {1, 12},
	"dayOfWeek":  {0, 6},
}

var fields = []string{
	"minute",
	"hour",
	"dayOfMonth",
	"month",
	"dayOfWeek",
}

/**
* - every
/ - 14/5 -> starting from 14th every 5 ex. 14, 19, 24
4-8 - range ex. 4,5,6,7,8
*/

var cronRe = regexp.MustCompile(`^((?:[^\s]+\s+){5}(?:\d{4})?)(?:\s+)?(.*)`)
var onlyDigit = regexp.MustCompile("^[0-9]*$")
var validChars = regexp.MustCompile(`^[\d|/|*|\-|,]+$`)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func parseRepeat(val string, field string) ([]int, error) {
	exp := strings.Split(val, "/")
	if len(exp) > 1 {
		if onlyDigit.MatchString(exp[0]) {
			// Convert to range
			exp[0] = exp[0] + "-" + strconv.Itoa(limits[field].max)
		}

		return parseRange(exp[0], exp[1], field)
	}
	return parseRange(val, "1", field)
}

func parseRange(val string, interval string, field string) ([]int, error) {
	exp := strings.Split(val, "-")
	if len(exp) > 1 {
		start, err := strconv.Atoi(exp[0])
		if err != nil {
			return nil, errors.New("invalid number value provided for range start: " + exp[0])
		}
		end, err := strconv.Atoi(exp[1])
		if err != nil {
			return nil, errors.New("invalid number value provided for range end: " + exp[1])
		}

		if start < limits[field].min || end > limits[field].max {
			return nil, errors.New("range error got " + exp[0] + "-" + exp[1] + " expected range " + strconv.Itoa(limits[field].min) + "-" + strconv.Itoa(limits[field].max))
		}
		if start >= end {
			return nil, errors.New("invalid range " + val)
		}

		inc, err := strconv.Atoi(interval)
		if err != nil {
			return nil, errors.New("invalid number value provided for interval")
		}
		if inc <= 0 {
			return nil, errors.New("invalid interval " + interval)
		}

		ans := []int{}
		for i := start; i <= end; i += inc {
			ans = append(ans, i)
		}
		return ans, nil
	}
	ival, err := strconv.Atoi(val)
	if err != nil {
		return nil, errors.New("invalid number value")
	}
	if ival < limits[field].min || ival > limits[field].max {
		return nil, errors.New("value error got " + val + " must be between " + strconv.Itoa(limits[field].min) + "-" + strconv.Itoa(limits[field].max))
	}
	return []int{ival}, nil
}

func converToStrArr(values []int) []string {
	str := []string{}
	for i := range values {
		number := values[i]
		text := strconv.Itoa(number)
		str = append(str, text)
	}
	return str
}

func removeDuplicates(arr []int) []int {
	set := make(map[int]bool)
	removed := []int{}

	for _, entry := range arr {
		if _, value := set[entry]; !value {
			set[entry] = true
			removed = append(removed, entry)
		}
	}
	return removed
}

func ParseField(field string, expression string) (string, error) {
	// Star means every number between min and max range
	maxRange := strconv.Itoa(limits[field].min) + "-" + strconv.Itoa(limits[field].max)
	expression = strings.ReplaceAll(expression, "*", maxRange)
	exp := strings.Split(expression, ",")

	if len(exp) == 0 {
		return "", errors.New("Wrong expression passed for " + field)
	}

	var vals []int
	var err error
	// Comma separated multiple values
	if len(exp) > 1 {
		for _, val := range exp {
			var temp []int
			temp, err = parseRepeat(val, field)
			if err != nil {
				vals = nil
				break
			}
			vals = append(vals, temp...)
		}
	} else {
		vals, err = parseRepeat(expression, field)
	}

	// Cleaning
	vals = removeDuplicates(vals)
	sort.Ints(vals)
	arr := converToStrArr(vals)

	return strings.Join(arr, " "), err
}

func run() error {
	args := os.Args[1:]
	if len(args) < 1 {
		return errors.New("not enough args")
	}

	match := cronRe.FindStringSubmatch(args[0])
	if len(match) < 2 {
		return errors.New("invalid cron expression")
	}
	blocks := strings.Split(args[0], " ")
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.Debug)
	for i, block := range blocks[:5] {
		if !validChars.MatchString(block) {
			return errors.New("Invalid characters: " + block)
		}
		vals, err := ParseField(fields[i], block)
		if err != nil {
			return errors.New(fields[i] + " error:(" + block + ") " + err.Error())
		}
		fmt.Fprintln(w, fields[i], "\t", vals)
	}
	fmt.Fprintln(w, "command", "\t", blocks[len(blocks)-1])
	w.Flush()

	return nil
}
