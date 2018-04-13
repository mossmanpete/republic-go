package hyper

import (
	"context"

	"github.com/republicprotocol/republic-go/dispatch"
)

type ChannelSet struct {
	BufferSize uint8
	ctx        context.Context
	Proposal   chan Proposal
	Prepare    chan Prepare
	Fault      chan Fault
	Commit     chan Commit
	Block      chan Block
}

func EmptyChannelSet(ctx context.Context, size uint8) ChannelSet {
	return ChannelSet{
		ctx:        ctx,
		BufferSize: size,
		Proposal:   make(chan Proposal, size),
		Prepare:    make(chan Prepare, size),
		Fault:      make(chan Fault, size),
		Commit:     make(chan Commit, size),
		Block:      make(chan Block, size),
	}
}

func (c *ChannelSet) Split(cs []ChannelSet) {

	proposals := make([]chan Proposal, len(cs))
	prepares := make([]chan Prepare, len(cs))
	commits := make([]chan Commit, len(cs))
	faults := make([]chan Fault, len(cs))
	blocks := make([]chan Block, len(cs))

	for i, chset := range cs {
		proposals[i] = chset.Proposal
		prepares[i] = chset.Prepare
		commits[i] = chset.Commit
		faults[i] = chset.Fault
		blocks[i] = chset.Block
	}

	go dispatch.Split(c.Proposal, proposals)
	go dispatch.Split(c.Prepare, prepares)
	go dispatch.Split(c.Commit, commits)
	go dispatch.Split(c.Fault, faults)
	go dispatch.Split(c.Block, blocks)

	func() {
		for {
			select {
			case <-c.ctx.Done():
				return
			}
		}
	}()
}

func (c *ChannelSet) Copy(cs ChannelSet) {
	for {
		select {
		case <-c.ctx.Done():
			return
		case proposal, ok := <-cs.Proposal:
			if !ok {
				return
			}
			c.Proposal <- proposal
		case prepare, ok := <-cs.Prepare:
			if !ok {
				return
			}
			// log.Println("Copying prepares in channelset")
			c.Prepare <- prepare
		case commit, ok := <-cs.Commit:
			if !ok {
				return
			}
			c.Commit <- commit
		case fault, ok := <-cs.Fault:
			if !ok {
				return
			}
			c.Fault <- fault

		case block, ok := <-cs.Block:
			if !ok {
				return
			}
			c.Block <- block
		}
	}
}
