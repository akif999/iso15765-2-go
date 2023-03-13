package iso15765

type nPCI struct {
	flowStatus    uint8
	blockSize     uint8
	squenceNumber uint8
	pCIType       PCIType
	dataLength    uint16
}

type nPDU struct {
	messageType MessageType
	size        uint16
	nAI         nAI
	nPCI        nPCI
	data        []uint8
}
