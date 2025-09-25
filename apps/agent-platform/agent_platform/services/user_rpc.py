import grpc

from agent_platform.config import get_settings
from agent_platform.generated import user_pb2, user_pb2_grpc


class UserRpcClient:
    def __init__(self, channel: grpc.Channel | None = None) -> None:
        settings = get_settings()
        self._target = settings.user_rpc_target
        self._timeout = settings.user_rpc_timeout_seconds
        self._channel = channel or grpc.insecure_channel(self._target)
        self._stub = user_pb2_grpc.UserStub(self._channel)

    def create_agent(self, nickname: str, avatar: str, phone: str, sex: int) -> str:
        request = user_pb2.CreateAgentReq(
            nickname=nickname,
            avatar=avatar,
            phone=phone,
            sex=sex,
        )
        response = self._stub.CreateAgent(request, timeout=self._timeout)
        return response.id
