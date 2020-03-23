/*
Copyright (c) 2020 AT&T. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/
package plugin

import (
	"fmt"
	"testing"

	"github.com/vmware-tanzu/octant/pkg/plugin/service"
)

func TestRegister(t *testing.T) {
	plugin, err := Register("openstack", "OpenStack test version")
	if err != nil {
		t.Fatalf("Registering the plugin returned an error")
	}

	err = plugin.Validate()
	if err != nil {
		t.Fatalf("Validating the plugin returned an error")
	}
}

func TestRoutes(t *testing.T) {
	router := service.NewRouter()
	NewOpenstackPlugin().initRoutes(router)

	tests := []struct {
		path   string
		exists bool
	}{
		{
			path:   "",
			exists: true,
		},
		{
			path:   "/not-real",
			exists: false,
		},
	}

	for _, test := range tests {
		test := test // pin the value so that the following function literal binds to it
		t.Run(fmt.Sprintf("Path='%s'", test.path), func(t *testing.T) {
			_, found := router.Match(test.path)

			if test.exists != found {
				if found {
					t.Errorf("Found path '%s' when it should not exist.", test.path)
				} else {
					t.Errorf("Didn't find path '%s' when it should exist.", test.path)
				}
			}
		})
	}
}
