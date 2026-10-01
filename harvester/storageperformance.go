package harvester

import (
	"fmt"
	"math"

	"github.com/harvester/harvester/pkg/builder"
	kubevirtv1 "kubevirt.io/api/core/v1"
)

var (
	validDiskCacheModes = map[kubevirtv1.DriverCache]bool{
		kubevirtv1.CacheNone:         true,
		kubevirtv1.CacheWriteThrough: true,
		kubevirtv1.CacheWriteBack:    true,
	}
	validDiskIOModes = map[kubevirtv1.DriverIO]bool{
		kubevirtv1.IONative:  true,
		kubevirtv1.IOThreads: true,
	}
	validIOThreadsPolicies = map[kubevirtv1.IOThreadsPolicy]bool{
		kubevirtv1.IOThreadsPolicyShared:           true,
		kubevirtv1.IOThreadsPolicyAuto:             true,
		kubevirtv1.IOThreadsPolicySupplementalPool: true,
	}
)

func (disk *Disk) hasPerformanceOptions() bool {
	return disk.Cache != "" || disk.IO != "" || disk.DedicatedIOThread
}

func checkDiskPerformance(disk *Disk) error {
	if !disk.hasPerformanceOptions() {
		return nil
	}
	if disk.Type == builder.DiskTypeCDRom {
		return fmt.Errorf("cache, io and dedicatedIOThread are not supported on %s disks", builder.DiskTypeCDRom)
	}
	if disk.Cache != "" && !validDiskCacheModes[kubevirtv1.DriverCache(disk.Cache)] {
		return fmt.Errorf("invalid disk cache mode %q, must be one of none, writethrough, writeback", disk.Cache)
	}
	if disk.IO != "" && !validDiskIOModes[kubevirtv1.DriverIO(disk.IO)] {
		return fmt.Errorf("invalid disk io mode %q, must be one of native, threads", disk.IO)
	}
	// Native AIO needs O_DIRECT, which libvirt only uses for an uncached disk
	if disk.IO == string(kubevirtv1.IONative) && disk.Cache != string(kubevirtv1.CacheNone) {
		return fmt.Errorf("disk io mode native requires cache mode none")
	}
	return nil
}

func checkIOThreads(policy string, count int) error {
	if policy != "" && !validIOThreadsPolicies[kubevirtv1.IOThreadsPolicy(policy)] {
		return fmt.Errorf("invalid io threads policy %q, must be one of shared, auto, supplementalPool", policy)
	}
	if count < 0 || count > math.MaxUint32 {
		return fmt.Errorf("io thread count must be between 0 and %d", uint32(math.MaxUint32))
	}
	if count > 0 && policy != string(kubevirtv1.IOThreadsPolicySupplementalPool) {
		return fmt.Errorf("io thread count requires io threads policy supplementalPool")
	}
	return nil
}

// ConfigureStoragePerformance applies the KubeVirt high-performance disk options, which the VM builder does not support
func (d *Driver) ConfigureStoragePerformance(vm *kubevirtv1.VirtualMachine) {
	domain := &vm.Spec.Template.Spec.Domain

	if d.BlockMultiQueue {
		domain.Devices.BlockMultiQueue = new(true)
	}
	if d.IOThreadsPolicy != "" {
		domain.IOThreadsPolicy = new(kubevirtv1.IOThreadsPolicy(d.IOThreadsPolicy))
		if d.IOThreadCount > 0 && d.IOThreadCount <= math.MaxUint32 {
			domain.IOThreads = &kubevirtv1.DiskIOThreads{SupplementalPoolThreadCount: new(uint32(d.IOThreadCount))}
		}
	}

	if d.DiskInfo == nil {
		return
	}
	for i, disk := range d.DiskInfo.Disks {
		if !disk.hasPerformanceOptions() {
			continue
		}
		diskName := fmt.Sprintf("%s-%d", diskNamePrefix, i)
		for j := range domain.Devices.Disks {
			if domain.Devices.Disks[j].Name != diskName {
				continue
			}
			domain.Devices.Disks[j].Cache = kubevirtv1.DriverCache(disk.Cache)
			domain.Devices.Disks[j].IO = kubevirtv1.DriverIO(disk.IO)
			if disk.DedicatedIOThread {
				domain.Devices.Disks[j].DedicatedIOThread = new(true)
			}
		}
	}
}
