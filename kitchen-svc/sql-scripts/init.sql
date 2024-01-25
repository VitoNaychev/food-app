DROP TABLE IF EXISTS ticket_items;
DROP TABLE IF EXISTS tickets;
DROP TABLE IF EXISTS menu_items;
DROP TABLE IF EXISTS ticket_items;

CREATE TABLE restaurants (
    id         int       PRIMARY KEY
);


CREATE TABLE menu_items (
    id              int                  PRIMARY KEY,
    restaurant_id   int                  REFERENCES restaurants(id),
    name            varchar(20)          NOT NULL,
    price           numeric(6, 2)        NOT NULL
);

CREATE TABLE tickets (
    id                        serial                       PRIMARY KEY,
    restaurant_id             int                          REFERENCES restaurants(id),
    total                     numeric(8, 2)                NOT NULL,
    state                     int                          NOT NULL,
    ready_by                  timestamp with time zone     NOT NULL
);

CREATE TABLE ticket_items (
    id           serial    PRIMARY KEY,
    ticket_id    int       REFERENCES tickets(id),
    menu_item_id int       REFERENCES menu_items(id),
    quantity     int       NOT NULL
);


