package strfmt

import (
	"bytes"
	"errors"
	"strings"
)

var (
	ErrInvalidMount = errors.New("invalid volume mount")
)

// data[/subPath]:mountPath[:ro]
func ParseVolumeMount(s string) (*VolumeMount, error) {
	if s == "" {
		return nil, ErrInvalidMount
	}

	parts := strings.Split(s, ":")
	vm := &VolumeMount{}

	if len(parts) > 3 {
		return nil, ErrInvalidMount
	}

	if len(parts) == 3 {
		n := len(parts)
		if parts[n-1] != "ro" {
			return nil, ErrInvalidMount
		}
		vm.ReadOnly = true
		parts = parts[0 : n-1]
	}

	if len(parts) != 2 {
		return nil, ErrInvalidMount
	}

	vm.MountPath = parts[1]

	volumeFrom := strings.Split(parts[0], "/")

	vm.Name = volumeFrom[0]

	if len(volumeFrom) == 2 {
		vm.SubPath = volumeFrom[1]
	}

	return vm, nil
}

// openapi:strfmt volume-mount
type VolumeMount struct {
	Name      string
	MountPath string
	SubPath   string
	ReadOnly  bool
}

func (v VolumeMount) String() string {
	buf := bytes.NewBufferString(v.Name)
	if v.SubPath != "" {
		buf.WriteByte('/')
		buf.WriteString(v.SubPath)
	}

	buf.WriteByte(':')
	buf.WriteString(v.MountPath)

	if v.ReadOnly {
		buf.WriteString(":ro")
	}

	return buf.String()
}

func (v VolumeMount) MarshalText() ([]byte, error) {
	return []byte(v.String()), nil
}

func (v *VolumeMount) UnmarshalText(data []byte) error {
	vm, err := ParseVolumeMount(string(data))
	if err != nil {
		return err
	}
	*v = *vm
	return nil
}
