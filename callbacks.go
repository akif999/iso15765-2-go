package iso15765

type CallBacks struct {
	Ind           func(nInd) error
	FFInd         func(nFFInd) error
	Cfm           func(nCfm) error
	CfgCfm        func(nCngParamCfm) error
	PduCustomPack func(nPDU) error
	GetMs         func() uint32
	SendFrame     func(IDType, uint32, FrameFormat, uint8, uint8) (uint8, error)
}
