-- 1 Order Service
CREATE TABLE IF NOT EXISTS orders(
    id bigserial PRIMARY KEY,
    customername VARCHAR(50) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',  --ending, completed, cancelled
    isdeleted BOOLEAN DEFAULT 'FALSE',
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS order_items (
    orderID INT,
    productID INT,
    quantity INT NOT NULL CHECK(quantity > 0),
    PRIMARY KEY (orderID, productID),
    FOREIGN KEY (orderID) REFERENCES orders(id) ON DELETE CASCADE
);


