package internal

import (
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func GetServers(url string) ([][2]string, error) {
	r, err := httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	rd := csv.NewReader(r.Body)
	rd.Comma = '\t'
	rd.FieldsPerRecord = 2

	var servers [][2]string
	for {
		record, err := rd.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		servers = append(servers, [2]string{record[1], record[0]})
	}

	return servers, nil
}

func GetString(configURL string) (string, error) {
	r, err := httpClient.Get(configURL)
	if err != nil {
		return "", err
	}
	defer r.Body.Close()

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func GetIP() (string, error) {
	var (
		s   string
		err error
	)
	// connection sometimes drops after VPN switch
	// retry requests to allow for this
	for i := 0; i < 10; i++ {
		if i > 0 {
			fmt.Printf("Error getting IP. Retrying in %d seconds...\n", i)
			time.Sleep(time.Duration(i) * time.Second)
		}
		if s, err = GetString("http://icanhazip.com/"); err == nil {
			break
		}
	}
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(s), nil
}

const defaultTimeout = 30 * time.Second

var httpClient = &http.Client{
	Timeout: defaultTimeout,
	Transport: &http.Transport{
		DisableKeepAlives: true,
	},
}
