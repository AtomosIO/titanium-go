package titanium

import (
	"fmt"
	"log"
	"os"
)

var _ = fmt.Print
var _ = log.New
var _ = os.Stderr

type DirectionType uint8
type TypeType uint8

const (
	ProjectSystemType = iota
	ProjectKernelType
)

var ProjectTypeToString = []string{
	"system",
	"kernel",
}

var ProjectTypeStrings = map[string]int8{
	ProjectTypeToString[ProjectSystemType]: ProjectSystemType,
	ProjectTypeToString[ProjectKernelType]: ProjectKernelType,
}

type ProjectInterface struct {
	// Name of the project interface
	Name string
	// Human readable description of the interface
	Description string

	// Alias used to connect the interface to an internal node
	Alias     string
	Type      TypeType
	Direction DirectionType
	Optional  bool
}

type ConfigurationEntity struct {
	// Name of the entity
	Name string
	// Long description of the entity
	Description string

	Kernel     string
	Interfaces []ConfigurationEntityInterface
}

type ConfigurationEntityInterface struct {
	Name  string
	Alias string
}

type Kernel struct {
	// File to run when starting kernel. When a command is not empty, the project
	// kernel is configured.
	Command string
	// If image is empty, it defaults to 'default'.
	Image string

	// If Command is not empty, there must be atleast one output-capable interface
	// No two interface can have the same Path or Name.
	Interfaces []KernelInterface
}

type KernelInterface struct {
	Name        string
	Description string

	Path      string
	Type      TypeType
	Direction DirectionType
	Optional  bool
}

// Interface Directions
const (
	InvalidDirection = iota
	InDirection
	OutDirection
	InOutDirection
)

var NodeDirectionStrings = map[string]DirectionType{
	"invalid": InvalidDirection,
	"in":      InDirection,
	"out":     OutDirection,
	"inout":   InOutDirection,
}
var NodeDirectionToStrings = []string{
	"invalid",
	"in",
	"out",
	"inout",
}

// Interface Types
const (
	InvalidType = iota
	FileType
)

var NodeTypeStrings = map[string]TypeType{
	"invalid": InvalidType,
	"file":    FileType,
}
var NodeTypeToStrings = []string{
	"invalid",
	"file",
}

type OutProjectInterface struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Alias       string `json:"alias"`
	Type        string `json:"type"`
	Direction   string `json:"direction"`
	Optional    bool   `json:"optional"`
}

type OutConfigurationEntity struct {
	Name        string                            `json:"name"`
	Description string                            `json:"description"`
	Kernel      string                            `json:"kernel"`
	Interfaces  []OutConfigurationEntityInterface `json:"interfaces"`
}

type OutConfigurationEntityInterface struct {
	Name  string `json:"name"`
	Alias string `json:"alias"`
}

type OutKernel struct {
	Command    string               `json:"command"`
	Image      string               `json:"image"`
	Interfaces []OutKernelInterface `json:"interfaces"`
}

type OutKernelInterface struct {
	Name        string `json:"name"`
	Description string `json:"description"`

	Path      string `json:"path"`
	Type      string `json:"type"`
	Direction string `json:"direction"`
	Optional  bool   `json:"optional"`
}

func ProjectInterfacesToOutProjectInterfaces(interfaces []ProjectInterface) []OutProjectInterface {
	output := make([]OutProjectInterface, len(interfaces))
	for index, pinterface := range interfaces {
		output[index] = OutProjectInterface{
			Name:        pinterface.Name,
			Description: pinterface.Description,
			Alias:       pinterface.Alias,
			Type:        NodeTypeToStrings[pinterface.Type],
			Direction:   NodeDirectionToStrings[pinterface.Direction],
			Optional:    pinterface.Optional,
		}
	}

	return output
}

func ConfigurationEntitiesToOutConfigurationEntities(entities []ConfigurationEntity) []OutConfigurationEntity {
	output := make([]OutConfigurationEntity, len(entities))
	for eindex, centities := range entities {
		config := OutConfigurationEntity{
			Name:        centities.Name,
			Description: centities.Description,
			Kernel:      centities.Kernel,
			Interfaces:  make([]OutConfigurationEntityInterface, len(centities.Interfaces)),
		}

		for cindex, cinterface := range centities.Interfaces {
			config.Interfaces[cindex] = OutConfigurationEntityInterface{
				Name:  cinterface.Name,
				Alias: cinterface.Alias,
			}
		}

		output[eindex] = config
	}
	return output
}

func KernelToOutKernel(kernel Kernel) OutKernel {
	output := OutKernel{
		Command:    kernel.Command,
		Image:      kernel.Image,
		Interfaces: make([]OutKernelInterface, len(kernel.Interfaces)),
	}

	for index, kinterface := range kernel.Interfaces {
		output.Interfaces[index] = OutKernelInterface{
			Name:        kinterface.Name,
			Description: kinterface.Description,
			Path:        kinterface.Path,
			Type:        NodeTypeToStrings[kinterface.Type],
			Direction:   NodeDirectionToStrings[kinterface.Direction],
			Optional:    kinterface.Optional,
		}
	}

	return output
}
