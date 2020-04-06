/*
Copyright (c) 2020 AT&T. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/
package airshipopenstack

import (
	"encoding/json"
	"log"
	"strconv"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/domains"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/projects"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/users"
	"github.com/gophercloud/gophercloud/pagination"
)

// gets OpenStack domains that are viewable by the tenant
// https://docs.openstack.org/security-guide/identity/domains.html
func GetDomains() []map[string]string {
	osp := NewOpenstackPlugin()
	m := make([]map[string]string, 0)

	err := domains.List(identityClientHelper(osp), domains.ListOpts{}).EachPage(
		func(page pagination.Page) (bool, error) {
			domainList, err := domains.ExtractDomains(page)

			if err != nil {
				log.Printf("Broken at domain list %v\n", err)
				return false, err
			}

			for _, domain := range domainList {
				name := domain.Name
				domainCache[domain.ID] = name
				m = append(m, map[string]string{
					"Name":        name,
					"Description": domain.Description,
					"Enabled":     strconv.FormatBool(domain.Enabled),
				})
			}

			return true, nil
		})

	if err != nil {
		log.Printf("identity domain list error: %s\n", err)
	}

	return m
}

// gets OpenStack projects that are viewable by the tenant
// https://docs.openstack.org/keystone/latest/admin/cli-manage-projects-users-and-roles.html
func GetProjects() []map[string]string {
	osp := NewOpenstackPlugin()
	m := make([]map[string]string, 0)

	err := projects.List(identityClientHelper(osp), projects.ListOpts{}).EachPage(
		func(page pagination.Page) (bool, error) {
			projectList, err := projects.ExtractProjects(page)

			if err != nil {
				log.Printf("Broken at project list %v\n", err)
				return false, err
			}

			for _, project := range projectList {
				name := project.Name
				projectCache[project.ID] = name

				m = append(m, map[string]string{
					"Enabled":     strconv.FormatBool(project.Enabled),
					"Description": project.Description,
					"Name":        name,
				})
			}

			return true, nil
		})

	if err != nil {
		log.Printf("identity project list error: %s\n", err)
	}

	return m
}

// gets OpenStack uesrs that are viewable by the tenant
// https://docs.openstack.org/keystone/latest/user/index.html
func GetUsers() []map[string]string {
	osp := NewOpenstackPlugin()
	m := make([]map[string]string, 0)

	err := users.List(identityClientHelper(osp), users.ListOpts{}).EachPage(
		func(page pagination.Page) (bool, error) {
			networkList, err := users.ExtractUsers(page)

			if err != nil {
				log.Printf("Broken at user list %v\n", err)
				return false, err
			}

			for _, user := range networkList {
				var email string
				emailInterface, ok := user.Extra["email"]
				if ok && emailInterface != nil {
					b, err := json.Marshal(emailInterface)
					if err != nil {
						log.Printf("Error getting email %v\n", err)
					}
					email = string(b)
				}

				if len(domainCache) == 0 {
					GetDomains()
				}
				m = append(m, map[string]string{
					"Name":        user.Name,
					"Email":       email,
					"Description": user.Description,
					"Enabled":     strconv.FormatBool(user.Enabled),
					"Domain Name": domainCache[user.DomainID],
				})
			}

			return true, nil
		})

	if err != nil {
		log.Printf("identity user list error: %s\n", err)
	}

	return m
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
