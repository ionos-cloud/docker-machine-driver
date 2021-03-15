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

// KubernetesNodePoolForPut struct for KubernetesNodePoolForPut
type KubernetesNodePoolForPut struct {
	// The resource's unique identifier.
	Id *string `json:"id,omitempty"`
	// The type of object
	Type *string `json:"type,omitempty"`
	// URL to the object representation (absolute path)
	Href *string `json:"href,omitempty"`
	Metadata *DatacenterElementMetadata `json:"metadata,omitempty"`
	Properties *KubernetesNodePoolPropertiesForPut `json:"properties"`
}



// GetId returns the Id field value
// If the value is explicit nil, the zero value for string will be returned
func (o *KubernetesNodePoolForPut) GetId() *string {
	if o == nil {
		return nil
	}


	return o.Id

}

// GetIdOk returns a tuple with the Id field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *KubernetesNodePoolForPut) GetIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}


	return o.Id, true
}

// SetId sets field value
func (o *KubernetesNodePoolForPut) SetId(v string) {


	o.Id = &v

}

// HasId returns a boolean if a field has been set.
func (o *KubernetesNodePoolForPut) HasId() bool {
	if o != nil && o.Id != nil {
		return true
	}

	return false
}



// GetType returns the Type field value
// If the value is explicit nil, the zero value for string will be returned
func (o *KubernetesNodePoolForPut) GetType() *string {
	if o == nil {
		return nil
	}


	return o.Type

}

// GetTypeOk returns a tuple with the Type field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *KubernetesNodePoolForPut) GetTypeOk() (*string, bool) {
	if o == nil {
		return nil, false
	}


	return o.Type, true
}

// SetType sets field value
func (o *KubernetesNodePoolForPut) SetType(v string) {


	o.Type = &v

}

// HasType returns a boolean if a field has been set.
func (o *KubernetesNodePoolForPut) HasType() bool {
	if o != nil && o.Type != nil {
		return true
	}

	return false
}



// GetHref returns the Href field value
// If the value is explicit nil, the zero value for string will be returned
func (o *KubernetesNodePoolForPut) GetHref() *string {
	if o == nil {
		return nil
	}


	return o.Href

}

// GetHrefOk returns a tuple with the Href field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *KubernetesNodePoolForPut) GetHrefOk() (*string, bool) {
	if o == nil {
		return nil, false
	}


	return o.Href, true
}

// SetHref sets field value
func (o *KubernetesNodePoolForPut) SetHref(v string) {


	o.Href = &v

}

// HasHref returns a boolean if a field has been set.
func (o *KubernetesNodePoolForPut) HasHref() bool {
	if o != nil && o.Href != nil {
		return true
	}

	return false
}



// GetMetadata returns the Metadata field value
// If the value is explicit nil, the zero value for DatacenterElementMetadata will be returned
func (o *KubernetesNodePoolForPut) GetMetadata() *DatacenterElementMetadata {
	if o == nil {
		return nil
	}


	return o.Metadata

}

// GetMetadataOk returns a tuple with the Metadata field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *KubernetesNodePoolForPut) GetMetadataOk() (*DatacenterElementMetadata, bool) {
	if o == nil {
		return nil, false
	}


	return o.Metadata, true
}

// SetMetadata sets field value
func (o *KubernetesNodePoolForPut) SetMetadata(v DatacenterElementMetadata) {


	o.Metadata = &v

}

// HasMetadata returns a boolean if a field has been set.
func (o *KubernetesNodePoolForPut) HasMetadata() bool {
	if o != nil && o.Metadata != nil {
		return true
	}

	return false
}



// GetProperties returns the Properties field value
// If the value is explicit nil, the zero value for KubernetesNodePoolPropertiesForPut will be returned
func (o *KubernetesNodePoolForPut) GetProperties() *KubernetesNodePoolPropertiesForPut {
	if o == nil {
		return nil
	}


	return o.Properties

}

// GetPropertiesOk returns a tuple with the Properties field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *KubernetesNodePoolForPut) GetPropertiesOk() (*KubernetesNodePoolPropertiesForPut, bool) {
	if o == nil {
		return nil, false
	}


	return o.Properties, true
}

// SetProperties sets field value
func (o *KubernetesNodePoolForPut) SetProperties(v KubernetesNodePoolPropertiesForPut) {


	o.Properties = &v

}

// HasProperties returns a boolean if a field has been set.
func (o *KubernetesNodePoolForPut) HasProperties() bool {
	if o != nil && o.Properties != nil {
		return true
	}

	return false
}


func (o KubernetesNodePoolForPut) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}

	if o.Id != nil {
		toSerialize["id"] = o.Id
	}
	

	if o.Type != nil {
		toSerialize["type"] = o.Type
	}
	

	if o.Href != nil {
		toSerialize["href"] = o.Href
	}
	

	if o.Metadata != nil {
		toSerialize["metadata"] = o.Metadata
	}
	

	if o.Properties != nil {
		toSerialize["properties"] = o.Properties
	}
	
	return json.Marshal(toSerialize)
}

type NullableKubernetesNodePoolForPut struct {
	value *KubernetesNodePoolForPut
	isSet bool
}

func (v NullableKubernetesNodePoolForPut) Get() *KubernetesNodePoolForPut {
	return v.value
}

func (v *NullableKubernetesNodePoolForPut) Set(val *KubernetesNodePoolForPut) {
	v.value = val
	v.isSet = true
}

func (v NullableKubernetesNodePoolForPut) IsSet() bool {
	return v.isSet
}

func (v *NullableKubernetesNodePoolForPut) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableKubernetesNodePoolForPut(val *KubernetesNodePoolForPut) *NullableKubernetesNodePoolForPut {
	return &NullableKubernetesNodePoolForPut{value: val, isSet: true}
}

func (v NullableKubernetesNodePoolForPut) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableKubernetesNodePoolForPut) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


