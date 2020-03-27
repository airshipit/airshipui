# AirshipUI OpenStack Plugin Documentation

## OpenStack Plugin Requirements
In order to make a connection to the OpenStack APIs the system needs enough information to attach to it.  There are 2 ways to accomplish this.

### File Based OpenStack Connection Details

Create an openstack.json file in the correct location for your system.

Windows:
```
C:\Users\<USER_ID>\AppData\Local\octant\etc\openstack.json
```

Linux:
```
~/.config/octant/etc/openstack.json
```

The file contents should look like this:
```
{
	"identityEndpoint": "http://<OPENSTACK_HOST>/identity/v3",
	"username": "<USERNAME>",
	"password": "<PASSWORD>",
	"tenantName": "<TENANT hint demo is the default>",
	"domainID": "<DOMAINID hint default is the default>"
}
```

### Environment Based OpenStack Connection Details

The following environment variables are required to be set in the same shell that is executing the Octant or AirshipUI binary:
```
	OS_USER_DOMAIN_ID
	OS_AUTH_URL
	OS_PROJECT_DOMAIN_ID
	OS_REGION_NAME
	OS_PROJECT_NAME
	OS_IDENTITY_API_VERSION
	OS_TENANT_NAME
	OS_TENANT_ID
	OS_AUTH_TYPE
	OS_PASSWORD
	OS_USERNAME
	OS_VOLUME_API_VERSION
	OS_TOKEN
	OS_USERID
```

This plugin should now have enough information to start and display data from your OpenStack instance

The next time you run Airship UI the OpenStack plugin will be available at http://127.0.0.1:7777/#/openstack