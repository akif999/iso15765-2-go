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
