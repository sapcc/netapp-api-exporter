package main

import (
	"github.com/pepabo/go-netapp/netapp"
	"strconv"
)

type NetappAggregate struct {
	AvailabilityZone    string
	FilerName           string
	Name                string
	OwnerName           string
	SizeUsed            float64
	SizeTotal           float64
	SizeAvailable       float64
	TotalReservedSpace  float64
	PercentUsedCapacity float64
	PhysicalUsed        float64
	PhysicalUsedPercent float64
}

func (f *FilerManager) GetNetappAggregate() (aggregates []*NetappAggregate, err error) {
	ff := new(bool)
	*ff = false
	opts := &netapp.AggrOptions{
		Query: &netapp.AggrInfo{
			AggrRaidAttributes: &netapp.AggrRaidAttributes{
				IsRootAggregate: ff,
			},
		},
		DesiredAttributes: &netapp.AggrInfo{
			AggrOwnershipAttributes: &netapp.AggrOwnershipAttributes{},
			AggrSpaceAttributes:     &netapp.AggrSpaceAttributes{},
		},
	}

	aggs, err := f.getAggrList(opts)

	if err == nil {
		logger.Printf("%s: %d aggregates fetched", f.Host, len(aggs))

		for _, n := range aggs {
			percentUsedCapacity, _ := strconv.ParseFloat(n.AggrSpaceAttributes.PercentUsedCapacity, 64)
			aggregates = append(aggregates, &NetappAggregate{
				AvailabilityZone:    f.AvailabilityZone,
				FilerName:           f.Name,
				Name:                n.AggregateName,
				OwnerName:           n.AggrOwnershipAttributes.OwnerName,
				SizeUsed:            float64(n.AggrSpaceAttributes.SizeUsed),
				SizeTotal:           float64(n.AggrSpaceAttributes.SizeTotal),
				SizeAvailable:       float64(n.AggrSpaceAttributes.SizeAvailable),
				TotalReservedSpace:  float64(n.AggrSpaceAttributes.TotalReservedSpace),
				PercentUsedCapacity: percentUsedCapacity,
				PhysicalUsed:        float64(n.AggrSpaceAttributes.PhysicalUsed),
				PhysicalUsedPercent: float64(n.AggrSpaceAttributes.PhysicalUsedPercent),
			})
		}
	}
	return
}

func (f *FilerManager) getAggrList(opts *netapp.AggrOptions) (res []netapp.AggrInfo, err error) {
	pageHandler := func(r netapp.AggrListPagesResponse) bool {
		if r.Error != nil {
			err = r.Error
			return false
		}
		res = append(res, r.Response.Results.AggrAttributes...)
		return true
	}
	f.NetappClient.Aggregate.ListPages(opts, pageHandler)
	return
}
