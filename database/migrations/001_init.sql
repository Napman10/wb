CREATE TABLE employee (
    id UUID NOT NULL PRIMARY KEY,
    fullname VARCHAR(255) NOT NULL DEFAULT '',
    gender INTEGER NOT NULL DEFAULT 0,
    age INTEGER NOT NULL DEFAULT 0,
    email VARCHAR(50) NOT NULL DEFAULT '',
    address TEXT NOT NULL DEFAULT '',
    vacation_days INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP(0),
    deleted_at TIMESTAMP
)