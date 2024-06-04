package utils

import (
	"testing"
)

const DefaultLen = 5

func TestGenShortUrl(t *testing.T) {
	type args struct {
		n int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{name: "normal length 0", args: args{n: DefaultLen}, want: DefaultLen}, {name: "normal length 1", args: args{n: DefaultLen}, want: DefaultLen},
		{name: "normal length 2", args: args{n: DefaultLen}, want: DefaultLen},
		{name: "normal length 3", args: args{n: DefaultLen}, want: DefaultLen},
		{name: "normal length 4", args: args{n: DefaultLen}, want: DefaultLen},
		{name: "normal length 5", args: args{n: DefaultLen}, want: DefaultLen},
		{name: "normal length 6", args: args{n: DefaultLen}, want: DefaultLen},
		{name: "spec length", args: args{n: 10}, want: 10},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GenShortUrl(tt.args.n); len(got) != tt.want {
				t.Errorf("GenShortUrl() = %v, want %v", len(got), tt.want)
			}
		})
	}
}

func TestToAddr(t *testing.T) {
	type args struct {
		baseUrl string
		str     string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{name: "default test to address", args: args{baseUrl: "localhost:8080", str: "12345"}, want: "localhost:8080/12345"},
		{name: "test empty address", args: args{baseUrl: "", str: ""}, want: "/"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToAddr(tt.args.baseUrl, tt.args.str); got != tt.want {
				t.Errorf("ToAddr() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestGetSymbols(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{name: "I dont know", want: symbols},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetSymbols(); got != tt.want {
				t.Errorf("GetSymbols() = %v, want %v", got, tt.want)
			}
		})
	}
}
