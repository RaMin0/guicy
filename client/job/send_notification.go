package job

// SendNotification job
var SendNotification sendNotification = "send-notification"

type sendNotification string

func (sendNotification) Name() string {
	return "Send Notification"
}
func (sendNotification) Description() string {
	return "Not at the GUC? No worries, send your students a notification right away!"
}
func (sendNotification) Exec(payload interface{}) (interface{}, error) {
	return map[string]interface{}{
		"status": "Done",
	}, nil
}
func (sendNotification) Inputs() []map[string]string {
	return []map[string]string{
		map[string]string{
			"id":    "subject",
			"type":  "text",
			"label": "Subject",
			"hint":  "Compensation",
		},
		map[string]string{
			"id":    "body",
			"type":  "text",
			"label": "Body",
			"hint":  "This week's lab will be compensated next week.",
		},
	}
}
func (sendNotification) Outputs() []map[string]string {
	return []map[string]string{
		map[string]string{
			"id":    "status",
			"type":  "text",
			"label": "Status",
		},
	}
}
