---
title: Basic Syntax of Postgres
date: 2026-10-26
tags:
  - postgres
  - database
---

# Overview

This note covers basic PostgreSQL syntax.

| # | Section |
| --- | --- |
| 1 | [Data Types](#data-types) |
| 2 | [PSQL Client](#psql---postgres-client) |
| 3 | [New User and Database Creation](#new-user-and-database-creation) |
| 4 | [Delete User and Database](#delete-user-and-database) |
| 5 | [Create a New Table](#create-a-new-table) |
| 6 | [Add a New Column to an Existing Table](#add-a-new-column-to-an-existing-table) |
| 7 | [Table Information Using psql](#table-information-using-psql) |
| 8 | [Remove a Column from a Table](#remove-a-column-from-a-table) |
| 9 | [Remove a Table](#remove-a-table) |
| 10 | [Query](#query)

# Data Types

| Type | Meaning |
| --- | --- |
| `SERIAL` | 32-bit integer that increases automatically (auto-increment) |
| `INT` / `INTEGER` | 32-bit integer (the two names are the same type) |
| `BIGINT` | 64-bit integer |
| `REAL` | 32-bit float |
| `DOUBLE PRECISION` | 64-bit float |
| `CHAR(n)` | Fixed length of n characters, padded with spaces if shorter |
| `VARCHAR(n)` | Variable length, up to n characters |
| `TEXT` | Variable length, no declared limit |
| `BOOLEAN` | `true` / `false` |
| `TIMESTAMP` | Date and time, no time zone (same as `TIMESTAMP WITHOUT TIME ZONE`) |
| `TIMESTAMPTZ` | Date and time with time zone (same as `TIMESTAMP WITH TIME ZONE`) |
| `DATE` | Calendar date only |
| `NUMERIC(15, 2)` | Exact decimal: up to 15 digits in total, 2 after the decimal point |
| `INT[]` | Array of integers (any type can be used, e.g. `TEXT[]`) |
| `JSON` | JSON stored as plain text, exactly as written |
| `JSONB` | JSON stored in binary form, faster to query and can be indexed |
| `SMALLINT`  |	16-bit integer |
| `BIGSERIAL` |	64-bit auto-increment (use with BIGINT) |
| `UUID` |	Universally unique ID, often used as a primary key with default gen_random_uuid() |
| `TIME` | Time of day only |
| `INTERVAL` |	A span of time, like '3 days' |
| `BYTEA` |	Raw binary data |

# PSQL - Postgres Client

Connect as your own user:

```bash
$ whoami
labib0x9
$ psql -U labib0x9 -d postgres

psql (18.3 (Homebrew))
Type "help" for help.

postgres=#
```

# New User and Database Creation

```bash
postgres=# CREATE ROLE new_user WITH LOGIN PASSWORD 'my_password';
CREATE ROLE

postgres=# CREATE DATABASE new_database OWNER new_user;
CREATE DATABASE

postgres=# GRANT ALL PRIVILEGES ON DATABASE new_database TO new_user;
GRANT
```

Then log in as the new user:

```bash
labib0x9@sarona (main) $ psql -U new_user -d new_database
psql (18.3 (Homebrew))
Type "help" for help.

new_database=>
```

# Delete User and Database

Drop the database first, then the user. A user who still owns a database cannot be dropped. You also cannot drop the database you are currently connected to, so connect to `postgres` first.

```bash
postgres=# DROP DATABASE new_database;
DROP DATABASE

postgres=# DROP DATABASE IF EXISTS new_database;
-- safe version: no error if the database does not exist

postgres=# DROP ROLE new_user;
DROP ROLE

postgres=# DROP ROLE IF EXISTS new_user;
-- safe version: no error if the role does not exist
```

```sql
alter role new_user with password 'new_password';
alter role new_user with createdb;
revoke all privileges on database new_database from new_user;
```

# Create a New Table

```sql
create table users (
    id serial primary key,
    email varchar(100) unique not null,
    name text not null,
    age int check (age >= 18),
    created_at timestamptz default now()
);

create table if not exists books (
    title text,
    author_id int references users(id) on delete cascade,
    category text,
    created_at timestamptz default now(),
    unique (title, author_id)
);

create table enrollments (
    student_id int references students(id) on delete cascade,
    course_id int references courses(id) on delete cascade,
    enrolled_at timestamptz default now(),
    primary key (student_id, course_id)
);

create table posts (
    id serial primary key,
    author_id int references users(id) on delete set null,   -- keep the post, clear author_id
    editor_id int references users(id) on delete restrict,   -- block deleting the user
    title text not null
);

create table products (
    id serial primary key,
    name text not null,
    price numeric(10, 2) not null,
    sale_price numeric(10, 2),
    constraint price_positive check (price > 0),
    constraint sale_less_than_price check (sale_price is null or sale_price < price),
    constraint uq_product_name unique (name, price)
);

create type order_status as enum ('pending', 'paid', 'shipped', 'cancelled');

create table shipments (
    id serial primary key,
    status order_status not null default 'pending'
);

create table articles (
    id serial primary key,
    tags text[] default '{}',
    metadata jsonb not null default '{}'::jsonb
);


create table categories (
    id serial primary key,
    name text not null,
    parent_id int references categories(id) on delete cascade
);

create table rectangles (
    width int not null,
    height int not null,
    area int generated always as (width * height) stored
);

create table users_backup (like users including all);

create table adult_users as
select * from users where age >= 18;

create temp table scratch (
    id int,
    note text
);


```

# Add a New Column to an Existing Table

```sql
alter table users add column phone varchar(20);
-- safe version: no error if the column already exists
alter table users add column if not exists phone varchar(20);
-- with a default value and a constraint
alter table users add column is_active boolean not null default true;
```

```sql
alter table users rename to customers;
alter table users rename column name to full_name;
alter table users alter column age type bigint;
alter table users alter column age set default 18;
alter table users alter column age drop default;
alter table users alter column name set not null;
alter table users alter column name drop not null;
alter table users add constraint age_adult check (age >= 18);
alter table users drop constraint age_adult;
```

# Remove a Column from a Table

```sql
alter table users drop column created_at;
-- safe version: no error if the column does not exist
alter table users drop column if exists created_at;
```

# Remove a Table

```sql
drop table users;
-- safe version: no error if the table does not exist
drop table if exists users;
-- also drops objects that depend on it (e.g. the foreign key in books)
drop table if exists users cascade;
drop database new_database with (force);

```

```sql
truncate table users;
truncate table users restart identity cascade;
```

# Table Information Using psql

```
\l              list databases
\c new_database connect to another database
\du             list users (roles)
\dt             list tables
\d users        describe a table
\d+ users       describe with extra details (size, storage)
\conninfo       show current connection
\q              quit psql
\?              help for psql commands
```

# Query