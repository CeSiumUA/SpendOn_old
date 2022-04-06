create table Transactions
(
	Id uuid primary key,
	Amount MONEY NOT NULL,
	SpentAt TIMESTAMP NOT NULL,
	Note TEXT,
	CategoryId INT REFERENCES Categories (Id) NOT NULL,
    UserId INT not null
)