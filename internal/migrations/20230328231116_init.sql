-- +goose Up
create table public.files_meta_data (
    user_id uuid not null primary key ,
    server_ip varchar(255) not null,
    file_name  varchar(255) not null,
    file_part_name uuid not null,
    file_part_number smallint not null,
    created_at timestamp not null default now(),
    updated_at timestamp
);

-- +goose Down
drop table if exists files_meta_data;
