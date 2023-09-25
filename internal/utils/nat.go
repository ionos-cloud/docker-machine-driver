package utils

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/docker/machine/libmachine/log"
	"github.com/ionos-cloud/docker-machine-driver/internal/pointer"
	"github.com/ionos-cloud/docker-machine-driver/pkg/sdk_utils"
	sdkgo "github.com/ionos-cloud/sdk-go/v6"
	"golang.org/x/exp/maps"
)

func (c *Client) GetNats(datacenterId string) (*sdkgo.NatGateways, error) {
	nats, _, err := c.NATGatewaysApi.DatacentersNatgatewaysGet(c.ctx, datacenterId).Depth(1).Execute()
	if err != nil {
		return nil, sdk_utils.ShortenOpenApiErr(err)
	}
	return &nats, nil
}

func (c *Client) GetNat(datacenterId, natId string) (*sdkgo.NatGateway, error) {
	nat, _, err := c.NATGatewaysApi.DatacentersNatgatewaysFindByNatGatewayId(c.ctx, datacenterId, natId).Execute()
	if err != nil {
		return nil, sdk_utils.ShortenOpenApiErr(err)
	}
	return &nat, nil
}

// NatRuleMaker 's objective is to easily create many NatGatewayRules with similar properties but different open ports
type NatRuleMaker struct {
	rules             []sdkgo.NatGatewayRule
	defaultProperties sdkgo.NatGatewayRuleProperties
}

func NewNRM(publicIp, srcSubnet, targetSubnet *string) NatRuleMaker {
	return NatRuleMaker{
		rules: make([]sdkgo.NatGatewayRule, 0),
		defaultProperties: sdkgo.NatGatewayRuleProperties{
			Name:         pointer.From("Docker Machine NAT Rule"),
			Type:         pointer.From(sdkgo.NatGatewayRuleType("SNAT")),
			SourceSubnet: srcSubnet,
			TargetSubnet: targetSubnet,
			PublicIp:     publicIp,
		},
	}
}

func (nrm *NatRuleMaker) Make() *[]sdkgo.NatGatewayRule {
	return &nrm.rules
}

func (nrm *NatRuleMaker) OpenPort(protocol string, port int32) *NatRuleMaker {
	return nrm.OpenPorts(protocol, port, port)
}

func (nrm *NatRuleMaker) OpenPorts(protocol string, start int32, end int32) *NatRuleMaker {
	properties := nrm.defaultProperties
	properties.Protocol = (*sdkgo.NatGatewayRuleProtocol)(&protocol)
	nameIdentifier := fmt.Sprintf("%s: ", protocol)
	if protocol != "ALL" && protocol != "ICMP" {
		properties.TargetPortRange = &sdkgo.TargetPortRange{Start: &start, End: &end}
		nameIdentifier = fmt.Sprintf("%s (%d - %d): ", protocol, start, end)
	}
	properties.Name = pointer.From(nameIdentifier + *properties.Name)
	nrm.rules = append(nrm.rules, sdkgo.NatGatewayRule{
		Properties: &properties,
	})
	return nrm
}

func flowlogStringToModel(flowlog string) sdkgo.FlowLog {
	var split = strings.Split(flowlog, ":")

	return sdkgo.FlowLog{
		Properties: &sdkgo.FlowLogProperties{
			Name:      &split[0],
			Action:    &split[1],
			Direction: &split[2],
			Bucket:    &split[3],
		},
	}
}

func flowlogsStringToModel(flowlogs []string) *sdkgo.FlowLogs {
	if len(flowlogs) == 0 {
		return nil
	}
	flowlog_models := sdkgo.NewFlowLogs()
	flowlog_models.Items = &[]sdkgo.FlowLog{}

	for _, flowlog := range flowlogs {
		*flowlog_models.Items = append(*flowlog_models.Items, flowlogStringToModel(flowlog))
	}

	return flowlog_models
}

func natRuleStringToModel(rule, natPublicIp, defaultsourceSubnet string) (*sdkgo.NatGatewayRule, error) {
	var split = strings.Split(rule, ":")
	ruleType := sdkgo.NatGatewayRuleType(split[1])
	ruleProtocol := sdkgo.NatGatewayRuleProtocol(split[2])
	publicIp := split[3]
	if publicIp == "" {
		publicIp = natPublicIp
	}
	sourceSubnet := split[4]
	if sourceSubnet == "" {
		sourceSubnet = defaultsourceSubnet
	}

	ruleModel := sdkgo.NatGatewayRule{
		Properties: &sdkgo.NatGatewayRuleProperties{
			Name:            &split[0],
			Type:            &ruleType,
			Protocol:        &ruleProtocol,
			PublicIp:        &publicIp,
			SourceSubnet:    &sourceSubnet,
			TargetPortRange: &sdkgo.TargetPortRange{},
		},
	}

	targetSubnet := split[5]
	if targetSubnet == "" {
		ruleModel.Properties.TargetSubnet = nil
	} else {
		ruleModel.Properties.TargetSubnet = &targetSubnet
	}

	if split[6] == "" {
		ruleModel.Properties.TargetPortRange.Start = nil
	} else {

		start, err := strconv.Atoi(split[6])
		start32 := int32(start)

		if err != nil {
			return nil, err
		}
		ruleModel.Properties.TargetPortRange.Start = &start32
	}

	if split[7] == "" {
		ruleModel.Properties.TargetPortRange.End = nil
	} else {

		end, err := strconv.Atoi(split[7])
		end32 := int32(end)

		if err != nil {
			return nil, err
		}
		ruleModel.Properties.TargetPortRange.End = &end32
	}

	return &ruleModel, nil
}

