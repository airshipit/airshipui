/*
Copyright (c) 2020 AT&T. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/
package plugin

import (
	"log"
	"strconv"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/domains"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/projects"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/users"
	"github.com/gophercloud/gophercloud/pagination"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

// gets OpenStack domains that are viewable by the tenant
// https://docs.openstack.org/security-guide/identity/domains.html
func getDomains(osp *OpenstackPlugin) component.Component {
	rows := []component.TableRow{}

	// TODO: Determine if the error needs to be handled from this function
	domains.List(identityClientHelper(osp), domains.ListOpts{}).EachPage(
		func(page pagination.Page) (bool, error) {
			domainList, err := domains.ExtractDomains(page)

			if err != nil {
				log.Fatalf("Broken at user list %v\n", err)
				return false, err
			}

			for _, domain := range domainList {
				name := domain.Name
				domainCache[domain.ID] = name
				rows = append(rows, component.TableRow{
					"Name":        component.NewText(name),
					"Description": component.NewText(domain.Description),
					"Enabled":     component.NewText(strconv.FormatBool(domain.Enabled)),
				})
			}

			return true, nil
		})

	return component.NewTableWithRows(
		"Domains",
		"No domains found",
		component.NewTableCols("Name", "Description", "Enabled"), rows)
}

// gets OpenStack projects that are viewable by the tenant
// https://docs.openstack.org/keystone/latest/admin/cli-manage-projects-users-and-roles.html
func getProjects(osp *OpenstackPlugin) component.Component {
	rows := []component.TableRow{}

	// TODO: Determine if the error needs to be handled from this function
	projects.List(identityClientHelper(osp), projects.ListOpts{}).EachPage(
		func(page pagination.Page) (bool, error) {
			projectList, err := projects.ExtractProjects(page)

			if err != nil {
				log.Fatalf("Broken at project list %v\n", err)
				return false, err
			}

			for _, project := range projectList {
				name := project.Name
				projectCache[project.ID] = name

				rows = append(rows, component.TableRow{
					"Enabled":     component.NewText(strconv.FormatBool(project.Enabled)),
					"Description": component.NewText(project.Description),
					"Name":        component.NewText(name),
				})
			}

			return true, nil
		})

	return component.NewTableWithRows(
		"Projects",
		"No projects found",
		component.NewTableCols("Name", "Enabled", "Description"), rows)
}

// gets OpenStack uesrs that are viewable by the tenant
// https://docs.openstack.org/keystone/latest/user/index.html
func getUsers(osp *OpenstackPlugin) component.Component {
	rows := []component.TableRow{}

	// TODO: Determine if the error needs to be handled from this function
	users.List(identityClientHelper(osp), users.ListOpts{}).EachPage(
		func(page pagination.Page) (bool, error) {
			networkList, err := users.ExtractUsers(page)

			if err != nil {
				log.Fatalf("Broken at user list %v\n", err)
				return false, err
			}

			for _, user := range networkList {
				email := ""
				if user.Extra["email"] != nil {
					email = user.Extra["email"].(string)
				}

				rows = append(rows, component.TableRow{
					"Name":        component.NewText(user.Name),
					"Email":       component.NewText(email),
					"Description": component.NewText(user.Description),
					"Enabled":     component.NewText(strconv.FormatBool(user.Enabled)),
					"Domain Name": component.NewText(domainCache[user.DomainID]),
				})
			}

			return true, nil
		})

	return component.NewTableWithRows(
		"Users",
		"No users found",
		component.NewTableCols("Name", "Email", "Description", "Enabled", "Domain Name"), rows)
}

// helper function to create an identity specific gophercloud client
func identityClientHelper(osp *OpenstackPlugin) *gophercloud.ServiceClient {
	client, err := openstack.NewIdentityV3(osp.provider, gophercloud.EndpointOpts{})

	if err != nil {
		log.Fatalf("Identity Client Error: %s\n", err)
		return nil
	}

	return client
}
