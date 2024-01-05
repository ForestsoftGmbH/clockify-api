package report

import (
	clockify_api "github.com/ForestsoftGmbH/clockify-api"
	"github.com/stretchr/testify/assert"
	"net/http"
	"os"
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

		api := NewReportApi()
		api.Credentials = clockify_api.Credentials{WorkspaceId: os.Getenv("CLOCKIFY_WORKSPACE"), ApiKey: os.Getenv("CLOCKIFY_API_KEY")}
		api.HttpClient = &MockClient{
			DoFunc: func(req *http.Request) (*http.Response, error) {
				// do whatever you want
				return &http.Response{
					StatusCode: http.StatusBadRequest,
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
	api.HttpClient = &MockClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			// do whatever you want
			return &http.Response{
				StatusCode: http.StatusBadRequest,
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
