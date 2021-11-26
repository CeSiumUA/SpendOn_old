create table Users
(
	Id INT primary key,
	Login nvarchar(max),
	PasswordHash nvarchar(256)
)

GO

alter table Transactions 
add UserId INT not null

GO

alter table Transactions
add foreign key(UserId) references Users(Id)