package gsm

import (
	"encoding/json"
	"os"
	"regexp"
)

import (
	"github.com/streadway/amqp"
	gsmlog "gsm/log"
)

type Instructions struct {
	Rev          string `json:"rev"`
	RepoPath     string `json:"repo_path"`
	RepoBasename string `json:"repo_basename"`
}

func Orchestrate(delivery amqp.Delivery, logger gsmlog.GsmLogger) (*Instructions, error) {
	defer func() {
		err := delivery.Ack(false)
		if err != nil {
			logger.Printf("Error acking delivery %+v: %+v\n", delivery, err)
			os.Exit(6)
		}
	}()

	var applicationJsonRegex = regexp.MustCompile(`application/json`)
	var instructions *Instructions
	var err error

	switch {
	case applicationJsonRegex.MatchString(delivery.ContentType):
		instructions, err = parseJson(delivery.Body)
	default:
		instructions, err = parseJson(delivery.Body)
	}

	if err != nil {
		return nil, err
	}

	return instructions, nil
}

func parseJson(rawBody []byte) (*Instructions, error) {
	var ret *Instructions
	err := json.Unmarshal(rawBody, ret)
	if err != nil {
		return nil, err
	}

	return ret, nil
}
