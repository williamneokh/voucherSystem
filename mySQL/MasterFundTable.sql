
USE voucher_db;
drop table Voucher;
drop table MasterFund;
drop table FloatFund;


CREATE TABLE Voucher (
Voucher_ID INT AUTO_INCREMENT,
PRIMARY KEY (Voucher_ID),
VID varchar (36),
UserID varchar (10),
UserPoints varchar(10),
CreatedDate DATETIME DEFAULT CURRENT_TIMESTAMP,
VoucherValue varchar(3),
RedeemedDate DATETIME,
MerchantID varchar (50)
);

CREATE TABLE MasterFund (
MFund_ID INT AUTO_INCREMENT,
PRIMARY KEY (MFund_ID),
TransactionType varchar(10),
SponsorIDOrVID varchar(36),
SponsorNameOrUserID varchar(30),
TransactionDate DATETIME DEFAULT CURRENT_TIMESTAMP,
Amount varchar(7),
BalancedFund varchar(8)
);

CREATE TABLE FloatFund (
FFund_ID INT AUTO_INCREMENT,
PRIMARY KEY (FFund_ID),
VID varchar (36),
FloatDate DATETIME DEFAULT CURRENT_TIMESTAMP,
FloatValue varchar(8),
WithdrawalDate DATETIME,
MerchantID varchar(50)
);

