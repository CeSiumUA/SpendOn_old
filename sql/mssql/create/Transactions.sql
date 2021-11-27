create table Transactions
(
	Id bigint primary key,
	Amount MONEY NOT NULL,
	SpentAt DATETIME2 NOT NULL,
	Note NVARCHAR(max),
	CategoryId INT REFERENCES Categories (Id) NOT NULL
)