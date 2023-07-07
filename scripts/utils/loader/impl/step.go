package impl

import (
	"errors"
	"fmt"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/samber/lo"
	vegeta "github.com/tsenart/vegeta/lib"
	"io"
	"strings"
	"time"
)

type StepOption func(*Step)

type Step struct {
	Name      string
	Requester *retryablehttp.Client
	BaseUrl   string
	Path      string
	Method    string
	Timeout   time.Duration
	Retry     bool
	Accept    func(int) bool
}

func (s *Step) Run(row []string) *vegeta.Result {
	start := time.Now()
	dataPath, _ := truncatePath(s.Path, lo.ToAnySlice(row)...)
	uri := fmt.Sprintf("%s%s", s.BaseUrl, dataPath)
	request, _ := retryablehttp.NewRequest(s.Method, uri, nil)

	metrics := &vegeta.Result{
		Attack:    s.Name,
		Timestamp: start,
	}

	defer func() {
		metrics.Latency = time.Since(metrics.Timestamp)
	}()

	resp, err := s.Requester.Do(request)
	if err != nil {
		metrics.Error = err.Error()
		return metrics
	}

	defer resp.Body.Close()
	body := io.Reader(resp.Body)
	if metrics.Body, err = io.ReadAll(body); err != nil {
		metrics.Error = err.Error()
		return metrics
	} else if _, err = io.Copy(io.Discard, body); err != nil {
		metrics.Error = err.Error()
		return metrics
	}

	metrics.BytesIn = uint64(len(metrics.Body))
	if request.ContentLength != -1 {
		metrics.BytesOut = uint64(request.ContentLength)
	}

	if metrics.Code = uint16(resp.StatusCode); s.Accept != nil && !s.Accept(resp.StatusCode) {
		metrics.Error = fmt.Sprintf("Unexpected status code: %d", resp.StatusCode)
		return metrics
	}

	return metrics
}

func NewStep(name, baseUrl string, options ...StepOption) *Step {
	step := &Step{
		Name:    name,
		BaseUrl: baseUrl,
	}

	for _, opt := range options {
		opt(step)
	}
	return step
}

func WithPathFormat(path string) StepOption {
	return func(step *Step) {
		step.Path = path
	}
}

func Withmethod(method string) StepOption {
	return func(step *Step) {
		step.Method = method
	}
}

func WithTimeout(timeout time.Duration) StepOption {
	return func(step *Step) {
		step.Timeout = timeout
	}
}

func WithRetry(retry bool) StepOption {
	return func(step *Step) {
		step.Retry = retry
	}
}

func WithAccept(accept func(int) bool) StepOption {
	return func(step *Step) {
		step.Accept = accept
	}
}

func truncatePath(str string, args ...any) (string, error) {
	n := strings.Count(str, "%s")
	if n > len(args) {
		return "", errors.New("Unexpected string:" + str)
	}
	return fmt.Sprintf(str, args[0:n]...), nil
}
