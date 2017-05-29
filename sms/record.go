package sms

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bcho/lc/client"
)

type RecentRecord struct {
	Phone    *string `json:"phone,omitempty"`
	Message  *string `json:"msg,omitempty"`
	Status   *string `json:"status,omitempty"`
	Created  *string `json:"created,omitempty"`
	Type     *string `json:"type,omitempty"`
	IsRecent *int    `json:"rcnt,omitempty"`
}

type GetRecentRecordsOptions struct {
	StartTime string
	EndTime   string
	Skip      int
	Limit     int
}

type recentRecordsResp struct {
	Count   *int            `json:"count,omitempty"`
	Results []*RecentRecord `json:"results,omitempty"`
}

// GetRecentRecords gets recent sms records.
func GetRecentRecords(c client.Client, appId string, opt *GetRecentRecordsOptions) (int, []*RecentRecord, *http.Response, error) {
	req, err := http.NewRequest("GET", client.UrlSMSRecords(appId), nil)
	if err != nil {
		return 0, nil, nil, err
	}

	q := req.URL.Query()
	q.Add("start_time", opt.StartTime)
	q.Add("end_time", opt.EndTime)
	q.Add("skip", fmt.Sprintf("%d", opt.Skip))
	q.Add("limit", fmt.Sprintf("%d", opt.Limit))
	req.URL.RawQuery = q.Encode()

	resp, err := c.Do(req)
	if err != nil {
		return 0, nil, resp, err
	}
	defer resp.Body.Close()

	var rv recentRecordsResp
	if err := json.NewDecoder(resp.Body).Decode(&rv); err != nil {
		return 0, nil, resp, err
	}

	var count int
	if rv.Count != nil {
		count = *rv.Count
	} else {
		count = 0
	}

	return count, rv.Results, resp, err
}
