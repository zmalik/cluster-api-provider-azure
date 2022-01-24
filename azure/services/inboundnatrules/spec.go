/*
Copyright 2021 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package inboundnatrules

import (
	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2021-02-01/network"
	"github.com/Azure/go-autorest/autorest/to"
	"github.com/pkg/errors"
)

// InboundNatSpec defines the specification for an inbound NAT rule.
type InboundNatSpec struct {
	Name                      string
	LoadBalancerName          string
	ResourceGroup             string
	FrontendIPConfigurationID *string
	PortsInUse                map[int32]struct{}
}

// ResourceName returns the name of the inbound NAT rule.
func (s *InboundNatSpec) ResourceName() string {
	return s.Name
}

// ResourceGroupName returns the name of the resource group.
func (s *InboundNatSpec) ResourceGroupName() string {
	return s.ResourceGroup
}

// OwnerResourceName returns the name of the load balancer associated with an inbound NAT rule.
func (s *InboundNatSpec) OwnerResourceName() string {
	return s.LoadBalancerName
}

// Parameters returns the parameters for the inbound NAT rule.
func (s *InboundNatSpec) Parameters(existing interface{}) (parameters interface{}, err error) {
	if existing != nil {
		if _, ok := existing.(network.InboundNatRule); !ok {
			return nil, errors.Errorf("%T is not a network.InboundNatRule", existing)
		}

		return nil, nil
	}

	if s.FrontendIPConfigurationID == nil {
		return nil, errors.Errorf("FrontendIPConfigurationID is not set")
	}

	sshFrontendPort, err := getAvailablePort(s.PortsInUse)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to find available SSH Frontend port for NAT Rule %s in load balancer %s", s.ResourceName(), s.OwnerResourceName())
	}

	rule := network.InboundNatRule{
		Name: to.StringPtr(s.ResourceName()),
		InboundNatRulePropertiesFormat: &network.InboundNatRulePropertiesFormat{
			BackendPort:          to.Int32Ptr(22),
			EnableFloatingIP:     to.BoolPtr(false),
			IdleTimeoutInMinutes: to.Int32Ptr(4),
			FrontendIPConfiguration: &network.SubResource{
				ID: s.FrontendIPConfigurationID,
			},
			Protocol:     network.TransportProtocolTCP,
			FrontendPort: &sshFrontendPort,
		},
	}

	return rule, nil
}

func getAvailablePort(portsInUse map[int32]struct{}) (int32, error) {
	// NAT rules need to use a unique port. Since we need one NAT rule per control plane and we expect to have 1, 3, 5, maybe 9 control planes, there should never be more than 9 ports in use.
	// This is an artificial limit of 20 ports that we can pick from, which should be plenty enough (in reality we should never reach that limit).
	// These NAT rules are used for SSH purposes which is why we start at 22 and then use 2201, 2202, etc.
	var i int32 = 22
	if _, ok := portsInUse[22]; ok {
		for i = 2201; i < 2220; i++ {
			if _, ok := portsInUse[i]; !ok {
				// Found available port
				return i, nil
			}
		}
		return i, errors.Errorf("No available SSH Frontend ports")
	}

	return i, nil
}
