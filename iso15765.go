package iso15765

import "fmt"

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
	inQueue  chan Frame
}

func New(
	addrMode AddressingMode,
	idType IDType,
	cb CallBacks,
	cfg Config,
) *ISO15765Node {
	return &ISO15765Node{
		addrMode: addrMode,
		idType:   idType,
		in:       IOStream{},
		out:      IOStream{},
		flPdu:    nPDU{},
		cb:       cb,
		cfg:      cfg,
		// TODO: There is sloopy value must fix or must be configurable.
		timers:  timer{nBs: 5000, nCr: 5000, nCs: 5000},
		inQueue: make(chan Frame, ISOQueueSize),
	}
}

func (n *ISO15765Node) Configure() error {
	if (n.idType != IDTypeStandard) && (n.idType != IDTypeExtended) {
		return fmt.Errorf("invalid ID type: %d", n.idType)
	}
	if (uint8(n.idType) & uint8(n.addrMode)) == 0 {
		return fmt.Errorf("invalid addressing mode: %d", n.addrMode)
	}

	if (n.cb.SendFrame == nil) || (n.cb.GetMs == nil) {
		return fmt.Errorf("SendFrame and GetMs must not be nil")
	}
	if n.cb.Ind == nil {
		n.cb.Ind = func(nInd) error { return nil }
	}
	if n.cb.FFInd == nil {
		n.cb.FFInd = func(nFFInd) error { return nil }
	}
	if n.cb.Cfm == nil {
		n.cb.Cfm = func(nCfm) error { return nil }
	}
	if n.cb.CfgCfm == nil {
		n.cb.CfgCfm = func(nCngParamCfm) error { return nil }
	}

	return nil
}

func (n *ISO15765Node) Enqueue(frame Frame) error {
	n.inQueue <- frame
	return nil
}

func (n *ISO15765Node) Send(frame nRequest) error {
	if n.out.status != IOStreamStatusIdle {
		return fmt.Errorf("send is busy")
	}
	if len(frame.msg) > int(ISOMsgSize) {
		return fmt.Errorf("msgSize is must be equal or less than %d", ISOMsgSize)
	}
	n.out.frameFormat = frame.FrameFormat
	n.out.msgSize = uint16(len(frame.msg))
	copy(n.out.msg, frame.msg)
	n.out.nPDU.nAI = frame.nAI
	n.out.status = IOStreamStatusTXBusy

	return nil
}

func (n *ISO15765Node) Process() error {
	if err := timeoutProcess(); err != nil {
		return err
	}
	for f := range n.inQueue {
		if err := n.processIn(f); err != nil {
			return err
		}
	}
	n.processOut()
	return nil
}

func (n *ISO15765Node) processIn(frame Frame) error {
	return nil
}

func (n *ISO15765Node) processOut() error {
	return nil
}
