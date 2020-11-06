package test

import "net/url"

func ParseUrlOrDie(s string) *url.URL {
	parse, err := url.Parse(s)
	if err != nil {
		panic(err)
	}
	return parse
}
