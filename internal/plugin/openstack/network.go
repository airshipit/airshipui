/*
Copyright (c) 2020 AT&T. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package plugin

import (
	"log"
	"strconv"
	"strings"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/networks"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/subnets"
	"github.com/gophercloud/gophercloud/pagination"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

// gets OpenStack networks available for the tenant
// https://docs.openstack.org/api-ref/network/v2/#networks
func getNetworks(osp *OpenstackPlugin) component.Component {
	rows := []component.TableRow{}

	// TODO: Determine if the error needs to be handled from this function
	networks.List(networkClientHelper(osp), networks.ListOpts{}).EachPage(
		func(page pagination.Page) (bool, error) {
			networkList, err := networks.ExtractNetworks(page)

			if err != nil {
				log.Fatalf("Network retrival error: %s\n", err)
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

				rows = append(rows, component.TableRow{
					"Name":               component.NewText(n.Name),
					"Project":            component.NewText(projectCache[n.ProjectID]),
					"Subnets":            component.NewText(subnetString),
					"Shared":             component.NewText(strconv.FormatBool(n.Shared)),
					"Status":             component.NewText(n.Status),
					"Admin State":        component.NewText(strconv.FormatBool(n.AdminStateUp)),
					"Availability Zones": component.NewText(strings.Join(n.AvailabilityZoneHints, ", ")),
				})
			}

			return true, nil
		})

	return component.NewTableWithRows(
		"Networks",
		"No networks found",
		component.NewTableCols("Name", "Project", "Subnets", "Shared", "Status",
			"Admin State", "Availability Zones"),
		rows)
}

// gets OpenStack subnets available for the tenant
// https://docs.openstack.org/api-ref/network/v2/#subnets
func getSubnets(osp *OpenstackPlugin) component.Component {
	rows := []component.TableRow{}

	// TODO: Determine if the error needs to be handled from this function
	subnets.List(networkClientHelper(osp), subnets.ListOpts{}).EachPage(
		func(page pagination.Page) (bool, error) {
			networkList, err := subnets.ExtractSubnets(page)

			if err != nil {
				log.Fatalf("Subnet list error: %s\n", err)
				return false, err
			}

			for _, subnet := range networkList {
				cidr := subnet.CIDR
				name := subnet.Name
				subnetCache[subnet.ID] = name + ": " + cidr

				rows = append(rows, component.TableRow{
					"Name":            component.NewText(name),
					"Network Address": component.NewText(cidr),
					"IP Version":      component.NewText("IPv" + strconv.Itoa(subnet.IPVersion)),
					"Gateway IP":      component.NewText(subnet.GatewayIP),
				})
			}

			return true, nil
		})

	return component.NewTableWithRows(
		"Subnets",
		"No subnets found",
		component.NewTableCols("Name", "Network Address", "IP Version", "Gateway IP"), rows)
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
