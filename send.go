package iso15765

import "fmt"

func (n *ISO15765Node) send() error {
	if (n.out.status != IOStreamStatusTXBusy) && (n.out.status != IOStreamStatusTXReady) {
		// if there is no pending action just return
		return nil
	}
	var id uint32

	n.out.nPDU.nPCI.pCIType = n.outFrameType()

	switch n.out.nPDU.nPCI.pCIType {
	case PCITypeSF:
		n.out.nPDU.nPCI.dataLength = n.out.msgSize
		n.out.nPDU.size = n.out.msgSize

		if err := pduPack(n.addrMode, &(n.out.nPDU), &id, n.out.msg); err != nil {
			n.out.status = IOStreamStatusIdle
			n.out.cfCnt = 0
			n.out.wfCnt = 0
			signaling(NConf, n.out, n.cb, 0)
			return err
		}
		_, err := n.cb.SendFrame(
			n.idType,
			id,
			n.out.frameFormat,
			getClosestCANDl(n.out.nPDU.size+uint16(getDataOffset(n.addrMode, PCITypeSF, n.out.nPDU.size)), n.out.frameFormat),
			n.out.nPDU.data,
		)
		if err != nil {
			return err
		}
	case PCITypeFF:
		n.out.nPDU.nPCI.dataLength = n.out.msgSize
		n.out.wfCnt = 0
		if n.out.frameFormat == FrameFormatStandard {
			if (n.addrMode & 0x01) == 0 {
				n.out.nPDU.size = 6
			} else {
				n.out.nPDU.size = 5
			}
		} else {
			if n.addrMode&0x01 == 0 {
				n.out.nPDU.size = 62
			} else {
				n.out.nPDU.size = 61
			}
		}
		if err := pduPack(n.addrMode, &(n.out.nPDU), &id, n.out.msg); err != nil {
			return err
		}
		n.out.cfCnt = 1

		n.out.status = IOStreamStatusTXWaitFC
		var dlc uint8 = 8
		if n.out.frameFormat != FrameFormatStandard {
			dlc = 64
		}
		_, err := n.cb.SendFrame(n.idType, id, n.out.frameFormat, dlc, n.out.nPDU.data)
		if err != nil {
			return err
		}
		n.out.lastUpdate.nBs = n.cb.GetMs()
	case PCITypeCF:
		if n.out.lastUpdate.nCs+uint32(n.cfg.stmin) > n.cb.GetMs() {
			// if the minimun difference between transmissions is not reached then skip
			return nil
		}
		n.out.nPDU.nPCI.squenceNumber = n.out.cfCnt & 0x0F
		if n.out.cfCnt == 0xFF {
			n.out.cfCnt = 0
		} else {
			n.out.cfCnt += 1
		}
		var maxPayload uint16
		if n.out.frameFormat == FrameFormatStandard {
			maxPayload = 6
			if (n.addrMode & 0x01) == 0x00 {
				maxPayload = 7
			}
		} else {
			maxPayload = 62
			if (n.addrMode & 0x01) == 0x00 {
				maxPayload = 63
			}
		}
		n.out.nPDU.size = n.out.msgSize - n.out.msgPos
		if n.out.nPDU.size >= maxPayload {
			n.out.nPDU.size = maxPayload
		}

		if err := pduPack(n.addrMode, &(n.out.nPDU), &id, n.out.msg[n.out.msgPos:]); err != nil {
			return err
		}

		n.out.msgPos += n.out.nPDU.size

		if n.out.nPDU.nPCI.squenceNumber == n.cfg.bs {
			n.out.status = IOStreamStatusTXWaitFC
			n.out.lastUpdate.nBs = n.cb.GetMs()
		}

		var offset uint16 = 2
		if (n.addrMode & 0x01) == 0 {
			offset = 1
		}
		_, err := n.cb.SendFrame(
			n.idType,
			id,
			n.out.frameFormat,
			getClosestCANDl(n.out.nPDU.size+offset, n.out.frameFormat),
			n.out.nPDU.data,
		)
		if err != nil {
			return err
		}
		n.out.lastUpdate.nCs = n.cb.GetMs()
		if n.out.msgPos >= n.out.msgSize {
			return fmt.Errorf("msgPos overflowed msgSize: %d", n.out.msgPos)
		}
	default:
		return fmt.Errorf("invalid PCI type: %d", n.out.nPDU.nPCI.pCIType)
	}
	return nil
}

