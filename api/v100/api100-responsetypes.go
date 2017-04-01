package api100

import (
	"time"

	"github.com/chaosvermittlung/funkloch-server/db/v100"
)

type authResponse struct {
	Token string `json:"token"`
}

type storeItemCountResponse struct {
	Name  string
	Count int
}

type equipmentCountResponse struct {
	Equipment db100.Equipment
	Store     db100.Store
	Count     int
}

type eventParticipiantsResponse struct {
	User      db100.User
	Arrival   time.Time
	Departure time.Time
}

type wishlistItemsResponse struct {
	Equipment db100.Equipment
	Count     int
}

type packinglistItemsResponse struct {
	StoreItemID int
	Equipment   db100.Equipment
	Store       db100.Store
}
