package bench

import (
	"runtime"

	"os/exec"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
)

type Info struct {
	CPUInfo
	SystemInfo
	MemoryInfo
}

type CPUInfo struct {
	Model string
	Cores int32
}

type SystemInfo struct {
	Hostname        string
	OS              string
	Platform        string
	PlatformVersion string
	KernelVersion   string
	Arch            string
	Virt            string
}

type MemoryInfo struct {
	Total     uint64
	Used      uint64
	TotalSwap uint64
	UsedSwap  uint64
}

func GetInfo() (info *Info, err error) {
	info = new(Info)

	//CPU
	ci, err := cpu.Info()
	if err != nil {
		return
	}
	info.Model = ci[0].ModelName
	info.Cores = ci[0].Cores

	//Memory info
	vmi, err := mem.VirtualMemory()
	if err != nil {
		return
	}
	info.Total = vmi.Total
	info.Used = vmi.Used
	smi, err := mem.SwapMemory()
	if err != nil {
		return
	}
	info.TotalSwap = smi.Total
	info.UsedSwap = smi.Used

	info.Arch = runtime.GOARCH

	//System
	si, err := host.Info()
	if err != nil {
		return
	}
	info.Hostname = si.Hostname
	info.OS = si.OS
	info.Platform = si.Platform
	info.PlatformVersion = si.PlatformVersion
	info.KernelVersion = si.KernelVersion

	//Virt
	virtB, err := exec.Command("virt-what").Output()
	if err != nil {
		info.Virt = "error"
		return
	}

	info.Virt = string(virtB[:len(virtB)-1])
	return
}
