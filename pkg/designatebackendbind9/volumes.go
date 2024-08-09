/*
Copyright 2024.

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

package designatebackendbind9

import (
	corev1 "k8s.io/api/core/v1"
)

const (
	scriptVolume       = "designatebackendbind9-scripts"
	dataVolume         = "designatebackendbind9-bind-volume"
	configVolume       = "designatebackendbind9-config-data"
	namedConfigVolume  = "designatebackendbind9-config-named"
	mergedConfigVolume = "designatebackendbind9-config-data-merged"
	logVolume          = "designatebackendbind9-log-volume"
)

func getServicePodVolumes(baseConfigMapName string) []corev1.Volume {
	var scriptMode int32 = 0755
	var configMode int32 = 0640
	return []corev1.Volume{
		{
			Name: scriptVolume,
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					DefaultMode: &scriptMode,
					SecretName:  baseConfigMapName + "-scripts",
				},
			},
		},
		{
			Name: configVolume,
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					DefaultMode: &configMode,
					SecretName:  baseConfigMapName + "-config-data",
				},
			},
		},
		{
			Name: namedConfigVolume,
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					DefaultMode: &configMode,
					SecretName:  baseConfigMapName + "-config-named",
				},
			},
		},
		{
			Name: mergedConfigVolume,
			VolumeSource: corev1.VolumeSource{
				EmptyDir: &corev1.EmptyDirVolumeSource{Medium: ""},
			},
		},
		{
			Name: logVolume,
			VolumeSource: corev1.VolumeSource{
				EmptyDir: &corev1.EmptyDirVolumeSource{Medium: ""},
			},
		},
	}
}

// TODO(beagles): we follow the old TripleO/kolla naming of these mounts, but do they really make sense here?

// getInitVolumesMounts - the init container will use the scripts mounted in the scriptVolume and create completed named
// configuration from the files in configVolume. The modified files will be stored in the mergedConfigVolume
func getInitVolumeMounts() []corev1.VolumeMount {
	return []corev1.VolumeMount{
		{
			Name:      configVolume,
			MountPath: "/var/lib/config-data/default",
			ReadOnly:  true,
		},
		{
			Name:      namedConfigVolume,
			MountPath: "/var/lib/config-data/default/named",
			ReadOnly:  false,
		},
		{
			Name:      mergedConfigVolume,
			MountPath: "/var/lib/config-data/merged",
			ReadOnly:  false,
		},
		{
			Name:      scriptVolume,
			MountPath: "/usr/local/bin/container-scripts",
			ReadOnly:  true,
		},
	}
}

func getServicePodVolumeMounts(persistentData string) []corev1.VolumeMount {
	return []corev1.VolumeMount{
		{
			Name:      mergedConfigVolume,
			MountPath: "/var/lib/config-data/merged",
			ReadOnly:  true,
		},
		{
			Name:      mergedConfigVolume,
			MountPath: "/var/lib/kolla/config_files/config.json",
			SubPath:   "designate-bind9-config.json",
			ReadOnly:  true,
		},
		{
			Name:      scriptVolume,
			MountPath: "/usr/local/bin/container-scripts",
			ReadOnly:  true,
		},
		{
			Name:      persistentData,
			MountPath: "/var/named-persistent",
			ReadOnly:  false,
		},
		{
			Name:      logVolume,
			MountPath: "/var/log/bind",
			ReadOnly:  false,
		},
	}
}
