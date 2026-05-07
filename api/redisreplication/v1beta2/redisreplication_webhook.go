/*
Copyright 2020 Opstree Solutions.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1beta2

import (
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

const (
	webhookPath = "/validate-redis-redis-opstreelabs-in-v1beta2-redisreplication"
)

// log is for logging in this package.
var redisreplicationlog = logf.Log.WithName("redisreplication-v1beta2-validation")

// +kubebuilder:webhook:path=/validate-redis-redis-opstreelabs-in-v1beta2-redisreplication,mutating=false,failurePolicy=fail,sideEffects=None,groups=redis.redis.opstreelabs.in,resources=redisreplications,verbs=create;update,versions=v1beta2,name=validate-redisreplication.redis.opstreelabs.in,admissionReviewVersions=v1

func (r *RedisReplication) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

var _ webhook.Validator = &RedisReplication{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *RedisReplication) ValidateCreate() (admission.Warnings, error) {
	redisreplicationlog.Info("validate create", "name", r.Name)

	return r.validate(nil)
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *RedisReplication) ValidateUpdate(old runtime.Object) (admission.Warnings, error) {
	redisreplicationlog.Info("validate update", "name", r.Name)

	return r.validate(old.(*RedisReplication))
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *RedisReplication) ValidateDelete() (admission.Warnings, error) {
	redisreplicationlog.Info("validate delete", "name", r.Name)

	return nil, nil
}

// validate validates the RedisReplication CR
func (r *RedisReplication) validate(_ *RedisReplication) (admission.Warnings, error) {
	var errors field.ErrorList

	// Validate ACL configuration
	if r.Spec.ACL != nil {
		if err := r.Spec.ACL.Validate(); err != nil {
			errors = append(errors, field.Invalid(
				field.NewPath("spec").Child("acl"),
				r.Spec.ACL,
				err.Error(),
			))
		}
	}

	// Validate ExternalMaster configuration.
	if r.Spec.ExternalMaster != nil {
		extPath := field.NewPath("spec").Child("externalMaster")

		if r.Spec.ExternalMaster.Host == "" {
			errors = append(errors, field.Required(
				extPath.Child("host"),
				"externalMaster.host must be set when externalMaster is specified",
			))
		}
		if r.Spec.ExternalMaster.Port != nil {
			p := *r.Spec.ExternalMaster.Port
			if p < 1 || p > 65535 {
				errors = append(errors, field.Invalid(
					extPath.Child("port"),
					p,
					"externalMaster.port must be between 1 and 65535",
				))
			}
		}

		// ExternalMaster is mutually exclusive with Sentinel: in slave-only
		// mode the controller skips all leader-election/failover logic, and
		// Sentinel-driven failover would conflict with that contract.
		if r.Spec.Sentinel != nil {
			errors = append(errors, field.Forbidden(
				extPath,
				"spec.externalMaster and spec.sentinel are mutually exclusive",
			))
		}
	}

	if len(errors) == 0 {
		return nil, nil
	}

	return nil, apierrors.NewInvalid(
		schema.GroupKind{Group: "redis.redis.opstreelabs.in", Kind: "RedisReplication"},
		r.Name,
		errors,
	)
}

func (r *RedisReplication) WebhookPath() string {
	return webhookPath
}
