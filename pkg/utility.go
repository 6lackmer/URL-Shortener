package pkg

import (
	"log"
	"math/rand"
	"net/url"
	"time"
)

func GenerateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func IsValidURL(urlStr string) bool {
	u, err := url.ParseRequestURI(urlStr)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func GetUniqueShortUrl(allShortUrls []string, length int) string {
	shortURL := GenerateRandomString(length)
	for _, value := range allShortUrls {
		if value == shortURL {
			log.Println("Short URL already exists")
			GetUniqueShortUrl(allShortUrls, length)
		}
	}
	return shortURL
}
