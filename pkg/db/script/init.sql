CREATE TABLE mas_personal_income_tax (
  id SERIAL PRIMARY KEY,
  level INTEGER NOT NULL UNIQUE,
  description VARCHAR(100),
  percent_rate INTEGER,
  min_amount INTEGER,
  max_amount INTEGER
);

INSERT INTO mas_personal_income_tax (level, description, percent_rate, min_amount, max_amount)
VALUES 
(1, '0-150,000', 0, 0, 150000),
(2, '150,001-500,000', 10, 150001, 500000),
(3, '500,001-1,000,000', 15, 500001, 1000000),
(4, '1,000,001-2,000,000', 20, 1000001, 2000000),
(5, '2,000,000 ขึ้นไป', 35, 2000001, 2000001);

CREATE TABLE mas_deductions (
  id SERIAL PRIMARY KEY,
  type VARCHAR(255) NOT NULL,
  amount INTEGER
);

INSERT INTO mas_deductions (type, amount)
VALUES
('personal', 60000),
('k-receipt', 50000)
;