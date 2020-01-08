package plugin

import (
	"fmt"
	"testing"

	"github.com/vmware-tanzu/octant/pkg/plugin/service"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

func TestRegister(t *testing.T) {
	p, err := Register("airship-ui", "Airship UI test version")
	if err != nil {
		t.Fatalf("Registering the plugin returned an error")
	}

	err = p.Validate()
	if err != nil {
		t.Fatalf("Validating the plugin returned an error")
	}
}

func TestRoutes(t *testing.T) {
	router := service.NewRouter()

	initRoutes(router)

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
func TestRouteHandles(t *testing.T) {
	router := service.NewRouter()

	initRoutes(router)

	tests := []struct {
		path  string
		title string
	}{
		{
			path:  "",
			title: "Argo UI",
		},
	}

	for _, test := range tests {
		test := test // pin the value so that the following function literal binds to it
		t.Run(fmt.Sprintf("Path='%s'", test.path), func(t *testing.T) {
			handleFunc, found := router.Match(test.path)

			if !found {
				t.Fatalf("Path '%s' was not found.", test.path)
			}

			request := &service.Request{
				Path: test.path,
			}
			contentResponse, err := handleFunc(request)
			if err != nil {
				t.Fatalf("handleFunc for path '%s' returned an error.", test.path)
			}

			title, err := component.TitleFromTitleComponent(contentResponse.Title)
			if err != nil {
				t.Fatalf("Getting the Title from the TitleComponents returned an error.")
			}

			if title != test.title {
				t.Errorf("Title is not correct. Got: '%s', Expected: '%s'", title, test.title)
			}
		})
	}
}
