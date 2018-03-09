package models

import (
	"github.com/Juniper/contrail/pkg/common"
)

//To skip import error.
var _ = common.OPERATION

// MakeVPNGroup makes VPNGroup
// nolint
func MakeVPNGroup() *VPNGroup {
	return &VPNGroup{
		//TODO(nati): Apply default
		ProvisioningLog:           "",
		ProvisioningProgress:      0,
		ProvisioningProgressStage: "",
		ProvisioningStartTime:     "",
		ProvisioningState:         "",
		UUID:                      "",
		ParentUUID:                "",
		ParentType:                "",
		FQName:                    []string{},
		IDPerms:                   MakeIdPermsType(),
		DisplayName:               "",
		Annotations:               MakeKeyValuePairs(),
		Perms2:                    MakePermType2(),
		ConfigurationVersion:      0,
		Type:                      "",
	}
}

// MakeVPNGroup makes VPNGroup
// nolint
func InterfaceToVPNGroup(i interface{}) *VPNGroup {
	m, ok := i.(map[string]interface{})
	_ = m
	if !ok {
		return nil
	}
	return &VPNGroup{
		//TODO(nati): Apply default
		ProvisioningLog:           common.InterfaceToString(m["provisioning_log"]),
		ProvisioningProgress:      common.InterfaceToInt64(m["provisioning_progress"]),
		ProvisioningProgressStage: common.InterfaceToString(m["provisioning_progress_stage"]),
		ProvisioningStartTime:     common.InterfaceToString(m["provisioning_start_time"]),
		ProvisioningState:         common.InterfaceToString(m["provisioning_state"]),
		UUID:                      common.InterfaceToString(m["uuid"]),
		ParentUUID:                common.InterfaceToString(m["parent_uuid"]),
		ParentType:                common.InterfaceToString(m["parent_type"]),
		FQName:                    common.InterfaceToStringList(m["fq_name"]),
		IDPerms:                   InterfaceToIdPermsType(m["id_perms"]),
		DisplayName:               common.InterfaceToString(m["display_name"]),
		Annotations:               InterfaceToKeyValuePairs(m["annotations"]),
		Perms2:                    InterfaceToPermType2(m["perms2"]),
		ConfigurationVersion:      common.InterfaceToInt64(m["configuration_version"]),
		Type:                      common.InterfaceToString(m["type"]),

		LocationRefs: InterfaceToVPNGroupLocationRefs(m["location_refs"]),
	}
}

func InterfaceToVPNGroupLocationRefs(i interface{}) []*VPNGroupLocationRef {
	list, ok := i.([]interface{})
	if !ok {
		return nil
	}
	result := []*VPNGroupLocationRef{}
	for _, item := range list {
		m, ok := item.(map[string]interface{})
		_ = m
		if !ok {
			return nil
		}
		result = append(result, &VPNGroupLocationRef{
			UUID: common.InterfaceToString(m["uuid"]),
			To:   common.InterfaceToStringList(m["to"]),
		})
	}

	return result
}

// MakeVPNGroupSlice() makes a slice of VPNGroup
// nolint
func MakeVPNGroupSlice() []*VPNGroup {
	return []*VPNGroup{}
}

// InterfaceToVPNGroupSlice() makes a slice of VPNGroup
// nolint
func InterfaceToVPNGroupSlice(i interface{}) []*VPNGroup {
	list := common.InterfaceToInterfaceList(i)
	if list == nil {
		return nil
	}
	result := []*VPNGroup{}
	for _, item := range list {
		result = append(result, InterfaceToVPNGroup(item))
	}
	return result
}