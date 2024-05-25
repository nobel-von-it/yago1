package handlers

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"nerd/shortener/flags"
	"nerd/shortener/storage"
	"nerd/shortener/utils"
	"net/http"
)

var (
	config = flags.GetConfig()
	sugar  = flags.GetSugar()
	events = storage.GetEvents()
)

const defaultLen = 5

func PostFormHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	url := chi.URLParam(r, "url")
	if len(url) == 0 {
		url = r.FormValue("url")
		if len(url) == 0 {
			http.Error(w, "url required", http.StatusBadRequest)
			return
		}
	}
	short := utils.GenShortUrl(defaultLen)
	addr := utils.ToAddr(config.BaseUrl, short)
	events.Add(short, url, config.StoragePath, config.BaseUrl, sugar)

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write([]byte("From: " + addr + " To: " + url + "\n"))

	sugar.Infoln("From:", addr, "To:", url)
}

func JsonPostFormHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Need post method", http.StatusMethodNotAllowed)
		return
	}

	var req storage.RequestData
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Parse storage error", http.StatusBadRequest)
		return
	}

	sugar.Infow("Got storage:", "url", req.Url)

	if req.Url == "" {
		http.Error(w, "Url is empty", http.StatusBadRequest)
		return
	}
	short := utils.GenShortUrl(defaultLen)
	addr := utils.ToAddr(config.BaseUrl, short)

	events.Add(short, req.Url, config.StoragePath, config.BaseUrl, sugar)

	w.Header().Set("content-type", "application/storage")
	w.WriteHeader(http.StatusCreated)

	res := storage.ResponseData{
		Result: addr,
	}

	if err := json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, "storage encode error", http.StatusInternalServerError)
		return
	}
	sugar.Infoln("From:", addr, "To:", req.Url)
}
