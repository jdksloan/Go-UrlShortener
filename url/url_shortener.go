package url

import (
	"fmt"
	"math"
)

var baseMap = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func ShortenURL(id int) (string, error) {
	url, err := baseConvert(id, baseMap)
	if err != nil {
		return "", fmt.Errorf("error converting URL: %v", err)
	}
	return url, nil
}

func baseConvert(i int, stringMap string) (string, error) {
	if i == 0 {
		return stringMap[:1], nil
	}
	short := ""
	base := len(stringMap)

	for i > 0 {
		r := i % base
		short += string(stringMap[r])
		i = int(math.Floor(float64(i / base)))
	}
	return short, nil
}
