drop database bookstore;
create database bookstore;
drop user storeap;
create user storeap with password '1qaz2wsx';
grant all privileges on database bookstore to storeap;