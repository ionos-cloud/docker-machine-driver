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

// Snapshots struct for Snapshots
type Snapshots struct {
	// The resource's unique identifier
	Id *string `json:"id,omitempty"`
	// The type of object that has been created
	Type *Type `json:"type,omitempty"`
	// URL to the object representation (absolute path)
	Href *string `json:"href,omitempty"`
	// Array of items in that collection
	Items *[]Snapshot `json:"items,omitempty"`
}

// GetId returns the Id field value
// If the value is explicit nil, the zero value for string will be returned
func (o *Snapshots) GetId() *string {
	if o == nil {
		return nil
	}

	return o.Id

}

// GetIdOk returns a tuple with the Id field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *Snapshots) GetIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.Id, true
}

// SetId sets field value
func (o *Snapshots) SetId(v string) {

	o.Id = &v

}

// HasId returns a boolean if a field has been set.
func (o *Snapshots) HasId() bool {
	if o != nil && o.Id != nil {
		return true
	}

	return false
}

// GetType returns the Type field value
// If the value is explicit nil, the zero value for Type will be returned
func (o *Snapshots) GetType() *Type {
	if o == nil {
		return nil
	}

	return o.Type

}

// GetTypeOk returns a tuple with the Type field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *Snapshots) GetTypeOk() (*Type, bool) {
	if o == nil {
		return nil, false
	}

	return o.Type, true
}

// SetType sets field value
func (o *Snapshots) SetType(v Type) {

	o.Type = &v

}

// HasType returns a boolean if a field has been set.
func (o *Snapshots) HasType() bool {
	if o != nil && o.Type != nil {
		return true
	}

	return false
}

// GetHref returns the Href field value
// If the value is explicit nil, the zero value for string will be returned
func (o *Snapshots) GetHref() *string {
	if o == nil {
		return nil
	}

	return o.Href

}

// GetHrefOk returns a tuple with the Href field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *Snapshots) GetHrefOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.Href, true
}

// SetHref sets field value
func (o *Snapshots) SetHref(v string) {

	o.Href = &v

}

// HasHref returns a boolean if a field has been set.
func (o *Snapshots) HasHref() bool {
	if o != nil && o.Href != nil {
		return true
	}

	return false
}

// GetItems returns the Items field value
// If the value is explicit nil, the zero value for []Snapshot will be returned
func (o *Snapshots) GetItems() *[]Snapshot {
	if o == nil {
		return nil
	}

	return o.Items

}

// GetItemsOk returns a tuple with the Items field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *Snapshots) GetItemsOk() (*[]Snapshot, bool) {
	if o == nil {
		return nil, false
	}

	return o.Items, true
}

// SetItems sets field value
func (o *Snapshots) SetItems(v []Snapshot) {

	o.Items = &v

}

// HasItems returns a boolean if a field has been set.
func (o *Snapshots) HasItems() bool {
	if o != nil && o.Items != nil {
		return true
	}

	return false
}

func (o Snapshots) MarshalJSON() ([]byte, error) {
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

	if o.Items != nil {
		toSerialize["items"] = o.Items
	}
	return json.Marshal(toSerialize)
}

type NullableSnapshots struct {
	value *Snapshots
	isSet bool
}

func (v NullableSnapshots) Get() *Snapshots {
	return v.value
}

func (v *NullableSnapshots) Set(val *Snapshots) {
	v.value = val
	v.isSet = true
}

func (v NullableSnapshots) IsSet() bool {
	return v.isSet
}

func (v *NullableSnapshots) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableSnapshots(val *Snapshots) *NullableSnapshots {
	return &NullableSnapshots{value: val, isSet: true}
}

func (v NullableSnapshots) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableSnapshots) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
