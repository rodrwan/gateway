-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied

WITH user1 AS (
    INSERT INTO users (email, first_name, last_name, phone, birthdate)
    VALUES ('mar.vidal@example.com', 'Mar', 'Vidal', '936-235-189', now() - interval '20 years')
    RETURNING id
)
INSERT INTO addresses (user_id, city, address_line, locality, administrative_area_level_1, country, postal_code)
select id, 'Bilbao', '2479 calle de alcalá', 'Bilbao', 'Región de Murcia', 'España', 78276 from user1;

WITH user2 AS (
    INSERT INTO users (email, first_name, last_name, phone, birthdate)
    VALUES ('margarita.hernandez@example.com', 'Margarita', 'Hernandez', '935-441-524', now() - interval '24 years')
    RETURNING id
)
INSERT INTO addresses (user_id, city, address_line, locality, administrative_area_level_1, country, postal_code)
select id, 'Pontevedra', '9950 calle de ferraz', 'Pontevedra', 'Melilla', 'España', 71670  from user2;

WITH user3 AS (
    INSERT INTO users (email, first_name, last_name, phone, birthdate)
    VALUES ('alfredo.vidal@example.com', 'Alfredo', 'Vidal', '914-644-766', now() - interval '33 years')
    RETURNING id
)
INSERT INTO addresses (user_id, city, address_line, locality, administrative_area_level_1, country, postal_code)
select id ,'Murcia', '7181 avenida de castilla', 'Murcia', 'Comunidad de Madrid', 'España', 97371 from user3;
-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

delete from addresses;
delete from users;
