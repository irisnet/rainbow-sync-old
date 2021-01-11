package gov

import (
	"github.com/irisnet/rainbow-sync/model"
	"time"
)

const (
	ProposalTypeSoftwareUpgrade       = "SoftwareUpgrade"
	ProposalTypeCancelSoftwareUpgrade = "CancelSoftwareUpgrade"
	ProposalTypeParameterChange       = "ParameterChange"
	ProposalTypeCommunityPoolSpend    = "CommunityPoolSpend"
	ProposalTypeText                  = "Text"
	ProposalTypeClientUpdate          = "ClientUpdate"
)

type (
	ContentTextProposal struct {
		Title       string `json:"title" bson:"title"`
		Description string `json:"description" bson:"description"`
	}
	ContentParameterChangeProposal struct {
		Title       string        `json:"title" bson:"title"`
		Description string        `json:"description" bson:"description"`
		Changes     []ParamChange `json:"changes" bson:"changes"`
	}
	ParamChange struct {
		Subspace string `json:"subspace" bson:"subspace"`
		Key      string `json:"key" bson:"key"`
		Value    string `json:"value" bson:"value"`
	}
	ContentCommunityPoolSpendProposal struct {
		Title       string       `json:"title" bson:"title"`
		Description string       `json:"description" bson:"description"`
		Recipient   string       `json:"recipient" bson:"recipient"`
		Amount      []model.Coin `json:"amount" bson:"amount"`
	}
	ContentSoftwareUpgradeProposal struct {
		Title       string `json:"title" bson:"title"`
		Description string `json:"description" bson:"description"`
		Plan        Plan   `json:"plan" bson:"plan"`
	}
	Plan struct {
		Name                string    `json:"name" bson:"name"`
		Time                time.Time `json:"time" bson:"time"`
		Height              int64     `json:"height" bson:"height"`
		Info                string    `json:"info" bson:"info"`
		UpgradedClientState string    `json:"upgraded_client_state" bson:"upgraded_client_state"`
	}
	ContentCancelSoftwareUpgradeProposal struct {
		Title       string `json:"title" bson:"title"`
		Description string `json:"description" bson:"description"`
	}
	ContentClientUpdateProposal struct {
		Title       string `json:"title" bson:"title"`
		Description string `json:"description" bson:"description"`
		ClientId    string `json:"client_id" bson:"client_id"`
		Header      string `json:"header" bson:"header"`
	}
)
