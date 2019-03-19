package main

import (
	"fmt"
	"log"
	"os"

	"github.com/pepabo/go-netapp/netapp"
)

type NetappAggr struct {
	Name                string
	SizeUsed            int
	SizeTotal           int
	SizeAvailable       int
	TotalReservedSpace  int
	PercentUsedCapacity string
	PhysicalUsedPercent int
}

func (f *Filer) GetAggrData() (r []NetappAggr) {
	opts := &netapp.AggrOptions{
		DesiredAttributes: &netapp.AggrInfo{
			AggrSpaceAttributes: &netapp.AggrSpaceAttributes{},
		},
	}

	l := f.getAggrList(opts)
	fmt.Printf("%+v\n", l[0])

	for _, n := range l {
		r = append(r, NetappAggr{
			Name:                n.AggregateName,
			SizeUsed:            n.AggrSpaceAttributes.SizeUsed,
			SizeTotal:           n.AggrSpaceAttributes.SizeTotal,
			SizeAvailable:       n.AggrSpaceAttributes.SizeAvailable,
			TotalReservedSpace:  n.AggrSpaceAttributes.TotalReservedSpace,
			PercentUsedCapacity: n.AggrSpaceAttributes.PercentUsedCapacity,
			PhysicalUsedPercent: n.AggrSpaceAttributes.PhysicalUsedPercent,
		})
	}
	return
}

func (f *Filer) getAggrList(opts *netapp.AggrOptions) (r []netapp.AggrInfo) {

	var pages []*netapp.AggrListResponse

	handler := func(r netapp.AggrListPagesResponse) bool {
		fmt.Println(r)
		if r.Error != nil {
			if os.Getenv("INFO") != "" {
				log.Printf("%s", r.Error)
			}
			return false
		}
		pages = append(pages, r.Response)
		return true
	}

	f.NetappClient.Aggregate.ListPages(opts, handler)

	for _, p := range pages {
		r = append(r, p.Results.AggrAttributes...)
	}

	return
}
