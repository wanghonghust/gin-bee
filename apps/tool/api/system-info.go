package api

import (
	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	"net/http"
	"time"
)

var CSystem = SystemController{}

type SysInfo struct {
	Host    map[string]any   `json:"host"`
	CpuInfo map[string]any   `json:"cpuInfo"`
	MemInfo map[string]any   `json:"memInfo"`
	Disk    []map[string]any `json:"disk"`
}

type SystemController struct {
}
type Disk struct {
	Info  disk.PartitionStat
	Usage *disk.UsageStat
}

// List
// @Summary
// @Schemes
// @Description 获取系统信息
// @Tags
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Success 200 {object} response.SystemInfoRes
// @Failure 400 {object} response.Response
// @Router /api/tool/system-info [get]
func (s *SystemController) List(c *gin.Context) {
	var sysInfo SysInfo
	cpuInfos, err := cpu.Info()
	if err != nil {
		c.JSONP(http.StatusBadRequest, gin.H{"msg": "查询系统信息失败"})
		return
	}
	percent, err := cpu.Percent(time.Duration(0), false)
	if err != nil {
		c.JSONP(http.StatusBadRequest, gin.H{"msg": "查询CPU使用率失败"})
		return
	}
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		c.JSONP(http.StatusBadRequest, gin.H{"msg": "查询内存信息失败"})
		return
	}
	hInfo, err := host.Info()
	if err != nil {
		c.JSONP(http.StatusBadRequest, gin.H{"msg": "查询主机信息失败"})
		return
	}
	parts, err := disk.Partitions(true)
	if err != nil {
		c.JSONP(http.StatusBadRequest, gin.H{"msg": "查询磁盘信息失败"})
		return
	}
	sysInfo.Host = make(map[string]any)
	sysInfo.Host = map[string]any{
		"hostname":        hInfo.Hostname,
		"os":              hInfo.OS,
		"platform":        hInfo.Platform,
		"kernelArch":      hInfo.KernelArch,
		"platformVersion": hInfo.PlatformVersion,
	}
	sysInfo.CpuInfo = make(map[string]any)
	sysInfo.CpuInfo = map[string]any{
		"modelName":  cpuInfos[0].ModelName,
		"cores":      cpuInfos[0].Cores,
		"cpuPercent": percent[0],
	}
	sysInfo.MemInfo = make(map[string]any)
	sysInfo.MemInfo = map[string]any{
		"usedPercent": memInfo.UsedPercent,
	}

	sysInfo.Disk = make([]map[string]any, 0)
	for _, part := range parts {
		var tmp map[string]any
		tmp = make(map[string]any)
		diskInfo, _ := disk.Usage(part.Mountpoint)
		tmp["device"] = part.Device
		tmp["total"] = diskInfo.Total
		tmp["fstype"] = part.Fstype
		tmp["free"] = diskInfo.Free
		tmp["used"] = diskInfo.Used
		tmp["usedPercent"] = diskInfo.UsedPercent
		sysInfo.Disk = append(sysInfo.Disk, tmp)
	}

	c.JSONP(http.StatusOK, gin.H{"data": sysInfo})
}
