local json = require ("json")
require ("packet")
C_A_LoginRequest = {
	Message,
	AccountName = "",
	BuildNo  = "",
	SocketId  = 0
}

function C_A_LoginRequest:new(o, id, destservertype, packetName)
    o = o or Message:new(o, id, destservertype, packetName)
    setmetatable(o, self)
    self.__index = self
    return o
end

aa = C_A_LoginRequest{}
aa.AccountName = test
aa.Message.Init(0, 1, "C_A_LoginRequest")
print(json.encode(aa))