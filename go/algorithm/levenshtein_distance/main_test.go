package main

import "testing"

func TestLevenshteinDistance(t *testing.T) {
	type args struct {
		s1 string
		s2 string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			args: args{s1: "lininon", s2: "linonon"},
			want: 1,
		},
		{
			args: args{s1: "sitting", s2: "kitten"},
			want: 3,
		},
		{
			args: args{s1: "kiching", s2: "kitten"},
			want: 4,
		},
		{
			args: args{s1: "abcdefg", s2: "bacdefg"},
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := LevenshteinDistanceV2(tt.args.s1, tt.args.s2); got != tt.want {
				t.Errorf("LevenshteinDistance() = %v, want %v", got, tt.want)
			}
		})
	}
}
