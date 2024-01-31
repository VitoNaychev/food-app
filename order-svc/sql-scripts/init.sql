DROP TABLE IF EXISTS order_items;
DROP TABLE IF EXISTS orders;
DROP TABLE IF EXISTS addresses;

CREATE TABLE addresses (
  id                  serial               PRIMARY KEY,
  lat                 numeric(20, 17)      NOT NULL,
  lon                 numeric(20, 17)      NOT NULL,
  address_line1       varchar(100)         NOT NULL,
  address_line2       varchar(100)                 ,
  city                varchar(70)          NOT NULL,
  country             varchar(60)          NOT NULL
  );
  
CREATE TABLE orders (
	id                        serial        PRIMARY KEY,
	customer_id               int           NOT NULL,
  restaurant_id             int           NOT NULL,
  total                     numeric(7, 2) NOT NULL,
  status                    int           NOT NULL,
  pickup_address            int           NOT NULL          REFERENCES addresses(id),
  delivery_address          int           NOT NULL          REFERENCES addresses(id)
);

CREATE TABLE order_items (
	id                 serial        PRIMARY KEY,
	order_id           int           NOT NULL       REFERENCES orders(id),
  menu_item_id       int           NOT NULL,
  quantity           int           NOT NULL
);