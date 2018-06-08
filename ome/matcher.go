package ome

import (
	"fmt"

	"github.com/republicprotocol/republic-go/logger"
	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/shamir"
	"github.com/republicprotocol/republic-go/smpc"
)

// ResolveStage defines the various stages that resolving can be in for any
// given Computation.
type ResolveStage = byte

// Values for ResolveStage.
const (
	ResolveStageNil ResolveStage = iota
	ResolveStagePriceExp
	ResolveStagePriceCo
	ResolveStageBuyVolumeExp
	ResolveStageBuyVolumeCo
	ResolveStageSellVolumeExp
	ResolveStageSellVolumeCo
	ResolveStageTokens
)

// A MatchCallback is called when a Computation is finished. The Computation
// can then be inspected to determine if the result is a match.
type MatchCallback func(Computation)

// A Matcher resolves Computations into a matched, or mismatched, result.
type Matcher interface {

	// Resolve a Computation to determine whether or not the orders involved
	// are a match. The ξ hash is used to define the ξ in which this
	// Computation exists, and the MatchCallback is called when a result has
	// be determined.
	Resolve(ξ [32]byte, com Computation, buyFragment, sellFragment order.Fragment, callback MatchCallback)
}

type matcher struct {
	storer Storer
	smpcer smpc.Smpcer
}

// NewMatcher returns a Matcher that will resolve Computations by resolving
// each component in a pipeline. If a mismatch is encounterd at any stage of
// the pipeline, the Computation is short circuited and the MatchCallback will
// be called immediately.
func NewMatcher(storer Storer, smpcer smpc.Smpcer) Matcher {
	return &matcher{
		storer: storer,
		smpcer: smpcer,
	}
}

// Resolve implements the Matcher interface.
func (matcher *matcher) Resolve(ξ [32]byte, com Computation, buyFragment, sellFragment order.Fragment, callback MatchCallback) {
	matcher.resolvePriceExp(smpc.NetworkID(ξ), buyFragment, sellFragment, com, callback)
}

func (matcher *matcher) resolvePriceExp(networkID smpc.NetworkID, buyFragment, sellFragment order.Fragment, com Computation, callback MatchCallback) {
	priceExpShare := buyFragment.Price.Exp.Sub(&sellFragment.Price.Exp)
	priceExpJoin := smpc.Join{
		ID:     smpc.JoinID(com.ID),
		Index:  smpc.JoinIndex(priceExpShare.Index),
		Shares: shamir.Shares{priceExpShare},
	}
	priceExpJoin.ID[31] = ResolveStagePriceExp

	err := matcher.smpcer.Join(networkID, priceExpJoin, func(joinID smpc.JoinID, values []uint64) {
		if len(values) != 1 {
			logger.Compute(logger.LevelError, fmt.Sprintf("cannot resolve priceExp: unexpected number of values: %v", len(values)))
			return
		}
		if isGreaterThanZero(values[0]) {
			logger.Compute(logger.LevelDebug, fmt.Sprintf("✔ priceExp => buy = %v, sell = %v", com.Buy, com.Sell))
			matcher.resolveBuyVolumeExp(networkID, buyFragment, sellFragment, com, callback)
			return
		}
		if isEqualToZero(values[0]) {
			logger.Compute(logger.LevelDebug, fmt.Sprintf("✔ priceExp => buy = %v, sell = %v", com.Buy, com.Sell))
			matcher.resolvePriceCo(networkID, buyFragment, sellFragment, com, callback)
			return
		}
		com.State = ComputationStateMismatched
		com.Match = false
		if err := matcher.storer.InsertComputation(com); err != nil {
			logger.Compute(logger.LevelError, fmt.Sprintf("cannot store mismatched computation buy = %v, sell = %v", com.Buy, com.Sell))
		}
		logger.Compute(logger.LevelDebugHigh, fmt.Sprintf("✗ priceExp => mismatch buy = %v, sell = %v", com.Buy, com.Sell))
		callback(com)
	})
	if err != nil {
		logger.Compute(logger.LevelError, fmt.Sprintf("cannot join priceExp: %v", err))
	}
}

