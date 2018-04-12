package hyper_test

import (
	"sync"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/hyperdrive"
)

var _ = Describe("Broadcast", func() {

	Context("Proposals", func() {

		It("should recieve all the broadcasted proposals", func() {
			chanSetIn := EmptyChannelSet(100)
			validator, _ := NewTestValidator(NewSharedBlocks(0, 0), 100)
			chanSetOut := ProcessBroadcast(chanSetIn, validator)

			statusMu := new(sync.RWMutex)
			status := map[[32]byte]uint8{}
			var wg sync.WaitGroup

			wg.Add(2)
			go func() {
				defer chanSetIn.Close()
				defer wg.Done()
				for i := 0; i < 100; i++ {
					proposal := Proposal{
						Height: uint64(i),
					}
					statusMu.Lock()
					status[ProposalHash(proposal)] = 1
					statusMu.Unlock()
					chanSetIn.Proposal <- proposal
				}
			}()

			go func() {
				defer wg.Done()
				for {
					select {
					case proposal, ok := <-chanSetOut.Proposal:
						if !ok {
							return
						}
						statusMu.RLock()
						Ω(status[ProposalHash(proposal)]).Should(Equal(uint8(1)))
						statusMu.RUnlock()

						statusMu.Lock()
						status[ProposalHash(proposal)]++
						statusMu.Unlock()
					}
				}
			}()

			wg.Wait()
		})

		It("should only return unique proposals", func() {
			chanSetIn := EmptyChannelSet(100)
			validator, _ := NewTestValidator(NewSharedBlocks(0, 0), 100)
			chanSetOut := ProcessBroadcast(chanSetIn, validator)

			statusMu := new(sync.RWMutex)
			status := map[[32]byte]uint8{}
			var wg sync.WaitGroup

			wg.Add(2)
			go func() {
				defer chanSetIn.Close()
				defer wg.Done()
				for i := 0; i < 100; i++ {
					for j := 0; j < i; j++ {
						proposal := Proposal{
							Height: uint64(i),
						}
						statusMu.Lock()
						status[ProposalHash(proposal)] = 1
						statusMu.Unlock()
						chanSetIn.Proposal <- proposal
					}
				}
			}()

			go func() {
				defer wg.Done()
				for {
					select {
					case proposal, ok := <-chanSetOut.Proposal:
						if !ok {
							return
						}
						statusMu.RLock()
						Ω(status[ProposalHash(proposal)]).Should(Equal(uint8(1)))
						statusMu.RUnlock()

						statusMu.Lock()
						status[ProposalHash(proposal)]++
						statusMu.Unlock()
					}
				}
			}()

			wg.Wait()
		})
	})
})
