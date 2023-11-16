DROP TABLE IF EXISTS orders;
DROP TABLE IF EXISTS addresses;
DROP TYPE IF EXISTS status;

CREATE TABLE addresses (
  id                  serial               PRIMARY KEY,
  lat                 numeric(10, 7)       NOT NULL,
  lon                 numeric(10, 7)       NOT NULL,
  address_line1       varchar(100)         NOT NULL,
  address_line2       varchar(100)                 ,
  city                varchar(40)          NOT NULL,
  country             varchar(40)          NOT NULL
  );
  
CREATE TYPE status AS ENUM ('APPROVAL_PENDING', 'REJECTED', 'DECLINED', 'APPROVED', 'CANCELED', 'PREPARING', 'PREPARED', 'PICKED_UP', 'COMPLETED');

CREATE TABLE orders (
	id                        serial        PRIMARY KEY,
	customer_id               int           UNIQUE NOT NULL,
  restaurant_id             int           UNIQUE NOT NULL,
  address_id                int           REFERENCES addresses(id),
  billing_id                int           UNIQUE NOT NULL,
  items                     integer[]     NOT NULL,
  total                     numeric(7, 2) NOT NULL,
  requested_delivery_time   timestamp             ,
  status                    status        NOT NULL
);