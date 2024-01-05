package report

import (
	"bytes"
	"encoding/json"
	"fmt"
	clockify_api "github.com/ForestsoftGmbH/clockify-api"
	"github.com/ForestsoftGmbH/clockify-api/timeentries"
	"io"
	"log"
	"net/http"
	url2 "net/url"
)

type ClientFilter struct {
	Ids []string `json:"ids,omitempty"`
}
type ReportRequest struct {
	Client         *ClientFilter        `json:"clients,omitempty"`
	Billable       bool                 `json:"billable"`
	InvoicingState string               `json:"invoicingState"`
	DateRangeStart string               `json:"dateRangeStart"`
	DateRangeEnd   string               `json:"dateRangeEnd"`
	DetailedFilter ReportDetailedFilter `json:"detailedFilter"`
}

type ReportDetailedFilter struct {
	Page         int    `json:"page"`
	PageSize     int    `json:"pageSize"`
	SortColumn   string `json:"sortColumn"`
	SortOrder    string `json:"sortOrder"`
	FilterColumn string `json:"filterColumn"`
	FilterQuery  string `json:"filterQuery"`
}

type ReportResponse struct {
	TimeEntries []timeentries.TimeEntry `json:"timeentries"`
}

type ReportApi struct {
	clockify_api.Api
}

func NewReportApi() *ReportApi {
	return &ReportApi{Api: *clockify_api.NewApi()}
}

func (a ReportApi) GetReport(request ReportRequest) (*ReportResponse, error) {
	//make a http request to https://reports.api.clockify.me/v1/workspaces/{workspaceId}/reports/detailed
	personJSON, err := json.Marshal(request)

	url := fmt.Sprintf("https://reports.api.clockify.me/v1/workspaces/%s/reports/detailed", a.Credentials.WorkspaceId)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(personJSON))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Api-Key", a.Credentials.ApiKey)

	if err != nil {
		fmt.Errorf("Could not make POST request to timeentries: %v", err)
		return nil, err
	}

	resp, err := a.HttpClient.Do(req)

	if err != nil {
		fmt.Errorf("Could not perform request to timeentries: %v", err)
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		answer, err := io.ReadAll(resp.Body)
		err = fmt.Errorf("Could not make POST request to timeentries: %v", resp.Status+" "+string(answer))
		return nil, err
	}

	report := &ReportResponse{}

	// But for good measure, let's look at the response body.
	derr := json.NewDecoder(resp.Body).Decode(report)

	if derr != nil {
		fmt.Errorf("Could not decode response body: %v", derr)
		return nil, derr
	}
	return report, nil
}

type SearchResult struct {
	Id           string `json:"id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	WorkspaceId  string `json:"workspaceId"`
	Archived     bool   `json:"archived"`
	Address      string `json:"address"`
	Note         string `json:"note"`
	CurrencyId   string `json:"currencyId"`
	CurrencyCode string `json:"currencyCode"`
}

func (a ReportApi) SearchClient(s string) (string, error) {
	url := fmt.Sprintf("https://api.clockify.me/api/v1/workspaces/%s/clients?name=%s", a.Credentials.WorkspaceId, url2.QueryEscape(s))

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("X-Api-Key", a.Credentials.ApiKey)
	if err != nil {
		return "", err
	}
	resp, err := a.HttpClient.Do(req)

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		answer, err := io.ReadAll(resp.Body)
		err = fmt.Errorf("Could not make request to timeentries for clientid: %v", resp.Status+" "+string(answer))
		return "", err
	}
	result := []SearchResult{}
	json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		log.Fatalln(err)
	}
	if len(result) == 0 {
		return "", fmt.Errorf("Could not find client %s", s)

	}
	return result[0].Id, nil
}
