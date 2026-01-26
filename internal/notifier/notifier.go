package notifier

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"
)

type Client struct {
	URL string
}

func (n Client) Send(event string, activeProver string) error {
	if n.URL == "" {
		return nil
	}
	payload := map[string]any{
		"event":        event,
		"activeProver": activeProver,
		"ts":           time.Now().UTC().Format(time.RFC3339Nano),
	}
	b, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", n.URL, bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}
