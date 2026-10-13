-- Runs once, when the mysql volume is first created.
-- goravel_test backs `make test-mysql`, kept apart from the dev database.
CREATE DATABASE IF NOT EXISTS goravel_test;
GRANT ALL ON goravel_test.* TO 'goravel'@'%';
