-- +goose Up
-- +goose StatementBegin
create table addresses (
	user_id                      uuid not null primary key references users(id),
	city                         text not null,
	address_line                 text not null,
	locality                     text not null,
	administrative_area_level_1  text not null,
	country                      text not null,
	postal_code                  integer not null
);

create index on addresses (user_id);

create trigger update_addresses_updated_at
before update on addresses for each row execute procedure update_updated_at_column();
-- +goose StatementEnd

-- +goose Down
drop table addresses;

