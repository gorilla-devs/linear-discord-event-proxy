// l'action deve contenere i seguenti steps:
// installazione di digitalocean,
// doctl sandbox deploy ./src --verbose-build
package main

import (
	"log"
	"os"

	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type HTTPRequest struct {
	OwHeaders struct {
		Accept         string `json:"accept"`
		AcceptCharset  string `json:"accept-charset"`
		AcceptEncoding string `json:"accept-encoding"`
		CdnLoop        string `json:"cdn-loop"`
		CfConnectingIP string `json:"cf-connecting-ip"`
		CfIpcountry    string `json:"cf-ipcountry"`
		CfRay          string `json:"cf-ray"`
		// CfVisitor      struct {
		// 	Scheme string `json:"scheme"`
		// } `json:"cf-visitor"`
		ContentType              string `json:"content-type"`
		Host                     string `json:"host"`
		LinearDelivery           string `json:"linear-delivery"`
		LinearEvent              string `json:"linear-event"`
		UserAgent                string `json:"user-agent"`
		XDatadogParentID         string `json:"x-datadog-parent-id"`
		XDatadogSampled          string `json:"x-datadog-sampled"`
		XDatadogSamplingPriority string `json:"x-datadog-sampling-priority"`
		XDatadogTraceID          string `json:"x-datadog-trace-id"`
		XForwardedFor            string `json:"x-forwarded-for"`
		XForwardedProto          string `json:"x-forwarded-proto"`
		XRequestID               string `json:"x-request-id"`
	} `json:"__ow_headers"`
	OwMethod  string    `json:"__ow_method"`
	OwPath    string    `json:"__ow_path"`
	Action    string    `json:"action"`
	CreatedAt time.Time `json:"createdAt"`
	Data      struct {
		Assignee struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"assignee"`
		AssigneeID          string        `json:"assigneeId"`
		BoardOrder          int           `json:"boardOrder"`
		CreatedAt           time.Time     `json:"createdAt"`
		CreatorID           string        `json:"creatorId"`
		ID                  string        `json:"id"`
		LabelIds            []interface{} `json:"labelIds"`
		Number              int           `json:"number"`
		PreviousIdentifiers []interface{} `json:"previousIdentifiers"`
		Priority            int           `json:"priority"`
		PriorityLabel       string        `json:"priorityLabel"`
		SortOrder           float64       `json:"sortOrder"`
		StartedAt           time.Time     `json:"startedAt"`
		State               struct {
			Color string `json:"color"`
			ID    string `json:"id"`
			Name  string `json:"name"`
			Type  string `json:"type"`
		} `json:"state"`
		StateID       string   `json:"stateId"`
		SubscriberIds []string `json:"subscriberIds"`
		Team          struct {
			ID   string `json:"id"`
			Key  string `json:"key"`
			Name string `json:"name"`
		} `json:"team"`
		TeamID    string    `json:"teamId"`
		Title     string    `json:"title"`
		UpdatedAt time.Time `json:"updatedAt"`
	} `json:"data"`
	OrganizationID string `json:"organizationId"`
	Type           string `json:"type"`
	UpdatedFrom    struct {
		Title     string    `json:"title"`
		UpdatedAt time.Time `json:"updatedAt"`
	} `json:"updatedFrom"`
	URL string `json:"url"`
}

type Request struct {
	Location string `json:"location"`
}

type Response struct {
	StatusCode int               `json:"statusCode,omitempty"`
	Headers    map[string]string `json:"headers,omitempty"`
	Body       string            `json:"body,omitempty"`
}

func Main(in map[string]interface{}) error {

	log.Println("Started")

	var inReq HTTPRequest
	b, err := json.Marshal(in)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(b, &inReq); err != nil {
		fmt.Println(err)
		return err
	}

	assignee := inReq.Data.Assignee.Name
	taskName := inReq.Data.Title
	updatedAt := inReq.Data.UpdatedAt
	teamName := inReq.Data.Team.Name
	action := inReq.Action
	url := inReq.URL

	var message string = "At " + updatedAt.String() + ", has been performed the " + action + " action to the " + taskName + " task\n assigned to " + assignee + " in the team " + teamName + "\n task URL: " + url

	err = sendRequest(message)
	if err != nil {
		return err
	}

	return nil
}

type DiscordMessage struct {
	Content string `json:"content"`
}

func sendRequest(message string) error {

	URL := os.Getenv("DISCORD_WEBHOOK")

	discordMessage := &DiscordMessage{
		Content: message,
	}

	b, err := json.Marshal(discordMessage)
	if err != nil {
		return err
	}

	client := &http.Client{}

	req, err := http.NewRequest("POST", URL, bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	_, err = client.Do(req)
	if err != nil {
		return err
	}

	return nil
}
