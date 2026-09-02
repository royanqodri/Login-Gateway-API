CREATE TABLE IF NOT EXISTS t_customer (
    id          BIGSERIAL PRIMARY KEY,
    customer_no VARCHAR(20) NOT NULL,
    name        VARCHAR(50) NOT NULL,
    description VARCHAR(300) NOT NULL,
    file_path   TEXT,
    status_data VARCHAR(2) NOT NULL,
    insert_by   TEXT NOT NULL,
    insert_time TIMESTAMP NOT NULL,
    update_by   TEXT NOT NULL,
    update_time TIMESTAMP NOT NULL
);

CREATE INDEX IDX_t_customer_1 ON t_customer USING btree (customer_no);
CREATE INDEX IDX_t_customer_2 ON t_customer USING btree (name);
CREATE INDEX IDX_t_customer_3 ON t_customer USING btree (customer_no, name);

CREATE UNIQUE INDEX UK_t_customer_1 ON t_customer(customer_no);
CREATE UNIQUE INDEX UK_t_customer_2 ON t_customer(name);

----------------------------------------------------------------------------------------------


CREATE TABLE IF NOT EXISTS t_user (
    id                  BIGSERIAL PRIMARY KEY,
    id_customer         BIGINT NOT NULL,
    username            VARCHAR(20) NOT NULL,
    password            VARCHAR(200) NOT NULL,
    name                VARCHAR(500) NOT NULL,
    phone_no            VARCHAR(15) NOT NULL,
    email_address       VARCHAR(100) NOT NULL,
    blood_type          VARCHAR(5) NOT NULL,
    status_data         VARCHAR(2) NOT NULL,
    insert_by           TEXT NOT NULL,
    insert_time         TIMESTAMP NOT NULL,
    update_by           TEXT NOT NULL,
    update_time         TIMESTAMP NOT NULL
);

CREATE INDEX IDX_t_user_1 ON t_user USING btree (id_customer);
CREATE INDEX IDX_t_user_2 ON t_user USING btree (id_customer, username);
CREATE INDEX IDX_t_user_3 ON t_user USING btree (id_customer, name);
CREATE INDEX IDX_t_user_4 ON t_user USING btree (id_customer, phone_no);
CREATE INDEX IDX_t_user_5 ON t_user USING btree (id_customer, email_address);

CREATE UNIQUE INDEX UK_t_user_1 ON t_user(id_customer, username);
CREATE UNIQUE INDEX UK_t_user_2 ON t_user(id_customer, email_address);
CREATE UNIQUE INDEX UK_t_user_3 ON t_user(id_customer, phone_no);


-----------------------------------------------------------

CREATE TABLE IF NOT EXISTS t_user_module (
    id                  BIGSERIAL PRIMARY KEY,
    id_customer         BIGINT NOT NULL,
    id_user             BIGINT NOT NULL,
    id_module           BIGINT NOT NULL,
    create_access       BOOLEAN DEFAULT FALSE NOT NULL,
    read_access         BOOLEAN DEFAULT FALSE NOT NULL,
    update_access       BOOLEAN DEFAULT FALSE NOT NULL,
    delete_access       BOOLEAN DEFAULT FALSE NOT NULL,
    export_access       BOOLEAN DEFAULT FALSE NOT NULL,
    status_data         VARCHAR(2) NOT NULL,
    insert_by           TEXT NOT NULL,
    insert_time         TIMESTAMP NOT NULL,
    update_by           TEXT NOT NULL,
    update_time         TIMESTAMP NOT NULL
);

CREATE INDEX IDX_t_user_module_1 ON t_user_module USING btree (id_customer);
CREATE INDEX IDX_t_user_module_2 ON t_user_module USING btree (id_customer, id_user);
CREATE INDEX IDX_t_user_module_3 ON t_user_module USING btree (id_customer, id_module);
CREATE INDEX IDX_t_user_module_4 ON t_user_module USING btree (id_customer, id_user, id_module);

CREATE UNIQUE INDEX UK_t_user_module_1 ON t_user_module(id_customer, id_user, id_module);

-------------------------------------------------------

CREATE TABLE IF NOT EXISTS t_module (
   id BIGSERIAL PRIMARY KEY,
   id_customer BIGINT NOT NULL,
   id_system int4 NOT NULL,
   code varchar(50) NOT NULL,
   name varchar(100) NOT NULL,
   status_data varchar(2) NOT NULL,
   insert_by TEXT NOT NULL,
   insert_time timestamp NOT NULL,
   update_by TEXT NOT NULL,
   update_time timestamp NOT NULL
);

CREATE INDEX IDX_t_module_1 ON t_module USING btree (id_customer);
CREATE INDEX IDX_t_module_2 ON t_module USING btree (id_customer, id_system);
CREATE INDEX IDX_t_module_3 ON t_module USING btree (id_customer, code);
CREATE INDEX IDX_t_module_4 ON t_module USING btree (id_customer, name);
CREATE INDEX IDX_t_module_5 ON t_module USING btree (id_customer, id_system, code);
CREATE INDEX IDX_t_module_6 ON t_module USING btree (id_customer, id_system, name);

CREATE UNIQUE INDEX UK_t_module_1 ON t_module(id_customer, id_system, code);
CREATE UNIQUE INDEX UK_t_module_2 ON t_module(id_customer, id_system, name);

--------------------------------------------------------

CREATE TABLE IF NOT EXISTS t_system (
	id BIGSERIAL PRIMARY KEY,
    id_customer BIGINT NOT NULL,
	code varchar(10) NOT NULL,
	name varchar(50) NOT NULL,
	status_data varchar(2) NOT NULL,
	insert_by TEXT NOT NULL,
	insert_time timestamp NOT NULL,
	update_by TEXT NOT NULL,
	update_time timestamp NOT NULL
);

CREATE INDEX IDX_t_system_1 ON t_system USING btree (id_customer);
CREATE INDEX IDX_t_system_2 ON t_system USING btree (id_customer, code);
CREATE INDEX IDX_t_system_3 ON t_system USING btree (id_customer, name);

CREATE UNIQUE INDEX UK_t_system_1 ON t_system(id_customer, code);
CREATE UNIQUE INDEX UK_t_system_2 ON t_system(id_customer, name);

-------------------------------