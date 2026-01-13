package domainmapping

import (
	"time"

	"github.com/JulienBreux/run-cli/internal/run/model/common/condition"
)

// DomainMapping represents a Cloud Run domain mapping.
type DomainMapping struct {
	Name       string           `json:"name"`
	RouteName  string           `json:"routeName"`
	Region     string           `json:"region"`
	Project    string           `json:"project"`
	Creator    string           `json:"creator"`
	Records    []ResourceRecord `json:"records"`
	CreateTime time.Time        `json:"createTime"`
	UpdateTime time.Time        `json:"updateTime"`
	Conditions []*condition.Condition `json:"conditions,omitempty"`
}

// ResourceRecord represents a DNS resource record.
type ResourceRecord struct {
	Type   string `json:"type"`
	Name   string `json:"name"`
	RRData string `json:"rrData"`
}
