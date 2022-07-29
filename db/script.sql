BEGIN;

CREATE TABLE wallet (
    wid bigserial PRIMARY KEY,
    username text NOT NULL,
    funds integer
);

INSERT INTO wallet (wid, username,funds) VALUES
    (1, 'john', 5),
    (2, 'jane', 10),
    (3, 'bob', 15),
    (4, 'rick', 100),
    (5, 'morty', 0);

COMMIT;
