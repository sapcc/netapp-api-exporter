package main

import (
	"testing"

	_ "github.com/motemen/go-loghttp/global"
	"github.com/stretchr/testify/assert"
)

func TestParseVolumeComment(t *testing.T) {
	var str, id, name, project string

	str = "share_id: 193b4209-2ef0-4752-a262-261b9fa27b25 in project: 631a3518e93d436fbdf57525babe8606"
	id, name, project, _ = parseVolumeComment(str)
	assert.Equal(t, "193b4209-2ef0-4752-a262-261b9fa27b25", id)
	assert.Equal(t, "", name)
	assert.Equal(t, "631a3518e93d436fbdf57525babe8606", project)

	str = "share_id: 69fe1228-360c-4063-8f29-3a5bfb6d9772, share_name: c_blackbox_1553028005, project: d940aae3f8084f15a9b67de5b3b39720"
	id, name, project, _ = parseVolumeComment(str)
	assert.Equal(t, "69fe1228-360c-4063-8f29-3a5bfb6d9772", id)
	assert.Equal(t, "c_blackbox_1553028005", name)
	assert.Equal(t, "d940aae3f8084f15a9b67de5b3b39720", project)
}
