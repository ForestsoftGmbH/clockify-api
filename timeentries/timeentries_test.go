package timeentries

import (
	"bytes"
	"errors"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"testing"
)

type MockClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *MockClient) Do(req *http.Request) (*http.Response, error) {
	if m.DoFunc != nil {
		return m.DoFunc(req)
	}
	// just in case you want default correct return value
	return &http.Response{}, nil
}

func TestMarkInvoiced(t *testing.T) {
	t.Run("Bad Request", func(t *testing.T) {
		request := ConvertTimeEntriesToInvoicedRequest([]TimeEntry{{Id: "65858d739e57ad1046b8bc36"}})
		api := NewTimeEntryApi()
		api.HttpClient = &MockClient{
			DoFunc: func(req *http.Request) (*http.Response, error) {
				// do whatever you want
				return &http.Response{
					StatusCode: http.StatusBadRequest,
					Status:     "Bad Request",
				}, nil
			},
		}
		result, err := api.MarkInvoiced(request)
		assert.NotNil(t, err, "Error should not nil")
		assert.Equal(t, errors.New("PATCH request failed to invoice timeentries: 400 Bad Request"), err)
		assert.False(t, result, "Result should be false")

	})

	t.Run("Mark Successful", func(t *testing.T) {
		request := ConvertTimeEntriesToInvoicedRequest([]TimeEntry{{Id: "65858d739e57ad1046b8bc36"}})
		api := NewTimeEntryApi()
		api.HttpClient = &MockClient{
			DoFunc: func(req *http.Request) (*http.Response, error) {
				// do whatever you want
				return &http.Response{
					StatusCode: http.StatusOK,
					Status:     "",
					Body:       io.NopCloser(bytes.NewReader([]byte(""))),
				}, nil
			},
		}
		result, err := api.MarkInvoiced(request)
		assert.Nil(t, err, "Error is not nil")
		assert.True(t, result, "Result should be false")

	})
}
