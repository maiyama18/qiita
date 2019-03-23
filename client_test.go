package qiita

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		desc        string
		accessToken string
		logger      *log.Logger

		expectedURL    string
		expectedLogger *log.Logger
	}{
		{
			desc:        "success",
			accessToken: "access_token",
			logger:      log.New(os.Stdout, "", log.LstdFlags),

			expectedURL:    BaseURL,
			expectedLogger: log.New(os.Stdout, "", log.LstdFlags),
		},
		{
			desc:        "success_with_no_logger",
			accessToken: "access_token",
			logger:      nil,

			expectedURL:    BaseURL,
			expectedLogger: log.New(ioutil.Discard, "", 0),
		},
		{
			desc:        "success_with_no_access_token",
			accessToken: "",
			logger:      log.New(os.Stdout, "", log.LstdFlags),

			expectedURL:    BaseURL,
			expectedLogger: log.New(os.Stdout, "", log.LstdFlags),
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			cli, err := New(tt.accessToken, tt.logger)
			if !assert.Nil(t, err) {
				t.FailNow()
			}

			assert.Equal(t, cli.URL.String(), tt.expectedURL)
			assert.Equal(t, cli.Logger, tt.expectedLogger)
		})
	}
}
