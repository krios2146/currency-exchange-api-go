CREATE TABLE IF NOT EXISTS Currencies (
    id          bigint  PRIMARY KEY NOT NULL,
    code        varchar UNIQUE NOT NULL,
    full_name   varchar NOT NULL,
    sign        varchar NOT NULL,

    CHECK (length(code) == 3)
);
