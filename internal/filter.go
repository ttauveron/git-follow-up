package internal

import (
	"fmt"
	"github.com/spf13/pflag"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"regexp"
	"time"
)

type Filter struct {
	From    time.Time
	Labels  []string
	Authors []string
	Display []string
}

var DisplayArgs = []string{"repo", "date", "hash", "message", "author"}
var FromArgs = []string{"ytd", "mtd", "wtd", "yesterday", "today"}

func NewFilter(flags *pflag.FlagSet) (f *Filter) {
	f = &Filter{}

	// Date filter
	from, err := flags.GetString("from")
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	f.SetFrom(from)

	// Labels filter
	labels, err := flags.GetStringSlice("label")
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	f.Labels = append(f.Labels, labels...)

	// Author filter
	authors, err := flags.GetStringSlice("author")
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	f.Authors = append(f.Authors, authors...)

	// Display filter
	if !flags.Changed("display") {
		f.Display = append(f.Display, DisplayArgs...)
	} else {
		displays, err := flags.GetStringSlice("display")
		if err != nil {
			fmt.Printf("%v\n", err)
		}
		f.Display = append(f.Display, displays...)
	}

	return
}

func (filter *Filter) SetFrom(from string) {
	// TODO handle timezone
	regexDate, _ := regexp.Compile("([12]\\d{3}-(0[1-9]|1[0-2])-(0[1-9]|[12]\\d|3[01]))")

	now := time.Now()
	currentYear, currentMonth, currentDay := now.Date()
	currentLocation := now.Location()

	switch {
	case from == "ytd":
		filter.From = time.Date(currentYear, time.January, 1, 0, 0, 0, 0, currentLocation)
		break
	case from == "mtd":
		filter.From = time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
		break
	case from == "wtd":
		y, m, d := now.AddDate(0, 0, -(int(now.Weekday()+1)%8 - 1)).Date()
		filter.From = time.Date(y, m, d, 0, 0, 0, 0, currentLocation)
		break
	case from == "today":
		filter.From = time.Date(currentYear, currentMonth, currentDay, 0, 0, 0, 0, currentLocation)
		break
	case from == "yesterday":
		filter.From = time.Date(currentYear, currentMonth, currentDay-1, 0, 0, 0, 0, currentLocation)
		break
	case regexDate.MatchString(from):
		filter.From, _ = time.Parse("2006-01-02", from)
		break
	default:
		fmt.Println("from flag not recognized")
		break
	}
}

func (filter Filter) Filter(c *object.Commit) (b bool) {

	b = true
	author := c.Author.Name + " " + c.Author.Email

	switch {
	// Filter by date
	case !c.Author.When.After(filter.From):
		b = false
		break
	// Filter by author
	case filter.Authors != nil && !MatchAny(author, filter.Authors):
		b = false
		break
	}

	return
}
