package v1beta2

import "fmt"

const defaultExternalMasterPort = int32(6379)

func (cr *RedisReplication) EnableSentinel() bool {
	return cr != nil && cr.Spec.Sentinel != nil && cr.Spec.Sentinel.Size > 0
}

// HasExternalMaster reports whether this RedisReplication is configured to
// replicate from an external Redis master (slave-only mode).
func (cr *RedisReplication) HasExternalMaster() bool {
	return cr != nil && cr.Spec.ExternalMaster != nil && cr.Spec.ExternalMaster.Host != ""
}

// ExternalMasterPort returns the configured external master port or the
// default Redis port (6379) if unset.
func (cr *RedisReplication) ExternalMasterPort() int32 {
	if cr == nil || cr.Spec.ExternalMaster == nil || cr.Spec.ExternalMaster.Port == nil {
		return defaultExternalMasterPort
	}
	return *cr.Spec.ExternalMaster.Port
}

func (cr *RedisReplication) SentinelStatefulSet() string {
	return cr.Name + "-s"
}

func (cr *RedisReplication) RedisStatefulSet() string {
	return cr.Name
}

func (cr *RedisReplication) SentinelHLService() string {
	return cr.Name + "-s-hl"
}

func (cr *RedisReplication) MasterService() string {
	return cr.Name + "-master"
}

// GetConnectionInfo returns connection info for clients based on the mode.
// The dnsDomain parameter should be the cluster DNS domain (e.g., "cluster.local").
func (cr *RedisReplication) GetConnectionInfo(dnsDomain string) *ConnectionInfo {
	if cr.HasExternalMaster() {
		return &ConnectionInfo{
			Host: cr.Spec.ExternalMaster.Host,
			Port: int(cr.ExternalMasterPort()),
		}
	}
	if cr.EnableSentinel() {
		return &ConnectionInfo{
			Host:       fmt.Sprintf("%s.%s.svc.%s", cr.SentinelHLService(), cr.Namespace, dnsDomain),
			Port:       26379,
			MasterName: "mymaster",
		}
	}
	return &ConnectionInfo{
		Host: fmt.Sprintf("%s.%s.svc.%s", cr.MasterService(), cr.Namespace, dnsDomain),
		Port: 6379,
	}
}

// ExternalMasterEndpoint returns a "host:port" string suitable for use as
// a synthetic master identifier in Status.MasterNode.
func (cr *RedisReplication) ExternalMasterEndpoint() string {
	if !cr.HasExternalMaster() {
		return ""
	}
	return fmt.Sprintf("%s:%d", cr.Spec.ExternalMaster.Host, cr.ExternalMasterPort())
}
