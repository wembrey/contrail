package models

import (
	"strconv"
	"strings"

	"github.com/Juniper/contrail/pkg/common"
	"github.com/Juniper/contrail/pkg/models/basemodels"
)

// FirewallRule constants.
const (
	DefaultMatchTagType               = "application"
	FirewallPolicyTagNameGlobalPrefix = "global:"
)

var protocolIDs = map[string]int64{
	"any":  0,
	"icmp": 1,
	"tcp":  6,
	"udp":  17,
}

// CheckAssociatedRefsInSameScope checks scope of firewallRule refs
// that global scoped firewall rule cannot reference scoped resources
func (fr *FirewallRule) CheckAssociatedRefsInSameScope() error {

	// this method is simply based on the fact global
	// firewall resource (draft or not) have a FQ name length equal to two
	// and scoped one (draft or not) have a FQ name longer than 2 to
	// distinguish a scoped firewall resource to a global one. If that
	// assumption disappear, all that method will need to be re-worked
	if len(fr.GetFQName()) != 2 {
		return nil
	}

	for _, ref := range fr.GetAddressGroupRefs() {
		if err := fr.checkRefInSameScope(ref); err != nil {
			return err
		}
	}

	for _, ref := range fr.GetServiceGroupRefs() {
		if err := fr.checkRefInSameScope(ref); err != nil {
			return err
		}
	}

	for _, ref := range fr.GetVirtualNetworkRefs() {
		if err := fr.checkRefInSameScope(ref); err != nil {
			return err
		}
	}

	return nil
}

func (fr *FirewallRule) checkRefInSameScope(ref basemodels.Reference) error {
	if len(ref.GetTo()) == 2 {
		return nil
	}

	refKind := strings.Replace(ref.GetReferredKind(), "-", " ", -1)
	return common.ErrorBadRequestf(
		"global Firewall Rule %s (%s) cannot reference a scoped %s %s (%s)",
		basemodels.FQNameToString(fr.GetFQName()),
		fr.GetUUID(),
		strings.Title(refKind),
		basemodels.FQNameToString(ref.GetTo()),
		ref.GetUUID(),
	)
}

// CheckServiceProperties checks for existence of service and serviceGroupRefs property
func (fr *FirewallRule) CheckServiceProperties() error {
	serviceGroupRefs := fr.GetServiceGroupRefs()
	service := fr.GetService()

	if service == nil && len(serviceGroupRefs) == 0 {
		return common.ErrorBadRequest("firewall Rule requires at least 'service' property or Service Group references(s)")
	}

	if service != nil && len(serviceGroupRefs) > 0 {
		return common.ErrorBadRequest(
			"firewall Rule cannot have both defined 'service' property and Service Group reference(s)",
		)
	}

	return nil
}

// AddDefaultMatchTag sets default matchTag if not defined in the request
func (fr *FirewallRule) AddDefaultMatchTag() {
	if fr.GetMatchTags() == nil {
		fr.MatchTags = &FirewallRuleMatchTagsType{
			TagList: []string{DefaultMatchTagType},
		}
	}
}

func (fr *FirewallRule) getProtocolID() (int64, error) {
	protocol := fr.GetService().GetProtocol()
	ok := true
	protocolID, err := strconv.ParseInt(protocol, 10, 64)
	if err != nil {
		protocolID, ok = protocolIDs[protocol]
	}

	if !ok || protocolID < 0 || protocolID > 255 {
		return 0, common.ErrorBadRequestf("rule with invalid protocol: %s", protocol)
	}

	return protocolID, nil
}

// SetProtocolID sets protocolID based on protocol property
func (fr *FirewallRule) SetProtocolID() error {
	protocolID, err := fr.getProtocolID()
	fr.Service.ProtocolID = protocolID
	return err
}

// CheckEndpoints validates endpoint_1 and endpoint_2
func (fr *FirewallRule) CheckEndpoints() error {

	//TODO Authorize only one condition per endpoint for the moment
	endpoints := []*FirewallRuleEndpointType{
		fr.GetEndpoint1(),
		fr.GetEndpoint2(),
	}

	for _, endpoint := range endpoints {
		if err := endpoint.ValidateEndpointType(); err != nil {
			return err
		}
	}

	//TODO check no ids present
	//TODO check endpoints exclusivity clause
	//TODO VN name in endpoints

	return nil
}

// GetTagFQName returns fqname based on a tag name and firewall rule
func (fr *FirewallRule) GetTagFQName(tagName string) ([]string, error) {
	if !strings.Contains(tagName, "=") {
		return nil, common.ErrorNotFoundf("invalid tag name '%s'", tagName)
	}

	if strings.HasPrefix(tagName, FirewallPolicyTagNameGlobalPrefix) {
		return []string{tagName[7:]}, nil
	}

	fqName := append([]string(nil), fr.GetFQName()...)
	if fr.GetParentType() == KindPolicyManagement {
		return append(fqName[:len(fqName)-2], tagName), nil
	}

	fqName = append([]string(nil), fr.GetFQName()...)
	if fr.GetParentType() == KindProject {
		return append(fqName[:len(fqName)-1], tagName), nil
	}

	return nil, common.ErrorBadRequestf(
		"Firewall rule %s (%s) parent type '%s' is not supported as security resource scope",
		fr.GetUUID(),
		fr.GetFQName(),
		fr.GetParentType(),
	)
}

// GetEndpoints returns request and database endpoints properties of firewall rules
func (fr *FirewallRule) GetEndpoints(
	databaseFR *FirewallRule,
) ([]*FirewallRuleEndpointType, []*FirewallRuleEndpointType) {
	endpoints := []*FirewallRuleEndpointType{
		fr.GetEndpoint1(),
		fr.GetEndpoint2(),
	}

	dbEndpoints := []*FirewallRuleEndpointType{
		databaseFR.GetEndpoint1(),
		databaseFR.GetEndpoint2(),
	}

	return endpoints, dbEndpoints
}

// FirewallRuleTagRef stub to be removed when tag refs implemented
type FirewallRuleTagRef struct {
	UUID string
	To   []string
}

// GetTagRefs to be removed when tag refs implemented
func (fr *FirewallRule) GetTagRefs() []*FirewallRuleTagRef {
	return []*FirewallRuleTagRef{}
}
