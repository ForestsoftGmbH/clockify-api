package report

import (
	"bytes"
	"encoding/json"
	"github.com/ForestsoftGmbH/clockify-api/timeentries"
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

func TestGetReport(t *testing.T) {
	t.Run("GetReport", func(t *testing.T) {
		req := ReportRequest{
			Client:         nil,
			Billable:       true,
			InvoicingState: "UNINVOICED",
			DateRangeStart: "2023-12-01T00:00:00.000Z",
			DateRangeEnd:   "2023-12-31T23:59:59.999Z",
			DetailedFilter: ReportDetailedFilter{
				Page:         1,
				PageSize:     500,
				SortColumn:   "DATE",
				SortOrder:    "ASCENDING",
				FilterColumn: "DATE",
				FilterQuery:  "2023-12-01",
			},
		}

		apiResponse := ReportResponse{
			TimeEntries: []timeentries.TimeEntry{
				{
					Id:          "5f8b7b4b9e57ad0b4c8b4567",
					Description: "Test",
				},
			},
		}

		api := NewReportApi()
		api.HttpClient = &MockClient{
			DoFunc: func(req *http.Request) (*http.Response, error) {
				jsonByte, _ := json.Marshal(apiResponse)
				// do whatever you want
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader(jsonByte)),
				}, nil
			},
		}

		report, err := api.GetReport(req)
		if err != nil {
			t.Errorf("Error: %v", err)
		}

		if report == nil {
			t.Errorf("report is nil")
		}

	})
}

func TestSearchClient(t *testing.T) {
	api := NewReportApi()
	searchResponse := []SearchResult{
		{
			Id:   "644d02af2d46af108391309e",
			Name: "Enbitcon GmbH",
		},
	}
	searchResponseJson, _ := json.Marshal(searchResponse)
	api.HttpClient = &MockClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			// do whatever you want
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader(searchResponseJson)),
			}, nil
		},
	}

	clientId, err := api.SearchClient("Enbitcon GmbH")
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	if clientId == "" {
		t.Errorf("client id could not found")
	}

	assert.Equal(t, "644d02af2d46af108391309e", clientId)

}
