create or replace procedure create_DB() as $$
begin
    create schema if not exists wb;

    create table if not exists wb.order (
        order_info jsonb not null,
        order_id serial primary key
    );
end;
$$
language plpgsql;

create or replace procedure drop_DB() as $$
begin
    drop schema wb cascade;
end;
$$
language plpgsql;

/*
call create_DB();
call drop_DB();

select * from wb.order;

drop procedure create_DB();
drop procedure drop_DB();
*/