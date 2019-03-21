package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNetappAggregate(t *testing.T) {
	host := os.Getenv("NETAPP_HOST")
	username := os.Getenv("NETAPP_USERNAME")
	password := os.Getenv("NETAPP_PASSWORD")
	region := os.Getenv("NETAPP_REGION")
	f := NewFiler("test", host, username, password, region)

	l := f.GetAggrData()
	n := l[2]

	fmt.Println("Aggregate Name:\t\t", n.Name)
	fmt.Println("Size Used:\t\t", n.SizeUsed)
	fmt.Println("Size Total:\t\t", n.SizeTotal)
	fmt.Println("Size Available:\t\t", n.SizeAvailable)
	fmt.Println("Size Used Percentage:\t", n.PercentUsedCapacity)
	fmt.Println("Physical Used Percent:\t", n.PhysicalUsedPercent)

	assert.Equal(t, "", n)
}
