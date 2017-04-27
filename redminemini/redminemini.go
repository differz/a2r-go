package redminemini

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"../config"
)

// TimeEntry is redmine API data struct
type TimeEntry struct {
	projectID  int
	issueID    int
	spentON    string
	hours      float64
	activityID int
	comments   string
}

var (
	t TimeEntry
)

// NewTimeEntry init TimeEntry struct
func NewTimeEntry(projectID int, spentON string, hours float64, comments string) *TimeEntry {
	t.projectID = projectID
	t.spentON = spentON
	t.hours = hours
	t.comments = comments

	return &t
}

// TimeEntryString prepear json string for POST
func (t *TimeEntry) TimeEntryString() string {

	ret := `{ "time_entry": {
        "project_id": ` + fmt.Sprintf("%d", t.projectID) + `,
        "spent_on": "` + t.spentON + `",
        "hours": ` + fmt.Sprintf("%.2f ", t.hours) + `,
        "activity_id": 11,
        "comments": "` + t.comments + `" }
}`
	return ret
}

// CreateTimeEntry creates time entry over API
//Parameters:
//    time_entry (required): a hash of the time entry attributes, including:
//        issue_id or project_id (only one is required): the issue id or project id to log time on
//        spent_on: the date the time was spent (default to the current date)
//        hours (required): the number of spent hours
//        activity_id: the id of the time activity. This parameter is required unless a default activity is defined in Redmine.
//        comments: short description for the entry (255 characters max)
//Response:
//    201 Created: time entry was created
//    422 Unprocessable Entity: time entry was not created due to validation failures (response body contains the error messages)
func (t *TimeEntry) CreateTimeEntry() bool {

	cfg := config.New()

	res := "/time_entries.json"
	url := "http://" + cfg.RedmineServer + res

	jsonStr := t.TimeEntryString()

	req, err := http.NewRequest("POST", url, bytes.NewBufferString(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "text/html, application/xhtml+xml, */*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.8,ru;q=0.5,uk;q=0.3")
	req.Header.Set("Accept-Encoding", "deflate") //gzip,
	req.Header.Set("X-Redmine-API-Key", cfg.RedmineKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	// fmt.Println("response Headers:", resp.Header)

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))

	return resp.Status == "201"
}
