package titanium

import (
	"errors"
	"fmt"
	"strconv"
	//"strings"
	"time"

	"github.com/atomosio/common"
)

var _ = fmt.Printf

type Instance struct {
	Code        int64  `json:"code"`
	Description string `json:"description"`

	Command      string `json:"command"`
	Stdout       int64
	Stderr       int64
	StdoutString string     `json:"stdout"`
	StderrString string     `json:"stderr"`
	Status       string     `json:"status"`
	Log          []LogEntry `json:"log"`
}

type LogEntry struct {
	Type string `json:"type"`
	// Unix timestamp
	Timestamp int64  `json:"timestamp"`
	Comment   string `json:"comment"`
}

type UpdateInstanceRequest struct {
	Status int16  `json:"status,omitempty"`
	Log    string `json:"log,omitempty"`
	Error  string `json:"error,omitempty"`
}

const (
	InstanceInvalidStatus = iota
	InstanceWaitingStatus // Waiting for some criteria to be met before starting
	InstanceQueuedStatus  // Instance has been queued to start up
	InstanceActiveStatus  // Instance is active
	InstanceStoppedStatus // Stopped

	SpinSleepDuration = time.Millisecond * 1000
)

var (
	ErrInstanceWaitForFinishTimeout = errors.New("Instance did not finish within timeout period")
)

var (
	InstanceStatusStrings = []string{
		"Invalid",
		"Waiting",
		"Queued",
		"Active",
		"Stopped",
	}
)

var (
	EventStrings = []string{
		"Invalid",
		"Waiting",
		"Queued",
		"Started",
		"Stopped",
		"Error",
		"Log",
		"Shutdown",
	}
)

// Retreives the instance information associated with the token this client was
// created with.
func (client *HttpClient) GetTokenInstance() (Instance, error) {
	return client.GetInstance(0)
}

func (client *HttpClient) SetInstanceActive(instanceId int64) error {
	// Get and unmarshal
	request := UpdateInstanceRequest{
		Status: InstanceActiveStatus,
	}
	response := &Response{}
	addr := fmt.Sprintf("%s%d", InstancesEndpoint, instanceId)
	err := client.DoMethodAndUnmarshal("PATCH", addr, request, response)
	if err != nil {
		return err
	}
	if response.Code != common.Success {
		return errors.New(response.Description)
	}

	return nil
}

func (client *HttpClient) SetInstanceStopped(instanceId int64) error {
	// Get and unmarshal
	request := UpdateInstanceRequest{
		Status: InstanceStoppedStatus,
	}
	response := &Response{}
	addr := fmt.Sprintf("%s%d", InstancesEndpoint, instanceId)
	err := client.DoMethodAndUnmarshal("PATCH", addr, request, response)
	if err != nil {
		return err
	}
	if response.Code != common.Success {
		return errors.New(response.Description)
	}

	return nil
}

func (client *HttpClient) LogInstanceComment(instanceId int64, comment string) error {
	// Get and unmarshal
	request := UpdateInstanceRequest{
		Log: comment,
	}
	response := &Response{}
	addr := fmt.Sprintf("%s%d", InstancesEndpoint, instanceId)
	err := client.DoMethodAndUnmarshal("PATCH", addr, request, response)
	if err != nil {
		return err
	}
	if response.Code != common.Success {
		return errors.New(response.Description)
	}

	return nil
}

func (client *HttpClient) LogInstanceError(instanceId int64, comment string) error {
	// Get and unmarshal
	request := UpdateInstanceRequest{
		Error: comment,
	}
	response := &Response{}
	addr := fmt.Sprintf("%s%d", InstancesEndpoint, instanceId)
	err := client.DoMethodAndUnmarshal("PATCH", addr, request, response)
	if err != nil {
		return err
	}
	if response.Code != common.Success {
		return errors.New(response.Description)
	}

	return nil
}

func (client *HttpClient) GetInstance(instanceId int64) (Instance, error) {
	var output Instance

	// Get and unmarshal
	addr := fmt.Sprintf("%s%d", InstancesEndpoint, instanceId)
	err := client.DoEmptyMethodAndUnmarshal("GET", addr, &output)
	if err != nil {
		return output, err
	}
	if output.Code != common.Success {
		return output, errors.New(output.Description)
	}
	// TODO Check reponse to make sure operation succeeded

	output.Stderr, err = strconv.ParseInt(output.StderrString, 10, 64)
	if err != nil {
		return output, err
	}
	output.Stdout, err = strconv.ParseInt(output.StdoutString, 10, 64)
	if err != nil {
		return output, err
	}
	return output, nil
}

func (client *HttpClient) WaitForInstanceToFinish(id int64, timeout time.Duration) error {
	waitTill := time.Now().Add(timeout)

	for {
		// Get cluster information
		instance, err := client.GetInstance(id)
		if err != nil {
			return err
		}

		// If we're done, exit function
		if instance.IsStopped() {
			return nil
		}

		if time.Now().After(waitTill) {
			return ErrInstanceWaitForFinishTimeout
		}

		// Go to sleep for a bit
		time.Sleep(SpinSleepDuration)
	}
}

func (instance Instance) IsWaiting() bool {
	return instance.Status == InstanceStatusStrings[InstanceWaitingStatus]
}

func (instance Instance) IsActive() bool {
	return instance.Status == InstanceStatusStrings[InstanceActiveStatus]
}

func (instance Instance) IsStopped() bool {
	return instance.Status == InstanceStatusStrings[InstanceStoppedStatus]
}

func (instance Instance) IsShuttingDown() bool {
	shutDownEventLast := false

	for _, entry := range instance.Log {
		switch entry.Type {
		case "Shutdown":
			shutDownEventLast = true
		case "Waiting", "Queued", "Started":
			// If we have had another event that causes an instance start, we aren't
			// shutting down.
			shutDownEventLast = false
		}
	}

	return shutDownEventLast
}
