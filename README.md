# Malina
Malina is a ecommerce project on Golang and PostgreSQL 9.6 with MVC pattern

Min. Go Version: 1.8

Installation

1. Download dump.sql file and import to your PostgreSql 9.6 database
2.  Installation
#### $ go get github.com/ilyaran/Malina 

or download the zip containing the source: 
https://github.com/ilyaran/Malina/archive/master.zip

3. Open /config/app.go file and set your settings
4. Open /assets/filemanager/conf.json file and set your configuration

Admin (Boss) account email: ilyaran@mail.ru, password: ilyaran
You can generate a password for any user, open main.go file and uncomment where generatePasswordForAnyAccount("your password").
You'll see crypted password in terminal, should replace password column in database

Enjoy, and contributions are more than welcome!
