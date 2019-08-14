// Copyright 2017 Intel Corp.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package usrsptypes

import (
	"net"

	"github.com/containernetworking/cni/pkg/types"
	"github.com/containernetworking/cni/pkg/types/current"
)

//
// Exported Types
//
type MemifConf struct {
	Role        string  `json:"role,omitempty"`  // Role of memif: master|slave
	Mode        string  `json:"mode,omitempty"`  // Mode of memif: ip|ethernet|inject-punt

	// Autogenerated as memif-<ContainerID[:12]>-<IfName>.sock i.e. memif-0958c8871b32-net1.sock
	// Filename only, no path. Will use if populated, but used to passed filename to container.
	Socketfile  string  `json:"socketfile,omitempty"`
}

type VhostConf struct {
	Mode        string  `json:"mode,omitempty"`  // vhost-user mode: client|server

	// Autogenerated as <ContainerID[:12]>-<IfName> i.e. 0958c8871b32-net1
	// Filename only, no path. Will use if populated, but used to passed filename to container.
	Socketfile  string  `json:"socketfile,omitempty"`
}

type BridgeConf struct {
	// ovs-dpdk specific note:
	//   ovs-dpdk requires a bridge to create an interfaces. So if 'NetType' is set
	//   to something other than 'bridge', a bridge is still need and this field will
	//   be inspected. For ovs-dpdk, if bridge data is not populated, it will default
	//   to 'br-0'. 
	BridgeName  string `json:"bridgeName,omitempty"` // Bridge Name
	BridgeId    int    `json:"bridgeId,omitempty"`   // Bridge Id - Deprecated in favor of BridgeName
	VlanId      int    `json:"vlanId,omitempty"`     // Optional VLAN Id
}

type UserSpaceConf struct {
	// The Container Instance will default to the Host Instance value if a given attribute
	// is not provided. However, they are not required to be the same and a Container
	// attribute can be provided to override. All values are listed as 'omitempty' to
	// allow the Container struct to be empty where desired.
	Engine     string     `json:"engine,omitempty"`  // CNI Implementation {vpp|ovs-dpdk}
	IfType     string     `json:"iftype,omitempty"`  // Type of interface {memif|vhostuser}
	NetType    string     `json:"netType,omitempty"` // Interface network type {none|bridge|interface}
	MemifConf  MemifConf  `json:"memif,omitempty"`
	VhostConf  VhostConf  `json:"vhost,omitempty"`
	BridgeConf BridgeConf `json:"bridge,omitempty"`
}

type NetConf struct {
	types.NetConf

	/*
	// Support chaining
	RawPrevResult *map[string]interface{} `json:"prevResult"`
	PrevResult    *current.Result         `json:"-"`
	*/

	// One of the following two must be provided: KubeConfig or SharedDir
	//
	// KubeConfig:
	//  Example: "kubeconfig": "/etc/cni/net.d/multus.d/multus.kubeconfig",
	//  Provides credentials for Userspace CNI to call KubeAPI to:
	//  - Read Volume Mounts:
	//    - "shared-dir": Directory on host socketfiles are created in
	//  - Write annotations:
	//    - "userspace-cni/configuration-data": Configuration data passed
	//      to containe in JSON format.
	//    - "userspace-cni/mapped-dir": Directory in container socketfiles
	//      are created in. Scraped from Volume Mounts above.
	//
	// SharedDir:
	//  Example: "sharedDir": "/usr/local/var/run/openvswitch/023bcd123/",
	//  Since credentials are not provided, Userspace CNI cannot call KubeAPI
	//  to read the Volume Mounts, so this is the same directory used in the
	//  "hostPath".
	//  Along the same lines, no annotations are written by Userspace CNI.
	//   1) Configuration data will be written to a file in the same
	//      directory as the socketfiles instead of to an annotation.
	//   2) The "userspace-cni/mapped-dir" annotation must be added to the
	//      pod spec so container know where to retrieve data.
	//      Example: userspace-cni/mappedDir: /var/lib/cni/usrspcni/
	KubeConfig    string        `json:"kubeconfig,omitempty"`
	SharedDir     string        `json:"sharedDir,omitempty"`

	LogFile       string        `json:"logFile,omitempty"`
	LogLevel      string        `json:"logLevel,omitempty"`

	Name          string        `json:"name"`
	HostConf      UserSpaceConf `json:"host,omitempty"`
	ContainerConf UserSpaceConf `json:"container,omitempty"`
}

// Defines the JSON data written to container. It is either written to:
//  1) Annotation - "userspace-cni/configuration-data"
//  -- OR --
//  2) a file in the directory designated by NetConf.SharedDir.
type ConfigurationData struct {
	ContainerId   string                   `json:"containerId"` // From args.ContainerId, used locally. Used in several place, namely in the socket filenames.
	IfName        string                   `json:"ifName"`      // From args.IfName, used locally. Used in several place, namely in the socket filenames.
	Name          string                   `json:"name"`        // From NetConf.Name
	Config        UserSpaceConf            `json:"config"`      // From NetConf.ContainerConf
	IPResult      current.Result           `json:"ipResult"`    // Network Status also has IP, but wrong format
}

// UnmarshallableString typedef for builtin string
type UnmarshallableString string

// K8sArgs is the valid CNI_ARGS used for Kubernetes
type K8sArgs struct {
	types.CommonArgs
	IP                         net.IP
	K8S_POD_NAME               types.UnmarshallableString
	K8S_POD_NAMESPACE          types.UnmarshallableString
	K8S_POD_INFRA_CONTAINER_ID types.UnmarshallableString
}
