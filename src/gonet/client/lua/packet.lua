require("crc")
local json = require ("json")
--server type
SERVICE_NONE           = 0
SERVICE_CLIENT         = 1
SERVICE_GATESERVER     = 2
SERVICE_ACCOUNTSERVER  = 3
SERVICE_WORLDSERVER    = 4
SERVICE_MONITORSERVER  = 5
--chat type
CHAT_MSG_TYPE_WORLD    = 0
CHAT_MSG_TYPE_PRIVATE  = 1
CHAT_MSG_TYPE_ORG      = 2
CHAT_MSG_TYPE_COUNT    = 3

Default_Ipacket_Stx  = 39
Default_Ipacket_Ckx  = 114

m_PacketCreateMap = {}
m_PacketMap = {}
Ipacket = {Stx = 0, DestServerType = 0, Ckx = 0, Id = 0}
Message = {PacketHead = Ipacket, MessageName=""}

function Message:new(o, id, destservertype, packetName)
    o = o or {}
    setmetatable(o, self)
    self.__index = self
    self:Init(id, destservertype, packetName)
    return o
end

function Message:Init(id, destservertype, packetName)
    self.PacketHead.Stx = Default_Ipacket_Stx
    self.PacketHead.DestServerType = destservertype
    self.PacketHead.Ckx = Default_Ipacket_Ckx
    self.PacketHead.Id = id
    self.MessageName = string.lower(packetName)
    return o
end

function RegisterPacket(packet, func)
    name = string.lower(packet.MessageName)
    id = CRC32.hash(name)
    packetFunc = function()
    		packet = new(packet)
    		packet:Init(0, 0, name)
    		return packet
    end
    m_PacketCreateMap[id] = packetFunc
    m_PacketMap[id] = func
end

function HandlePacket(dat)
    id = bytes_to_int(string.sub(dat, 0, 4))
    buff = string.sub(dat, 4)
    packet, bEx = m_PacketCreateMap[id]
    if bEx then
        json.decode(buff, packet)
        m_PacketMap[id](packet)
    end
end

--前四位为包头名
function Encode(packet)
    name = string.lower(packet.MessageName)
    packetId = CRC32.hash(name)
	buff = json.encode(packet)
	data = int_to_bytes(packetId) .. buff
	return data
end

function bytes_to_int(str,endian,signed) -- use length of string to determine 8,16,32,64 bits
    local t={str:byte(1,-1)}
    if endian=="big" then --reverse bytes
        local tt={}
        for k=1,#t do
            tt[#t-k+1]=t[k]
        end
        t=tt
    end
    local n=0
    for k=1,#t do
        n=n+t[k]*2^((k-1)*8)
    end
    if signed then
        n = (n > 2^(#t-1) -1) and (n - 2^#t) or n -- if last bit set, negative.
    end
    return n
end

function int_to_bytes(num,endian,signed)
    if num<0 and not signed then num=-num print"warning, dropping sign from number converting to unsigned" end
    local res={}
    local n = math.ceil(select(2,math.frexp(num))/8) -- number of bytes to be used.
    if signed and num < 0 then
        num = num + 2^n
    end
    for k=n,1,-1 do -- 256 = 2^8 bits per char.
        local mul=2^(8*(k-1))
        res[k]=math.floor(num/mul)
        num=num-res[k]*mul
    end
    assert(num==0)
    if endian == "big" then
        local t={}
        for k=1,n do
            t[k]=res[n-k+1]
        end
        res=t
    end
    return string.char(unpack(res))
end