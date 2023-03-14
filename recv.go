package iso15765

import "fmt"

func (n *ISO15765Node) processRecv(frame Frame) error {
	n.recv.frameFormat = frame.FFormat
	return nil
}

func nPduUnpack(mode AddressingMode, nPdu *nPDU, id uint32, dlc uint8, data []uint8) error {
	if (nPdu == nil) || (data == nil) {
		return fmt.Errorf("nPdu and data must not be nil")
	}
	switch mode {
	case AddressingModeMixed11:
		nPdu.nAI.addrPriority = data[0]
	case AddressingModeNormal:
		nPdu.nAI.addrPriority = uint8((id & 0x00000700) >> 8)
		nPdu.nAI.targetAddress = uint8((id & 0x00000038) >> 3)
		nPdu.nAI.sourceAddress = uint8(id & 0x00000007 >> 0)
		if uint8((id&0x00000040)>>6) == 1 {
			nPdu.nAI.targetAddressType = TargetAddressTypePhysical
		} else {
			nPdu.nAI.targetAddressType = TargetAddressTypeFunctional
		}
		break
	case AddressingModeMixed29:
		nPdu.nAI.addressExtension = data[0]
		nPdu.nAI.addrPriority = uint8((id & 0x1C000000) >> 26)
		if uint8((id&0x00FF0000)>>16) == 0xCE {
			nPdu.nAI.targetAddressType = TargetAddressTypePhysical
		} else {
			nPdu.nAI.targetAddressType = TargetAddressTypeFunctional
		}
		nPdu.nAI.targetAddress = uint8((id & 0x0000FF00) >> 8)
		nPdu.nAI.sourceAddress = uint8((id & 0x000000FF) >> 0)
		break
	case AddressingModeFixed:
		nPdu.nAI.addrPriority = uint8((id & 0x1C000000) >> 26)
		if uint8((id&0x00FF0000)>>16) == 0xDA {
			nPdu.nAI.targetAddressType = TargetAddressTypePhysical
		} else {
			nPdu.nAI.targetAddressType = TargetAddressTypeFunctional
		}
		nPdu.nAI.targetAddress = uint8((id & 0x0000FF00) >> 8)
		nPdu.nAI.sourceAddress = uint8((id & 0x000000FF) >> 0)
		break
	case AddressingModeExtended:
		nPdu.nAI.addrPriority = uint8((id & 0x700) >> 8)
		nPdu.nAI.targetAddress = uint8((id & 0x38) >> 3)
		nPdu.nAI.sourceAddress = uint8((id & 0x07) >> 0)
		if uint8((id&0x40)>>6) == 0x01 {
			nPdu.nAI.targetAddress = TargetAddressTypePhysical
		} else {
			nPdu.nAI.targetAddress = TargetAddressTypeFunctional
		}
		nPdu.nAI.addressExtension = data[0]
		break
	default:
		return fmt.Errorf("invalid addressing mode: %d", mode)
	}

	return nil
}

func nPduUnpackData(mode AddressingMode, nPdu *nPDU, data []uint8) error {
	if (nPdu == nil) || (data == nil) {
		return fmt.Errorf("nPdu and data must not be nil")
	}

	switch nPdu.nPCI.pCIType {
	case PCITypeSF:
		copy(nPdu.data, data[getDataOffset(mode, PCITypeSF, nPdu.nPCI.dataLength):])
	case PCITypeFF:
		copy(nPdu.data, data[getDataOffset(mode, PCITypeFF, nPdu.size):])
	case PCITypeCF:
		copy(nPdu.data, data[getDataOffset(mode, PCITypeCF, nPdu.size):])
	case PCITypeFC:
		copy(nPdu.data, data[getDataOffset(mode, PCITypeFC, nPdu.size):])
	default:
		return fmt.Errorf("invalid PCI type: %d", nPdu.nPCI.pCIType)
	}

	return nil
}
