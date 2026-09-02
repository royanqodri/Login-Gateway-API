-- Drop tables in reverse order (karena ada ketergantungan)


DROP TABLE IF EXISTS t_customer;
DROP TABLE IF EXISTS t_user_module;
DROP TABLE IF EXISTS t_user;
DROP TABLE IF EXISTS t_module;
DROP TABLE IF EXISTS t_system;