create table Users
(
	Id uniqueidentifier primary key,
	Login nvarchar(max),
	PasswordHash nvarchar(256)
)

GO

alter table Transactions 
add UserId uniqueidentifier not null

GO

alter table Transactions
add foreign key(UserId) references Users(Id)