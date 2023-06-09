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
			Condition: "Category == 'Bangumi' || any(Tags, {# == 'VCB' || # == 'Public'})",
		},
		torrent: Torrent{
			Category: "Bangumi",
		},
		want: true,
	}, {
		name: "test2",
		rule: Rule{
			Condition: "Category == 'Bangumi' || any(Tags, {# == 'VCB' || # == 'Public'})",
		},
		torrent: Torrent{
			Tags: []string{
				"VCB",
				"Normal",
			},
		},
		want: true,
	}, {
		name: "test2",
		rule: Rule{
			Condition: "Category == 'Bangumi' || any(Tags, {# == 'VCB' || # == 'Public'})",
		},
		torrent: Torrent{
			Tags: []string{
				"Music",
				"Movie",
			},
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
