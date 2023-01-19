package extflag

import (
	"reflect"
	"testing"
)

func TestToMapOfStringToStringSlice(t *testing.T) {
	tests := []struct {
		name string
		val  string
		want map[string][]string
	}{
		{
			name: "Nat Gateways example 1",
			val:  "1=10.0.0.1,10.0.0.2:2=10.0.0.10",
			want: map[string][]string{"1": {"10.0.0.1", "10.0.0.2"}, "2": {"10.0.0.10"}},
		},
		{
			name: "Nat Gateways example 2",
			val:  "1=10.0.0.1,10.0.0.2:2=10.0.0.10:3=1,2,3,4,5,6,7,8,99:11=Foo,Bar",
			want: map[string][]string{"1": {"10.0.0.1", "10.0.0.2"}, "2": {"10.0.0.10"}, "3": {"1", "2", "3", "4", "5", "6", "7", "8", "99"}, "11": {"Foo", "Bar"}},
		},
		{
			name: "nil",
			val:  "",
			want: nil,
		},
		{
			name: "bogus 1",
			val:  "1,2,3",
			want: nil,
		},
		{
			name: "bogus 2",
			val:  "1:1,2,3",
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToMapOfStringToStringSlice(tt.val); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToMapOfStringToStringSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKebabCaseToCamelCase(t *testing.T) {
	tests := []struct {
		name  string
		kebab string
		want  string
	}{
		{
			name:  "flag example",
			kebab: "ionoscloud-endpoint",
			want:  "IONOSCLOUD_ENDPOINT",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := KebabCaseToCamelCase(tt.kebab); got != tt.want {
				t.Errorf("KebabCaseToCamelCase() = %v, want %v", got, tt.want)
			}
		})
	}
}
