package msgparser

import (
	"fmt"
	"github.com/irisnet/rainbow-sync/lib/logger"
	"github.com/kaifei-bianjie/msg-parser/types"
	"regexp"
)

var (
	// IsAlphaNumeric defines a regular expression for matching against alpha-numeric
	// values.
	IsAlphaNumeric = regexp.MustCompile(`^[a-zA-Z0-9]+$`).MatchString
)

type Handler func(v types.SdkMsg) CustomMsgDocInfo

var _ Router = (*router)(nil)

// Router implements a msg-parser Handler router.
type Router interface {
	AddRoute(r string, h Handler) (rtr Router)
	HasRoute(r string) bool
	GetRoute(path string) (h Handler, err error)
	GetRoutesLen() int
}

type router struct {
	routes map[string]Handler
}

// NewRouter creates a new Router interface instance
func NewRouter() Router {
	return &router{
		routes: make(map[string]Handler),
	}
}

// AddRoute adds a governance handler for a given path. It returns the Router
// so AddRoute calls can be linked. It will panic if the router is sealed.
func (rtr *router) AddRoute(path string, h Handler) Router {

	if !IsAlphaNumeric(path) {
		logger.Warn("addroute failed for route expressions can only contain alphanumeric characters")
		return rtr
	}
	if rtr.HasRoute(path) {
		logger.Warn(fmt.Sprintf("route %s has already been initialized", path))
		return rtr
	}

	rtr.routes[path] = h
	return rtr
}

// HasRoute returns true if the router has a path registered or false otherwise.
func (rtr *router) HasRoute(path string) bool {
	return rtr.routes[path] != nil
}

// GetRoute returns a Handler for a given path.
func (rtr *router) GetRoute(path string) (Handler, error) {
	if !rtr.HasRoute(path) {
		return nil, fmt.Errorf("route \"%s\" does not exist", path)
	}

	return rtr.routes[path], nil
}

func (rtr *router) GetRoutesLen() int {
	return len(rtr.routes)
}

func RegisteRouter() Router {
	msgRoute := NewRouter()
	msgRoute.AddRoute(BankRouteKey, handleBank).
		AddRoute(StakingRouteKey, handleStaking).
		AddRoute(DistributionRouteKey, handleDistribution).
		AddRoute(CrisisRouteKey, handleCrisis).
		AddRoute(EvidenceRouteKey, handleEvidence).
		AddRoute(GovRouteKey, handleGov).
		AddRoute(SlashingRouteKey, handleSlashing).
		AddRoute(IbcRouteKey, handleIbc).
		AddRoute(IbcTransferRouteKey, handleIbc).
		AddRoute(NftRouteKey, handleNft).
		AddRoute(ServiceRouteKey, handleService).
		AddRoute(TokenRouteKey, handleToken).
		AddRoute(HtlcRouteKey, handleHtlc).
		AddRoute(CoinswapRouteKey, handleCoinswap).
		AddRoute(RandomRouteKey, handleRandom).
		AddRoute(OracleRouteKey, handleOracle).
		AddRoute(RecordRouteKey, handleRecord).
		AddRoute(FarmRouteKey, handleFarm).
		AddRoute(TIbcRouteKey, handleTIbc).
		AddRoute(TIbcTransferRouteKey, handleTIbc)
	return msgRoute
}
