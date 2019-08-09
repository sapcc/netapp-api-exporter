package main

import (
	"github.com/pepabo/go-netapp/netapp"
)

type Aggregate struct {
	AvailabilityZone    string
	FilerName           string
	Name                string
	OwnerName           string
	SizeUsed            int
	SizeTotal           int
	SizeAvailable       int
	TotalReservedSpace  int
	PercentUsedCapacity string
	PhysicalUsed        int
	PhysicalUsedPercent int
}

func (f *Filer) GetNetappAggregate(r chan<- *Aggregate, done chan<- struct{}) {
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

	aggrs := f.getAggrList(opts)
	logger.Printf("%s: %d aggregates fetched", f.Host, len(aggrs))

	for _, n := range aggrs {
		r <- &Aggregate{
			FilerName:           f.Name,
			AvailabilityZone:    f.AvailabilityZone,
			Name:                n.AggregateName,
			OwnerName:           n.AggrOwnershipAttributes.OwnerName,
			SizeUsed:            n.AggrSpaceAttributes.SizeUsed,
			SizeTotal:           n.AggrSpaceAttributes.SizeTotal,
			SizeAvailable:       n.AggrSpaceAttributes.SizeAvailable,
			TotalReservedSpace:  n.AggrSpaceAttributes.TotalReservedSpace,
			PercentUsedCapacity: n.AggrSpaceAttributes.PercentUsedCapacity,
			PhysicalUsed:        n.AggrSpaceAttributes.PhysicalUsed,
			PhysicalUsedPercent: n.AggrSpaceAttributes.PhysicalUsedPercent,
		}
	}

	if len(aggrs) != 0 {
		done <- struct{}{}
	}
}

func (f *Filer) getAggrList(opts *netapp.AggrOptions) (res []netapp.AggrInfo) {
	pageHandler := func(r netapp.AggrListPagesResponse) bool {
		if r.Error != nil {
			logger.Warnf("%s", r.Error)
			return false
		}
		res = append(res, r.Response.Results.AggrAttributes...)
		return true
	}
	f.NetappClient.Aggregate.ListPages(opts, pageHandler)
	return
}
