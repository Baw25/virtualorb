package status_report

import (
	"fmt"
	"math"
	"math/rand"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/distatus/battery"
	"github.com/shirou/gopsutil/cpu"

	"github.com/ssimunic/gosensors"

	"github.com/shirou/gopsutil/disk"
)

const (
	Darwin     = "darwin"
	Windows    = "windows"
	Linux      = "linux"
	Mount      = "/System/Volumes/Data"
	FileSystem = "apfs"
)

type StatusReport struct {
	BatteryLevelPercent             string `json:"battery_level_percent"`
	BatteryVoltage                  string `json:"battery_voltage"`
	DeviceCPUPercent                string `json:"device_cpu_percent"`
	DeviceCPUTemp                   string `json:"device_cpu_temp"`
	DeviceDiskSpaceUsedPercent      string `json:"device_disk_space_used_percent"`
	DeviceDiskSpaceAvailablePercent string `json:"device_disk_space_available_percent"`
}

// Generate the entire status report: battery, cpu usage, cpu temp, disk space information
func GenerateSingleStatusReport() (StatusReport, error) {
	var newStatusReport StatusReport
	var reportError error

	batteryLevel, batterErr := GetDeviceBatteryLevelPercent()
	if batterErr != nil {
		reportError = batterErr
	}

	batteryVoltage, voltErr := GetDeviceBatteryVoltage()
	if voltErr != nil {
		reportError = voltErr
	}

	deviceCPUPercent, cpuPercentErr := GetDeviceCPU()
	if cpuPercentErr != nil {
		reportError = cpuPercentErr
	}

	deviceCPUTemp, cpuTempErr := GetCPUTemp()
	if cpuTempErr != nil {
		reportError = cpuTempErr
	}

	deviceDiskSpaceUsed, deviceDiskSpaceAvailable, diskErr := GetDeviceDiskSpace()
	if cpuTempErr != nil && (deviceDiskSpaceUsed != 0 && deviceDiskSpaceAvailable != 0) {
		reportError = diskErr
	}

	if reportError != nil {
		return newStatusReport, reportError
	}

	newStatusReport.BatteryLevelPercent = fmt.Sprintf("%f", batteryLevel)
	newStatusReport.BatteryVoltage = fmt.Sprintf("%fV", batteryVoltage)
	newStatusReport.DeviceCPUPercent = fmt.Sprintf("%f", deviceCPUPercent)
	newStatusReport.DeviceCPUTemp = deviceCPUTemp
	newStatusReport.DeviceDiskSpaceUsedPercent = fmt.Sprintf("%f", deviceDiskSpaceUsed)
	newStatusReport.DeviceDiskSpaceAvailablePercent = fmt.Sprintf("%f", deviceDiskSpaceAvailable)

	return newStatusReport, nil
}

// To simulate grabbing the orb's device information,
// we will acquire the current machine's battery, cpu usage, cpu temp, disk space information
func GetDeviceBatteryLevelPercent() (float64, error) {
	var deviceBatteryLevel float64
	firstBattery, err := battery.Get(0)
	if err != nil {
		return 0, err
	}
	deviceBatteryLevel = toFixedFloat((firstBattery.Current/firstBattery.Full)*100, 2)
	return deviceBatteryLevel, nil
}

func GetDeviceBatteryVoltage() (float64, error) {
	var deviceVoltage float64
	firstBattery, err := battery.Get(0)
	if err != nil {
		return 0, err
	}
	deviceVoltage = toFixedFloat(firstBattery.Voltage, 2)
	return deviceVoltage, nil
}

// Only collects cpu percentage from user level
func GetDeviceCPU() (float64, error) {
	percent, err := cpu.Percent(time.Second, true)
	duration := time.Duration(1) * time.Second
	time.Sleep(duration)
	if err != nil {
		return 0, err
	}

	resultPercent := percent[cpu.CPUser]
	return resultPercent, nil
}

func GetCPUTemp() (string, error) {
	var cpuTemp string

	if runtime.GOOS == Windows {
		// TODO: need to implement for windows - generate randomly for now
		cpuTemp = strconv.Itoa(randInt(40, 50)) + "°C"
		return cpuTemp, nil
	} else if runtime.GOOS == Darwin {
		// TODO: find clean way across various OS. In reality, this would be written for the OS of the actual orb
		out, err := exec.Command("./osx-cpu-temp/osx-cpu-temp").Output()
		if err != nil {
			return "0.0 °C", err
		}
		outStr := string(out[:])
		cpuTemp = outStr

		// If cputemp results in 0.0, then spoof temperature
		// osx-cpu-temp utility does not work for Apple M1 chips
		if strings.Contains(cpuTemp, "0.0") {
			cpuTemp = strconv.Itoa(randInt(40, 50)) + "°C"
		}

		return cpuTemp, nil
	} else if runtime.GOOS == Linux {
		// Gosensors assumes the presence of lm-sensors available on linux machines
		// See: https://github.com/ssimunic/gosensors
		sensors, err := gosensors.NewFromSystem()
		var cpuTemp string
		if err != nil {
			return "0.0 °C", err
		}
		for chip := range sensors.Chips {
			for key, value := range sensors.Chips[chip] {
				if key == "CPU" {
					cpuTemp = value
					break
				}
			}
		}
		return cpuTemp, nil
	}

	return "0.0 °C", nil
}

// For the sake virtualorb, only grab the disk space from the system space or app file system
func GetDeviceDiskSpace() (float64, float64, error) {
	var deviceDiskSpaceUsed float64
	var deviceDiskSpaceAvailable float64
	parts, partErr := disk.Partitions(true)
	if partErr != nil {
		return 0, 0, partErr
	}

	for _, p := range parts {
		device := p.Mountpoint
		s, err := disk.Usage(device)
		if err != nil {
			return 0, 0, err
		}
		if p.Mountpoint == Mount || s.Fstype == FileSystem {
			if s.Total == 0 {
				continue
			}

			deviceDiskSpaceUsed = s.UsedPercent
			deviceDiskSpaceAvailable = 100.0 - s.UsedPercent

		}
	}

	return deviceDiskSpaceUsed, deviceDiskSpaceAvailable, nil
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

func toFixedFloat(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}
