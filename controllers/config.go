package controllers

import "github.com/suffiks/suffiks/extension"

type Config struct {
	extension.ConfigSpec `json:",inline"`

	// AllowedDomains is a list of domains that are allowed to access the extension.
	AllowedDomains []string `json:"allowedDomains"`
}
