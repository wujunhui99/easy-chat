package constants

// 设备类型常量（简化集合）
const (
	DeviceMobile  = "mobile"  // 主设备，可踢其它设备
	DeviceWeb     = "web"     // 次级设备
	DeviceDesktop = "desktop" // 次级设备
	DevicePad     = "pad"     // 次级设备（平板）
)

// PrimaryDevices 主设备集合
var PrimaryDevices = map[string]struct{}{
	DeviceMobile: {},
}

// AllowedDevices 允许登录设备集合
var AllowedDevices = map[string]struct{}{
	DeviceMobile:  {},
	DeviceWeb:     {},
	DeviceDesktop: {},
	DevicePad:     {},
}

func IsAllowedDevice(dt string) bool {
	_, ok := AllowedDevices[dt]
	return ok
}

func IsPrimaryDevice(dt string) bool {
	_, ok := PrimaryDevices[dt]
	return ok
}
