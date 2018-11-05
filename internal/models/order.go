package models

import (
	"hameid.net/cdex/dex/internal/store"
	"hameid.net/cdex/dex/internal/wrappers"
)

// Order record
type Order struct {
	Hash         *wrappers.Hash    `json:"order_hash"`
	Token        *wrappers.Address `json:"token"`
	Base         *wrappers.Address `json:"base"`
	Price        *wrappers.BigInt  `json:"price"`
	Quantity     *wrappers.BigInt  `json:"quantity"`
	IsBid        bool              `json:"is_bid"`
	CreatedAt    uint64            `json:"created_at"`
	CreatedBy    *wrappers.Address `json:"created_by"`
	Volume       *wrappers.BigInt  `json:"volume"`
	VolumeFilled *wrappers.BigInt  `json:"volume_filled"`
	IsOpen       bool              `json:"is_open"`
}

// Save inserts Order
func (order *Order) Save(store *store.DataStore) error {
	query := `INSERT INTO orders (
		order_hash, token, base, price, quantity, is_bid, created_at, created_by, volume)
		VALUES ($1, $2, $3, $4, $5, $6, to_timestamp($7), $8, $9)`

	_, err := store.DB.Exec(
		query,
		order.Hash,
		order.Token,
		order.Base,
		order.Price.String(),
		order.Quantity.String(),
		order.IsBid,
		order.CreatedAt,
		order.CreatedBy,
		order.Volume.String(),
	)

	return err
}

// Get scans the order by hash from database
func (order *Order) Get(store *store.DataStore) error {
	query := `SELECT order_hash, token, base, price, quantity, is_bid, trunc(extract(epoch from created_at::timestamp with time zone)), created_by, volume, volume_filled FROM orders WHERE order_hash=$1`

	row := store.DB.QueryRow(query, order.Hash)

	var t float64

	err := row.Scan(
		&order.Hash,
		&order.Token,
		&order.Base,
		&order.Price,
		&order.Quantity,
		&order.IsBid,
		// &order.CreatedAt,
		&t,
		&order.CreatedBy,
		&order.Volume,
		&order.VolumeFilled,
	)

	order.CreatedAt = uint64(t)

	return err
}

// Update order details
func (order *Order) Update(store *store.DataStore) error {
	query := `UPDATE orders SET volume_filled=$2, is_open=$3 WHERE order_hash=$1`

	if order.Volume.Cmp(&order.VolumeFilled.Int) == 0 {
		order.IsOpen = false
	}

	_, err := store.DB.Exec(
		query,
		order.Hash,
		order.VolumeFilled.String(),
		order.IsOpen,
	)

	return err
}

// StoreFilledVolume updates order's filled volume
func (order *Order) StoreFilledVolume(store *store.DataStore) error {
	query := `UPDATE orders SET volume_filled=$1 WHERE order_hash=$2`

	_, err := store.DB.Exec(
		query,
		order.VolumeFilled.String(),
		order.Hash,
	)

	return err
}

// Close updates order's filled volume
func (order *Order) Close(store *store.DataStore) error {
	query := `UPDATE orders SET is_open=$1 WHERE order_hash=$2`

	_, err := store.DB.Exec(
		query,
		false,
		order.Hash,
	)

	return err
}