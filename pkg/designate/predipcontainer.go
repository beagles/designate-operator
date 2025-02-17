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

package designate

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/utils/ptr"
)

type PredIPContainerDetails struct {
	ContainerImage string
	VolumeMounts   []corev1.VolumeMount
	Command        string
	EnvVars        []corev1.EnvVar
}

func PredictableIPContainer(init PredIPContainerDetails) corev1.Container {

	args := []string{
		"-c",
		init.Command,
	}

	capabilities := []corev1.Capability{"NET_ADMIN", "SYS_ADMIN", "SYS_NICE"}
	return corev1.Container{
		Name:  "predictableips",
		Image: init.ContainerImage,
		SecurityContext: &corev1.SecurityContext{
			Capabilities: &corev1.Capabilities{
				Add:  capabilities,
				Drop: []corev1.Capability{},
			},
			RunAsUser: ptr.To(int64(0)),
		},
		Command: []string{
			"/bin/bash",
		},
		Args:         args,
		Env:          init.EnvVars,
		VolumeMounts: init.VolumeMounts,
	}
}
