package iso15765

const (
	ISOMsgSize   uint16 = 4095
	ISOQueueSize uint8  = 64
)

type FrameFormat uint8

const (
	FrameFormatStandard FrameFormat = iota
	FrameFormatFD
)

type IDType uint8

const (
	IDTypeStandard IDType = iota
	IDTypeExtended
)

type AddressingMode uint8

const (
	AddressingModeUnkown AddressingMode = iota
	AddressingModeNormal
	AddressingModeFixed
	AddressingModeFixed11
	AddressingModeExtended
	AddressingModeExtended29
)

type PCIType uint8

const (
	PCITypeSF PCIType = 0x00
	PCITypeFF PCIType = 0x01
	PCITypeCF PCIType = 0x02
	PCITypeFC PCIType = 0x03
	PCITypeUN PCIType = 0xFF
)

type TargetAddressType uint8

const (
	TargetAddressTypePhysical = iota
	TargetAddressTypeFunctional
)

type MessageType uint8

const (
	MessageTypeDiag = iota
	MessageTypeRemoteDiag
)

type FlowControlStatus uint8

const (
	FlowControlStatusContinue FlowControlStatus = 0x00
	FlowControlStatusWait     FlowControlStatus = 0x01
	FlowControlStatusOverFlow FlowControlStatus = 0x02
)

type IOStreamStatus uint8

const (
	IOStreamStatusIdle = iota
	IOStreamStatusRXBusy
	IOStreamStatusTXBusy
	IOStreamStatusTXReady
	IOStreamStatusTXWaitFC
)

type FlowParam uint8

const (
	Stmin FlowParam = iota
	BS
)

type SignalTP uint8

const (
	NIndi SignalTP = iota
	NFFIndi
	NConf
	NCngPConf
)

type nAI struct {
	addrPriority      uint8
	sourceAddress     uint8
	targetAddress     uint8
	addressExtension  uint8
	targetAddressType TargetAddressType
}

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
	nPCI        nPCI
	data        []uint8
}

type nCfm struct {
	nAI  nAI
	nPCI nPCI
}

type nFFInd struct {
	FrameFormat FrameFormat
	nAI         nAI
	nPCI        nPCI
	msgSize     uint16
}

type nRequest struct {
	FrameFormat FrameFormat
	nAI         nAI
	nPCI        nPCI
	msgSize     uint16
	msg         []uint8
}

type nInd struct {
	FrameFormat FrameFormat
	nAI         nAI
	nPCI        nPCI
	msgSize     uint16
	msg         []uint8
}

type nCngParamReq struct {
	nAI         nAI
	nPCI        nPCI
	flParameter FlowParam
	val         uint8
}

type nCngParamCfm struct {
	nAI         nAI
	nPCI        nPCI
	flParameter FlowParam
	val         uint8
}

type CallBacks struct {
	Ind           func(nInd) error
	FFInd         func(nFFInd) error
	Cfm           func(nCfm) error
	CfgCfm        func(nCngParamCfm) error
	PduCustomPack func(nPDU) error
	GetMs         func() uint32
	SendFrame     func(IDType, uint32, FrameFormat, uint8, uint8) (uint8, error)
}

type timer struct {
	nBs  uint32
	nCr  uint32
	n_cs uint32
}

type IOStream struct {
	frameFormat FrameFormat
	nPDU        nPDU
	cfCnt       uint8
	wfCnt       uint8
	cfgWf       uint8
	Stmin       uint8
	cfgBS       uint8
	status      IOStreamStatus
	msgSize     uint16
	msgPos      uint16
	lastUpdate  timer
	msg         []uint8
}

type Config struct {
	stmin uint8
	bs    uint8
	wf    uint8
	nBs   uint16
	nCr   uint16
}

type Frame struct {
	ID      uint32
	IDType  IDType
	FFormat FrameFormat
	dlc     uint16
	data    []uint8
}

const (
	maxDataLength uint8 = 64
)

type ISO15765Node struct {
	addrMode AddressingMode
	idType   IDType
	in       IOStream
	out      IOStream
	flPdu    nPDU
	cb       CallBacks
	cfg      Config
	timers   timer
	inQueue  []Frame
}

func New() *ISO15765Node {
	return &ISO15765Node{}
}

func (n *ISO15765Node) Init() error {
	return nil
}

func (n *ISO15765Node) Send(frame nRequest) error {
	return nil
}

func (n *ISO15765Node) Enqueue(frame Frame) error {
	return nil
}

func (n *ISO15765Node) Process() error {
	return nil
}
