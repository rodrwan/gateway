-- +goose Up
-- +goose StatementBegin
create table users (
	id              uuid primary key default gen_random_uuid(),
	email           citext not null unique,

	first_name      text not null,
	last_name       text not null,
	phone           text not null,
	birthdate       date not null
);

create trigger update_users_updated_at
before update on users for each row execute procedure update_updated_at_column();
-- +goose StatementEnd

-- +goose Down
drop table users cascade;

