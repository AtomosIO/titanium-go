package titanium

import (
	"encoding/json"
	"fmt"
	"strings"
)

var _ = fmt.Printf

type Instance struct {
	Command    string `json:"command"`
	Executable string
	Arguments  []string
	Directory  string
}

// Retreives the instance information associated with the token this client was
// created with.
func (client *HttpClient) GetTokenInstance() (Instance, error) {
	var output Instance

	// Get and unmarshal
	data, err := client.get(InstancesEndpoint)
	if err != nil {
		return output, err
	}
	err = json.Unmarshal(data, &output)
	if err != nil {
		return output, err
	}

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
