/*

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

package designateunbound

import (
	designatev1beta1 "github.com/openstack-k8s-operators/designate-operator/api/v1beta1"
	designate "github.com/openstack-k8s-operators/designate-operator/pkg/designate"
	common "github.com/openstack-k8s-operators/lib-common/modules/common"
	"github.com/openstack-k8s-operators/lib-common/modules/common/affinity"
	"github.com/openstack-k8s-operators/lib-common/modules/common/env"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	configVolume = "designateunbound-config"
)

// Deployment func
func Deployment(instance *designatev1beta1.DesignateUnbound,
	configHash string,
	labels map[string]string,
	annotations map[string]string,
) *appsv1.Deployment {
	var rootUser int64 = 0
	var configMode int32 = 0640

	volumes := []corev1.Volume{
		{
			Name: configVolume,
			VolumeSource: corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{
					DefaultMode: &configMode,
					LocalObjectReference: corev1.LocalObjectReference{
						Name: "designate-unbound-config-data",
					},
				},
			},
		},
	}
	mounts := []corev1.VolumeMount{
		{
			Name:      configVolume,
			MountPath: "/etc/unbound/conf.d",
			ReadOnly:  true,
		},
	}

	livenessProbe := &corev1.Probe{
		// TODO might need tuning
		TimeoutSeconds:      15,
		PeriodSeconds:       13,
		InitialDelaySeconds: 10,
	}
	readinessProbe := &corev1.Probe{
		// TODO might need tuning
		TimeoutSeconds:      15,
		PeriodSeconds:       15,
		InitialDelaySeconds: 10,
	}

	// TODO(beagles): use equivalent's of healthcheck's in tripleo which
	// seem to largely based on connections to database. The pgrep's
	// could be tightened up too but they seem to be a bit tricky.

	livenessProbe.Exec = &corev1.ExecAction{
		Command: []string{
			"/usr/sbin/unbound-streamtcp", "-u", ".", "SOA", "IN",
		},
	}

	readinessProbe.Exec = &corev1.ExecAction{
		Command: []string{
			"/usr/sbin/unbound-streamtcp", "-u", ".", "SOA", "IN",
		},
	}

	envVars := map[string]env.Setter{}
	envVars["KOLLA_CONFIG_STRATEGY"] = env.SetValue("COPY_ALWAYS")
	envVars["CONFIG_HASH"] = env.SetValue(configHash)

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.Name,
			Namespace: instance.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: instance.Spec.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: annotations,
					Labels:      labels,
				},
				Spec: corev1.PodSpec{
					ServiceAccountName: instance.Spec.ServiceAccount,
					// Unbound doesn't use any config in common with the other
					// designate services so just give it it's own config
					// volume.
					Volumes: volumes,
					Containers: []corev1.Container{{
						Name:    ServiceName,
						Image:   instance.Spec.ContainerImage,
						Command: []string{"/usr/sbin/unbound"},
						Args: []string{
							"-d",
							"-d",
							"-p",
						},
						SecurityContext: &corev1.SecurityContext{
							RunAsUser: &rootUser,
						},
						Env:            env.MergeEnvs([]corev1.EnvVar{}, envVars),
						VolumeMounts:   mounts,
						Resources:      instance.Spec.Resources,
						ReadinessProbe: readinessProbe,
						LivenessProbe:  livenessProbe,
					}},
				},
			},
		},
	}

	// If possible two pods of the same service should not
	// run on the same worker node. If this is not possible
	// the get still created on the same worker node.
	deployment.Spec.Template.Spec.Affinity = affinity.DistributePods(
		common.AppSelector,
		[]string{
			designate.ServiceName,
		},
		corev1.LabelHostname,
	)
	if instance.Spec.NodeSelector != nil && len(instance.Spec.NodeSelector) > 0 {
		deployment.Spec.Template.Spec.NodeSelector = instance.Spec.NodeSelector
	}

	return deployment
}
