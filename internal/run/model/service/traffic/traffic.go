package traffic

const (
	TrafficTargetAllocationTypeLatest = "TRAFFIC_TARGET_ALLOCATION_TYPE_LATEST"
)

// TrafficTargetStatus represents the actual traffic allocated to a revision.
type TrafficTargetStatus struct {
	Type     string `json:"type,omitempty"`
	Revision string `json:"revision,omitempty"`
	Percent  int32  `json:"percent,omitempty"`
	Tag      string `json:"tag,omitempty"`
	URI      string `json:"uri,omitempty"`
}
