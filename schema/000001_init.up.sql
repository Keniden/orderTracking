CREATE TABLE users (
    id serial PRIMARY KEY,
    name varchar(255) NOT NULL,
    username varchar(255) NOT NULL UNIQUE,
    password_hash varchar(255) NOT NULL
);

CREATE TABLE order_lists(
    id serial PRIMARY KEY,
    title varchar(255) NOT NULL,
    description varchar(255)
);

CREATE TABLE users_lists(
    id serial PRIMARY KEY,
    user_id int REFERENCES users (id) ON DELETE CASCADE NOT NULL,
    list_id int REFERENCES order_lists (id) ON DELETE CASCADE NOT NULL
);

CREATE TABLE order_items (
    id serial PRIMARY KEY,
    title varchar(255) NOT NULL,
    description varchar(255),
    done boolean NOT NULL DEFAULT false
);

CREATE TABLE lists_items (
    id serial PRIMARY KEY,
    user_id int REFERENCES users (id) ON DELETE CASCADE NOT NULL,
    list_id int REFERENCES order_lists (id) ON DELETE CASCADE NOT NULL
);