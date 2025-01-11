package schema

import (
	"cendit.io/garage/primer/enum"
	"cendit.io/garage/schema"
	"github.com/uptrace/bun"
)

type KYC struct {
	bun.BaseModel `bun:"table:kycs" rsf:"false"`

	ID            string            `json:"id" bun:"id"`           // Required
	UserID        string            `json:"user_id" bun:"user_id"` // Required
	BVN           string            `json:"bvn"`                   // Not required
	IDDocument    enum.IDDocument   `json:"id_document"`           // Not required, Enum
	IDURL         string            `json:"id_url"`                // Not required
	Status        enum.KYCStatus    `json:"status"`                // Enum
	FaceMap       string            `json:"face_map"`              // Not required
	MaritalStatus string            `json:"marital_status"`        // Not required
	NextOfKin     *schema.NextOfKin `json:"next_of_kin"`           // Not required
	Address       *Address          `json:"address"`               // Non-modifiable after confirmation
	EarningRange  string            `json:"earning_range"`         // Not required
	IncomeRange   string            `json:"income_range"`          // Not required
	WorkIndustry  string            `json:"work_industry"`         // Not required
	UseCase       string            `json:"use_case"`              // Not required

	schema.Datum
}
