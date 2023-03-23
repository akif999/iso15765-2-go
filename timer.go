package iso15765

import "fmt"

type timer struct {
	nBs uint32
	nCr uint32
	nCs uint32
}

func (n *ISO15765Node) timeout() error {
	if (n.out.status != IOStreamStatusTXWaitFC) || (n.out.lastUpdate.nBs == 0) ||
		(n.cfg.nBs == 0) || ((n.out.lastUpdate.nBs + uint32(n.cfg.nBs)) >= n.cb.GetMs()) {
		return nil
	} else {
		// TODO: Imporove message
		n.out.cfCnt = 0x00
		signaling(NIndi, n.out, n.cb, n.out.msgSize)
		return fmt.Errorf("session timed out")
	}
}