func (matcher *matcher) resolvePriceCo(networkID smpc.NetworkID, buyFragment, sellFragment order.Fragment, com Computation, callback MatchCallback) {
	priceCoShare := buyFragment.Price.Co.Sub(&sellFragment.Price.Co)
	priceCoJoin := smpc.Join{
		ID:     smpc.JoinID(com.ID),
		Index:  smpc.JoinIndex(priceCoShare.Index),
		Shares: shamir.Shares{priceCoShare},
	}
	priceCoJoin.ID[31] = ResolveStagePriceCo

	err := matcher.smpcer.Join(networkID, priceCoJoin, func(joinID smpc.JoinID, values []uint64) {
		if len(values) != 1 {
			logger.Compute(logger.LevelError, fmt.Sprintf("cannot resolve priceCo: unexpected number of values: %v", len(values)))
			return
		}
		if isGreaterThanOrEqualToZero(values[0]) {
			logger.Compute(logger.LevelDebug, fmt.Sprintf("✔ priceCo => buy = %v, sell = %v", com.Buy, com.Sell))
			matcher.resolveBuyVolumeExp(networkID, buyFragment, sellFragment, com, callback)
			return
		}
		com.State = ComputationStateMismatched
		com.Match = false
		if err := matcher.storer.InsertComputation(com); err != nil {
			logger.Compute(logger.LevelError, fmt.Sprintf("cannot store mismatched computation buy = %v, sell = %v", com.Buy, com.Sell))
		}
		logger.Compute(logger.LevelDebugHigh, fmt.Sprintf("✗ priceCo => mismatch buy = %v, sell = %v", com.Buy, com.Sell))
		callback(com)
	})
	if err != nil {
		logger.Compute(logger.LevelError, fmt.Sprintf("cannot join priceCo: %v", err))
	}
}

func (matcher *matcher) resolveBuyVolumeExp(networkID smpc.NetworkID, buyFragment, sellFragment order.Fragment, com Computation, callback MatchCallback) {
	buyVolumeExpShare := buyFragment.Volume.Exp.Sub(&sellFragment.MinimumVolume.Exp)
	buyVolumeExpJoin := smpc.Join{
		ID:     smpc.JoinID(com.ID),
		Index:  smpc.JoinIndex(buyVolumeExpShare.Index),
		Shares: shamir.Shares{buyVolumeExpShare},
	}
	buyVolumeExpJoin.ID[31] = ResolveStageBuyVolumeExp

	err := matcher.smpcer.Join(networkID, buyVolumeExpJoin, func(joinID smpc.JoinID, values []uint64) {
		if len(values) != 1 {
			logger.Compute(logger.LevelError, fmt.Sprintf("cannot resolve buyVolumeExp: unexpected number of values: %v", len(values)))
			return
		}
		if isGreaterThanZero(values[0]) {
			logger.Compute(logger.LevelDebug, fmt.Sprintf("✔ buyVolumeExp => buy = %v, sell = %v", com.Buy, com.Sell))
			matcher.resolveSellVolumeExp(networkID, buyFragment, sellFragment, com, callback)
			return
		}
		if isEqualToZero(values[0]) {
			logger.Compute(logger.LevelDebug, fmt.Sprintf("✔ buyVolumeExp => buy = %v, sell = %v", com.Buy, com.Sell))
			matcher.resolveBuyVolumeCo(networkID, buyFragment, sellFragment, com, callback)
			return
		}
		com.State = ComputationStateMismatched
		com.Match = false
		if err := matcher.storer.InsertComputation(com); err != nil {
			logger.Compute(logger.LevelError, fmt.Sprintf("cannot store mismatched computation buy = %v, sell = %v", com.Buy, com.Sell))
		}
		logger.Compute(logger.LevelDebugHigh, fmt.Sprintf("✗ buyVolumeExp => mismatch buy = %v, sell = %v", com.Buy, com.Sell))
		callback(com)
	})
	if err != nil {
		logger.Compute(logger.LevelError, fmt.Sprintf("cannot join buyVolumeExp: %v", err))
	}
}

