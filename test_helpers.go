package qiita

import (
	"bufio"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"
)

func newTestServer(t *testing.T, responseFile string, expectedMethod, expectedRequestPath, expectedRawQuery string) *httptest.Server {
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

		headerPath := fmt.Sprintf("./testdata/responses/%s-header", responseFile)
		h, err := os.Open(headerPath)
		if !assert.Nil(t, err) {
			t.FailNow()
		}
		sc := bufio.NewScanner(h)

		var statusCode int
		for sc.Scan() {
			line := sc.Text()
			if len(line) == 0 {
				continue
			}

			if strings.HasPrefix(line, "HTTP/2") {
				codeStr := strings.Split(line, " ")[1]
				statusCode, _ = strconv.Atoi(codeStr)
			} else {
				key := strings.Split(line, ": ")[0]
				value := strings.Split(line, ": ")[1]
				w.Header().Set(key, value)
			}
		}
		w.WriteHeader(statusCode)

		bodyPath := fmt.Sprintf("./testdata/responses/%s-body", responseFile)
		f, err := os.Open(bodyPath)
		if !assert.Nil(t, err) {
			t.FailNow()
		}
		b, err := ioutil.ReadAll(f)
		if !assert.Nil(t, err) {
			t.FailNow()
		}
		_, err = w.Write(b)
		if !assert.Nil(t, err) {
			t.FailNow()
		}
	}))
}
