package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestApplianceFullName(t *testing.T) {
	cases := []struct {
		a        Appliance
		expected string
	}{
		{a: Appliance{Proto: "app", OS: "os", Arch: "arch", Name: "name"},
			expected: "os-arch-name:latest",
		},
		{a: Appliance{Proto: "vm", OS: "os", Arch: "arch", Name: "name"},
			expected: "name:latest",
		},
		{a: Appliance{Proto: "vm", OS: "os", Arch: "arch", Name: "name", Tag: "1.0"},
			expected: "name:1.0",
		},
	}

	for _, c := range cases {
		assert.Equal(t, c.expected, c.a.FullName(), "Fail to get correct fullname")
	}
}
