package postgres

import (
	"context"
	"fmt"
	"order-service/internal/models"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Order struct {
	db *pgxpool.Pool
}

func NewOrderRepository(db *pgxpool.Pool) *Order {
	return &Order{db: db}
}

func (r *Order) Create(ctx context.Context, order models.Order) (int64, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)

	query := `
		INSERT INTO orders (customername, status) 
		VALUES ($1, $2)
		RETURNING ID;
	`

	var orderID int64
	err = tx.QueryRow(ctx, query, order.CustomerName, order.Status).Scan(&orderID)
	if err != nil {
		return 0, err
	}

	// Inserting order items. in case when same product id is given, it check on conflict, if so it's just adding quantity for previus row.
	queryOrderItems := `
		INSERT INTO order_items (OrderID, ProductID, Quantity) VALUES
		($1, $2, $3)
		ON CONFLICT (OrderID, ProductID)
		DO UPDATE SET Quantity = order_items.Quantity + EXCLUDED.Quantity;
	`

	for _, v := range order.OrderItems {
		_, err = tx.Exec(ctx, queryOrderItems, orderID, v.ProductID, v.Quantity)
		if err != nil {
			return 0, err
		}
	}

	return orderID, tx.Commit(ctx)
}

func (r *Order) GetWithFilter(ctx context.Context, filter models.OrderFilter) (models.Order, error) {
	query := `
		SELECT id, customername, status, created_at 
		FROM orders 
		WHERE id = $1 AND isdeleted = FALSE
	`

	var order models.Order
	err := r.db.QueryRow(ctx, query, filter.ID).Scan(
		&order.ID,
		&order.CustomerName,
		&order.Status,
		&order.Created_at,
	)
	if err != nil {
		return models.Order{}, err
	}

	// Get order items
	itemsQuery := `
	SELECT orderID, productID, quantity 
	FROM order_items 
	WHERE orderID = $1
`

	rows, err := r.db.Query(ctx, itemsQuery, filter.ID)
	if err != nil {
		return models.Order{}, err
	}
	defer rows.Close()

	var orderItems []models.OrderItem
	for rows.Next() {
		var item models.OrderItem
		err := rows.Scan(&item.OrderID, &item.ProductID, &item.Quantity)
		if err != nil {
			return models.Order{}, err
		}
		orderItems = append(orderItems, item)
	}

	if err = rows.Err(); err != nil {
		return models.Order{}, err
	}

	order.OrderItems = orderItems
	return order, nil
}

func (r *Order) GetListWithFilter(ctx context.Context, filter models.OrderFilter) ([]models.Order, error) {
	// First, get all orders
	ordersQuery := `
        SELECT id, customername, status, created_at 
        FROM orders 
        WHERE isdeleted = FALSE
        ORDER BY created_at DESC
    `

	rows, err := r.db.Query(ctx, ordersQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var order models.Order
		err := rows.Scan(&order.ID, &order.CustomerName, &order.Status, &order.Created_at)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	// If no orders found, return empty slice
	if len(orders) == 0 {
		return orders, nil
	}

	// Get all order items for the fetched orders
	itemsQuery := `
        SELECT orderID, productID, quantity 
        FROM order_items 
        WHERE orderID = ANY($1)
    `

	// Prepare list of order IDs
	orderIDs := make([]int64, len(orders))
	for i, order := range orders {
		orderIDs[i] = order.ID
	}

	itemRows, err := r.db.Query(ctx, itemsQuery, orderIDs)
	if err != nil {
		return nil, err
	}
	defer itemRows.Close()

	// Create a map to store items for each order
	itemsMap := make(map[int64][]models.OrderItem)
	for itemRows.Next() {
		var item models.OrderItem
		err := itemRows.Scan(&item.OrderID, &item.ProductID, &item.Quantity)
		if err != nil {
			return nil, err
		}
		itemsMap[item.OrderID] = append(itemsMap[item.OrderID], item)
	}

	if err = itemRows.Err(); err != nil {
		return nil, err
	}

	// Assign order items to each order
	for i, order := range orders {
		if items, ok := itemsMap[order.ID]; ok {
			orders[i].OrderItems = items
		} else {
			orders[i].OrderItems = []models.OrderItem{} // empty slice if no items
		}
	}

	return orders, nil
}

func (r *Order) Update(ctx context.Context, update models.OrderUpdateData) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Build dynamic update query for order fields
	query := "UPDATE orders SET "
	params := []any{}
	paramCount := 1

	if update.CustomerName != nil {
		query += fmt.Sprintf("customername = $%d, ", paramCount)
		params = append(params, *update.CustomerName)
		paramCount++
	}

	if update.Status != nil {
		query += fmt.Sprintf("status = $%d, ", paramCount)
		params = append(params, *update.Status)
		paramCount++
	}

	if update.IsDeleted != nil {
		query += fmt.Sprintf("isdeleted = $%d, ", paramCount)
		params = append(params, *update.IsDeleted)
		paramCount++
	}

	// Remove trailing comma and add WHERE clause
	query = strings.TrimSuffix(query, ", ")
	query += fmt.Sprintf(" WHERE id = $%d", paramCount)
	params = append(params, update.ID)
	paramCount++

	// Execute order update
	_, err = tx.Exec(ctx, query, params...)
	if err != nil {
		return fmt.Errorf("failed to update order: %w", err)
	}

	// Handle order items update if provided
	if update.OrderItems != nil {
		// First delete existing items
		_, err = tx.Exec(ctx, "DELETE FROM order_items WHERE orderID = $1", update.ID)
		if err != nil {
			return fmt.Errorf("failed to clear order items: %w", err)
		}

		// Insert new items if any
		items := *update.OrderItems
		if len(items) > 0 {
			itemQuery := "INSERT INTO order_items (orderID, productID, quantity) VALUES "
			itemParams := []any{}
			itemParamCount := 1

			valueClauses := []string{}
			for _, item := range items {
				valueClause := fmt.Sprintf("($%d, $%d, $%d)",
					itemParamCount, itemParamCount+1, itemParamCount+2)
				valueClauses = append(valueClauses, valueClause)
				itemParams = append(itemParams, update.ID, item.ProductID, item.Quantity)
				itemParamCount += 3
			}

			itemQuery += strings.Join(valueClauses, ", ")
			_, err = tx.Exec(ctx, itemQuery, itemParams...)
			if err != nil {
				return fmt.Errorf("failed to insert order items: %w", err)
			}
		}
	}

	return tx.Commit(ctx)
}
