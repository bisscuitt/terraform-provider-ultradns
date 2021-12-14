package zone

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ultradns/ultradns-go-sdk/pkg/zone"
)

func flattenTsig(t *zone.Tsig) *schema.Set {
	set := &schema.Set{F: zeroIndexHash}
	tsig := make(map[string]interface{})

	tsig["tsig_key_name"] = t.TsigKeyName
	tsig["tsig_key_value"] = t.TsigKeyValue
	tsig["tsig_algorithm"] = t.TsigAlgorithm
	tsig["description"] = t.Description

	set.Add(tsig)

	return set
}

func flattenRestrictIP(ri []*zone.RestrictIP) *schema.Set {
	set := &schema.Set{F: schema.HashResource(restrictIPResource())}

	for _, restrictIPData := range ri {
		restrictIP := make(map[string]interface{})

		restrictIP["start_ip"] = restrictIPData.StartIP
		restrictIP["end_ip"] = restrictIPData.EndIP
		restrictIP["cidr"] = restrictIPData.Cidr
		restrictIP["single_ip"] = restrictIPData.SingleIP
		restrictIP["comment"] = restrictIPData.Comment

		set.Add(restrictIP)
	}

	return set
}

func flattenNotifyAddresses(na []*zone.NotifyAddress) *schema.Set {
	set := &schema.Set{F: schema.HashResource(notifyAddressResource())}

	for _, notifyAddressData := range na {
		notifyAddress := make(map[string]interface{})

		notifyAddress["notify_address"] = notifyAddressData.NotifyAddress
		notifyAddress["description"] = notifyAddressData.Description

		set.Add(notifyAddress)
	}

	return set
}

func flattenRegistrarInfo(ri *zone.RegistrarInfo) *schema.Set {
	set := &schema.Set{F: zeroIndexHash}

	registrarInfo := make(map[string]interface{})

	registrarInfo["registrar"] = ri.Registrar
	registrarInfo["who_is_expiration"] = ri.WhoIsExpiration
	registrarInfo["name_servers"] = flattenRegistrarInfoNameServer(ri.NameServers)

	set.Add(registrarInfo)

	return set
}

func flattenRegistrarInfoNameServer(nsl *zone.NameServersList) *schema.Set {
	set := &schema.Set{F: zeroIndexHash}

	RegistrarInfoNameServersList := make(map[string]interface{})

	RegistrarInfoNameServersList["ok"] = nsl.Ok
	RegistrarInfoNameServersList["unknown"] = nsl.Unknown
	RegistrarInfoNameServersList["missing"] = nsl.Missing
	RegistrarInfoNameServersList["incorrect"] = nsl.Incorrect

	set.Add(RegistrarInfoNameServersList)

	return set
}

func flattenNameServer(ns *zone.NameServer) *schema.Set {
	set := &schema.Set{F: zeroIndexHash}

	nameServer := make(map[string]interface{})

	nameServer["ip"] = ns.IP
	nameServer["tsig_key"] = ns.TsigKey
	nameServer["tsig_key_value"] = ns.TsigKeyValue
	nameServer["tsig_algorithm"] = ns.TsigAlgorithm

	set.Add(nameServer)

	return set
}

func flattenTransferStatusDetails(tsd *zone.TransferStatusDetails) *schema.Set {
	set := &schema.Set{F: zeroIndexHash}
	transferDetails := make(map[string]interface{})

	transferDetails["last_refresh"] = tsd.LastRefresh
	transferDetails["next_refresh"] = tsd.NextRefresh
	transferDetails["last_refresh_status"] = tsd.LastRefreshStatus
	transferDetails["last_refresh_status_message"] = tsd.LastRefreshStatusMessage

	set.Add(transferDetails)

	return set
}

func zeroIndexHash(i interface{}) int {
	return 0
}
