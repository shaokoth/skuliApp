package inventory

import "time"

// Item is a stock-tracked asset or consumable owned by a school.
type Item struct {
	ID        int64     `json:"id"`
	SchoolID  int64     `json:"school_id"`
	Name      string    `json:"name"`
	Category  string    `json:"category"`
	SKU       string    `json:"sku"`
	Quantity  int       `json:"quantity"`
	Unit      string    `json:"unit"`
	UnitCost  int64     `json:"unit_cost"` // minor units
	Location  string    `json:"location"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Version   int32     `json:"-"`
}