func natRulesStringToModel(rules []string, natPublicIp, sourceSubnet string) (*sdkgo.NatGatewayRules, error) {
	if len(rules) == 0 {
		return nil, nil
	}
	rule_models := sdkgo.NewNatGatewayRules()
	rule_models.Items = &[]sdkgo.NatGatewayRule{}

	for _, rule := range rules {
		rule_model, err := natRuleStringToModel(rule, natPublicIp, sourceSubnet)

		if err != nil {
			return nil, err
		}

		*rule_models.Items = append(*rule_models.Items, *rule_model)
	}

	return rule_models, nil
}

func (c *Client) CreateNat(datacenterId, name string, publicIps, flowlogs, natRules []string, lansToGateways map[string][]string, sourceSubnet string, skipDefaultRules bool) (*sdkgo.NatGateway, error) {
	var lans []sdkgo.NatGatewayLanProperties
	publicIp := publicIps[0]

	err := c.createLansIfNotExist(datacenterId, maps.Keys(lansToGateways))
	if err != nil {
		return nil, err
	}
	time.Sleep(5 * time.Second)
	for lanId, gatewayIps := range lansToGateways {
		id, err := strconv.ParseInt(lanId, 10, 32)
		if err != nil {
			return nil, err
		}
		// Unpack the map into NatGatewayLanProperties objects. https://api.ionos.com/docs/cloud/v6/#tag/NAT-Gateways/operation/datacentersNatgatewaysPost
		var ptrGatewayIps *[]string = nil
		if len(gatewayIps) > 1 || gatewayIps[0] != "" {
			// We do this check so that we don't set the GatewayIps property if it's empty. If the property is empty, a gateway IP is generated by the API.
			ptrGatewayIps = &gatewayIps
		}
		lans = append(lans, sdkgo.NatGatewayLanProperties{Id: pointer.From(int32(id)), GatewayIps: ptrGatewayIps})
	}

	rules := &[]sdkgo.NatGatewayRule{}
	if !skipDefaultRules {
		nrm := NewNRM(&publicIp, &sourceSubnet, nil)
		nrm.
			OpenPort("TCP", 22).  // SSH
			OpenPort("UDP", 53).  // DNS
			OpenPort("TCP", 80).  // HTTP
			OpenPort("TCP", 179). // Calico BGP Port
			OpenPort("TCP", 443). //

			OpenPort("TCP", 2376). // Node driver Docker daemon TLS port
			OpenPort("UDP", 4789). // Flannel VXLAN overlay networking on Windows cluster
			OpenPort("TCP", 6443). // Rancher Webhook
			OpenPort("TCP", 6783). // Weave Port
			OpenPort("TCP", 8443). // Rancher webhook

			OpenPort("UDP", 8472). // Canal/Flannel VXLAN overlay networking
			OpenPort("TCP", 9099). // Canal/Flannel livenessProbe/readinessProbe
			OpenPort("TCP", 9100). // Default port required by Monitoring to scrape metrics from Linux node-exporters
			OpenPort("TCP", 9443). // Rancher webhook
			OpenPort("TCP", 9796). // Default port required by Monitoring to scrape metrics from Windows node-exporters

			OpenPort("TCP", 10254).         // Ingress controller livenessProbe/readinessProbe
			OpenPort("TCP", 10256).         //
			OpenPorts("TCP", 2379, 2380).   // etcd
			OpenPorts("UDP", 6783, 6784).   // Weave Port (UDP)
			OpenPorts("TCP", 10250, 10252). // Metrics server communication with all nodes API

			OpenPorts("TCP", 30000, 32767). //
			OpenPorts("UDP", 30000, 32767). //
			OpenPort("ALL", 0)              // Outbound
		default_rules := nrm.Make()

		*rules = append(*rules, *default_rules...)
	}

	new_rules, err := natRulesStringToModel(natRules, publicIps[0], sourceSubnet)

	if err != nil {
		return nil, err
	}

	if new_rules != nil {
		*rules = append(*rules, *new_rules.Items...)
	}

	nat, resp, err := c.NATGatewaysApi.DatacentersNatgatewaysPost(c.ctx, datacenterId).NatGateway(
		sdkgo.NatGateway{
			Properties: &sdkgo.NatGatewayProperties{
				Name:      &name,
				PublicIps: &publicIps,
				Lans:      &lans,
			},
			Entities: &sdkgo.NatGatewayEntities{
				Rules:    &sdkgo.NatGatewayRules{Items: rules},
				Flowlogs: flowlogsStringToModel(flowlogs),
			},
		},
	).Execute()
	if err != nil {
		return nil, err
	}

	err = c.waitTillProvisioned(resp.Header.Get("location"))
	return &nat, err
}

func (c *Client) createLansIfNotExist(datacenterId string, lanIds []string) error {
	for _, lanid := range lanIds {
		_, resp, err := c.LANsApi.DatacentersLansFindById(c.ctx, datacenterId, lanid).Execute()
		if resp.HttpNotFound() {
			// Before err check as 404s throw an err.
			log.Infof("Creating LAN %s for NAT\n", lanid)
			_, err := c.CreateLan(datacenterId, "Docker Machine LAN (NAT)", false)
			if err != nil {
				return err
			}
			continue // breakpoint
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) RemoveNat(datacenterId, natId string) error {
	_, err := c.NATGatewaysApi.DatacentersNatgatewaysDelete(c.ctx, datacenterId, natId).Execute()
	if err != nil {
		return sdk_utils.ShortenOpenApiErr(err)
	}
	return nil
}
