/*
 * CLOUD API
 *
 * An enterprise-grade Infrastructure is provided as a Service (IaaS) solution that can be managed through a browser-based \"Data Center Designer\" (DCD) tool or via an easy to use API.   The API allows you to perform a variety of management tasks such as spinning up additional servers, adding volumes, adjusting networking, and so forth. It is designed to allow users to leverage the same power and flexibility found within the DCD visual tool. Both tools are consistent with their concepts and lend well to making the experience smooth and intuitive.
 *
 * API version: 5.0
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package ionoscloud

import (
	"encoding/json"
)

// KubernetesNodePoolLan struct for KubernetesNodePoolLan
type KubernetesNodePoolLan struct {
	// The LAN ID of an existing LAN at the related datacenter
	Id *int32 `json:"id"`
}



// GetId returns the Id field value
// If the value is explicit nil, the zero value for int32 will be returned
func (o *KubernetesNodePoolLan) GetId() *int32 {
	if o == nil {
		return nil
	}


	return o.Id

}

// GetIdOk returns a tuple with the Id field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *KubernetesNodePoolLan) GetIdOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}


	return o.Id, true
}

// SetId sets field value
func (o *KubernetesNodePoolLan) SetId(v int32) {


	o.Id = &v

}

// HasId returns a boolean if a field has been set.
func (o *KubernetesNodePoolLan) HasId() bool {
	if o != nil && o.Id != nil {
		return true
	}

	return false
}


func (o KubernetesNodePoolLan) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}

	if o.Id != nil {
		toSerialize["id"] = o.Id
	}
	
	return json.Marshal(toSerialize)
}

type NullableKubernetesNodePoolLan struct {
	value *KubernetesNodePoolLan
	isSet bool
}

func (v NullableKubernetesNodePoolLan) Get() *KubernetesNodePoolLan {
	return v.value
}

func (v *NullableKubernetesNodePoolLan) Set(val *KubernetesNodePoolLan) {
	v.value = val
	v.isSet = true
}

func (v NullableKubernetesNodePoolLan) IsSet() bool {
	return v.isSet
}

func (v *NullableKubernetesNodePoolLan) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableKubernetesNodePoolLan(val *KubernetesNodePoolLan) *NullableKubernetesNodePoolLan {
	return &NullableKubernetesNodePoolLan{value: val, isSet: true}
}

func (v NullableKubernetesNodePoolLan) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableKubernetesNodePoolLan) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