func (n *ISO15765Node) outFrameType() PCIType {
	result := PCITypeCF

	if n.out.cfCnt == 0 {
		if (n.addrMode & 0x01) == 1 {
			if n.out.frameFormat == FrameFormatStandard {
				if n.out.msgSize <= 6 {
					result = PCITypeSF
				} else {
					result = PCITypeFF
				}
			} else {
				if n.out.msgSize <= 61 {
					result = PCITypeSF
				} else {
					result = PCITypeFF
				}
			}
		} else {
			if n.out.frameFormat == FrameFormatStandard {
				if n.out.msgSize <= 7 {
					result = PCITypeSF
				} else {
					result = PCITypeFF
				}
			} else {
				if n.out.msgSize <= 62 {
					result = PCITypeSF
				} else {
					result = PCITypeFF
				}
			}
		}
	}

	return result
}

func pduPack(mode AddressingMode, pdu *nPDU, id *uint32, data []uint8) error {
	if (data == nil) || (id == nil) {
		return fmt.Errorf("data and id must not be nil")
	}

	switch mode {
	case AddressingModeExtended:
		pdu.data[0] = pdu.nAI.targetAddress
	case AddressingModeNormal:
		*id = 0x80 |
			(uint32(pdu.nAI.addrPriority) << 8) |
			(uint32(pdu.nAI.targetAddress) << 3) |
			(uint32(pdu.nAI.sourceAddress) << 0)
		if pdu.nAI.targetAddressType == TargetAddressTypePhysical {
			*id = *id | 0x00000040
		} else {
			*id = *id | 0x00000000
		}
	case AddressingModeMixed29:
		*id = (uint32(pdu.nAI.addrPriority) << 26) |
			(uint32(pdu.nAI.targetAddress) << 8) |
			(uint32(pdu.nAI.sourceAddress) << 0)
		if pdu.nAI.targetAddressType == TargetAddressTypePhysical {
			*id = *id | (0x000000CE << 16)
		} else {
			*id = *id | (0x000000CD << 16)
		}
	case AddressingModeFixed:
		*id = (uint32(pdu.nAI.addrPriority) << 26) |
			(uint32(pdu.nAI.targetAddress) << 8) |
			(uint32(pdu.nAI.sourceAddress) << 0)
		if pdu.nAI.targetAddressType == TargetAddressTypePhysical {
			*id = *id | (0x000000DA << 16)
		} else {
			*id = *id | (0x000000DB << 16)
		}
	case AddressingModeMixed11:
		*id = 0x80 |
			(uint32(pdu.nAI.addrPriority) << 8) |
			(uint32(pdu.nAI.targetAddress) << 3) |
			(uint32(pdu.nAI.sourceAddress) << 0)
		if pdu.nAI.targetAddressType == TargetAddressTypePhysical {
			*id = *id | 0x00000040
		} else {
			*id = *id | 0x00000000
		}
	default:
		fmt.Errorf("invalid addressing mode: %d", mode)
	}

	return nil
}

func pduPackData(mode AddressingMode, pdu *nPDU, data []uint8) error {
	if data == nil {
		return fmt.Errorf("data must not be nil")
	}

	var offset uint8

	switch pdu.nPCI.pCIType {
	case PCITypeSF:
		offset = getDataOffset(mode, PCITypeSF, pdu.size)
	case PCITypeFF:
		offset = getDataOffset(mode, PCITypeFF, pdu.size)
	case PCITypeCF:
		offset = getDataOffset(mode, PCITypeCF, pdu.size)
	case PCITypeFC:
		offset = getDataOffset(mode, PCITypeFC, pdu.size)
	default:
		fmt.Errorf("invalid PCI type: %d", pdu.nPCI.pCIType)
	}
	copy(data[offset:], data)

	return nil
}
