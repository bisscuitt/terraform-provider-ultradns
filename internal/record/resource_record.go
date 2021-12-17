package record

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ultradns/terraform-provider-ultradns/internal/service"
	"github.com/ultradns/ultradns-go-sdk/pkg/rrset"
)

func ResourceRecord() *schema.Resource {
	return &schema.Resource{

		CreateContext: resourceRecordCreate,
		ReadContext:   resourceRecordRead,
		UpdateContext: resourceRecordUpdate,
		DeleteContext: resourceRecordDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: resourceRecordSchema(),
	}
}

func resourceRecordCreate(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	services := meta.(*service.Service)
	rrSetData := newRRSet(rd)
	rrSetKeyData := newRRSetKey(rd)

	_, err := services.RecordService.CreateRecord(rrSetKeyData, rrSetData)

	if err != nil {
		return diag.FromErr(err)
	}

	rd.SetId(rrSetKeyData.ID())

	return resourceRecordRead(ctx, rd, meta)
}

func resourceRecordRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	services := meta.(*service.Service)
	rrSetKeyID := rd.Id()
	rrSetKey := getRRSetKey(rrSetKeyID)
	_, resList, err := services.RecordService.ReadRecord(rrSetKey)

	if resList != nil && resList.ResultInfo != nil && resList.ResultInfo.ReturnedCount > 0 && len(resList.RRSets) > 0 {
		if err := rd.Set("zone_name", resList.ZoneName); err != nil {
			return diag.FromErr(err)
		}

		ownerNameBefore := rd.Get("owner_name").(string)
		ownerNameAfter := resList.RRSets[0].OwnerName

		if ownerNameBefore != ownerNameAfter && ownerNameBefore+"."+resList.ZoneName != ownerNameAfter {
			if err := rd.Set("owner_name", ownerNameAfter); err != nil {
				return diag.FromErr(err)
			}
		}

		recordTypeBefore := rd.Get("record_type").(string)
		recordTypeAfter := resList.RRSets[0].RRType

		if rrset.GetRRTypeFullString(recordTypeBefore) != recordTypeAfter {
			if err := rd.Set("record_type", getRecordTypeString(resList.RRSets[0].RRType)); err != nil {
				return diag.FromErr(err)
			}
		}

		if err := rd.Set("ttl", resList.RRSets[0].TTL); err != nil {
			return diag.FromErr(err)
		}

		set := &schema.Set{F: schema.HashString}

		for _, val := range resList.RRSets[0].RData {
			set.Add(val)
		}

		if err := rd.Set("record_data", set); err != nil {
			return diag.FromErr(err)
		}
	}

	if err != nil {
		rd.SetId("")

		return nil
	}

	return diags
}

func resourceRecordUpdate(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	services := meta.(*service.Service)
	rrSetData := newRRSet(rd)
	rrSetKeyData := getRRSetKey(rd.Id())

	_, err := services.RecordService.UpdateRecord(rrSetKeyData, rrSetData)

	if err != nil {
		return diag.FromErr(err)
	}

	return resourceRecordRead(ctx, rd, meta)
}

func resourceRecordDelete(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	services := meta.(*service.Service)
	rrSetKeyData := getRRSetKey(rd.Id())

	_, err := services.RecordService.DeleteRecord(rrSetKeyData)

	if err != nil {
		rd.SetId("")

		return diag.FromErr(err)
	}

	rd.SetId("")

	return diags
}

func newRRSet(rd *schema.ResourceData) *rrset.RRSet {
	rrSetData := &rrset.RRSet{}

	if val, ok := rd.GetOk("owner_name"); ok {
		rrSetData.OwnerName = val.(string)
	}

	if val, ok := rd.GetOk("record_type"); ok {
		rrSetData.RRType = val.(string)
	}

	if val, ok := rd.GetOk("ttl"); ok {
		rrSetData.TTL = val.(int)
	}

	if val, ok := rd.GetOk("record_data"); ok {
		recordData := val.(*schema.Set).List()
		rrSetData.RData = make([]string, len(recordData))

		for i, record := range recordData {
			rrSetData.RData[i] = record.(string)
		}
	}

	return rrSetData
}

func newRRSetKey(rd *schema.ResourceData) *rrset.RRSetKey {
	rrSetKeyData := &rrset.RRSetKey{}

	if val, ok := rd.GetOk("zone_name"); ok {
		rrSetKeyData.Zone = val.(string)
	}

	if val, ok := rd.GetOk("owner_name"); ok {
		rrSetKeyData.Name = val.(string)
	}

	if val, ok := rd.GetOk("record_type"); ok {
		rrSetKeyData.Type = val.(string)
	}

	return rrSetKeyData
}

func getRRSetKey(id string) *rrset.RRSetKey {
	rrSetKeyData := &rrset.RRSetKey{}
	splitStringData := strings.Split(id, ":")

	if len(splitStringData) == 3 {
		rrSetKeyData.Name = splitStringData[0]
		rrSetKeyData.Zone = splitStringData[1]
		rrSetKeyData.Type = getRecordTypeString(splitStringData[2])
	}

	return rrSetKeyData
}
