package dao

import (
	"context"
	"fmt"

	"titan-container-platform/core"

	"github.com/Masterminds/squirrel"
)

// CreateOrder creates a new order in the database.
func CreateOrder(ctx context.Context, order *core.Order) error {
	query := fmt.Sprintf(`INSERT INTO %s (id, account, cpu, ram, storage, duration, status)
			VALUES (:id, :account, :cpu, :ram, :storage, :duration, :status);`, orderInfoTable)
	_, err := mDB.NamedExec(query, order)

	return err
}

// LoadOrdersByStatus retrieves orders based on their status.
func LoadOrdersByStatus(status core.OrderStatus) ([]*core.Order, error) {
	var infos []*core.Order

	query := fmt.Sprintf("SELECT * FROM %s WHERE status=?", orderInfoTable)
	err := mDB.Select(&infos, query, status)
	if err != nil {
		return nil, err
	}

	return infos, nil

	// var infos []*core.Order
	// hQuery := fmt.Sprintf(`SELECT * FROM %s WHERE status in (?) `, orderInfoTable)
	// shQuery, args, err := sqlx.In(hQuery, statuses)
	// if err != nil {
	// 	return infos, err
	// }

	// shQuery = mDB.Rebind(shQuery)
	// err = mDB.Select(&infos, shQuery, args...)

	// return infos, err
}

// UpdateOrderStatus updates the status of an order in the database.
func UpdateOrderStatus(id string, status core.OrderStatus) error {
	query := fmt.Sprintf(`UPDATE %s SET status=? WHERE id=? `, orderInfoTable)
	_, err := mDB.Exec(query, status, id)

	return err
}

// LoadAccountOrdersByStatus retrieves a list of orders for a given account, with pagination.
func LoadAccountOrdersByStatus(ctx context.Context, account string, status core.OrderStatus, page, size int) ([]*core.Order, int64, error) {
	out := make([]*core.Order, 0)

	var count int64
	if page < 1 {
		page = 1
	}

	query, args, err := squirrel.Select("*").From(orderInfoTable).Where(squirrel.Eq{"account": account}).Where(squirrel.Eq{"status": status}).Offset(uint64((page - 1) * size)).Limit(uint64(size)).ToSql()
	if err != nil {
		return nil, 0, err
	}

	if err := mDB.SelectContext(ctx, &out, query, args...); err != nil {
		return nil, 0, err
	}

	sq2 := squirrel.Select("COUNT(*)").From(orderInfoTable).Where(squirrel.Eq{"account": account}).Where(squirrel.Eq{"status": status})

	query2, args2, err := sq2.ToSql()
	if err != nil {
		return nil, 0, err
	}

	err = mDB.Get(&count, query2, args2...)
	if err != nil {
		return nil, 0, err
	}

	return out, count, nil
}

// LoadAccountOrders retrieves a list of orders for a given account, with pagination.
func LoadAccountOrders(ctx context.Context, account string, page, size int) ([]*core.Order, int64, error) {
	out := make([]*core.Order, 0)

	var count int64
	if page < 1 {
		page = 1
	}

	query, args, err := squirrel.Select("*").From(orderInfoTable).Where(squirrel.Eq{"account": account}).Offset(uint64((page - 1) * size)).Limit(uint64(size)).ToSql()
	if err != nil {
		return nil, 0, err
	}

	if err := mDB.SelectContext(ctx, &out, query, args...); err != nil {
		return nil, 0, err
	}

	sq2 := squirrel.Select("COUNT(*)").From(orderInfoTable).Where(squirrel.Eq{"account": account})

	query2, args2, err := sq2.ToSql()
	if err != nil {
		return nil, 0, err
	}

	err = mDB.Get(&count, query2, args2...)
	if err != nil {
		return nil, 0, err
	}

	return out, count, nil
}
