DROP TABLE IF EXISTS restaurants;
DROP TABLE IF EXISTS tickets;

CREATE TABLE tickets (
  id                        serial        PRIMARY KEY,
  order_id                  int           UNIQUE NOT NULL,
  restaurant_id             int           NOT NULL,
  ETC                       timestamp     NOT NULL, -- Estimated Time to Complete
  RTA                       timestamp             , -- Requested Time of Arival
  items                     integer[]     NOT NULL,
  order_total               numeric(8, 2) NOT NULL,
  status                    int           NOT NULL
  );

CREATE TABLE restaurants (
  id                  serial               PRIMARY KEY
  );