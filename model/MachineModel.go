package model

type Machine struct {
	InstanceName string  `json:"instance_name,omitempty"`
	PrivateIP    string  `json:"private_ip,omitempty"`
	PublicIP     *string `json:"public_ip,omitempty"`
	Password     string  `json:"password,omitempty"`
}
