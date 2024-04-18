INSERT INTO public."Users" (
        first_name,
        last_name,
        username,
        "password",
        uni_id,
        email,
        user_type
    )
VALUES ('%s', '%s', '%s', '%s', '%d', '%s', '%s');