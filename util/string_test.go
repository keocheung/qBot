package util

import "testing"

func TestStringsContain(t *testing.T) {
	type args struct {
		s   []string
		str string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test1",
			args: args{
				s:   []string{"a", "b", "c"},
				str: "c",
			},
			want: true,
		},
		{
			name: "test2",
			args: args{
				s:   []string{"a", "b", "c"},
				str: "d",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StringsContain(tt.args.s, tt.args.str); got != tt.want {
				t.Errorf("StringsContain() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringArraysHasCommon(t *testing.T) {
	type args struct {
		s1 []string
		s2 []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test1",
			args: args{
				s1: []string{"a", "b", "c"},
				s2: []string{"a", "b", "c"},
			},
			want: true,
		},
		{
			name: "test2",
			args: args{
				s1: []string{"a", "b", "c"},
				s2: []string{"d", "e", "f"},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StringArraysHasCommon(tt.args.s1, tt.args.s2); got != tt.want {
				t.Errorf("StringArraysHasCommon() = %v, want %v", got, tt.want)
			}
		})
	}
}
