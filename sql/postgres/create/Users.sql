create table Users
(
	Id INT primary key,
	Login varchar(500),
	PasswordHash varchar(256)
);

alter table Transactions
add foreign key(UserId) references Users(Id);

alter table Users
add Currency varchar(3)