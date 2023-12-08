DROP TABLE IF EXISTS addresses;
DROP TABLE IF EXISTS customers;

CREATE TABLE customers (
  id                  serial               PRIMARY KEY,
	first_name          varchar(20)          NOT NULL,
  last_name           varchar(20)          NOT NULL,
  phone_number        varchar(20)          UNIQUE NOT NULL,
  email               varchar(60)          UNIQUE NOT NULL,
  password            varchar(72)          NOT NULL
  );

CREATE TABLE addresses (
  id                  serial               PRIMARY KEY,
  customer_id         int                  REFERENCES customers(id),
  lat                 numeric(20, 17)      NOT NULL,
  lon                 numeric(20, 17)      NOT NULL,
  address_line1       varchar(100)         NOT NULL,
  address_line2       varchar(100)                 ,
  city                varchar(70)          NOT NULL,
  country             varchar(60)          NOT NULL
  );