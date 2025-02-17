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
package designatemdns

import (
	corev1 "k8s.io/api/core/v1"

	"github.com/openstack-k8s-operators/designate-operator/pkg/designate"
)

func GetVolumes(name string) []corev1.Volume {
	var config0640AccessMode int32 = 0640
	return append(
		designate.GetVolumes(name),
		corev1.Volume{
			Name: "bind-ips",
			VolumeSource: corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: designate.MdnsPredIPConfigMap,
					},
					DefaultMode: &config0640AccessMode,
				},
			},
		},
	)
}

func getPredIPVolumeMounts() []corev1.VolumeMount {
	return []corev1.VolumeMount{
		{
			Name:      "scripts",
			MountPath: "/usr/local/bin/container-scripts",
			ReadOnly:  true,
		},
		{
			Name:      "bind-ips",
			MountPath: "/var/lib/predictableips",
		},
	}
}
