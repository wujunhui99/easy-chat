/**
 * @author: dn-jinmin/dn-jinmin
 * @doc:
 */

package xerr

var codeText = map[int]string{
	SERVER_COMMON_ERROR: "服务器异常，稍后再尝试",
	REQUEST_PARAM_ERROR: "请求参数有误",
	DB_ERROR:            "数据库繁忙，稍后再尝试",
	UNAUTHORIZED_ERROR:  "未授权或登录状态已失效",
	DEVICE_KICKED_ERROR: "当前设备已被其它登录挤下线",
	INVALID_DEVICE_TYPE: "设备类型非法",
	NO_PERMISSION:       "当前设备无权限执行该操作",
	USER_NOT_FOUND:      "用户不存在",
}

func ErrMsg(errcode int) string {
	if msg, ok := codeText[errcode]; ok {
		return msg
	}
	return codeText[SERVER_COMMON_ERROR]
}
