package iso15765

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

func New(
	addrMode AddressingMode,
	idType IDType,
	frameFormat FrameFormat,
	messageType MessageType,
	cb CallBacks,
	cfg Config,
) *ISO15765Node {
	return &ISO15765Node{
		addrMode: addrMode,
		idType:   idType,
		in:       IOStream{frameFormat: frameFormat},
		out:      IOStream{frameFormat: frameFormat},
		flPdu:    nPDU{},
		cb:       cb,
		cfg:      cfg,
		timers:   timer{nBs: 5000, nCr: 5000, nCs: 5000},
		inQueue:  make([]Frame, ISOQueueSize),
	}
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
