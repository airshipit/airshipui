package plugin

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/vmware-tanzu/octant/pkg/plugin/service"
	"github.com/vmware-tanzu/octant/pkg/plugin/service/fake"
)

type fakeRequest struct {
	ctrl        *gomock.Controller
	path, title string
}

func (f *fakeRequest) DashboardClient() service.Dashboard {
	return fake.NewMockDashboard(f.ctrl)
}

func (f *fakeRequest) Path() string {
	return ""
}

func (f *fakeRequest) Context() context.Context {
	return context.TODO()
}

func TestRegister(t *testing.T) {
	plugin, err := Register("argo-ui", "Argo UI test version")
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
	NewArgoPlugin().initRoutes(router)

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
