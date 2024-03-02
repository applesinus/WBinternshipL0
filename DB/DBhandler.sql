create procedure create_DB() as $$

create schema wb;

create table wb.order (
    order_uid varchar(255) not null,
    track_number varchar(255) not null,
    entry varchar(255) not null,
    delivery_id int not null,
    payment_id varchar(255) not null,

    locate varchar(255) not null,
    internal_signature varchar(255) not null,
    customer_id varchar(255) not null,
    delivery_service varchar(255) not null,
    shardkey varchar(255) not null,
    sm_id int not null,
    date_created varchar(255) not null,
    oof_shard varchar(255) not null,

    foreign key (delivery_id) references wb.delivery (delivery_id),
    foreign key (payment_id) references wb.payment (payment_id),

    primary key (order_uid)
)

create table wb.delivery (
    delivery_id numeric not null,
    name varchar(255) not null,
    phone varchar(255) not null,
    zip varchar(255) not null,
    city varchar(255) not null,
    address varchar(255) not null,
    region varchar(255) not null,
    email varchar(255) not null,
    primary key (delivery_id)
)

create table wb.payment (
    transaction varchar(255) not null,
    request_id varchar(255) not null,
    currency varchar(255) not null,
    provider varchar(255) not null,
    amount int not null,
    payment_dt bigint not null,
    bank varchar(255) not null,
    delivery_cost int not null,
    goods_total int not null,
    custom_fee int not null,
    primary key (transaction)
)

create table wb.item (
    chrt_id int not null,
    order_uid varchar(255) not null,
    track_number varchar(255) not null,
    price int not null,
    rid varchar(255) not null,
    name varchar(255) not null,
    sale int not null,
    size varchar(255) not null,
    total_price int not null,
    nm_id int not null,
    brand varchar(255) not null,
    status int not null,
    primary key (chrt_id),
    foreign key (order_uid) references wb.order (order_uid)
)
$$
language plpgsql

create procedure drop_DB() as $$

drop schema wb cascade
$$
language plpgsql