package utils

import (
	"math/rand"

	"go.uber.org/zap"
)

const symbols = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func GenShortUrl(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = symbols[rand.Intn(len(symbols))]
	}
	return string(b)
}

func ToAddr(baseUrl, str string) string {
	return baseUrl + "/" + str
}

// test sugar is work
func Info(sugar *zap.SugaredLogger) {
	sugar.Info("hello", "world")
}

func GetSymbols() string {
	return symbols
}