func (matcher *matcher) resolveBuyVolumeCo(networkID smpc.NetworkID, buyFragment, sellFragment order.Fragment, com Computation, callback MatchCallback) {
	buyVolumeCoShare := buyFragment.Volume.Co.Sub(&sellFragment.MinimumVolume.Co)
	buyVolumeCoJoin := smpc.Join{
		ID:     smpc.JoinID(com.ID),
		Index:  smpc.JoinIndex(buyVolumeCoShare.Index),
		Shares: shamir.Shares{buyVolumeCoShare},
	}
	buyVolumeCoJoin.ID[31] = ResolveStageBuyVolumeCo

	err := matcher.smpcer.Join(networkID, buyVolumeCoJoin, func(joinID smpc.JoinID, values []uint64) {
		if len(values) != 1 {
			logger.Compute(logger.LevelError, fmt.Sprintf("cannot resolve buyVolumeCo: unexpected number of values: %v", len(values)))
			return
		}
		if isGreaterThanOrEqualToZero(values[0]) {
			logger.Compute(logger.LevelDebug, fmt.Sprintf("✔ buyVolumeCo => buy = %v, sell = %v", com.Buy, com.Sell))
			matcher.resolveSellVolumeExp(networkID, buyFragment, sellFragment, com, callback)
			return
		}
		com.State = ComputationStateMismatched
		com.Match = false
		if err := matcher.storer.InsertComputation(com); err != nil {
			logger.Compute(logger.LevelError, fmt.Sprintf("cannot store mismatched computation buy = %v, sell = %v", com.Buy, com.Sell))
		}
		logger.Compute(logger.LevelDebugHigh, fmt.Sprintf("✗ buyVolumeCo => mismatch buy = %v, sell = %v", com.Buy, com.Sell))
		callback(com)
	})
	if err != nil {
		logger.Compute(logger.LevelError, fmt.Sprintf("cannot join buyVolumeCo: %v", err))
	}
}

func (matcher *matcher) resolveSellVolumeExp(networkID smpc.NetworkID, buyFragment, sellFragment order.Fragment, com Computation, callback MatchCallback) {
	sellVolumeExpShare := sellFragment.Volume.Exp.Sub(&buyFragment.MinimumVolume.Exp)
	sellVolumeExpJoin := smpc.Join{
		ID:     smpc.JoinID(com.ID),
		Index:  smpc.JoinIndex(sellVolumeExpShare.Index),
		Shares: shamir.Shares{sellVolumeExpShare},
	}
	sellVolumeExpJoin.ID[31] = ResolveStageSellVolumeExp

	err := matcher.smpcer.Join(networkID, sellVolumeExpJoin, func(joinID smpc.JoinID, values []uint64) {
		if len(values) != 1 {
			logger.Compute(logger.LevelError, fmt.Sprintf("cannot resolve sellVolumeExp: unexpected number of values: %v", len(values)))
			return
		}
		if isGreaterThanZero(values[0]) {
			logger.Compute(logger.LevelDebug, fmt.Sprintf("✔ sellVolumeExp => buy = %v, sell = %v", com.Buy, com.Sell))
			matcher.resolveTokens(networkID, buyFragment, sellFragment, com, callback)
			return
		}
		if isEqualToZero(values[0]) {
			logger.Compute(logger.LevelDebug, fmt.Sprintf("✔ sellVolumeExp => buy = %v, sell = %v", com.Buy, com.Sell))
			matcher.resolveSellVolumeCo(networkID, buyFragment, sellFragment, com, callback)
			return
		}
		com.State = ComputationStateMismatched
		com.Match = false
		if err := matcher.storer.InsertComputation(com); err != nil {
			logger.Compute(logger.LevelError, fmt.Sprintf("cannot store mismatched computation buy = %v, sell = %v", com.Buy, com.Sell))
		}
		logger.Compute(logger.LevelDebugHigh, fmt.Sprintf("✗ sellVolumeExp => mismatch buy = %v, sell = %v", com.Buy, com.Sell))
		callback(com)
	})
	if err != nil {
		logger.Compute(logger.LevelError, fmt.Sprintf("cannot join sellVolumeExp: %v", err))
	}
}

