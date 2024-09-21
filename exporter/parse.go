package exporter

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
)

type subscriptionUserinfo struct {
	upload   int64
	download int64
	total    int64
	expire   int64
}

func parse(url string) (*subscriptionUserinfo, error) {
	// Verify if the actual URL starts with 'http'
	if !strings.HasPrefix(url, "http") {
		return nil, errors.New("invalid URL: must start with 'http'")
	}

	// Send HEAD request
	resp, err := http.Head(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("HTTP status code is " + strconv.Itoa(resp.StatusCode) + ", not 200")
	}

	// Get Subscription-Userinfo header
	userInfo := resp.Header.Get("Subscription-Userinfo")
	if userInfo == "" {
		return nil, errors.New("Subscription-Userinfo header not found")
	}

	// Parse Subscription-Userinfo
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
