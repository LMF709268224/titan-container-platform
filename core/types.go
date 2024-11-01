package core

import "time"

// User represents a user in the system.
type User struct {
	Account   string    `db:"account" json:"account"`
	Avatar    string    `db:"avatar" json:"avatar"`
	Username  string    `db:"user_name" json:"user_name"`
	UserEmail string    `db:"user_email" json:"user_email"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

// ResponseUser represents a user response structure.
type ResponseUser struct {
	Account   string    `db:"account" json:"account"`
	Username  string    `db:"user_name" json:"user_name"`
	UserEmail string    `db:"user_email" json:"user_email"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

// OrderReq represents a request to create or get price an order with specified resources.
type OrderReq struct {
	CPUCores    int `json:"cpu"`
	RAMSize     int `json:"ram"`      // in GB
	StorageSize int `json:"storage"`  // in GB
	Duration    int `json:"duration"` // Hour
}

// Order represents a customer's order in the system.
type Order struct {
	ID          string      `db:"id" json:"id"`
	Account     string      `db:"account" json:"account"`
	CPUCores    int         `db:"cpu" json:"cpu"`
	RAMSize     int         `db:"ram" json:"ram"`
	StorageSize int         `db:"storage" json:"storage"`
	Duration    int         `db:"duration" json:"duration"` // Hour
	Status      OrderStatus `db:"status" json:"status"`
	CreatedAt   time.Time   `db:"created_at" json:"created_at"`
}

// OrderStatus represents the status of an order.
type OrderStatus int

const (
	// OrderStatusCreated indicates that the order has been created.
	OrderStatusCreated OrderStatus = iota
	// OrderStatusPaid indicates that the order has been paid.
	OrderStatusPaid
	// OrderStatusDone indicates that the order has been completed.
	OrderStatusDone
	// OrderStatusExpired indicates that the order has expired.
	OrderStatusExpired
	// OrderStatusFailed indicates that the order has creation failed.
	OrderStatusFailed
	// OrderStatusTimeout indicates that the order has payment timeout.
	OrderStatusTimeout
)
