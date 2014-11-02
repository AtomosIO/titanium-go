package titanium

import (
	"errors"
	"fmt"
	"time"

	"github.com/atomosio/common"
)

var _ = fmt.Printf

const (
	ClusterInvalidStatus = iota
	ClusterWaitingStatus // Waiting for some criteria to be met before starting
	ClusterActiveStatus  // Instances being scheduled and/or running
	ClusterStoppedStatus // Stopped
)

const (
	InvalidClusterType = iota
	BatchClusterType
)

var (
	ClusterStatusStrings = []string{
		"Invalid",
		"Waiting",
		"Active",
		"Stopped",
	}
	TypeStrings = []string{
		"Invalid",
		"Batch",
	}
	TypeStringsTo = map[string]int{
		TypeStrings[InvalidClusterType]: InvalidClusterType,
		TypeStrings[BatchClusterType]:   BatchClusterType,
	}

	ErrClusterWaitForFinishTimeout = errors.New("Cluster did not finish within timeout period")
)

type Cluster struct {
	Response

	Id        int64   `json:"cluster_id"`
	Status    string  `json:"status"`
	Clusters  []int64 `json:"clusters"`
	Instances []int64 `json:"instances"`
}

type CreateClusterRequest struct {
	Type       string            `json:"type"`
	Name       string            `json:"name"`
	Project    string            `json:"project"`
	Interfaces map[string]string `json:"interfaces"`
}

type CreateClusterResponse struct {
	Response
	ClusterId int64 `json:"cluster_id,omitempty"`
}

// Retreives information related to the cluster
func (client *HttpClient) GetCluster(id int64) (Cluster, error) {
	var output Cluster

	// Get and unmarshal
	addr := fmt.Sprintf("%s%d", ClustersEndpoint, id)
	err := client.getAndUnmarshal(addr, &output)
	if err != nil {
		return output, err
	}

	if output.Response.Code != common.Success {
		return output, errors.New("Failed to get cluster information: " + output.Response.Description)
	}

	output.Id = id
	return output, nil
}

func (client *HttpClient) CreateBatchCluster(name, project string, interfaces map[string]string) (Cluster, error) {
	request := CreateClusterRequest{
		Type:       TypeStrings[BatchClusterType],
		Name:       name,
		Project:    project,
		Interfaces: interfaces,
	}

	// Post and unmarshal response
	var response CreateClusterResponse
	err := client.postAndUnmarshal(ClustersEndpoint, &request, &response)
	if err != nil {
		return Cluster{}, err
	}

	if response.Code != common.Success {
		return Cluster{}, errors.New(response.Description)
	}

	// Get information on the newly created cluster
	return client.GetCluster(response.ClusterId)
}

func (client *HttpClient) WaitForClusterToFinish(id int64, seconds time.Duration) error {
	waitTill := time.Now().Add(seconds)

	for {
		// Get cluster information
		cluster, err := client.GetCluster(id)
		if err != nil {
			return err
		}

		//fmt.Printf("Waiting for cluster -> %+v\n", cluster)
		// If we're done, exit function
		if cluster.IsStopped() {
			return nil
		}

		if time.Now().After(waitTill) {
			return ErrClusterWaitForFinishTimeout
		}

		// Go to sleep for a bit
		time.Sleep(SpinSleepDuration)
	}
}

func (cluster Cluster) IsWaiting() bool {
	return cluster.Status == ClusterStatusStrings[ClusterWaitingStatus]
}

func (cluster Cluster) IsActive() bool {
	return cluster.Status == ClusterStatusStrings[ClusterActiveStatus]
}

func (cluster Cluster) IsStopped() bool {
	return cluster.Status == ClusterStatusStrings[ClusterStoppedStatus]
}
