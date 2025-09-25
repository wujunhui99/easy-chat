/**
 * @author: dn-jinmin/dn-jinmin
 * @doc:
 */

package ctxdata

import "context"

func GetUid(ctx context.Context) string {
	if u, ok := ctx.Value(Identify).(string); ok {
		return u
	}
	return ""
}

func GetDevicetype(ctx context.Context) string {
	if u, ok := ctx.Value(DeveiceType).(string); ok {
		return u
	}
	return ""
}

func GetDeviceName(ctx context.Context) string {
	if u, ok := ctx.Value(DeveiceName).(string); ok {
		return u
	}
	return ""
}

func GetDeviceId(ctx context.Context) string {
	if u, ok := ctx.Value(DeviceID).(string); ok {
		return u
	}
	return ""
}