func (matcher *matcher) resolveSellVolumeCo(networkID smpc.NetworkID, buyFragment, sellFragment order.Fragment, com Computation, callback MatchCallback) {
	sellVolumeCoShare := sellFragment.Volume.Co.Sub(&buyFragment.MinimumVolume.Co)
	sellVolumeCoJoin := smpc.Join{
		ID:     smpc.JoinID(com.ID),
		Index:  smpc.JoinIndex(sellVolumeCoShare.Index),
		Shares: shamir.Shares{sellVolumeCoShare},
	}
	sellVolumeCoJoin.ID[31] = ResolveStageSellVolumeCo

	err := matcher.smpcer.Join(networkID, sellVolumeCoJoin, func(joinID smpc.JoinID, values []uint64) {
		if len(values) != 1 {
			logger.Compute(logger.LevelError, fmt.Sprintf("cannot resolve sellVolumeCo: unexpected number of values: %v", len(values)))
			return
		}
		if isGreaterThanOrEqualToZero(values[0]) {
			logger.Compute(logger.LevelDebug, fmt.Sprintf("✔ sellVolumeCo => buy = %v, sell = %v", com.Buy, com.Sell))
			matcher.resolveTokens(networkID, buyFragment, sellFragment, com, callback)
			return
		}
		com.State = ComputationStateMismatched
		com.Match = false
		if err := matcher.storer.InsertComputation(com); err != nil {
			logger.Compute(logger.LevelError, fmt.Sprintf("cannot store mismatched computation buy = %v, sell = %v", com.Buy, com.Sell))
		}
		logger.Compute(logger.LevelDebugHigh, fmt.Sprintf("✗ sellVolumeCo => mismatch buy = %v, sell = %v", com.Buy, com.Sell))
		callback(com)
	})
	if err != nil {
		logger.Compute(logger.LevelError, fmt.Sprintf("cannot join sellVolumeCo: %v", err))
	}
}

func (matcher *matcher) resolveTokens(networkID smpc.NetworkID, buyFragment, sellFragment order.Fragment, com Computation, callback MatchCallback) {
	tokensShare := buyFragment.Tokens.Sub(&sellFragment.Tokens)
	tokensJoin := smpc.Join{
		ID:     smpc.JoinID(com.ID),
		Index:  smpc.JoinIndex(tokensShare.Index),
		Shares: shamir.Shares{tokensShare},
	}
	tokensJoin.ID[31] = ResolveStageTokens

	err := matcher.smpcer.Join(networkID, tokensJoin, func(joinID smpc.JoinID, values []uint64) {
		if len(values) != 1 {
			logger.Compute(logger.LevelError, fmt.Sprintf("cannot resolve tokens: unexpected number of values: %v", len(values)))
			return
		}
		if isEqualToZero(values[0]) {
			com.State = ComputationStateMatched
			com.Match = true
			if err := matcher.storer.InsertComputation(com); err != nil {
				logger.Compute(logger.LevelError, fmt.Sprintf("cannot store matched computation buy = %v, sell = %v", com.Buy, com.Sell))
			}
			logger.Compute(logger.LevelDebug, fmt.Sprintf("✔ tokens => buy = %v, sell = %v", com.Buy, com.Sell))
			callback(com)
			return
		}
		com.State = ComputationStateMismatched
		com.Match = false
		if err := matcher.storer.InsertComputation(com); err != nil {
			logger.Compute(logger.LevelError, fmt.Sprintf("cannot store mismatched computation buy = %v, sell = %v", com.Buy, com.Sell))
		}
		logger.Compute(logger.LevelDebugHigh, fmt.Sprintf("✗ tokens => mismatch buy = %v, sell = %v", com.Buy, com.Sell))
		callback(com)
	})
	if err != nil {
		logger.Compute(logger.LevelError, fmt.Sprintf("cannot join tokens: %v", err))
	}
}

func isGreaterThanOrEqualToZero(value uint64) bool {
	return value >= 0 && value < shamir.Prime/2
}

func isGreaterThanZero(value uint64) bool {
	return value > 0 && value < shamir.Prime/2
}

func isEqualToZero(value uint64) bool {
	return value == 0 || value == shamir.Prime
}
