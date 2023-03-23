package iso15765

import "fmt"

func (n *ISO15765Node) recv(frame Frame) error {
	n.in.frameFormat = frame.FFormat
	if err := nPduUnpack(
		n.addrMode, &(n.in.nPDU), frame.ID, uint8(frame.dlc), frame.data,
	); err != nil {
		return err
	}
	switch n.in.nPDU.nPCI.pCIType {
	case PCITypeFC:
		n.recvFC(frame)
	case PCITypeCF:
		n.recvCF(frame)
	case PCITypeSF:
		n.recvSF(frame)
	case PCITypeFF:
		n.recvFF(frame)
	}

	return nil
}

func (n *ISO15765Node) recvSF(frame Frame) error {
	if (n.in.status & IOStreamStatusRXBusy) != 0 {
		return fmt.Errorf("rx is busy")
	}
	copy(n.in.msg, n.in.nPDU.data[:n.in.nPDU.size])
	n.in.status = IOStreamStatusIdle
	err := signaling(NIndi, n.in, n.cb, n.in.nPDU.size)
	if err != nil {
		return err
	}

	return nil
}

func (n *ISO15765Node) recvFF(frame Frame) error {
	if len(n.in.msg) < int(ISOMsgSize) {
		return fmt.Errorf("FF data length is must be equal or less than 4095: got=%d", len(n.in.msg))
	}
	if (n.in.status & IOStreamStatusRXBusy) != 0 {
		return fmt.Errorf("rx is busy")
	}
	copy(n.in.msg, n.in.nPDU.data[:n.in.nPDU.size])
	n.in.msgPos = n.in.nPDU.size
	n.in.cfCnt = 0
	n.in.wfCnt = 0
	err := signaling(NFFIndi, n.in, n.cb, n.in.msgSize)
	if err != nil {
		return err
	}
	err = sendFC()
	if err != nil {
		return err
	}

	return nil
}

func (n *ISO15765Node) recvCF(frame Frame) error {
	if (n.in.status & IOStreamStatusRXBusy) == 0 {
		return fmt.Errorf("consecutiveFrame was received without being expected")
	}
	// According to (ref: iso15765-2 p.26) if we are not in progress of reception we should ignore it
	if (n.in.cfCnt + 1) > 0xFF {
		n.in.cfCnt = 0
	} else {
		n.in.cfCnt += 1
	}
	if (n.in.cfCnt & 0x0F) != n.in.nPDU.nPCI.squenceNumber {
		return fmt.Errorf("consecutiveFrame sequence number was invalid")
	}

	copy(n.in.msg[n.in.msgPos:], n.in.nPDU.data[:n.in.nPDU.size])
	n.in.msgPos += n.in.nPDU.size

	if n.in.msgPos >= n.in.msgSize {
		err := signaling(NIndi, n.in, n.cb, n.in.msgSize)
		if err != nil {
			return err
		}
		n.in = IOStream{}
		return nil
	}
	if n.cfg.bs != 0 {
		if n.in.cfCnt == n.cfg.bs {
			n.in.cfCnt = 0
			err := sendFC()
			if err != nil {
				return err
			}
		}
	}
	n.in.lastUpdate.nCr = n.cb.GetMs()

	return nil
}

func (n *ISO15765Node) recvFC(frame Frame) error {
	if n.out.status != IOStreamStatusTXWaitFC {
		return fmt.Errorf("upon reception of an unexpected protocol data unit")
	}
	switch n.in.nPDU.nPCI.flowStatus {
	case FlowControlStatusWait:
		n.out.wfCnt += 1
		// TODO: fix it
		if err := n.checkMaxWfCapacity(); err != nil {
			n.out.lastUpdate.nBs = n.cb.GetMs()
			return err
		}
	case FlowControlStatusOverFlow:
		return fmt.Errorf("flow control status is overflow")
	case FlowControlStatusContinue:
		n.out.cfgBS = n.in.nPDU.nPCI.blockSize
		n.out.Stmin = n.in.nPDU.nPCI.stmin
		setStreamData(n.out, 1, 0, IOStreamStatusTXReady)
	default:
		return fmt.Errorf("invalid flow status: %d", n.in.nPDU.nPCI.flowStatus)
	}
	return nil
}

func (n *ISO15765Node) checkMaxWfCapacity() error {
	if n.out.wfCnt >= n.cfg.wf {
		setStreamData(n.out, 0, 0, IOStreamStatusIdle)
		return fmt.Errorf("wf capacity overflowed: %d", n.out.wfCnt)
	}
	return nil
}

func signaling(tp SignalTP, stream IOStream, cb CallBacks, msgSize uint16) error {
	// if cb == nil {
	// 	return fmt.Errorf("callback is must not be nil")
	// }
	switch tp {
	case NIndi:
		var indn nInd
		indn.FrameFormat = stream.frameFormat
		indn.nAI = stream.nPDU.nAI
		indn.nPCI = stream.nPDU.nPCI
		stream.status = IOStreamStatusIdle
		cb.Ind(indn)
	case NFFIndi:
		var indn nFFInd
		indn.FrameFormat = stream.frameFormat
		indn.nAI = stream.nPDU.nAI
		indn.nPCI = stream.nPDU.nPCI
		stream.status = stream.status | IOStreamStatusRXBusy
		cb.FFInd(indn)
	case NConf:
		var cnf nCfm
		cnf.nAI = stream.nPDU.nAI
		cnf.nPCI = stream.nPDU.nPCI
		cb.Cfm(cnf)
	default:
		return nil
	}

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
