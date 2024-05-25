package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"nerd/shortener/handlers"
	"nerd/shortener/storage"
	"nerd/shortener/utils"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const (
	DefaultLen = 5
	TestUrl    = "http://testingurl.com"
)

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
	symbols := utils.GetSymbols()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			short := utils.GenShortUrl(tt.args.n)
			assert.Equalf(t, tt.want, len(short), "GenShortUrl(%v)", tt.args.n)
			for _, c := range short {
				assert.NotContains(t, string(c), symbols)
			}
		})
	}
}

func TestInfo(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "1"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			utils.Info(sugar)
		})
	}
}

func TestPostHandler(t *testing.T) {
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
			args: args{method: http.MethodPost, addr: TestUrl},
			want: want{code: http.StatusCreated},
		},
		{
			name: "try to get form page",
			args: args{method: http.MethodGet, addr: ""},
			want: want{code: http.StatusMethodNotAllowed},
		},
		{
			name: "try another method",
			args: args{method: http.MethodConnect, addr: ""},
			want: want{code: http.StatusMethodNotAllowed},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.args.method, "/api/shorten", strings.NewReader(fmt.Sprintf("url=%s", tt.args.addr)))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			rw := httptest.NewRecorder()

			handlers.PostFormHandler(rw, req)

			assert.Equal(t, tt.want.code, rw.Code)
		})
	}
}

func TestGetHandler(t *testing.T) {

	ev, err := events.Get(0)
	if err != nil {
		return
	}
	type args struct {
		method string
		event  *storage.Event
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
			args: args{method: http.MethodPost, event: ev},
			want: want{code: http.StatusMethodNotAllowed},
		},
		{
			name: "try to get non-existing address",
			args: args{method: http.MethodGet, event: &storage.Event{
				Uuid:        "SLT4d",
				ShortUrl:    "",
				OriginalUrl: "",
			}},
			want: want{code: http.StatusBadRequest},
		},
		{
			name: "try to get existing address",
			args: args{method: http.MethodGet, event: ev},
			want: want{code: http.StatusTemporaryRedirect},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.args.method, fmt.Sprintf("/%s", tt.args.event.Uuid), nil)
			rw := httptest.NewRecorder()

			handlers.GetHandler(rw, req)

			assert.Equal(t, tt.want.code, rw.Code)
		})
	}
}

func TestJsonPostFormHandler(t *testing.T) {
	type args struct {
		method     string
		pseudoJson string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "test add address with storage",
			args: args{method: http.MethodPost, pseudoJson: `{"url": "https://testingurl.com"}`},
			want: http.StatusCreated,
		},
		{
			name: "test add address with storage method get",
			args: args{method: http.MethodGet, pseudoJson: `{"url": "https://testingurl.com"}`},
			want: http.StatusMethodNotAllowed,
		},
		{
			name: "test add address with storage but without storage",
			args: args{method: http.MethodPost, pseudoJson: "lskdjfslkdjfsldkfj"},
			want: http.StatusBadRequest,
		},
		{
			name: "test add address with storage but url is empty",
			args: args{method: http.MethodPost, pseudoJson: `"url": ""`},
			want: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.args.method, "/api/shorten", strings.NewReader(tt.args.pseudoJson))
			rw := httptest.NewRecorder()

			handlers.JsonPostFormHandler(rw, req)

			assert.Equal(t, tt.want, rw.Code)
		})
	}
	AfterTest()
}

func TestEvents_Load_Find(t *testing.T) {

	ev, err := events.Get(0)
	if err != nil {
		return
	}
	type args struct {
		uuid string
	}
	tests := []struct {
		name string
		args args
		want *storage.Event
	}{
		{
			name: "load events and find existing event",
			args: args{uuid: ev.Uuid},
			want: ev,
		},
		{
			name: "load events and find non-existing event",
			args: args{uuid: "DrM23"},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ev := events.Find(tt.args.uuid)
			assert.Equal(t, tt.want, ev)
		})
	}
}

func AfterTest() {
	counter := 0
	for _, e := range events.Events {
		if e.OriginalUrl == TestUrl {
			events.Delete(e.Uuid)
			counter += 1
		}
	}
	sugar.Infow("AfterTest", "clean", counter)
}
