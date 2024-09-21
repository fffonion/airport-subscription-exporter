package exporter

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

type subscriptionUserinfo struct {
	upload   int64
	download int64
	total    int64
	expire   int64
}

type cacheEntry struct {
	info      *subscriptionUserinfo
	timestamp time.Time
}

var (
	cache      = make(map[string]cacheEntry)
	cacheMutex sync.RWMutex
)

func parse(url string, updateIntervalSeconds int) (*subscriptionUserinfo, error) {
	cacheMutex.RLock()
	entry, exists := cache[url]
	cacheMutex.RUnlock()

	if exists && time.Since(entry.timestamp) < time.Duration(updateIntervalSeconds)*time.Second {
		log.Printf("Debug: Returning cached result for URL: %s (%d seconds before)\n",
			url, time.Since(entry.timestamp)/time.Second)
		return entry.info, nil
	}

	info, err := fetchAndParse(url)
	if err != nil {
		return nil, err
	}

	cacheMutex.Lock()
	cache[url] = cacheEntry{
		info:      info,
		timestamp: time.Now(),
	}
	cacheMutex.Unlock()

	return info, nil
}

func fetchAndParse(url string) (*subscriptionUserinfo, error) {
	if !strings.HasPrefix(url, "http") {
		return nil, errors.New("invalid URL: must start with 'http'")
	}

	resp, err := http.Head(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("HTTP status code is " + strconv.Itoa(resp.StatusCode) + ", not 200")
	}

	userInfo := resp.Header.Get("Subscription-Userinfo")
	if userInfo == "" {
		return nil, errors.New("Subscription-Userinfo header not found")
	}

	info := &subscriptionUserinfo{}
	pairs := strings.Split(userInfo, "; ")
	for _, pair := range pairs {
		kv := strings.SplitN(pair, "=", 2)
		if len(kv) != 2 {
			continue
		}
		key, value := kv[0], kv[1]
		intValue, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			continue
		}
		switch key {
		case "upload":
			info.upload = intValue
		case "download":
			info.download = intValue
		case "total":
			info.total = intValue
		case "expire":
			info.expire = intValue
		}
	}

	return info, nil
}
