package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"golang.org/x/sys/unix"
)

func main() {
	calculateCPUUsage()
	calculateMemUsage()
	calculateStorageUsage()
}

func cpu() (idleAllTime float64, nonIdleAllTime float64, totalTime float64) {
	var cpuRawData string
	file, err := ioutil.ReadFile("/proc/stat")
	if err != nil {
		fmt.Println(err)
	}
	lines := strings.Split(string(file), "\n")
	// for i, lineData := range lines {
	// 	if i == 0 {
	// 		cpuRawData = lineData
	// 		continue
	// 	}
	// }

	cpuRawData = lines[0]
	cpuData := strings.Fields(cpuRawData)
	// fmt.Println(cpuData, cpuRawData)
	// fmt.Println(cpuData)
	user, err := strconv.ParseFloat(cpuData[1], 32)
	nice, err := strconv.ParseFloat(cpuData[2], 32)
	system, err := strconv.ParseFloat(cpuData[3], 32)
	idle, err := strconv.ParseFloat(cpuData[4], 32)
	iowait, err := strconv.ParseFloat(cpuData[5], 32)
	irq, err := strconv.ParseFloat(cpuData[6], 32)
	softirq, err := strconv.ParseFloat(cpuData[7], 32)
	steal, err := strconv.ParseFloat(cpuData[8], 32)
	guest, err := strconv.ParseFloat(cpuData[9], 32)
	guest_nice, err := strconv.ParseFloat(cpuData[10], 32)

	// total := user + nice + system + idle + iowait + irq + softirq + steal + guest + guest_nice
	// idleFloat := (1 - (idle / total)) * 100
	// // idlePercent := truncate((1-(idle/total))*100, 0.01)
	// // idlePercent := (1 - (idle / total)) * 100
	// idlePercent := strconv.FormatFloat(idleFloat, 'f', 3, 64)
	// // fmt.Printf("%.2f", idlePercent)
	// fmt.Println(idlePercent)
	// fmt.Println(idleFloat)

	allUserTime := user + guest
	allNiceTime := nice + guest_nice
	idleAllTime = idle + iowait
	systemAllTime := system + irq + softirq
	nonIdleAllTime = allUserTime + allNiceTime + systemAllTime + steal
	totalTime = idleAllTime + nonIdleAllTime

	return idleAllTime, nonIdleAllTime, totalTime
}

func calculateCPUUsage() {
	prevIdleTime, _, prevTotal := cpu()
	// fmt.Println(prevIdleTime, prevTotal)
	time.Sleep(1 * time.Second)
	idleTime, _, total := cpu()
	// fmt.Println(idleTime, total)
	totalCalc := total - prevTotal
	idleCalc := idleTime - prevIdleTime

	usageCalc := (totalCalc - idleCalc) / totalCalc
	fmt.Println(strconv.FormatFloat(usageCalc, 'f', 2, 64))
}

func memory() (memTotal float64, memFree float64) {
	var memoryRawData []string
	file, err := ioutil.ReadFile("/proc/meminfo")
	if err != nil {
		fmt.Println(err)
	}
	lines := strings.Split(string(file), "\n")

	for i, lineData := range lines {
		re := regexp.MustCompile(`\W*(kB)\W*|:`)
		removed := re.ReplaceAllString(lineData, "")
		slicer := strings.Fields(removed)
		data := slicer[1]
		// fmt.Println(slicer)
		memoryRawData = append(memoryRawData, data)
		if i == 49 {
			break
		}
	}
	memTotal, _ = strconv.ParseFloat(memoryRawData[0], 32)
	memFree, _ = strconv.ParseFloat(memoryRawData[2], 32)

	return memTotal, memFree
	// fmt.Println(memoryRawData)
}

func calculateMemUsage() {
	memTotal, memFree := memory()
	memUsage := ((memTotal - memFree) / memTotal) * 100
	fmt.Println(strconv.FormatFloat(memUsage, 'f', 2, 64))
}

func storage() (storageFree float64, storageTotal float64) {
	var stat unix.Statfs_t

	wd, err := os.Getwd()

	if err != nil {
		return
	}

	unix.Statfs(wd, &stat)

	storageFree = float64(stat.Bfree * uint64(stat.Bsize))
	storageTotal = float64(stat.Blocks * uint64(stat.Bsize))
	// fmt.Println(storageFree, storageTotal)
	return storageFree, storageTotal

}

func calculateStorageUsage() {
	storageFree, storageTotal := storage()

	storageUsage := ((storageTotal - storageFree) / storageTotal) * 100

	fmt.Println(strconv.FormatFloat(storageUsage, 'f', 2, 64))
}
