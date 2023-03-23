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

func getClosestCANDl(size uint16, fFormat FrameFormat) uint8 {
	var dataLength uint8 = 0

	if fFormat == FrameFormatStandard {
		if size <= 8 {
			dataLength = uint8(size)
		} else {
			dataLength = 8
		}
	} else {
		if size <= 8 {
			dataLength = uint8(size)
		} else if size <= 12 {
			dataLength = 12
		} else if size <= 16 {
			dataLength = 16
		} else if size <= 20 {
			dataLength = 20
		} else if size <= 24 {
			dataLength = 24
		} else if size <= 32 {
			dataLength = 32
		} else if size <= 48 {
			dataLength = 48
		} else {
			dataLength = 64
		}
	}

	return dataLength
}

func setStreamData(io IOStream, cf, wf uint8, status IOStreamStatus) {
	io.cfCnt = cf
	io.wfCnt = wf
	io.status = status
}
