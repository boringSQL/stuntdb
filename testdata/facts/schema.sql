create table customers(id bigint primary key, email text not null);
create table orders(id bigint primary key, customer_id bigint not null references customers, total numeric(12,2) not null);
create table events(id bigint primary key, customer_id bigint not null references customers, kind text not null);
