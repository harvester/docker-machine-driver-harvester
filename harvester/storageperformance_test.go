package harvester

import (
	"testing"

	"github.com/harvester/harvester/pkg/builder"
	"github.com/stretchr/testify/require"
	kubevirtv1 "kubevirt.io/api/core/v1"
)

func Test_checkDiskPerformance(t *testing.T) {
	testcases := []struct {
		description string
		disk        Disk
		expectError bool
	}{
		{"no options", Disk{Type: builder.DiskTypeDisk}, false},
		{"high performance", Disk{Type: builder.DiskTypeDisk, Cache: "none", IO: "native", DedicatedIOThread: true}, false},
		{"writeback with threads", Disk{Type: builder.DiskTypeDisk, Cache: "writeback", IO: "threads"}, false},
		{"cd-rom without options", Disk{Type: builder.DiskTypeCDRom}, false},
		{"cd-rom with options", Disk{Type: builder.DiskTypeCDRom, Cache: "none"}, true},
		{"invalid cache", Disk{Type: builder.DiskTypeDisk, Cache: "unsafe"}, true},
		{"invalid io", Disk{Type: builder.DiskTypeDisk, IO: "io_uring"}, true},
		{"native without cache none", Disk{Type: builder.DiskTypeDisk, Cache: "writeback", IO: "native"}, true},
		{"native without cache", Disk{Type: builder.DiskTypeDisk, IO: "native"}, true},
	}

	assert := require.New(t)

	for _, tc := range testcases {
		err := checkDiskPerformance(&tc.disk)
		if tc.expectError {
			assert.Error(err, tc.description)
		} else {
			assert.NoError(err, tc.description)
		}
	}
}

func Test_checkIOThreads(t *testing.T) {
	testcases := []struct {
		description string
		policy      string
		count       int
		expectError bool
	}{
		{"unset", "", 0, false},
		{"shared", "shared", 0, false},
		{"auto", "auto", 0, false},
		{"supplemental pool with count", "supplementalPool", 4, false},
		{"invalid policy", "dedicated", 0, true},
		{"negative count", "supplementalPool", -1, true},
		{"count without supplemental pool", "auto", 2, true},
		{"count without policy", "", 2, true},
	}

	assert := require.New(t)

	for _, tc := range testcases {
		err := checkIOThreads(tc.policy, tc.count)
		if tc.expectError {
			assert.Error(err, tc.description)
		} else {
			assert.NoError(err, tc.description)
		}
	}
}

func newTestVM(diskNames ...string) *kubevirtv1.VirtualMachine {
	vm := &kubevirtv1.VirtualMachine{Spec: kubevirtv1.VirtualMachineSpec{Template: &kubevirtv1.VirtualMachineInstanceTemplateSpec{}}}
	for _, name := range diskNames {
		vm.Spec.Template.Spec.Domain.Devices.Disks = append(vm.Spec.Template.Spec.Domain.Devices.Disks, kubevirtv1.Disk{Name: name})
	}
	return vm
}

func Test_ConfigureStoragePerformance(t *testing.T) {
	assert := require.New(t)

	d := &Driver{
		BlockMultiQueue: true,
		IOThreadsPolicy: "supplementalPool",
		IOThreadCount:   4,
		DiskInfo: &DiskInfo{Disks: []Disk{
			{Type: builder.DiskTypeDisk},
			{Type: builder.DiskTypeDisk, Cache: "none", IO: "native", DedicatedIOThread: true},
		}},
	}
	vm := newTestVM("disk-0", "disk-1", "cloudinitdisk")
	d.ConfigureStoragePerformance(vm)

	domain := vm.Spec.Template.Spec.Domain
	assert.Equal(new(true), domain.Devices.BlockMultiQueue)
	assert.Equal(new(kubevirtv1.IOThreadsPolicySupplementalPool), domain.IOThreadsPolicy)
	assert.Equal(&kubevirtv1.DiskIOThreads{SupplementalPoolThreadCount: new(uint32(4))}, domain.IOThreads)

	assert.Equal(kubevirtv1.Disk{Name: "disk-0"}, domain.Devices.Disks[0], "disk without options must be untouched")
	assert.Equal(kubevirtv1.CacheNone, domain.Devices.Disks[1].Cache)
	assert.Equal(kubevirtv1.IONative, domain.Devices.Disks[1].IO)
	assert.Equal(new(true), domain.Devices.Disks[1].DedicatedIOThread)
	assert.Equal(kubevirtv1.Disk{Name: "cloudinitdisk"}, domain.Devices.Disks[2])
}

func Test_ConfigureStoragePerformance_Unset(t *testing.T) {
	assert := require.New(t)

	d := &Driver{}
	vm := newTestVM("disk-1")
	d.ConfigureStoragePerformance(vm)

	assert.Equal(newTestVM("disk-1"), vm, "VM must be unchanged when no options are set")
}
