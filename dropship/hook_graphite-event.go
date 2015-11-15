package dropship

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"text/template"
	"time"
)

type GraphiteEventHook struct{}

func (h GraphiteEventHook) Execute(config HookConfig, service Config) error {
	host, ok := config["host"]
	if !ok {
		return fmt.Errorf("Graphite Hook: unable to call graphite invalid host provided %v", config["host"])
	}
	delete(config, "host")

	config["when"] = string(time.Now().Unix())

	if w, ok := config["what"]; ok {
		what, err := parseTemplate(w, service)
		if err != nil {
			return err
		}

		config["what"] = what
	}

	if d, ok := config["data"]; ok {
		data, err := parseTemplate(d, service)
		if err != nil {
			return err
		}
		config["data"] = data
	}

	config["tags"] = config["tags"] + " " + service.Name

	body, err := json.Marshal(config)
	if err != nil {
		return fmt.Errorf("Graphite Hook: %s", err)
	}

	resp, err := http.Post(host+"/events/", "application/json", bytes.NewReader(body))

	if err != nil {
		return fmt.Errorf("Graphite Hook: %s", err)
	}

	if resp.StatusCode >= 400 {
		return fmt.Errorf("Graphite Hook: unable to post to events. responded with %d", resp.StatusCode)
	}

	return nil
}

func parseTemplate(temp string, service Config) (string, error) {
	tmpl, err := template.New("data").Parse(temp)
	if err != nil {
		return "", err
	}

	hostname, err := os.Hostname()
	if err != nil {
		return "", err
	}

	data := TemplateData{service, hostname}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	return buf.String(), err
}
