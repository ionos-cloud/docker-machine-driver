package extraflag

type IntToStringSlice struct {
	Name   string
	Usage  string
	EnvVar string
	Value  map[int][]string
}

func (f IntToStringSlice) String() string {
	return f.Name
}

func (f IntToStringSlice) Default() interface{} {
	return f.Value
}
