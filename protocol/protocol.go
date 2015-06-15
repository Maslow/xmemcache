package protocol

type Packet struct {
	PacketHeader
	extras []byte
	key    []byte
	value  []byte
}

type PacketHeader struct {
	magic             byte    //Magic number identifying the package
	opcode            byte    //Command code
	keyLength         uint16  //Length in bytes of the text key that follows the command extras
	extrasLength      uint8   //Lenght in bytes of the command extras
	dataType          byte    //Reserved for future use
	vbucketIdOrStatus [2]byte //Vbucket id in request: The virtual bucket for this command. Or status in response: Status of the response (non-zero on error)
	totalBodyLength   uint32  //Length in bytes of extra + key + value
	opaque            [4]byte //Will be copied back to you in the response
	CAS               [8]byte //Data version check
}

func (p *PacketHeader) parseHeader(data []byte) bool {
	if len(data) < 24 {
		return false
	}
	if 0x80 != data[0] && 0x81 != data[0] {
		return false
	}

	p.magic = data[0]
	p.opcode = data[1]
	p.keyLength = uint16(data[2]<<8 | data[3])

	p.extrasLength = uint8(data[4])
	p.dataType = data[5]
	p.vbucketIdOrStatus[0], p.vbucketIdOrStatus[1] = data[6], data[7]

	p.totalBodyLength = uint32(data[8]<<24 | data[9]<<16 | data[10]<<8 | data[11])

	p.opaque = [4]byte{data[12], data[13], data[14], data[15]}

	for j, i := 0, 16; i < 24; i++ {
		p.CAS[j] = data[i]
		j++
	}
	return true
}

func (p *Packet) Parse(data []byte) bool {
	d := make([]byte, len(data))
	copy(d, data)

	if false == p.parseHeader(d) {
		return false
	}

	if 0 != p.extrasLength {
		p.extras = d[24 : 24+p.extrasLength]
	}

	keyStart := uint16(24 + p.extrasLength)
	if 0 != p.keyLength {
		p.key = d[keyStart : keyStart+p.keyLength]
	}

	valueStart := uint32(keyStart + p.keyLength)
	valueLength := p.totalBodyLength - uint32(p.extrasLength) - uint32(p.keyLength)
	if 0 != valueLength {
		p.value = d[valueStart : valueStart+valueLength]
	}
	return true
}

func (p *Packet) GetKey() (key []byte) {
	key = p.key
	return
}
