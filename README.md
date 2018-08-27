# Gateway

Service that connect GraphQL queries with a microservice architecture.

```
query {
    user(email:"test@finciero.com") {
        id,
        first_name,
        last_name,
        email,
        phone,
        birthdate,
        dni,
        dni_type,
        family_status,
        gender,
        title,
        address {
            address_line,
            city,
            locality,
            administrative_area_level_1,
            country,
            place,
            postal_code
        },
        cards {
            id,
            user_id,
            product_id,
            card_number,
            reference_email,
            reference_user_id,
            reference_id,
            deposits {
                id,
                amount,
                status,
                payment_id,
                total,
                fee,
                created_at,
                dollar,
                total
            }
        }
    }
}
```