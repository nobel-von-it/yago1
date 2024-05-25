package storage

import (
	"encoding/json"
	"errors"
	"go.uber.org/zap"
	"nerd/shortener/utils"
	"os"
)

var events Events

type (
	Event struct {
		Uuid        string `storage:"uuid"`
		ShortUrl    string `storage:"short_url"`
		OriginalUrl string `storage:"original_url"`
	}

	Events struct {
		Events []Event `storage:"events"`
	}

	RequestData struct {
		Url string `storage:"url"`
	}

	ResponseData struct {
		Result string `storage:"result"`
	}
)

func (e *Event) String() string {
	return e.Uuid
}

func (es *Events) Get(id int) (*Event, error) {
	if id < 0 || id >= len(events.Events) {
		return nil, errors.New("invalid id")
	}
	return &events.Events[id], nil
}

func (es *Events) Save(storagePath string) error {
	data, err := json.MarshalIndent(es, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(storagePath, data, 0666)
}

func (es *Events) Load(storagePath string) error {
	data, err := os.ReadFile(storagePath)
	if err != nil {
		return err
	}
	if len(data) == 0 {
		es.Events = make([]Event, 0)
		return nil
	}
	if err = json.Unmarshal(data, &es); err != nil {
		return err
	}
	return nil
}

func (es *Events) Add(short, url, storagePath, baseUrl string, sugar *zap.SugaredLogger) {
	es.Events = append(es.Events, Event{
		Uuid:        short,
		ShortUrl:    utils.ToAddr(baseUrl, short),
		OriginalUrl: url,
	})
	if err := es.Save(storagePath); err != nil {
		sugar.Infow("error on save", "err", err)
	}
}

func (es *Events) Find(uuid string) *Event {
	for _, e := range es.Events {
		if e.Uuid == uuid {
			return &e
		}
	}
	return nil
}

func (es *Events) Delete(uuid string) {
	for i, e := range es.Events {
		if e.Uuid == uuid {
			es.Events = append(es.Events[:i], es.Events[i+1:]...)
		}
	}
}

func GetEvents() *Events {
	return &events
}
