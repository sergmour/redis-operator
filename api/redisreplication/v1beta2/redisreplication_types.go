package v1beta2

import (
	common "github.com/OT-CONTAINER-KIT/redis-operator/api/common/v1beta2"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type RedisReplicationSpec struct {
	Size                          *int32                            `json:"clusterSize"`
	KubernetesConfig              common.KubernetesConfig           `json:"kubernetesConfig"`
	RedisExporter                 *common.RedisExporter             `json:"redisExporter,omitempty"`
	RedisConfig                   *common.RedisConfig               `json:"redisConfig,omitempty"`
	Storage                       *common.Storage                   `json:"storage,omitempty"`
	NodeSelector                  map[string]string                 `json:"nodeSelector,omitempty"`
	PodSecurityContext            *corev1.PodSecurityContext        `json:"podSecurityContext,omitempty"`
	SecurityContext               *corev1.SecurityContext           `json:"securityContext,omitempty"`
	PriorityClassName             string                            `json:"priorityClassName,omitempty"`
	Affinity                      *corev1.Affinity                  `json:"affinity,omitempty"`
	Tolerations                   *[]corev1.Toleration              `json:"tolerations,omitempty"`
	TLS                           *common.TLSConfig                 `json:"TLS,omitempty"`
	PodDisruptionBudget           *common.RedisPodDisruptionBudget  `json:"pdb,omitempty"`
	ACL                           *common.ACLConfig                 `json:"acl,omitempty"`
	ReadinessProbe                *corev1.Probe                     `json:"readinessProbe,omitempty" protobuf:"bytes,11,opt,name=readinessProbe"`
	LivenessProbe                 *corev1.Probe                     `json:"livenessProbe,omitempty" protobuf:"bytes,12,opt,name=livenessProbe"`
	InitContainer                 *common.InitContainer             `json:"initContainer,omitempty"`
	Sidecars                      *[]common.Sidecar                 `json:"sidecars,omitempty"`
	ServiceAccountName            *string                           `json:"serviceAccountName,omitempty"`
	TerminationGracePeriodSeconds *int64                            `json:"terminationGracePeriodSeconds,omitempty" protobuf:"varint,4,opt,name=terminationGracePeriodSeconds"`
	EnvVars                       *[]corev1.EnvVar                  `json:"env,omitempty"`
	TopologySpreadConstrains      []corev1.TopologySpreadConstraint `json:"topologySpreadConstraints,omitempty"`
	HostPort                      *int                              `json:"hostPort,omitempty"`
	Sentinel                      *Sentinel                         `json:"sentinel,omitempty"`
	// ExternalMaster, when set, configures all pods of this RedisReplication
	// to act as read-replicas connecting to a Redis master that lives outside
	// of this Kubernetes cluster (typically in a primary K8s cluster).
	// When this field is set:
	//   - No local master is elected; the controller skips leader-election
	//     and failover logic.
	//   - All pods are configured with `replicaof <host> <port>` and
	//     `replica-read-only yes`.
	//   - Authentication to the external master uses the password from
	//     KubernetesConfig.ExistingPasswordSecret (masterauth).
	//   - When Spec.TLS is also set, the slave-to-master replication link
	//     uses TLS (`tls-replication yes`).
	//   - This field is mutually exclusive with Spec.Sentinel.
	// +optional
	ExternalMaster *ExternalMasterConfig `json:"externalMaster,omitempty"`
}

// ExternalMasterConfig defines the static endpoint of a Redis master that
// lives outside of this Kubernetes cluster.
type ExternalMasterConfig struct {
	// Host is the DNS name or IP address of the external Redis master.
	// +kubebuilder:validation:MinLength=1
	Host string `json:"host"`
	// Port is the TCP port of the external Redis master.
	// Defaults to 6379 when not specified.
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=65535
	// +optional
	Port *int32 `json:"port,omitempty"`
}

type Sentinel struct {
	common.KubernetesConfig `json:",inline"`
	common.SentinelConfig   `json:",inline"`
	Size                    int32 `json:"size"`
}

func (cr *RedisReplicationSpec) GetReplicationCounts(t string) int32 {
	replica := cr.Size
	return *replica
}

// ConnectionInfo provides connection details for clients to connect to Redis
type ConnectionInfo struct {
	// Host is the service FQDN
	Host string `json:"host,omitempty"`
	// Port is the service port
	Port int `json:"port,omitempty"`
	// MasterName is the Sentinel master group name, only set when Sentinel mode is enabled
	// +optional
	MasterName string `json:"masterName,omitempty"`
}

// RedisStatus defines the observed state of Redis
type RedisReplicationStatus struct {
	MasterNode string `json:"masterNode,omitempty"`
	// ConnectionInfo provides connection details for clients to connect to Redis
	// +optional
	ConnectionInfo *ConnectionInfo `json:"connectionInfo,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:storageversion
// +kubebuilder:printcolumn:name="Master",type="string",JSONPath=".status.masterNode"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// Redis is the Schema for the redis API
type RedisReplication struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RedisReplicationSpec   `json:"spec"`
	Status RedisReplicationStatus `json:"status,omitempty"`
}

func (rr *RedisReplication) GetStatefulSetName() string {
	return rr.Name
}

// +kubebuilder:object:root=true

// RedisList contains a list of Redis
type RedisReplicationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []RedisReplication `json:"items"`
}

//nolint:gochecknoinits
func init() {
	SchemeBuilder.Register(&RedisReplication{}, &RedisReplicationList{})
}
