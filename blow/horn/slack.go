package horn

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/slack-go/slack"
)

// EventData carries the event data
type EventData map[string]interface{}

// Slack posts the EventData to the webhook url
func Slack(webHook string, ed EventData) (bool, error) {
	msgStr, err := formMessageString(ed)
	if err != nil {
		return false, err
	}

	res, err := postMessage(webHook, msgStr)
	if err != nil {
		return false, err
	}
	defer res.Body.Close()

	return true, nil

}

func formMessageString(ed EventData) (string, error) {
	indentStr := "    "
	eventTitle := slack.NewTextBlockObject("mrkdwn", ed["source"].(string), false, false)
	sourceSection := slack.NewContextBlock(
		"",
		[]slack.MixedElement{eventTitle}...,
	)
	eventDescription := slack.NewTextBlockObject("mrkdwn", "*"+ed["detail-type"].(string)+"*", false, false)
	descriptionSection := slack.NewSectionBlock(eventDescription, nil, nil)

	timeField := slack.NewTextBlockObject("mrkdwn", "*Time:*\n"+ed["time"].(string), false, false)
	regionField := slack.NewTextBlockObject("mrkdwn", "*Region:*\n"+ed["region"].(string), false, false)
	// fieldSlice := make([]*slack.TextBlockObject, 0)
	// fieldSlice = append(fieldSlice, timeField)
	// fieldSlice = append(fieldSlice, regionField)
	// fieldsSection := slack.NewSectionBlock(nil, fieldSlice, nil)
	fieldsSection := slack.NewContextBlock(
		"",
		[]slack.MixedElement{timeField, regionField}...,
	)

	resourceList, err := json.MarshalIndent(ed["resources"], "", indentStr)
	if err != nil {
		return "", err
	}
	resourceText := slack.NewTextBlockObject("mrkdwn", "*Resources:*\n"+string(resourceList), false, false)
	// resourceSection := slack.NewSectionBlock(resourceText, nil, nil)
	resourceSection := slack.NewContextBlock(
		"",
		[]slack.MixedElement{resourceText}...,
	)

	divSection := slack.NewDividerBlock()

	detailMap, err := json.MarshalIndent(ed["detail"], "", indentStr)
	if err != nil {
		return "", err
	}
	detailText := slack.NewTextBlockObject("mrkdwn", "```"+string(detailMap)+"```", false, false)
	// detailSection := slack.NewSectionBlock(detailText, nil, nil)
	detailSection := slack.NewContextBlock(
		"",
		[]slack.MixedElement{detailText}...,
	)

	mentionText := slack.NewTextBlockObject("mrkdwn", "@raam", false, false)
	mentionSection := slack.NewSectionBlock(mentionText, nil, nil)

	msg := slack.NewBlockMessage(
		sourceSection,
		descriptionSection,
		divSection,
		fieldsSection,
		resourceSection,
		detailSection,
		mentionSection,
	)

	msgStr, err := json.MarshalIndent(msg, "", indentStr)
	if err != nil {
		return "", err
	}

	return string(msgStr), nil
}

func postMessage(webhook string, msg string) (http.Response, error) {
	var msgBytes = []byte(msg)
	var res *http.Response
	req, err := http.NewRequest("post", webhook, bytes.NewBuffer(msgBytes))
	if err != nil {
		return *res, err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	res, err = client.Do(req)
	if err != nil {
		return *res, err
	}
	defer res.Body.Close()

	return *res, nil
}
