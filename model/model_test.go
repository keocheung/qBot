package model

import (
	"testing"

	"github.com/antonmedv/expr"
)

func TestRules(t *testing.T) {
	tests := []struct {
		name    string
		rule    Rule
		torrent Torrent
		want    bool
	}{{
		name: "test1",
		rule: Rule{
			Condition: "(Category == 'Bangumi' || any(Tags, {# == 'VCB' || # == 'Public'})) && MaxRatio == -1",
		},
		torrent: Torrent{
			Category: "Bangumi",
			MaxRatio: -1,
		},
		want: true,
	}, {
		name: "test2",
		rule: Rule{
			Condition: "(Category == 'Bangumi' || any(Tags, {# == 'VCB' || # == 'Public'})) && MaxRatio == -1",
		},
		torrent: Torrent{
			Category: "Music",
			Tags: []string{
				"VCB",
				"Normal",
			},
			MaxRatio: -1,
		},
		want: true,
	}, {
		name: "test2",
		rule: Rule{
			Condition: "(Category == 'Bangumi' || any(Tags, {# == 'VCB' || # == 'Public'})) && MaxRatio == -1",
		},
		torrent: Torrent{
			Category: "Music",
			Tags: []string{
				"Music",
				"Movie",
			},
			MaxRatio: -1,
		},
		want: false,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := expr.Eval(tt.rule.Condition, tt.torrent)
			if err != nil {
				t.Errorf("Eval error = %v", err)
				return
			}
			if got != tt.want {
				t.Errorf("Rules() = %v, want %v", got, tt.want)
			}
		})
	}

}
