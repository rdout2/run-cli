package scaling

// Scaling represents the scaling configuration of a Cloud Run service.
type Scaling struct {
	ScalingMode         string `json:"scalingMode,omitempty"`
	ManualInstanceCount int32  `json:"manualInstanceCount,omitempty"`
	MinInstances        int32  `json:"minInstances,omitempty"`
	MaxInstances        int32  `json:"maxInstances,omitempty"`
}
