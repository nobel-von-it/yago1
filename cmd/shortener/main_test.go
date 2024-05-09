package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"nerd/yago1/cmd/shortener/config"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAddMap(t *testing.T) {
	mp := make(map[string]string)
	type args struct {
		key   string
		value string
	}
	tests := []struct {
		name string
		args args
	}{
		{name: "add new k/v in empty map", args: args{key: "hello", value: "world"}},
		{name: "try to rewrite existing value", args: args{key: "hello", value: "world!!!!"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			AddMap(mp, tt.args.key, tt.args.value)
			assert.Equal(t, tt.args.value, mp[tt.args.key])
		})
	}
}

func TestFindVal(t *testing.T) {
	mp := map[string]string{"k1": "v1", "k2": "v2"}
	type args struct {
		val string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "just try to find key by value 1", args: args{val: "v1"}, want: "k1"},
		{name: "just try to find key by value 2", args: args{val: "v2"}, want: "k2"},
		{name: "try to find key by non-existing value", args: args{val: "asdfsdf"}, want: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, FindVal(mp, tt.args.val), "FindVal(%v, %v)", mp, tt.args.val)
		})
	}
}

func TestGenShortUrl(t *testing.T) {
	type args struct {
		n int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{name: "normal length 0", args: args{n: defaultLen}, want: defaultLen},
		{name: "normal length 1", args: args{n: defaultLen}, want: defaultLen},
		{name: "normal length 2", args: args{n: defaultLen}, want: defaultLen},
		{name: "normal length 3", args: args{n: defaultLen}, want: defaultLen},
		{name: "normal length 4", args: args{n: defaultLen}, want: defaultLen},
		{name: "normal length 5", args: args{n: defaultLen}, want: defaultLen},
		{name: "normal length 6", args: args{n: defaultLen}, want: defaultLen},
		{name: "spec length", args: args{n: 10}, want: 10},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			short := GenShortUrl(tt.args.n)
			assert.Equalf(t, tt.want, len(short), "GenShortUrl(%v)", tt.args.n)
			for _, c := range short {
				assert.NotContains(t, string(c), symbols)
			}
		})
	}
}

func TestPostHandler(t *testing.T) {
	_ = config.ParseArgs()

	type args struct {
		method string
		addr   string
	}
	type want struct {
		code int
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "post new address",
			args: args{method: http.MethodPost, addr: "https://www.youtube.com"},
			want: want{code: http.StatusCreated},
		},
		{
			name: "try to get form page",
			args: args{method: http.MethodGet, addr: ""},
			want: want{code: http.StatusOK},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.args.method, "/", strings.NewReader(fmt.Sprintf("url=%s", tt.args.addr)))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			rw := httptest.NewRecorder()

			PostFormHandler(rw, req)

			assert.Equal(t, tt.want.code, rw.Code)
		})
	}
}

func TestGetHandler(t *testing.T) {
	type addr struct {
		key   string
		value string
	}
	type args struct {
		method string
		addr   addr
	}
	type want struct {
		code int
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "test another method",
			args: args{method: http.MethodPost, addr: addr{}},
			want: want{code: http.StatusBadRequest},
		},
		{
			name: "try to get non-existing address",
			args: args{method: http.MethodGet, addr: addr{key: "https://www.youtube.com", value: "SndTL"}},
			want: want{code: http.StatusBadRequest},
		},
		{
			name: "try to get existing address",
			args: args{method: http.MethodGet, addr: addr{key: "https://www.youtube.com", value: "SndTL"}},
			want: want{code: http.StatusTemporaryRedirect},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.args.method, fmt.Sprintf("/%s", tt.args.addr.value), nil)
			rw := httptest.NewRecorder()

			GetHandler(rw, req)

			assert.Equal(t, tt.want.code, rw.Code)

			AddMap(shoring, tt.args.addr.key, tt.args.addr.value)
		})
	}
}
