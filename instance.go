package titanium

import (
	"errors"
	"fmt"
	"github.com/atomosio/common"
	"strings"
)

var _ = fmt.Printf

type Instance struct {
	Code        int64  `json:"code"`
	Description string `json:"description"`

	Command string `json:"command"`
	Stdout  int64  `json:"stdout"`
	Stderr  int64  `json:"stderr"`
	Status  string `json:"status"`

	Executable string
	Arguments  []string
	Directory  string
}

type UpdateInstanceRequest struct {
	Status        int16  `json:"status,omitempty"`
	StatusComment string `json:"statusComment,omitempty"`
	Error         string `json:"error,omitempty"`
}

const (
	InstanceInvalidStatus = iota
	InstanceWaitingStatus // Waiting for some criteria to be met before starting
	InstanceQueuedStatus  // Instance has been queued to start up
	InstanceActiveStatus  // Instance is active
	InstanceStoppedStatus // Stopped
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
	err := client.patchAndUnmarshal(addr, request, response)
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
	err := client.patchAndUnmarshal(addr, request, response)
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
	err := client.patchAndUnmarshal(addr, request, response)
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
	err := client.getAndUnmarshal(addr, &output)
	if err != nil {
		return output, err
	}
	if output.Code != common.Success {
		return output, errors.New(output.Description)
	}
	// TODO Check reponse to make sure operation succeeded

	// Sample command:
	// /atomos/user/project/directory/executable arguments and more arugments
	commandSplits := strings.SplitN(output.Command, " ", 2)
	output.Executable = commandSplits[0]

	lastSeperatorIndex := strings.LastIndex(output.Executable, "/")
	if lastSeperatorIndex != -1 {
		output.Directory = output.Executable[:lastSeperatorIndex]
	}

	if len(commandSplits) == 2 {
		output.Arguments = strings.Split(commandSplits[1], " ")
	}

	return output, nil
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
