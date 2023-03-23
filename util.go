package iso15765

func getDataOffset(mode AddressingMode, pci PCIType, dataSize uint16) uint8 {
	var offset uint8
	offset = uint8(mode) & 0x01

	switch pci {
	case PCITypeSF:
		offset += 1
		if !(dataSize <= uint16(8-offset)) {
			offset += 1
		}
	case PCITypeFF:
		offset += 2
		if dataSize > 4095 {
			offset += 4
		}
	case PCITypeCF:
		offset += 1
	case PCITypeFC:
		offset += 3
	default:
		offset = 0
	}
	return offset
}

func setStreamData(io IOStream, cf, wf uint8, status IOStreamStatus) {
	io.cfCnt = cf
	io.wfCnt = wf
	io.status = status
}
