package event

import (
	"encoding/json"
	"errors"
	"strings"
	"toppr/gjallarhorn/blow/horn"
)

// eventData carries the event data
type eventData map[string]interface{}

// Legal AWS Events
var eventSource map[string]map[string]bool

func init() {
	eventSource = make(map[string]map[string]bool)
	eventSource["aws.elasticache"] = map[string]bool{"": true}
	eventSource["aws.es"] = map[string]bool{"": true}
	eventSource["aws.guardduty"] = map[string]bool{"": true}
	eventSource["aws.health"] = map[string]bool{"": true}
	eventSource["aws.inspector"] = map[string]bool{"": true}
	eventSource["aws.rds"] = map[string]bool{"": true}
	eventSource["aws.secretsmanager"] = map[string]bool{"": true}
	eventSource["aws.securityhub"] = map[string]bool{"": true}
	eventSource["aws.tag"] = map[string]bool{"": true}
	eventSource["aws.waf"] = map[string]bool{"": true}
	eventSource["aws.applicationinsights"] = map[string]bool{"": true}
	eventSource["aws.elasticloadbalancing"] = map[string]bool{"": true}

}

// ProcessMessage processes the raw event string received from AWS events
//
// Inputs:
// 		message is the raw json string extracted from AWS event message
// Output:
// 		If success, nil
// 		otherwise, error
func ProcessMessage(message string) error {
	eventData, err := unmarshalMessageString(message)
	if err != nil {
		return err
	}

	isLegal := isLegalEvent(*eventData)
	if !isLegal {
		return errors.New("Illegal event data")
	}

	// privateHook := "https://hooks.slack.com/services/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"
	postHook := "https://hooks.slack.com/services/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"
	webHook := postHook
	_, err = horn.Slack(webHook, horn.EventData(*eventData))
	if err != nil {
		return err
	}

	return nil
}

func isLegalEvent(ed eventData) bool {
	_, ok := eventSource[ed["source"].(string)]
	if !ok {
		return false
	}
	if strings.Contains(ed["source"].(string), "aws.") {
		return true
		// if _, ok := eventSource[ed["source"].(string)][ed["detail-type"].(string)]; ok {
		// 	return true
		// }
	}
	return false
}

func unmarshalMessageString(message string) (*eventData, error) {
	// eventData := make(map[string]interface{})
	eventData := make(eventData)
	if err := json.Unmarshal([]byte(message), &eventData); err != nil {
		return nil, err
	}
	return &eventData, nil
}
