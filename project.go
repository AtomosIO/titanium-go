package titanium

import (
	"errors"
	"fmt"
	"github.com/atomosio/common"
)

var _ = fmt.Printf

type CreateProjectRequest struct {
	Name   string `json:"name"`
	Public bool   `json:"public"`
}

// Create a project
func (client *HttpClient) CreateProject(projectName string, public bool) error {
	response := Response{}
	// Post and unmarshal
	request := CreateProjectRequest{
		Name:   projectName,
		Public: public,
	}
	err := client.postAndUnmarshal(ProjectsEndpoint, &request, &response)
	if err != nil {
		return err
	}

	if response.Code != common.Success {
		return errors.New("Failed to get cluster information: " + response.Description)
	}

	return nil
}

type UpdateProjectRequest struct {
	Title         string                   `json:"title,omitempty"`
	Description   string                   `json:"description,omitempty"`
	Interfaces    []OutProjectInterface    `json:"interfaces,omitempty"`
	Configuration []OutConfigurationEntity `json:"configuration,omitempty"`
	Kernel        *OutKernel               `json:"kernel,omitempty"`
	Type          string                   `json:"type"`
}

func (client *HttpClient) SetSystem(project string, interfaces []ProjectInterface, entities []ConfigurationEntity) {
	// Convert from ProjectInterface to OutProjectInterface
	outInterfaces := ProjectInterfacesToOutProjectInterfaces(interfaces)
	outConfigurations := ConfigurationEntitiesToOutConfigurationEntities(entities)
	request := UpdateProjectRequest{
		Type:          ProjectTypeToString[ProjectSystemType],
		Interfaces:    outInterfaces,
		Configuration: outConfigurations,
	}

	//send request
	response := Response{}
	addr := fmt.Sprintf("%s/%s", ProjectsEndpoint, project)
	err := client.patchAndUnmarshal(addr, &request, &response)
	if err != nil {
		//TODO Change panic to return error
		panic(err)
	}

	if response.Code != common.Success {
		panic(response.Description)
	}
}

func (client *HttpClient) SetProjectKernel(project string, kernel Kernel) {
	// Convert from ProjectInterface to OutProjectInterface
	outKernel := KernelToOutKernel(kernel)
	request := UpdateProjectRequest{
		Kernel: &outKernel,
	}

	//send request
	response := Response{}
	addr := fmt.Sprintf("%s/%s", ProjectsEndpoint, project)
	err := client.patchAndUnmarshal(addr, &request, &response)
	if err != nil {
		panic(err)
	}

	if response.Code != common.Success {
		panic(response.Description)
	}
}