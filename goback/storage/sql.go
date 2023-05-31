package storage

const sqlGetRoomUsers = `
with ranked as (
    select email,
           created_at,
           first_VALUE(connected) over w as connected,
           first_VALUE(active) over w as active,
           row_number() over w as row_number
    from active_room_users
    where room_id = $1
        WINDOW w AS (PARTITION BY email order by created_at desc RANGE BETWEEN UNBOUNDED PRECEDING AND UNBOUNDED FOLLOWING)
) select email, created_at, connected, active from ranked where row_number = 1 order by email
`
