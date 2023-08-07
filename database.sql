/**
  This is the SQL script that will be used to initialize the database schema.
  We will evaluate you based on how well you design your database.
  1. How you design the tables.
  2. How you choose the data types and keys.
  3. How you name the fields.
  In this assignment we will use PostgreSQL as the database.
  */

/** This is test table. Remove this table and replace with your own tables. */
CREATE TABLE IF NOT EXISTS users (
  id UUID PRIMARY KEY,
  full_name VARCHAR(60) NOT NULL,
  phone_number VARCHAR(13) UNIQUE NOT NULL,
  password_hash TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS login_logs (
  id UUID PRIMARY KEY,
  user_id UUID REFERENCES users(id) NOT NULL,
  login_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX users_id_index ON users(id);
CREATE INDEX users_phone_number_index ON users(phone_number);
