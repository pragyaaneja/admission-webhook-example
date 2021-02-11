// Package mutate deals with AdmissionReview requests and responses, it takes in the request body and returns a readily converted JSON []byte that can be
// returned from a http Handler w/o needing to further convert or modify it, it also makes testing Mutate() kind of easy w/o need for a fake http server, etc.
package main

import (
	"encoding/json"
	"fmt"

	v1beta1 "k8s.io/api/admission/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Mutate mutates
func Mutate(body []byte) ([]byte, error) {
	// if verbose {
	// 	log.Printf("recieved object: %s\n", string(body)) // untested section
	// }

	// unmarshal request into AdmissionReview struct
	admReview := v1beta1.AdmissionReview{}

	if err := json.Unmarshal(body, &admReview); err != nil {
		return nil, fmt.Errorf("unmarshaling request failed with %s", err)
	}

	admRequest := admReview.Request
	admResponse := v1beta1.AdmissionResponse{}

	var err error
	var pod *corev1.Pod

	responseBody := []byte{}

	// get the Pod object and unmarshal it into its struct, if we cannot, we might as well stop here
	if err := json.Unmarshal(admRequest.Object.Raw, &pod); err != nil {
		return nil, fmt.Errorf("unable unmarshal pod json object %v", err)
	}

	// the actual mutation is done by a string in JSONPatch style, i.e. we don't _actually_ modify the object, but
	// tell K8S how it should modifiy it
	p := make([]map[string]interface{}, 0)
	for i := range pod.Spec.Containers {
		limitValOld := pod.Spec.Containers[i].Resources.Limits["kubernetes.azure.com/sgx_epc_mem_in_MiB"]
		requestValOld := pod.Spec.Containers[i].Resources.Requests["kubernetes.azure.com/sgx_epc_mem_in_MiB"]
		requestValNew := pod.Spec.Containers[i].Resources.Requests["sgx.intel.com/epc"]

		if !requestValNew.IsZero() {
			// the deployed pod configuration has the new EPC resource name
			pod.Spec.Containers[i].Resources.Requests["sgx.intel.com/enclave"] = *resource.NewQuantity(1, resource.DecimalSI)
			pod.Spec.Containers[i].Resources.Requests["sgx.intel.com/provision"] = *resource.NewQuantity(1, resource.DecimalSI)
			pod.Spec.Containers[i].Resources.Limits["sgx.intel.com/enclave"] = *resource.NewQuantity(1, resource.DecimalSI)
			pod.Spec.Containers[i].Resources.Limits["sgx.intel.com/provision"] = *resource.NewQuantity(1, resource.DecimalSI)

		} else if !requestValOld.IsZero() {
			// the deployed pod configuration has the old EPC resource name
			requestValInt, _ := requestValOld.AsInt64()
			newRequestVal := resource.NewQuantity(requestValInt*1024*1024, resource.BinarySI)
			pod.Spec.Containers[i].Resources.Requests = make(corev1.ResourceList)
			pod.Spec.Containers[i].Resources.Requests["sgx.intel.com/enclave"] = *resource.NewQuantity(1, resource.DecimalSI)
			pod.Spec.Containers[i].Resources.Requests["sgx.intel.com/provision"] = *resource.NewQuantity(1, resource.DecimalSI)
			pod.Spec.Containers[i].Resources.Requests["sgx.intel.com/epc"] = *newRequestVal

			limitValInt, _ := limitValOld.AsInt64()
			newLimitVal := resource.NewQuantity(limitValInt*1024*1024, resource.BinarySI)
			pod.Spec.Containers[i].Resources.Limits = make(corev1.ResourceList)
			pod.Spec.Containers[i].Resources.Limits["sgx.intel.com/enclave"] = *resource.NewQuantity(1, resource.DecimalSI)
			pod.Spec.Containers[i].Resources.Limits["sgx.intel.com/provision"] = *resource.NewQuantity(1, resource.DecimalSI)
			pod.Spec.Containers[i].Resources.Limits["sgx.intel.com/epc"] = *newLimitVal

		} else {
			continue
		}

	}

	patch := make(map[string]interface{})
	patch["op"] = "replace"
	patch["path"] = "/spec/containers"
	patch["value"] = pod.Spec.Containers

	p = append(p, patch)
	// log.Printf("patch resp: %v\n", p)

	// set response options
	admResponse.Allowed = true
	admResponse.UID = admRequest.UID
	pT := v1beta1.PatchTypeJSONPatch
	admResponse.PatchType = &pT // it's annoying that this needs to be a pointer as you cannot give a pointer to a constant?

	// parse the []map into JSON
	marshalPatch, err := json.Marshal(p)
	admResponse.Patch = []byte(marshalPatch)
	// log.Printf("admission resp patch: %v\n", admResponse.Patch)

	// Success
	admResponse.Result = &metav1.Status{
		Status: "Success",
	}

	admReview.Response = &admResponse
	// back into JSON so we can return the finished AdmissionReview w/ Response directly
	// w/o needing to convert things in the http handler
	responseBody, err = json.Marshal(admReview)
	if err != nil {
		return nil, err // untested section
	}

	// if verbose {
	// 	log.Printf("resp: %s\n", string(responseBody)) // untested section
	// }

	return responseBody, nil
}

// func createResourceList() corev1.ResourceList {

// }
