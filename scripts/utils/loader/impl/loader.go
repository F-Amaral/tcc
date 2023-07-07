package impl

import (
	"github.com/hashicorp/go-retryablehttp"
	"github.com/samber/lo"
	vegeta "github.com/tsenart/vegeta/lib"
	"log"
)

type LoaderOption func(*Loader)

type Loader struct {
	Dataset   [][]string
	Steps     []*Step
	Requester *retryablehttp.Client
	results   *vegeta.Results
}

func NewLoader(dataset [][]string, options ...LoaderOption) *Loader {
	requester := retryablehttp.NewClient()
	requester.RetryMax = 5
	loader := &Loader{
		Steps:     nil,
		Requester: requester,
		results:   &vegeta.Results{},
	}
	withDataset(dataset)(loader)
	for _, option := range options {
		option(loader)
	}
	return loader
}

func AddStep(step *Step) LoaderOption {
	return func(l *Loader) {
		step.Requester = l.Requester
		l.Steps = append(l.Steps, step)
	}
}

func withDataset(dataset [][]string) LoaderOption {
	if dataset == nil {
		log.Fatal("dataset is nil")
	}
	return func(l *Loader) {
		l.Dataset = dataset
	}
}

func (l *Loader) Run() {
	for _, row := range l.Dataset[1:] {
		for _, step := range l.Steps {
			metricsCh := lo.Async1(func() *vegeta.Result {
				return step.Run(row)
			})

			metrics := <-metricsCh
			if metrics != nil {
				l.results.Add(metrics)
			}
		}
	}
}

func (l *Loader) Results() *vegeta.Results {
	return l.results
}
