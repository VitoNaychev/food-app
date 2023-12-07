DROP TABLE IF EXISTS working_hours;
DROP TABLE IF EXISTS menu_items;
DROP TABLE IF EXISTS restaurants;
DROP TABLE IF EXISTS addresses;



CREATE TABLE restaurants (
  id                  serial               PRIMARY KEY,
  name                varchar(30)          NOT NULL,
  phone_number        varchar(16)          UNIQUE NOT NULL,
  email               varchar(30)          UNIQUE NOT NULL,
  password            varchar(72)          NOT NULL,
  IBAN                varchar(34)          UNIQUE NOT NULL,
  status              int                  NOT NULL
  );

CREATE TABLE addresses (
  id                  serial                PRIMARY KEY,
  restaurant_id       int                   NOT NULL      UNIQUE    REFERENCES restaurants(id),
  lat                 numeric(20, 17)       NOT NULL,
  lon                 numeric(20, 17)       NOT NULL,
  address_line1       varchar(100)          NOT NULL,
  address_line2       varchar(100)                 ,
  city                varchar(40)           NOT NULL,
  country             varchar(40)           NOT NULL
  );

    
CREATE TABLE working_hours (
  id                  serial               PRIMARY KEY,
  restaurant_id       int                  NOT NULL                 REFERENCES restaurants(id),
  day                 int                  NOT NULL,
  opening             time                 NOT NULL,
  closing             time                 NOT NULL
  );
  
CREATE TABLE menu_items (
  id                  serial               PRIMARY KEY,
  restaurant_id       int                  REFERENCES restaurants(id),
  name                varchar(20)          NOT NULL,
  price               numeric(6, 2)        NOT NULL,
  details             text                 
  );