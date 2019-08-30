/*
Copyright 2018 The Kubernetes Authors.

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

package converter

import (
	"fmt"
	"strings"

	"k8s.io/klog"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func convertExampleCRD(Object *unstructured.Unstructured, toVersion string) (*unstructured.Unstructured, metav1.Status) {
	klog.V(2).Info("converting crd")

	convertedObject := Object.DeepCopy()
	fromVersion := Object.GetAPIVersion()

	if toVersion == fromVersion {
		return nil, statusErrorWithMessage("conversion from a version to itself should not call the webhook: %s", toVersion)
	}

	switch Object.GetAPIVersion() {
	case "stable.example.com/v1beta1":
		switch toVersion {
		case "stable.example.com/v1":
			hostPort, ok, _ := unstructured.NestedString(convertedObject.Object, "spec", "hostPort")
			if ok {
				delete(convertedObject.Object, "spec")
				parts := strings.Split(hostPort, ":")
				if len(parts) != 2 {
					return nil, statusErrorWithMessage("invalid hostPort value `%v`", hostPort)
				}
				host := parts[0]
				port := parts[1]
				unstructured.SetNestedField(convertedObject.Object, host, "spec", "host")
				unstructured.SetNestedField(convertedObject.Object, port, "spec", "port")
			}
		default:
			return nil, statusErrorWithMessage("unexpected conversion version %q", toVersion)
		}
	case "stable.example.com/v1":
		switch toVersion {
		case "stable.example.com/v1beta1":
			host, hasHost, _ := unstructured.NestedString(convertedObject.Object, "spec", "host")
			port, hasPort, _ := unstructured.NestedString(convertedObject.Object, "spec", "port")
			if hasHost || hasPort {
				if !hasHost {
					host = ""
				}
				if !hasPort {
					port = ""
				}
				hostPort := fmt.Sprintf("%s:%s", host, port)
				unstructured.SetNestedField(convertedObject.Object, hostPort, "spec", "hostPort")
			}
		default:
			return nil, statusErrorWithMessage("unexpected conversion version %q", toVersion)
		}
	default:
		return nil, statusErrorWithMessage("unexpected conversion version %q", fromVersion)
	}
	return convertedObject, statusSucceed()
}
