DROP TABLE IF EXISTS deliveries;
DROP TABLE IF EXISTS addresses;
DROP TABLE IF EXISTS locations;

DROP TABLE IF EXISTS couriers;

CREATE TABLE couriers (
    id                  int                  PRIMARY KEY,
    name                varchar(20)          NOT NULL
);

CREATE TABLE addresses (
  id                  serial               PRIMARY KEY,
  lat                 numeric(10, 7)       NOT NULL,
  lon                 numeric(10, 7)       NOT NULL,
  address_line1       varchar(100)         NOT NULL,
  address_line2       varchar(100)                 ,
  city                varchar(40)          NOT NULL,
  country             varchar(40)          NOT NULL
  );

CREATE TABLE deliveries (
  id                  serial    PRIMARY KEY,
  courier_id          int       REFERENCES couriers(id),
  pickup_address_id   int       REFERENCES addresses(id),
  delivery_address_id int       REFERENCES addresses(id),
  prepared_by_time    timestamp NOT NULL,
  state               int       NOT NULL
  );

CREATE TABLE locations (
  id                   int             PRIMARY KEY,
  courier_id           int             NOT NULL,
  lat                  numeric(10, 7)  NOT NULL,
  lon                  numeric(10, 7)  NOT NULL
  );

