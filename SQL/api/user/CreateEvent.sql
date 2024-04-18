INSERT INTO public."Events"
(event_id, "name", tags, description, start_time, end_time, loc_id, contact_phone, contact_email, visibility, uni_id, rso_id, superadmin_approval)
VALUES(nextval('"Events_event_id_seq"'::regclass), '', ?, '', '', '', nextval('"Events_loc_id_seq"'::regclass), '', '', '', nextval('"Events_uni_id_seq"'::regclass), nextval('"Events_rso_id_seq"'::regclass), false);