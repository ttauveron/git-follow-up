package internal

import (
	"testing"
	"time"

	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

func TestFilter_setFrom(t *testing.T) {
	type fields struct {
		From    time.Time
		Labels  []string
		Authors []string
		Display []string
	}
	type args struct {
		from string
		now  time.Time
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   time.Time
	}{
		{
			name: "Week to date on sunday",
			args: args{
				from: "wtd",
				now:  time.Date(2019, time.June, 30, 0, 0, 0, 0, time.UTC),
			},
			want: time.Date(2019, time.June, 24, 0, 0, 0, 0, time.UTC),
		},
		{
			name: "Week to date on thursday",
			args: args{
				from: "wtd",
				now:  time.Date(2019, time.June, 27, 0, 0, 0, 0, time.UTC),
			},
			want: time.Date(2019, time.June, 24, 0, 0, 0, 0, time.UTC),
		},
		{
			name: "Week to date on monday",
			args: args{
				from: "wtd",
				now:  time.Date(2019, time.June, 24, 0, 0, 0, 0, time.UTC),
			},
			want: time.Date(2019, time.June, 24, 0, 0, 0, 0, time.UTC),
		},
		{
			name: "Month to date on last day of month",
			args: args{
				from: "mtd",
				now:  time.Date(2019, time.February, 24, 0, 0, 0, 0, time.UTC),
			},
			want: time.Date(2019, time.February, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			name: "Month to date on first day of month",
			args: args{
				from: "mtd",
				now:  time.Date(2019, time.February, 1, 0, 0, 0, 0, time.UTC),
			},
			want: time.Date(2019, time.February, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			name: "Year to date",
			args: args{
				from: "ytd",
				now:  time.Date(2019, time.February, 1, 0, 0, 0, 0, time.UTC),
			},
			want: time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			name: "Year to date on first day of the year",
			args: args{
				from: "ytd",
				now:  time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC),
			},
			want: time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			name: "Today",
			args: args{
				from: "today",
				now:  time.Date(2019, time.April, 1, 0, 0, 0, 0, time.UTC),
			},
			want: time.Date(2019, time.April, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			name: "Yesterday",
			args: args{
				from: "yesterday",
				now:  time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC),
			},
			want: time.Date(2018, time.December, 31, 0, 0, 0, 0, time.UTC),
		},
		{
			name: "Arbitrary date",
			args: args{
				from: "2019-05-05",
				now:  time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC),
			},
			want: time.Date(2019, time.May, 5, 0, 0, 0, 0, time.UTC),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filter := &Filter{}
			filter.setFrom(tt.args.from, tt.args.now)
			if got := filter.From; got != tt.want {
				t.Errorf("%q. Run() =\n%v, want\n%v", tt.name, got, tt.want)
			}
		})
	}
}

func TestFilter_Filter(t *testing.T) {
	type fields struct {
		From    time.Time
		Labels  []string
		Authors []string
		Display []string
	}
	type args struct {
		c *object.Commit
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		wantB  bool
	}{
		{
			name: "filtering by author, author found",
			args: args{
				c: &object.Commit{
					Author: object.Signature{
						Name:  "jean",
						Email: "test@test.te",
						When:  time.Date(2019, time.May, 5, 0, 0, 0, 0, time.UTC),
					},
				},
			},
			fields: fields{
				Authors: []string{"jean"},
			},

			wantB: true,
		},
		{
			name: "filtering by author, author not found",
			args: args{
				c: &object.Commit{
					Author: object.Signature{
						Name:  "jack",
						Email: "test@test.te",
						When:  time.Date(2019, time.May, 5, 0, 0, 0, 0, time.UTC),
					},
				},
			},
			fields: fields{
				Authors: []string{"jean"},
			},

			wantB: false,
		},
		{
			name: "filtering by date, matching",
			args: args{
				c: &object.Commit{
					Author: object.Signature{
						Name:  "jack",
						Email: "test@test.te",
						When:  time.Date(2019, time.May, 5, 0, 0, 0, 0, time.UTC),
					},
				},
			},
			fields: fields{
				From: time.Date(2019, time.May, 5, 0, 0, 0, 0, time.UTC),
			},

			wantB: true,
		},
		{
			name: "filtering by date, not matching",
			args: args{
				c: &object.Commit{
					Author: object.Signature{
						Name:  "jack",
						Email: "test@test.te",
						When:  time.Date(2019, time.May, 4, 0, 0, 0, 0, time.UTC),
					},
				},
			},
			fields: fields{
				From: time.Date(2019, time.May, 5, 0, 0, 0, 0, time.UTC),
			},

			wantB: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filter := Filter{
				From:    tt.fields.From,
				Labels:  tt.fields.Labels,
				Authors: tt.fields.Authors,
				Display: tt.fields.Display,
			}
			if gotB := filter.Filter(tt.args.c); gotB != tt.wantB {
				t.Errorf("Filter.Filter() = %v, want %v", gotB, tt.wantB)
			}
		})
	}
}
