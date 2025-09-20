/**
 * @author: dn-jinmin/dn-jinmin
 * @doc:
 */

package xerr

const (
	SERVER_COMMON_ERROR = 100001
	REQUEST_PARAM_ERROR = 100002
	DB_ERROR            = 100003
	UNAUTHORIZED_ERROR  = 100004 // 未授权或登录态失效
	DEVICE_KICKED_ERROR = 100005 // 设备被踢下线
	INVALID_DEVICE_TYPE = 100006 // 非法设备类型
	NO_PERMISSION       = 100007 // 无权限操作
	USER_NOT_FOUND      = 100101 // 用户不存在/未注册
)
