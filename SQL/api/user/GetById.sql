SELECT user_id,
    first_name,
    last_name,
    email,
    auth,
    is_affiliated_with_rso
FROM users
WHERE user_id = %s;