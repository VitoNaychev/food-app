DROP TABLE IF EXISTS couriers;

CREATE TABLE couriers (
  id                  serial               PRIMARY KEY,
  first_name          varchar(20)          NOT NULL,
  last_name           varchar(20)          NOT NULL,
  phone_number        varchar(20)          NOT NULL,
  email               varchar(60)          NOT NULL,
  password            varchar(72)          NOT NULL,
  IBAN                varchar(34)          NOT NULL
  );
  