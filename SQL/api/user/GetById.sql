SELECT user_id, 
    first_name, 
    last_name, 
    username, 
    email
FROM public."Users"
WHERE user_id = %s