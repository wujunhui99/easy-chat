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
	if u, ok := ctx.Value(Deveicetype).(string); ok {
		return u
	}
	return ""
}