/*
 * CLOUD API
 *
 *  IONOS Enterprise-grade Infrastructure as a Service (IaaS) solutions can be managed through the Cloud API, in addition or as an alternative to the \"Data Center Designer\" (DCD) browser-based tool.    Both methods employ consistent concepts and features, deliver similar power and flexibility, and can be used to perform a multitude of management tasks, including adding servers, volumes, configuring networks, and so on.
 *
 * API version: 6.0
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package ionoscloud

import (
	"encoding/json"
)

// SecurityGroupEntities struct for SecurityGroupEntities
type SecurityGroupEntities struct {
	Rules   *FirewallRules `json:"rules,omitempty"`
	Nics    *Nics          `json:"nics,omitempty"`
	Servers *Servers       `json:"servers,omitempty"`
}

// NewSecurityGroupEntities instantiates a new SecurityGroupEntities object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewSecurityGroupEntities() *SecurityGroupEntities {
	this := SecurityGroupEntities{}

	return &this
}

// NewSecurityGroupEntitiesWithDefaults instantiates a new SecurityGroupEntities object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewSecurityGroupEntitiesWithDefaults() *SecurityGroupEntities {
	this := SecurityGroupEntities{}
	return &this
}

// GetRules returns the Rules field value
// If the value is explicit nil, nil is returned
func (o *SecurityGroupEntities) GetRules() *FirewallRules {
	if o == nil {
		return nil
	}

	return o.Rules

}

// GetRulesOk returns a tuple with the Rules field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *SecurityGroupEntities) GetRulesOk() (*FirewallRules, bool) {
	if o == nil {
		return nil, false
	}

	return o.Rules, true
}

// SetRules sets field value
func (o *SecurityGroupEntities) SetRules(v FirewallRules) {

	o.Rules = &v

}

// HasRules returns a boolean if a field has been set.
func (o *SecurityGroupEntities) HasRules() bool {
	if o != nil && o.Rules != nil {
		return true
	}

	return false
}

// GetNics returns the Nics field value
// If the value is explicit nil, nil is returned
func (o *SecurityGroupEntities) GetNics() *Nics {
	if o == nil {
		return nil
	}

	return o.Nics

}

// GetNicsOk returns a tuple with the Nics field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *SecurityGroupEntities) GetNicsOk() (*Nics, bool) {
	if o == nil {
		return nil, false
	}

	return o.Nics, true
}

// SetNics sets field value
func (o *SecurityGroupEntities) SetNics(v Nics) {

	o.Nics = &v

}

// HasNics returns a boolean if a field has been set.
func (o *SecurityGroupEntities) HasNics() bool {
	if o != nil && o.Nics != nil {
		return true
	}

	return false
}

// GetServers returns the Servers field value
// If the value is explicit nil, nil is returned
func (o *SecurityGroupEntities) GetServers() *Servers {
	if o == nil {
		return nil
	}

	return o.Servers

}

// GetServersOk returns a tuple with the Servers field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *SecurityGroupEntities) GetServersOk() (*Servers, bool) {
	if o == nil {
		return nil, false
	}

	return o.Servers, true
}

// SetServers sets field value
func (o *SecurityGroupEntities) SetServers(v Servers) {

	o.Servers = &v

}

// HasServers returns a boolean if a field has been set.
func (o *SecurityGroupEntities) HasServers() bool {
	if o != nil && o.Servers != nil {
		return true
	}

	return false
}

func (o SecurityGroupEntities) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Rules != nil {
		toSerialize["rules"] = o.Rules
	}

	if o.Nics != nil {
		toSerialize["nics"] = o.Nics
	}

	if o.Servers != nil {
		toSerialize["servers"] = o.Servers
	}

	return json.Marshal(toSerialize)
}

type NullableSecurityGroupEntities struct {
	value *SecurityGroupEntities
	isSet bool
}

func (v NullableSecurityGroupEntities) Get() *SecurityGroupEntities {
	return v.value
}

func (v *NullableSecurityGroupEntities) Set(val *SecurityGroupEntities) {
	v.value = val
	v.isSet = true
}

func (v NullableSecurityGroupEntities) IsSet() bool {
	return v.isSet
}

func (v *NullableSecurityGroupEntities) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableSecurityGroupEntities(val *SecurityGroupEntities) *NullableSecurityGroupEntities {
	return &NullableSecurityGroupEntities{value: val, isSet: true}
}

func (v NullableSecurityGroupEntities) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableSecurityGroupEntities) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
