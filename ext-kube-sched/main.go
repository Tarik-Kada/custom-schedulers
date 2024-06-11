package main

import (
    "encoding/json"
    "net/http"
    "log"
    "time"
)

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
    NodeName          string            `json:"nodeName"`
    Status            string            `json:"nodeStatus"`
    CpuCapacity       int64             `json:"cpuCapacity"`
    MemoryCapacity    int64             `json:"memoryCapacity"`
    EphemeralCapacity int64             `json:"ephemeralCapacity"`
    CpuUsage          int64             `json:"cpuUsage"`
    MemoryUsage       int64             `json:"memoryUsage"`
    EphemeralUsage    int64             `json:"ephemeralUsage"`
    ScalarResources   map[string]int64  `json:"scalarResources"`
    RunningPods       []FilteredPod     `json:"runningPods"`
}

type FilteredPod struct {
    Name              string            `json:"name"`
    Namespace         string            `json:"namespace"`
    Labels            map[string]string `json:"labels"`
    CpuRequests       int64             `json:"cpuRequests"`
    MemoryRequests    int64             `json:"memoryRequests"`
    EphemeralRequests int64             `json:"ephemeralRequests"`
    ScalarRequests    map[string]int64  `json:"scalarRequests"`
    Containers        []ContainerInfo   `json:"containers"`
}

type ContainerInfo struct {
    Name  string `json:"name"`
    Image string `json:"image"`
}

type SchedulerResponse struct {
    Node string `json:"node"`
}

func schedulePod(w http.ResponseWriter, r *http.Request) {
    log.Println("Received scheduling request")
    log.Printf("Request: %v\n", r)
    var req SchedulerRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    start := time.Now()

    bestNode := ""
    bestScore := int64(-1)

    for _, node := range req.ClusterInfo.Nodes {
        if fitsRequest(&req.Pod, &node) {
            score := scoreNode(&req.Pod, &node)
            log.Printf("Node: %s, Score: %d\n", node.NodeName, score)
            if score > bestScore {
                bestNode = node.NodeName
                bestScore = score
            }
        }
    }

    duration := time.Since(start)
    log.Printf("Scheduling request processed in %v\n", duration)
    log.Printf("Selected Node: %s\n", bestNode)

    resp := SchedulerResponse{Node: bestNode}
    if err := json.NewEncoder(w).Encode(&resp); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

func fitsRequest(pod *FilteredPod, node *NodeInfo) bool {
    if node.CpuCapacity-node.CpuUsage < pod.CpuRequests {
        return false
    }
    if node.MemoryCapacity-node.MemoryUsage < pod.MemoryRequests {
        return false
    }
    if node.EphemeralCapacity-node.EphemeralUsage < pod.EphemeralRequests {
        return false
    }

    for resourceName, requested := range pod.ScalarRequests {
        if nodeCapacity, exists := node.ScalarResources[resourceName]; exists {
            if nodeUsage, exists := node.ScalarResources[resourceName]; exists {
                if nodeCapacity-nodeUsage < requested {
                    return false
                }
            } else {
                return false
            }
        } else {
            return false
        }
    }

    return true
}

func scoreNode(pod *FilteredPod, node *NodeInfo) int64 {
    cpuFree := node.CpuCapacity - node.CpuUsage - pod.CpuRequests
    memoryFree := node.MemoryCapacity - node.MemoryUsage - pod.MemoryRequests
    ephemeralFree := node.EphemeralCapacity - node.EphemeralUsage - pod.EphemeralRequests

    score := cpuFree + memoryFree + ephemeralFree

    for resourceName, requested := range pod.ScalarRequests {
        if nodeCapacity, exists := node.ScalarResources[resourceName]; exists {
            if nodeUsage, exists := node.ScalarResources[resourceName]; exists {
                score += nodeCapacity - nodeUsage - requested
            }
        }
    }

    return score
}

func main() {
    http.HandleFunc("/", schedulePod)
    log.Println("Serving scheduler on port 8080")
    err := http.ListenAndServe(":8080", nil)
    if err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}
