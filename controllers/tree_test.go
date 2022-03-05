package controllers

import "testing"

func TestTree(t *testing.T) {
	tests := map[string]struct {
		hosts []string
		find  Host
		found bool
	}{
		"empty": {
			hosts: []string{},
			find:  "",
			found: false,
		},
		"single": {
			hosts: []string{"com"},
			find:  "com",
			found: true,
		},
		"single_not_found": {
			hosts: []string{"com"},
			find:  "io",
			found: false,
		},
		"two_level": {
			hosts: []string{"domain.com"},
			find:  "domain.com",
			found: true,
		},
		"two_level_not_found": {
			hosts: []string{"domain.com"},
			find:  "domain.io",
			found: false,
		},
		"two_level_not_found_three_level": {
			hosts: []string{"domain.com"},
			find:  "sub.domain.io",
			found: false,
		},
		"three_level": {
			hosts: []string{"sub.domain.com"},
			find:  "sub.domain.com",
			found: true,
		},
		"three_level_not_found": {
			hosts: []string{"domain.com"},
			find:  "sub.com.non",
			found: false,
		},
		"wildcard": {
			hosts: []string{"*.domain.com", "some.otherdomain.com"},
			find:  "sub.domain.com",
			found: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			tree := &tree{
				node: make(node),
			}
			for _, host := range test.hosts {
				tree.Add(host)
			}

			if test.found != tree.Contains(test.find) {
				t.Errorf("expected %v, got %v", test.found, tree.Contains(test.find))
			}
		})
	}
}
