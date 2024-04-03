INSERT INTO Users (
        first_name,
        last_name,
        email,
        password,
        auth,
        is_affiliated_with_rso
    )
VALUES (%s, %s, %s, %s, %s, %s);