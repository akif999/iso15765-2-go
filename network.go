package iso15765

type nAI struct {
	addrPriority      uint8
	sourceAddress     uint8
	targetAddress     uint8
	addressExtension  uint8
	targetAddressType TargetAddressType
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
	msg         []uint8
}

type nInd struct {
	FrameFormat FrameFormat
	nAI         nAI
	nPCI        nPCI
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
