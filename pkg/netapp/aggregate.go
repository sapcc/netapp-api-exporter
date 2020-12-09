package netapp

import (
	"strconv"

	n "github.com/pepabo/go-netapp/netapp"
)

type Aggregate struct {
	Name                string
	OwnerName           string
	SizeUsed            float64
	SizeTotal           float64
	SizeAvailable       float64
	TotalReservedSpace  float64
	PercentUsedCapacity float64
	PhysicalUsed        float64
	PhysicalUsedPercent float64
	IsEncrypted         bool
}

func (c *Client) ListAggregates() (aggregates []*Aggregate, err error) {
	aggrInfos, err := c.listAggregates()
	if err != nil {
		return nil, err
	}
	for _, aggr := range aggrInfos {
		aggregates = append(aggregates, parseAggregate(aggr))
	}
	return
}

func (c *Client) listAggregates() (res []n.AggrInfo, err error) {
	opts := newAggrOpts(false)
	pageHandler := func(r n.AggrListPagesResponse) bool {
		if r.Error != nil {
			err = r.Error
			return false
		}
		res = append(res, r.Response.Results.AggrAttributes...)
		return true
	}
	c.Aggregate.ListPages(opts, pageHandler)
	return
}

func newAggrOpts(isRootAggregate bool) *n.AggrOptions {
	return &n.AggrOptions{
		Query: &n.AggrInfo{
			AggrRaidAttributes: &n.AggrRaidAttributes{
				IsRootAggregate: &isRootAggregate,
			},
		},
		DesiredAttributes: &n.AggrInfo{
			AggrRaidAttributes:      &n.AggrRaidAttributes{},
			AggrOwnershipAttributes: &n.AggrOwnershipAttributes{},
			AggrSpaceAttributes:     &n.AggrSpaceAttributes{},
		},
	}
}

func parseAggregate(aggrInfo n.AggrInfo) *Aggregate {
	var isEncrypted bool
	percentUsedCapacity, _ := strconv.ParseFloat(aggrInfo.AggrSpaceAttributes.PercentUsedCapacity, 64)
	if aggrInfo.AggrRaidAttributes.IsEncrypted != nil && *aggrInfo.AggrRaidAttributes.IsEncrypted {
		isEncrypted = true
	}
	return &Aggregate{
		Name:                aggrInfo.AggregateName,
		OwnerName:           aggrInfo.AggrOwnershipAttributes.OwnerName,
		SizeUsed:            float64(aggrInfo.AggrSpaceAttributes.SizeUsed),
		SizeTotal:           float64(aggrInfo.AggrSpaceAttributes.SizeTotal),
		SizeAvailable:       float64(aggrInfo.AggrSpaceAttributes.SizeAvailable),
		TotalReservedSpace:  float64(aggrInfo.AggrSpaceAttributes.TotalReservedSpace),
		PercentUsedCapacity: percentUsedCapacity,
		PhysicalUsed:        float64(aggrInfo.AggrSpaceAttributes.PhysicalUsed),
		PhysicalUsedPercent: float64(aggrInfo.AggrSpaceAttributes.PhysicalUsedPercent),
		IsEncrypted:         isEncrypted,
	}
}
