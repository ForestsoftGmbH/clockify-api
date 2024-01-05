package timeentries

import (
	"bytes"
	"encoding/json"
	"fmt"
	clockify_api "github.com/ForestsoftGmbH/clockify-api"
	"io"
	"net/http"
	"strconv"
)

type Api clockify_api.Api

type TimeEntry struct {
	Id          string   `json:"_id"`
	Description string   `json:"description,omitempty"`
	ClientName  string   `json:"clientName,omitempty"`
	ProjectName string   `json:"projectName,omitempty"`
	Rate        float32  `json:"rate,omitempty"`
	Amount      float32  `json:"amount,omitempty"`
	Interval    Interval `json:"timeInterval,omitempty"`
}

type Interval struct {
	Start    string  `json:"start,omitempty"`
	End      string  `json:"end,omitempty"`
	Duration float64 `json:"duration,omitempty"`
}

type TimeEntryApi struct {
	clockify_api.Api
}

func ConvertTimeEntriesToInvoicedRequest(timeEntries []TimeEntry) InvoicedRequest {
	var timeEntryIds []string
	for _, timeEntry := range timeEntries {
		timeEntryIds = append(timeEntryIds, timeEntry.Id)
	}
	return InvoicedRequest{
		Invoiced:     true,
		TimeEntryIds: timeEntryIds,
	}
}

func NewTimeEntryApi() *TimeEntryApi {
	return &TimeEntryApi{Api: *clockify_api.NewApi()}
}

func (a TimeEntryApi) MarkInvoiced(request InvoicedRequest) (bool, error) {
	invoicedJSON, err := json.Marshal(request)
	//https://api.clockify.me/api/v1/workspaces/{workspaceId}/time-entries/invoiced
	url := fmt.Sprintf("https://api.clockify.me/api/v1/workspaces/%s/time-entries/invoiced", a.Credentials.WorkspaceId)

	req, err := http.NewRequest(http.MethodPatch, url, bytes.NewBuffer(invoicedJSON))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Api-Key", a.Credentials.ApiKey)

	if err != nil {
		fmt.Errorf("Could not make PATCH request to invoice timeentries: %v", err)
		return false, err
	}

	resp, err := a.HttpClient.Do(req)

	if resp.Body != nil {
		defer resp.Body.Close()
	}

	if resp.StatusCode != 200 {
		var answer string
		if resp.Body != nil {
			answerBytes, _ := io.ReadAll(resp.Body)
			answer = string(answerBytes)
		} else {
			answer = resp.Status
		}
		err = fmt.Errorf("PATCH request failed to invoice timeentries: %v", strconv.Itoa(resp.StatusCode)+" "+string(answer))
		return false, err
	}

	return true, nil
}

type InvoicedRequest struct {
	Invoiced     bool     `json:"invoiced"`
	TimeEntryIds []string `json:"timeEntryIds"`
}
