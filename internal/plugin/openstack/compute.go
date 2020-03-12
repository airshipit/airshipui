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
	"github.com/gophercloud/gophercloud/openstack/compute/v2/flavors"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/images"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	"github.com/gophercloud/gophercloud/pagination"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

// gets OpenStack flavors that are available for the tenant
// https://docs.openstack.org/nova/latest/user/flavors.html
func getFlavors(osp *OpenstackPlugin) component.Component {
	rows := []component.TableRow{}

	// TODO: Determine if the error needs to be handled from this function
	flavors.ListDetail(computeClientHelper(osp), flavors.ListOpts{AccessType: flavors.AllAccess}).EachPage(
		func(page pagination.Page) (bool, error) {
			flavorList, err := flavors.ExtractFlavors(page)

			if err != nil {
				log.Fatalf("compute flavor Error: %s\n", err)
			}

			for _, flavor := range flavorList {
				name := flavor.Name
				flavorCache[flavor.ID] = name

				rows = append(rows, component.TableRow{
					"Name":           component.NewText(name),
					"VCPUs":          component.NewText(strconv.Itoa(flavor.VCPUs)),
					"RAM":            component.NewText(strconv.Itoa(flavor.RAM)),
					"Root Disk":      component.NewText(strconv.Itoa(flavor.Disk)),
					"Ephemeral Disk": component.NewText(strconv.Itoa(flavor.Ephemeral)),
					"Swap Disk":      component.NewText(strconv.Itoa(flavor.Swap)),
					"RX/TX factor":   component.NewText(strconv.FormatFloat(flavor.RxTxFactor, 'f', 1, 64)),
					"Public":         component.NewText(strconv.FormatBool(flavor.IsPublic)),
				})
			}

			return true, nil
		})

	return component.NewTableWithRows(
		"Flavors",
		"No flavors found",
		component.NewTableCols("Name", "VCPUs", "RAM", "Root Disk", "Ephemeral Disk",
			"Swap Disk", "RX/TX factor", "Public"),
		rows)
}

// gets OpenStack images that are available for the tenant
// https://docs.openstack.org/image-guide/
func getImages(osp *OpenstackPlugin) component.Component {
	rows := []component.TableRow{}

	// TODO: Determine if the error needs to be handled from this function
	images.ListDetail(computeClientHelper(osp), images.ListOpts{}).EachPage(
		func(page pagination.Page) (bool, error) {
			list, err := images.ExtractImages(page)

			if err != nil {
				log.Fatalf("Image list error: %s\n", err)
				return false, err
			}

			for _, image := range list {
				name := image.Name
				imageCache[image.ID] = name

				rows = append(rows, component.TableRow{
					"Name":   component.NewText(name),
					"Type":   component.NewText(strconv.Itoa(image.Progress)),
					"Status": component.NewText(image.Status),
				})
			}

			return true, nil
		})

	return component.NewTableWithRows(
		"Images",
		"No images found",
		component.NewTableCols("Name", "Type", "Status"), rows)
}

// gets OpenStack vms that are configured for the tenant
// https://docs.openstack.org/python-openstackclient/latest/cli/command-objects/server.html
func getVMs(osp *OpenstackPlugin) component.Component {
	rows := []component.TableRow{}

	// TODO: Determine if the error needs to be handled from this function
	servers.List(computeClientHelper(osp), servers.ListOpts{AllTenants: true}).EachPage(
		func(page pagination.Page) (bool, error) {
			list, err := servers.ExtractServers(page)

			if err != nil {
				log.Fatalf("Server list error: %s\n", err)
				return false, err
			}

			for _, server := range list {
				name := server.Name
				vmCache[server.ID] = name

				// pull together the addresses from the structure
				var addresses string
				var tmp []string
				for subnetName := range server.Addresses {
					for _, v := range server.Addresses[subnetName].([]interface{}) {
						tmp = append(tmp, v.(map[string]interface{})["addr"].(string))
					}
				}

				if len(tmp) > 1 {
					addresses = strings.Join(tmp, ", ")
				} else {
					addresses = tmp[0]
				}

				rows = append(rows, component.TableRow{
					"Instance Name": component.NewText(name),
					"Image Name":    component.NewText(imageCache[server.Image["id"].(string)]),
					"IP Address":    component.NewText(addresses),
					"Flavor":        component.NewText(flavorCache[server.Flavor["id"].(string)]),
					"Key Pair":      component.NewText(server.KeyName),
					"Status":        component.NewText(server.Status),
					"Created":       component.NewText(server.Created.Local().String()),
				})
			}

			return true, nil
		})

	return component.NewTableWithRows(
		"Servers",
		"No servers found",
		component.NewTableCols("Instance Name", "Image Name", "IP Address", "Flavor", "Key Pair", "Status", "Created"), rows)
}

// helper function to create a compute specific gophercloud client
func computeClientHelper(osp *OpenstackPlugin) *gophercloud.ServiceClient {
	client, err := openstack.NewComputeV2(osp.provider, gophercloud.EndpointOpts{
		Type: "compute",
	})

	if err != nil {
		log.Fatalf("Compute Client Error: %s\n", err)
		return nil
	}

	return client
}
