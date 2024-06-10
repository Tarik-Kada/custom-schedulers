package main

import (
    "encoding/json"
    "net/http"
    "log"
    "fmt"
)

package main

type SchedulerRequest struct {
    Parameters      map[string]interface{} `json:"parameters"`
    Pod             FilteredPod            `json:"pod"`
    ClusterInfo     ClusterInfo            `json:"clusterInfo"`
    Metrics         map[string]interface{} `json:"metrics"`
    PrometheusError string                 `json:"prometheusError,omitempty"`
}

type ClusterInfo struct {
    Nodes []NodeInfo `json:"nodes"`
}

type NodeInfo struct {
    NodeName          string              `json:"nodeName"`
    Status            string              `json:"nodeStatus"`
    CpuCapacity       int64               `json:"cpuCapacity"`
    MemoryCapacity    int64               `json":"memoryCapacity"`
    EphemeralCapacity int64               `json":"ephemeralCapacity"`
    CpuUsage          int64               `json:"cpuUsage"`
    MemoryUsage       int64               `json":"memoryUsage"`
    EphemeralUsage    int64               `json":"ephemeralUsage"`
    ScalarResources   map[string]int64    `json:"scalarResources"`
    RunningPods       []FilteredPod       `json:"runningPods"`
}

type FilteredPod struct {
    Name             string            `json:"name"`
    Namespace        string            `json:"namespace"`
    Labels           map[string]string `json:"labels"`
    ServingService   string            `json:"servingService"`
    ServingRevision  string            `json:"servingRevision"`
    CpuRequests      int64             `json:"cpuRequests"`
    MemoryRequests   int64             `json:"memoryRequests"`
    EphemeralRequests int64            `json:"ephemeralRequests"`
    ScalarRequests   map[string]int64  `json:"scalarRequests"`
    Containers       []ContainerInfo   `json:"containers"`
}

type ContainerInfo struct {
    Name  string `json:"name"`
    Image string `json:"image"`
}

type SchedulerResponse struct {
    Node string `json:"node"`
}

func schedulePod(w http.ResponseWriter, r *http.Request) {
    var req SchedulerRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    bestNode := ""
    bestScore := int64(-1)

    for _, node := range req.ClusterInfo.Nodes {
        if fitsRequest(&req.Pod, &node) {
            score := scoreNode(&req.Pod, &node)
            if score > bestScore {
                bestNode = node.NodeName
                bestScore = score
            }
        }
    }

    resp := SchedulerResponse{Node: bestNode}
    if err := json.NewEncoder(w).Encode(&resp); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

func fitsRequest(pod *FilteredPod, node *NodeInfo) bool {
    // Implement the logic to check if the node can fit the pod
    // This should include checking CPU, memory, ephemeral storage, and scalar resources
    return true // Placeholder, implement actual logic
}

func scoreNode(pod *FilteredPod, node *NodeInfo) int64 {
    // Implement the scoring logic
    return 0 // Placeholder, implement actual logic
}

func main() {
    http.HandleFunc("/schedule", schedulePod)
    log.Fatal(http.ListenAndServe(":8080", nil))
}