-- REGISTER
INSERT INTO users (email, password, pin_hash, balance) VALUES 
('john@example.com', 'password123', '123456', 0);

-- LOGIN
SELECT u.id, u.email, u.password
FROM users u
WHERE u.email = 'john@example.com' AND u.password = 'password123';

--TOPUP
INSERT INTO transactions (receiver_id, transaction_type, amount, status, description)
VALUES ('11111111-1111-1111-1111-111111111111', 'TOP_UP', 100000, 'PENDING', 'Top up via Bank Transfer');
UPDATE users SET balance = balance + 100000 WHERE id = '11111111-1111-1111-1111-111111111111';
UPDATE transactions SET status = 'COMPLETED', completed_at = CURRENT_TIMESTAMP WHERE id = 'transaction_id';