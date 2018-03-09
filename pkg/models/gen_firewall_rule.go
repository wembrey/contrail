package models

import (
	"github.com/Juniper/contrail/pkg/common"
)

//To skip import error.
var _ = common.OPERATION

// MakeFirewallRule makes FirewallRule
// nolint
func MakeFirewallRule() *FirewallRule {
	return &FirewallRule{
		//TODO(nati): Apply default
		UUID:                 "",
		ParentUUID:           "",
		ParentType:           "",
		FQName:               []string{},
		IDPerms:              MakeIdPermsType(),
		DisplayName:          "",
		Annotations:          MakeKeyValuePairs(),
		Perms2:               MakePermType2(),
		ConfigurationVersion: 0,
		Endpoint1:            MakeFirewallRuleEndpointType(),
		Endpoint2:            MakeFirewallRuleEndpointType(),
		ActionList:           MakeActionListType(),
		Service:              MakeFirewallServiceType(),
		Direction:            "",
		MatchTagTypes:        MakeFirewallRuleMatchTagsTypeIdList(),
		MatchTags:            MakeFirewallRuleMatchTagsType(),
	}
}

// MakeFirewallRule makes FirewallRule
// nolint
func InterfaceToFirewallRule(i interface{}) *FirewallRule {
	m, ok := i.(map[string]interface{})
	_ = m
	if !ok {
		return nil
	}
	return &FirewallRule{
		//TODO(nati): Apply default
		UUID:                 common.InterfaceToString(m["uuid"]),
		ParentUUID:           common.InterfaceToString(m["parent_uuid"]),
		ParentType:           common.InterfaceToString(m["parent_type"]),
		FQName:               common.InterfaceToStringList(m["fq_name"]),
		IDPerms:              InterfaceToIdPermsType(m["id_perms"]),
		DisplayName:          common.InterfaceToString(m["display_name"]),
		Annotations:          InterfaceToKeyValuePairs(m["annotations"]),
		Perms2:               InterfaceToPermType2(m["perms2"]),
		ConfigurationVersion: common.InterfaceToInt64(m["configuration_version"]),
		Endpoint1:            InterfaceToFirewallRuleEndpointType(m["endpoint_1"]),
		Endpoint2:            InterfaceToFirewallRuleEndpointType(m["endpoint_2"]),
		ActionList:           InterfaceToActionListType(m["action_list"]),
		Service:              InterfaceToFirewallServiceType(m["service"]),
		Direction:            common.InterfaceToString(m["direction"]),
		MatchTagTypes:        InterfaceToFirewallRuleMatchTagsTypeIdList(m["match_tag_types"]),
		MatchTags:            InterfaceToFirewallRuleMatchTagsType(m["match_tags"]),

		ServiceGroupRefs: InterfaceToFirewallRuleServiceGroupRefs(m["service_group_refs"]),

		AddressGroupRefs: InterfaceToFirewallRuleAddressGroupRefs(m["address_group_refs"]),

		SecurityLoggingObjectRefs: InterfaceToFirewallRuleSecurityLoggingObjectRefs(m["security_logging_object_refs"]),

		VirtualNetworkRefs: InterfaceToFirewallRuleVirtualNetworkRefs(m["virtual_network_refs"]),
	}
}

func InterfaceToFirewallRuleSecurityLoggingObjectRefs(i interface{}) []*FirewallRuleSecurityLoggingObjectRef {
	list, ok := i.([]interface{})
	if !ok {
		return nil
	}
	result := []*FirewallRuleSecurityLoggingObjectRef{}
	for _, item := range list {
		m, ok := item.(map[string]interface{})
		_ = m
		if !ok {
			return nil
		}
		result = append(result, &FirewallRuleSecurityLoggingObjectRef{
			UUID: common.InterfaceToString(m["uuid"]),
			To:   common.InterfaceToStringList(m["to"]),
		})
	}

	return result
}

func InterfaceToFirewallRuleVirtualNetworkRefs(i interface{}) []*FirewallRuleVirtualNetworkRef {
	list, ok := i.([]interface{})
	if !ok {
		return nil
	}
	result := []*FirewallRuleVirtualNetworkRef{}
	for _, item := range list {
		m, ok := item.(map[string]interface{})
		_ = m
		if !ok {
			return nil
		}
		result = append(result, &FirewallRuleVirtualNetworkRef{
			UUID: common.InterfaceToString(m["uuid"]),
			To:   common.InterfaceToStringList(m["to"]),
		})
	}

	return result
}

func InterfaceToFirewallRuleServiceGroupRefs(i interface{}) []*FirewallRuleServiceGroupRef {
	list, ok := i.([]interface{})
	if !ok {
		return nil
	}
	result := []*FirewallRuleServiceGroupRef{}
	for _, item := range list {
		m, ok := item.(map[string]interface{})
		_ = m
		if !ok {
			return nil
		}
		result = append(result, &FirewallRuleServiceGroupRef{
			UUID: common.InterfaceToString(m["uuid"]),
			To:   common.InterfaceToStringList(m["to"]),
		})
	}

	return result
}

func InterfaceToFirewallRuleAddressGroupRefs(i interface{}) []*FirewallRuleAddressGroupRef {
	list, ok := i.([]interface{})
	if !ok {
		return nil
	}
	result := []*FirewallRuleAddressGroupRef{}
	for _, item := range list {
		m, ok := item.(map[string]interface{})
		_ = m
		if !ok {
			return nil
		}
		result = append(result, &FirewallRuleAddressGroupRef{
			UUID: common.InterfaceToString(m["uuid"]),
			To:   common.InterfaceToStringList(m["to"]),
		})
	}

	return result
}

// MakeFirewallRuleSlice() makes a slice of FirewallRule
// nolint
func MakeFirewallRuleSlice() []*FirewallRule {
	return []*FirewallRule{}
}

// InterfaceToFirewallRuleSlice() makes a slice of FirewallRule
// nolint
func InterfaceToFirewallRuleSlice(i interface{}) []*FirewallRule {
	list := common.InterfaceToInterfaceList(i)
	if list == nil {
		return nil
	}
	result := []*FirewallRule{}
	for _, item := range list {
		result = append(result, InterfaceToFirewallRule(item))
	}
	return result
}