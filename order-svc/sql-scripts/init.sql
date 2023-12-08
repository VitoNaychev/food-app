DROP TABLE IF EXISTS orders;
DROP TABLE IF EXISTS addresses;
DROP TYPE IF EXISTS status;

CREATE TABLE addresses (
  id                  serial               PRIMARY KEY,
  lat                 numeric(20, 17)      NOT NULL,
  lon                 numeric(20, 17)      NOT NULL,
  address_line1       varchar(100)         NOT NULL,
  address_line2       varchar(100)                 ,
  city                varchar(70)          NOT NULL,
  country             varchar(60)          NOT NULL
  );
  
-- CREATE TYPE status AS ENUM ('APPROVAL_PENDING', 'REJECTED', 'DECLINED', 'APPROVED', 'CANCELED', 'PREPARING', 'PREPARED', 'PICKED_UP', 'COMPLETED');

CREATE TABLE orders (
	id                        serial        PRIMARY KEY,
	customer_id               int           UNIQUE NOT NULL,
  restaurant_id             int           UNIQUE NOT NULL,
  items                     integer[]     NOT NULL,
  total                     numeric(7, 2) NOT NULL,
  delivery_time             timestamp             ,
  status                    int           NOT NULL,
  pickup_address            int           NOT NULL          REFERENCES addresses(id),
  delivery_address          int           NOT NULL          REFERENCES addresses(id)
);