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
	"github.com/gophercloud/gophercloud/openstack/blockstorage/v2/volumes"
	"github.com/gophercloud/gophercloud/pagination"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

// gets OpenStack volumes viewable by the tenant
// https://docs.openstack.org/cinder/latest/
func getVolumes(osp *OpenstackPlugin) component.Component {
	rows := []component.TableRow{}

	err := volumes.List(volumeClientHelper(osp), volumes.ListOpts{}).EachPage(
		func(page pagination.Page) (bool, error) {
			volumeList, err := volumes.ExtractVolumes(page)

			if err != nil {
				log.Printf("Broken at volume list %v\n", err)
				return false, err
			}

			for _, volume := range volumeList {
				// extract the potentially multiple devices
				var attachedDevices []string
				for index := range volume.Attachments {
					attachment := volume.Attachments[index]
					attachedDevices = append(attachedDevices, vmCache[attachment.ServerID]+" at "+attachment.Device)
				}

				rows = append(rows, component.TableRow{
					"Name":              component.NewText(volume.Name),
					"Description":       component.NewText(volume.Description),
					"Size":              component.NewText(strconv.Itoa(volume.Size) + "GiB"),
					"Status":            component.NewText(volume.Status),
					"Type":              component.NewText(volume.VolumeType),
					"Attached To":       component.NewText(strings.Join(attachedDevices, ", ")),
					"Availability Zone": component.NewText(volume.AvailabilityZone),
					"Bootable":          component.NewText(volume.Bootable),
					"Encrypted":         component.NewText(strconv.FormatBool(volume.Encrypted)),
				})
			}

			return true, nil
		})

	if err != nil {
		log.Printf("volume list error: %s\n", err)
	}

	return component.NewTableWithRows(
		"Volumes",
		"No volumes found",
		component.NewTableCols("Name", "Description", "Size", "Status", "Group", "Type", "Attached To",
			"Availability Zone", "Bootable", "Encrypted"), rows)
}

// helper function to create an volume specific gophercloud client
func volumeClientHelper(osp *OpenstackPlugin) *gophercloud.ServiceClient {
	client, err := openstack.NewBlockStorageV2(osp.provider, gophercloud.EndpointOpts{})

	if err != nil {
		log.Fatalf("NewIdentityV3 error: %v", err)
	}

	return client
}
