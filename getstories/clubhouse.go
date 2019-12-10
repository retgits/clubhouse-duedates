package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	ClubhouseURL          = "https://api.clubhouse.io/api/v3/"
	SearchStoriesEndpoint = "search/stories"
)

type SearchStoriesByOwnerInput struct {
	APIToken string
	Days     int
	Owner    string
}

func SearchStoriesByOwner(s SearchStoriesByOwnerInput) (SearchStoriesOutput, error) {
	startdate := time.Now().Format("2006-01-02")
	enddate := time.Now().AddDate(0, 0, s.Days).Format("2006-01-02")

	url := fmt.Sprintf("%s%s?token=%s", ClubhouseURL, SearchStoriesEndpoint, s.APIToken)

	ssi := SearchStoriesInput{
		PageSize: 25,
		Query:    fmt.Sprintf("due:%s..%s and owner:%s and !is:done", startdate, enddate, s.Owner),
	}

	payload, err := ssi.Marshal()
	if err != nil {
		return SearchStoriesOutput{}, err
	}

	req, err := http.NewRequest("GET", url, bytes.NewReader(payload))
	if err != nil {
		return SearchStoriesOutput{}, err
	}

	req.Header.Add("content-type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return SearchStoriesOutput{}, err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return SearchStoriesOutput{}, err
	}

	sso, err := UnmarshalSearchStoriesOutput(body)
	if err != nil {
		return SearchStoriesOutput{}, err
	}

	return sso, nil
}

func UnmarshalSearchStoriesInput(data []byte) (SearchStoriesInput, error) {
	var r SearchStoriesInput
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *SearchStoriesInput) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type SearchStoriesInput struct {
	PageSize int64  `json:"page_size"`
	Query    string `json:"query"`
}

func UnmarshalSearchStoriesOutput(data []byte) (SearchStoriesOutput, error) {
	var r SearchStoriesOutput
	err := json.Unmarshal(data, &r)
	return r, err
}

type SearchStoriesOutput struct {
	Next    interface{} `json:"next"`
	Stories []Story     `json:"data"`
	Total   int64       `json:"total"`
}

type Story struct {
	AppURL               string        `json:"app_url"`
	Description          string        `json:"description"`
	Archived             bool          `json:"archived"`
	Started              bool          `json:"started"`
	StoryLinks           []interface{} `json:"story_links"`
	EntityType           string        `json:"entity_type"`
	Labels               []interface{} `json:"labels"`
	ExternalTickets      []interface{} `json:"external_tickets"`
	MentionIDS           []interface{} `json:"mention_ids"`
	MemberMentionIDS     []interface{} `json:"member_mention_ids"`
	StoryType            string        `json:"story_type"`
	LinkedFiles          []interface{} `json:"linked_files"`
	CompletedAtOverride  interface{}   `json:"completed_at_override"`
	StartedAt            *string       `json:"started_at"`
	CompletedAt          interface{}   `json:"completed_at"`
	Name                 string        `json:"name"`
	Completed            bool          `json:"completed"`
	Comments             []Comment     `json:"comments"`
	Blocker              bool          `json:"blocker"`
	Branches             []interface{} `json:"branches"`
	EpicID               *int64        `json:"epic_id"`
	PreviousIterationIDS []interface{} `json:"previous_iteration_ids"`
	RequestedByID        string        `json:"requested_by_id"`
	IterationID          interface{}   `json:"iteration_id"`
	Tasks                []Task        `json:"tasks"`
	StartedAtOverride    interface{}   `json:"started_at_override"`
	WorkflowStateID      int64         `json:"workflow_state_id"`
	UpdatedAt            string        `json:"updated_at"`
	GroupMentionIDS      []interface{} `json:"group_mention_ids"`
	SupportTickets       []interface{} `json:"support_tickets"`
	FollowerIDS          []string      `json:"follower_ids"`
	OwnerIDS             []string      `json:"owner_ids"`
	ExternalID           interface{}   `json:"external_id"`
	ID                   int64         `json:"id"`
	Estimate             interface{}   `json:"estimate"`
	Commits              []interface{} `json:"commits"`
	Files                []interface{} `json:"files"`
	Position             int64         `json:"position"`
	Blocked              bool          `json:"blocked"`
	ProjectID            int64         `json:"project_id"`
	Deadline             string        `json:"deadline"`
	CreatedAt            string        `json:"created_at"`
	MovedAt              string        `json:"moved_at"`
}

type Comment struct {
	EntityType       string        `json:"entity_type"`
	StoryID          int64         `json:"story_id"`
	MentionIDS       []string      `json:"mention_ids"`
	AuthorID         string        `json:"author_id"`
	MemberMentionIDS []string      `json:"member_mention_ids"`
	UpdatedAt        string        `json:"updated_at"`
	GroupMentionIDS  []interface{} `json:"group_mention_ids"`
	ExternalID       interface{}   `json:"external_id"`
	ID               int64         `json:"id"`
	Position         int64         `json:"position"`
	Reactions        []interface{} `json:"reactions"`
	CreatedAt        string        `json:"created_at"`
	Text             string        `json:"text"`
}

type Task struct {
	Description      string        `json:"description"`
	EntityType       string        `json:"entity_type"`
	StoryID          int64         `json:"story_id"`
	MentionIDS       []string      `json:"mention_ids"`
	MemberMentionIDS []string      `json:"member_mention_ids"`
	CompletedAt      *string       `json:"completed_at"`
	UpdatedAt        string        `json:"updated_at"`
	GroupMentionIDS  []interface{} `json:"group_mention_ids"`
	OwnerIDS         []interface{} `json:"owner_ids"`
	ExternalID       interface{}   `json:"external_id"`
	ID               int64         `json:"id"`
	Position         int64         `json:"position"`
	Complete         bool          `json:"complete"`
	CreatedAt        string        `json:"created_at"`
}
