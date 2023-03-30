package main

import (
	"fmt"
	"io/ioutil"
	"math/big"
	"strconv"
	"strings"
)

func main() {
	cpu()
}

func cpu() {
	var cpuRawData string
	file, err := ioutil.ReadFile("/proc/stat")
	if err != nil {
		fmt.Println(err)
	}
	lines := strings.Split(string(file), "\n")
	for i, lineData := range lines {
		if i == 0 {
			cpuRawData = lineData
		}
	}
	cpuData := strings.Fields(cpuRawData)
	fmt.Println(cpuData)
	user, err := strconv.ParseFloat(cpuData[2], 32)
	nice, err := strconv.ParseFloat(cpuData[3], 32)
	idle, err := strconv.ParseFloat(cpuData[4], 32)
	iowait, err := strconv.ParseFloat(cpuData[5], 32)
	irq, err := strconv.ParseFloat(cpuData[6], 32)
	softirq, err := strconv.ParseFloat(cpuData[7], 32)
	steal, err := strconv.ParseFloat(cpuData[8], 32)
	guest, err := strconv.ParseFloat(cpuData[9], 32)
	guest_nice, err := strconv.ParseFloat(cpuData[10], 32)

	total := user + nice + idle + iowait + irq + softirq + steal + guest + guest_nice
	idlePercent := truncate((1-(idle/total))*100, 0.001)
	// fmt.Printf("%.2f", idlePercent)
	fmt.Println(idlePercent)
}
func truncate(f float64, unit float64) float64 {
	bf := big.NewFloat(0).SetPrec(1000).SetFloat64(f)
	bu := big.NewFloat(0).SetPrec(1000).SetFloat64(unit)

	bf.Quo(bf, bu)

	// Truncate:
	i := big.NewInt(0)
	bf.Int(i)
	bf.SetInt(i)

	f, _ = bf.Mul(bf, bu).Float64()
	return f
}

// func Round(x, unit float64) float64 {
// 	return math.Round(x/unit) * unit
// }
