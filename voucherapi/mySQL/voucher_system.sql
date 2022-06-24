DROP DATABASE IF EXISTS voucher_system_db;

CREATE DATABASE voucher_system_db;


USE voucher_system_db;
DROP TABLE IF EXISTS Voucher;
DROP TABLE IF EXISTS MasterFund;
DROP TABLE IF EXISTS FloatFund;


CREATE TABLE Voucher (
Voucher_ID INT AUTO_INCREMENT,
PRIMARY KEY (Voucher_ID),
VID varchar (36),
UserID varchar (10),
UserPoints varchar(10),
CreatedDate DATETIME DEFAULT CURRENT_TIMESTAMP,
VoucherValue varchar(3),
RedeemedDate DATETIME,
MerchantID varchar (50),
Branch varchar(50)
);

CREATE TABLE MasterFund (
MFund_ID INT AUTO_INCREMENT,
PRIMARY KEY (MFund_ID),
TransactionType varchar(10),
SponsorIDOrVID varchar(36),
SponsorNameOrUserID varchar(30),
TransactionDate DATETIME DEFAULT CURRENT_TIMESTAMP,
Amount varchar(8),
BalancedFund varchar(8)
);

CREATE TABLE FloatFund (
FFund_ID INT AUTO_INCREMENT,
PRIMARY KEY (FFund_ID),
VID varchar (36),
FloatDate DATETIME DEFAULT CURRENT_TIMESTAMP,
FloatValue varchar(8),
WithdrawalDate DATETIME,
Branch varchar(50)
);

