package schema

import (
	"cendit.io/garage/primer/enum"
	"github.com/uptrace/bun"
)

type PlatformConfig struct {
	bun.BaseModel `bun:"table:platform_configs"`

	Datum

	Tiers             KYCTier        `json:"tiers" bun:"tiers"`
	Rates             Rate           `json:"rates" bun:"rates"`
	FreeTransfer      int            `json:"free_transfer" bun:"free_transfer"`
	FreeTransferCycle enum.Frequency `json:"free_transfer_cycle" bun:"free_transfer_cycle"`
}

type Rate struct {
	// TODO define rates

	bun.BaseModel `bun:"table:rates"`

	Datum
}
