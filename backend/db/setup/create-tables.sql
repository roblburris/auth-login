CREATE TABLE USERS (
    aud VARCHAR PRIMARY KEY,
    name VARCHAR,
    email VARCHAR,
    salt BYTEA,
    pw BYTEA,
    role VARCHAR
);

