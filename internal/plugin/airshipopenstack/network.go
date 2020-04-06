/*
Copyright (c) 2020 AT&T. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package airshipopenstack

import (
	"log"
	"strconv"
	"strings"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/networks"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/subnets"
	"github.com/gophercloud/gophercloud/pagination"
)

// gets OpenStack networks available for the tenant
// https://docs.openstack.org/api-ref/network/v2/#networks
func GetNetworks() []map[string]string {
	osp := NewOpenstackPlugin()
	m := make([]map[string]string, 0)

	err := networks.List(networkClientHelper(osp), networks.ListOpts{}).EachPage(
		func(page pagination.Page) (bool, error) {
			networkList, err := networks.ExtractNetworks(page)

			if err != nil {
				log.Printf("Network retrival error: %s\n", err)
				return false, err
			}

			for _, n := range networkList {
				// pull together the subnet names from the cache
				var subnetString string
				if len(n.Subnets) > 1 {
					var tmp []string
					for subnetID := range n.Subnets {
						name := subnetCache[n.Subnets[subnetID]]
						if len(name) > 0 {
							tmp = append(tmp, name)
						}
					}
					subnetString = strings.Join(tmp, ", ")
				} else {
					subnetString = subnetCache[n.Subnets[0]]
				}

				if len(projectCache) == 0 {
					GetProjects()
				}
				m = append(m, map[string]string{
					"Name":               n.Name,
					"Project":            projectCache[n.ProjectID],
					"Subnets":            subnetString,
					"Shared":             strconv.FormatBool(n.Shared),
					"Status":             n.Status,
					"Admin State":        strconv.FormatBool(n.AdminStateUp),
					"Availability Zones": strings.Join(n.AvailabilityZoneHints, ", "),
				})
			}

			return true, nil
		})

	if err != nil {
		log.Printf("network list error: %s\n", err)
	}

	return m
}

// gets OpenStack subnets available for the tenant
// https://docs.openstack.org/api-ref/network/v2/#subnets
func GetSubnets() []map[string]string {
	osp := NewOpenstackPlugin()
	m := make([]map[string]string, 0)

	err := subnets.List(networkClientHelper(osp), subnets.ListOpts{}).EachPage(
		func(page pagination.Page) (bool, error) {
			networkList, err := subnets.ExtractSubnets(page)

			if err != nil {
				log.Printf("Subnet list error: %s\n", err)
				return false, err
			}

			for _, subnet := range networkList {
				cidr := subnet.CIDR
				name := subnet.Name
				subnetCache[subnet.ID] = name + ": " + cidr

				m = append(m, map[string]string{
					"Name":            name,
					"Network Address": cidr,
					"IP Version":      "IPv" + strconv.Itoa(subnet.IPVersion),
					"Gateway IP":      subnet.GatewayIP,
				})
			}

			return true, nil
		})

	if err != nil {
		log.Printf("network subnet list error: %s\n", err)
	}
	return m
}

// helper function to create a network specific gophercloud client
func networkClientHelper(osp *OpenstackPlugin) *gophercloud.ServiceClient {
	client, err := openstack.NewNetworkV2(osp.provider, gophercloud.EndpointOpts{
		Name:   "neutron",
		Region: "RegionOne",
	})

	if err != nil {
		log.Fatalf("Network Client Error: %s\n", err)
		return nil
	}

	return client
}
