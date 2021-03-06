package fs

import (
	"testing"

	"github.com/dotcloud/docker/pkg/libcontainer/cgroups"
)

const (
	memoryStatContents = `cache 512
rss 1024`
	memoryUsageContents    = "2048\n"
	memoryMaxUsageContents = "4096\n"
)

func TestMemoryStats(t *testing.T) {
	helper := NewCgroupTestUtil("memory", t)
	defer helper.cleanup()
	helper.writeFileContents(map[string]string{
		"memory.stat":               memoryStatContents,
		"memory.usage_in_bytes":     memoryUsageContents,
		"memory.max_usage_in_bytes": memoryMaxUsageContents,
	})

	memory := &memoryGroup{}
	err := memory.GetStats(helper.CgroupData, &actualStats)
	if err != nil {
		t.Fatal(err)
	}
	expectedStats := cgroups.MemoryStats{Usage: 2048, MaxUsage: 4096, Stats: map[string]uint64{"cache": 512, "rss": 1024}}
	expectMemoryStatEquals(t, expectedStats, actualStats.MemoryStats)
}

func TestMemoryStatsNoStatFile(t *testing.T) {
	helper := NewCgroupTestUtil("memory", t)
	defer helper.cleanup()
	helper.writeFileContents(map[string]string{
		"memory.usage_in_bytes":     memoryUsageContents,
		"memory.max_usage_in_bytes": memoryMaxUsageContents,
	})

	memory := &memoryGroup{}
	err := memory.GetStats(helper.CgroupData, &actualStats)
	if err == nil {
		t.Fatal("Expected failure")
	}
}

func TestMemoryStatsNoUsageFile(t *testing.T) {
	helper := NewCgroupTestUtil("memory", t)
	defer helper.cleanup()
	helper.writeFileContents(map[string]string{
		"memory.stat":               memoryStatContents,
		"memory.max_usage_in_bytes": memoryMaxUsageContents,
	})

	memory := &memoryGroup{}
	err := memory.GetStats(helper.CgroupData, &actualStats)
	if err == nil {
		t.Fatal("Expected failure")
	}
}

func TestMemoryStatsNoMaxUsageFile(t *testing.T) {
	helper := NewCgroupTestUtil("memory", t)
	defer helper.cleanup()
	helper.writeFileContents(map[string]string{
		"memory.stat":           memoryStatContents,
		"memory.usage_in_bytes": memoryUsageContents,
	})

	memory := &memoryGroup{}
	err := memory.GetStats(helper.CgroupData, &actualStats)
	if err == nil {
		t.Fatal("Expected failure")
	}
}

func TestMemoryStatsBadStatFile(t *testing.T) {
	helper := NewCgroupTestUtil("memory", t)
	defer helper.cleanup()
	helper.writeFileContents(map[string]string{
		"memory.stat":               "rss rss",
		"memory.usage_in_bytes":     memoryUsageContents,
		"memory.max_usage_in_bytes": memoryMaxUsageContents,
	})

	memory := &memoryGroup{}
	err := memory.GetStats(helper.CgroupData, &actualStats)
	if err == nil {
		t.Fatal("Expected failure")
	}
}

func TestMemoryStatsBadUsageFile(t *testing.T) {
	helper := NewCgroupTestUtil("memory", t)
	defer helper.cleanup()
	helper.writeFileContents(map[string]string{
		"memory.stat":               memoryStatContents,
		"memory.usage_in_bytes":     "bad",
		"memory.max_usage_in_bytes": memoryMaxUsageContents,
	})

	memory := &memoryGroup{}
	err := memory.GetStats(helper.CgroupData, &actualStats)
	if err == nil {
		t.Fatal("Expected failure")
	}
}

func TestMemoryStatsBadMaxUsageFile(t *testing.T) {
	helper := NewCgroupTestUtil("memory", t)
	defer helper.cleanup()
	helper.writeFileContents(map[string]string{
		"memory.stat":               memoryStatContents,
		"memory.usage_in_bytes":     memoryUsageContents,
		"memory.max_usage_in_bytes": "bad",
	})

	memory := &memoryGroup{}
	err := memory.GetStats(helper.CgroupData, &actualStats)
	if err == nil {
		t.Fatal("Expected failure")
	}
}
