package titanium

import (
	"errors"
	"fmt"
	"github.com/atomosio/common"
)

var _ = fmt.Printf

type Cluster struct {
	Response

	Status    string  `json:"status"`
	Clusters  []int64 `json:"clusters"`
	Instances []int64 `json:"instances"`
}

const (
	ClusterInvalidStatus = iota
	ClusterWaitingStatus // Waiting for some criteria to be met before starting
	ClusterActiveStatus  // Instances being scheduled and/or running
	ClusterStoppedStatus // Stopped

	//TransformingStatus // The structure of the cluster is changing, will go back to Running when all done.
)

var (
	ClusterStatusStrings = []string{
		"Invalid",
		"Waiting",
		"Active",
		"Stopped",
	}
)

// Retreives information related to the cluster
func (client *HttpClient) GetCluster(id int64) (Cluster, error) {
	var output Cluster

	// Get and unmarshal
	addr := fmt.Sprintf("%s%d", ClustersEndpoint, id)
	//fmt.Println(addr)
	err := client.getAndUnmarshal(addr, &output)
	if err != nil {
		return output, err
	}

	if output.Response.Code != common.Success {
		return output, errors.New("Failed to get cluster information: " + output.Response.Description)
	}

	return output, nil
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
