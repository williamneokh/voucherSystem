
USE voucher_db;


CREATE TABLE Voucher (
VID varchar (10) NOT NULL PRIMARY KEY,
UserID varchar (10),
UserPoints varchar(10),
CreatedDate DATETIME DEFAULT CURRENT_TIMESTAMP,
VoucherValue varchar(3),
RedeemedDate DATETIME,
MerchantID varchar (10)
);

CREATE TABLE MasterFund (
MFund_ID INT AUTO_INCREMENT,
PRIMARY KEY (MFund_ID),
TransactionType varchar(10),
SponsorIDOrVID varchar(8),
SponsorNameOrUserID varchar(30),
TransactionDate DATETIME DEFAULT CURRENT_TIMESTAMP,
Amount varchar(7),
BalancedFund varchar(8)
);
