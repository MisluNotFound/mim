package logic

import (
	"mim/internal/logic/rpc"
)

func InitLogic() {
    go rpc.InitLogicRpc()
}
