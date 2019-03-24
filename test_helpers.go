package qiita

import (
	"bufio"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
	"testing"
)

func setup(t *testing.T, mockFilesBaseDir, mockResponseHeaderFile, mockResponseBodyFile, expectedMethod, expectedRequestPath, expectedRawQuery string) (*Client, func()) {
	server := newTestServer(t, mockFilesBaseDir, mockResponseHeaderFile, mockResponseBodyFile, expectedMethod, expectedRequestPath, expectedRawQuery)

	serverURL, err := url.Parse(server.URL)
	if !assert.Nil(t, err) {
		t.FailNow()
	}
	cli := &Client{
		URL:        serverURL,
		HTTPClient: server.Client(),
		Logger:     log.New(ioutil.Discard, "", 0),
	}

	teardown := func() {
		server.Close()
	}

	return cli, teardown
}

func newTestServer(t *testing.T, mockFilesBaseDir, mockResponseHeaderFile, mockResponseBodyFile, expectedMethod, expectedRequestPath, expectedRawQuery string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if !assert.Equal(t, expectedMethod, req.Method) {
			t.FailNow()
		}
		if !assert.Equal(t, expectedRequestPath, req.URL.Path) {
			t.FailNow()
		}
		if !assert.Equal(t, expectedRawQuery, req.URL.RawQuery) {
			t.FailNow()
		}

		headerPath := path.Join(mockFilesBaseDir, mockResponseHeaderFile)
		statusCode, kvs := parseHeader(t, headerPath)
		for k, v := range kvs {
			w.Header().Set(k, v)
		}
		w.WriteHeader(statusCode)

		bodyPath := path.Join(mockFilesBaseDir, mockResponseBodyFile)
		body := parseBody(t, bodyPath)

		_, _ = w.Write(body)
	}))
}

func parseHeader(t *testing.T, headerPath string) (int, map[string]string) {
	t.Helper()

	h, err := os.Open(headerPath)
	if !assert.Nil(t, err) {
		t.FailNow()
	}

	sc := bufio.NewScanner(h)

	kvs := make(map[string]string)
	var statusCode int
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "HTTP/2") {
			codeStr := strings.Split(line, " ")[1]
			statusCode, _ = strconv.Atoi(codeStr)
		} else {
			key := strings.Split(line, ": ")[0]
			value := strings.Split(line, ": ")[1]

			kvs[key] = value
		}
	}

	return statusCode, kvs
}

func parseBody(t *testing.T, bodyPath string) []byte {
	f, err := os.Open(bodyPath)
	if !assert.Nil(t, err) {
		t.FailNow()
	}
	b, err := ioutil.ReadAll(f)
	if !assert.Nil(t, err) {
		t.FailNow()
	}

	return b
}
