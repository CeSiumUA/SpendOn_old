create table Users
(
	Id int identity primary key,
	Login nvarchar(max),
	PasswordHash nvarchar(256)
)

alter table Transactions 
add UserId int not null

alter table Transactions
add foreign key(UserId) references Users(Id)