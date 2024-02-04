SET TIMEZONE='Asia/Dushanbe';

CREATE TABLE limits (
    id SERIAL PRIMARY KEY NOT NULL,
    name VARCHAR(100) NOT NULL,
    max_amount  BIGINT NOT NULL
);

CREATE TABLE wallets (
    id SERIAL PRIMARY KEY NOT NULL,
    balance BIGINT NOT NULL DEFAULT 0 CHECK (balance >= 0),
    type INT NOT NULL DEFAULT 1,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    user_id CHAR(36) NOT NULL UNIQUE,

    FOREIGN KEY (type) REFERENCES limits(id)
);

CREATE TABLE transactions (
    id SERIAL PRIMARY KEY NOT NULL,
    wallet_id INT NOT NULL,
    amount BIGINT NOT NULL CHECK (amount > 0),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),

    FOREIGN KEY (wallet_id) REFERENCES wallets(id)
);

INSERT INTO limits (name, max_amount)
VALUES
    ('unidentified wallet', 1000000),
    ('identified wallet', 10000000);

INSERT INTO wallets (balance, type, user_id)
VALUES
    (50000, 1, '36764dc2-2653-4e7f-b24c-430deca66b88'),
    (150000, 2, 'c76fdd66-3d0c-4633-8274-c12f67e4fa2a'),
    (510000, 1, '1c6287a0-7071-4b63-af89-24a87ce89599'),
    (30000, 2, '69bccb14-69f8-48c8-b123-f80d65e6927f'),
    (0, 2, 'd136f61a-6a4c-4029-8bc6-6b722b80e0b3');

INSERT INTO transactions (wallet_id, amount)
VALUES
    (1, 50000),
    (2, 150000),
    (3, 510000),
    (4, 30000);