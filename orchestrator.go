package gsm

import (
	"encoding/json"
	"os"
	"regexp"
)

import (
	"github.com/streadway/amqp"
)

var ApplicationJsonRegex = regexp.MustCompile(`^application/json`)

type Instructions struct {
	Rev       string `json:"rev"`
	RepoName  string `json:"repo_name"`
	RepoOrg   string `json:"repo_org"`
	AuthToken string
}

type Orchestrator struct {
	Configuration
}

func NewOrchestrator(config Configuration) *Orchestrator {
	return &Orchestrator{
		Configuration: config,
	}
}

func (me *Orchestrator) Orchestrate(delivery amqp.Delivery) (*Instructions, error) {
	defer func() {
		err := delivery.Ack(false)
		if err != nil {
			me.Logger.Printf("Error acking delivery %+v: %+v\n", delivery, err)
			os.Exit(6)
		}
	}()

	var instructions *Instructions
	var err error

	switch {
	case ApplicationJsonRegex.MatchString(delivery.ContentType):
		instructions, err = me.parseJson(delivery.Body)
	default:
		instructions, err = me.parseJson(delivery.Body)
	}

	if err != nil {
		return nil, err
	}

	instructions.AuthToken = me.AuthToken

	return instructions, nil
}

func (me *Orchestrator) parseJson(rawBody []byte) (*Instructions, error) {
	var ret *Instructions
	err := json.Unmarshal(rawBody, ret)
	if err != nil {
		return nil, err
	}

	return ret, nil
}
