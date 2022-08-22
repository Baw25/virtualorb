package status_report

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateSingleStatusReport(t *testing.T) {
	testStatusReport, _ := GenerateSingleStatusReport()

	assert.NotEqual(t, testStatusReport.BatteryLevelPercent, "", "BatteryLevelPercent should not be an empty string")
	assert.NotEqual(t, testStatusReport.BatteryVoltage, "", "BatteryVoltage should not be an empty string")
	assert.NotEqual(t, testStatusReport.DeviceCPUPercent, "", "DeviceCPUPercent should not be an empty string")
	assert.NotEqual(t, testStatusReport.DeviceCPUTemp, "", "DeviceCPUTemp should not be an empty string")
	assert.NotEqual(t, testStatusReport.DeviceDiskSpaceUsedPercent, "", "DeviceDiskSpaceUsedPercent should not be an empty string")
	assert.NotEqual(t, testStatusReport.DeviceDiskSpaceAvailablePercent, "", "DeviceDiskSpaceAvailablePercent should not be an empty string")
}

func TestGetDeviceBatteryLevelPercent(t *testing.T) {
	percent, _ := GetDeviceBatteryLevelPercent()

	assert.Equal(t, (percent <= 100.0 && percent > 0.0), true, "percent must be between 0 and 100")
	assert.Equal(t, reflect.TypeOf(percent).Kind(), reflect.Float64, "percent should be a float")
}

func TestGetDeviceBatteryVoltage(t *testing.T) {
	voltage, _ := GetDeviceBatteryVoltage()

	assert.Equal(t, reflect.TypeOf(voltage).Kind(), reflect.Float64, "voltage should be a float")
}

func TestGetDeviceCPU(t *testing.T) {
	percent, _ := GetDeviceCPU()

	assert.Equal(t, (percent <= 100.0 && percent > 0.0), true, "percent must be between 0 and 100")
	assert.Equal(t, reflect.TypeOf(percent).Kind(), reflect.Float64, "percent should be a float")
}

func TestGetCPUTemp(t *testing.T) {
	temp, _ := GetCPUTemp()

	assert.Equal(t, reflect.TypeOf(temp).Kind(), reflect.String, "temperature should be a string")
}

func TestGetDeviceDiskSpace(t *testing.T) {
	used, available, _ := GetDeviceDiskSpace()

	assert.Equal(t, (used <= 100.0 && used > 0.0), true, "used must be between 0 and 100")
	assert.Equal(t, (available <= 100.0 && available > 0.0), true, "available must be between 0 and 100")
	assert.Equal(t, reflect.TypeOf(used).Kind(), reflect.Float64, "used should be a float")
	assert.Equal(t, reflect.TypeOf(available).Kind(), reflect.Float64, "available should be a float")
	assert.Equal(t, 100.0-used, available, "available and used difference must make sense")
}
