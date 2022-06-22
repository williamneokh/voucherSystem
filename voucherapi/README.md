# Voucher System API

### Adding Sponsored Fund
Method: "POST"

localhost:3000/api/sponsor/{sponsorID}

example JSON Input

{
"SponsorNameOrUserID":"Seng Siong",
"Amount":"2000"
}

example JSON Return
{
"ok": true,
"msg": "[MS-VOUCHERS]: Sponsor fund deposit, successful",
"data": "Code: ss008"
}


### List Deposit/Withdrawal from MasterFund table
Method: "POST"

localhost:3000/api/masterfund

example JSON Return

{
"ok": true,
"msg": "[MS-VOUCHERS]: Result: 1",
"data": {
"Mfund_ID": 1,
"TransactionType": "Deposit",
"SponsorIDOrVID": "SS13",
"SponsorNameOrUserID": "Seng Siong",
"TransactionDate": "2022-06-19 22:12:06",
"Amount": "12345678",
"BalancedFund": "12345678"
}
}
{
"ok": true,
"msg": "[MS-VOUCHERS]: Result: 2",
"data": {
"Mfund_ID": 2,
"TransactionType": "Deposit",
"SponsorIDOrVID": "d",
"SponsorNameOrUserID": "Seng Siong",
"TransactionDate": "2022-06-19 22:14:11",
"Amount": "12345678",
"BalancedFund": "24691356"
}
}

### User exchange Voucher
Method: "POST"

localhost:3000/api/getvoucher

example JSON Input

{
"UserID":"User001",
"Points":"1000",
"Value":"5"
}

example JSON Return

{
"ok": true,
"msg": "[MS-VOUCHERS]: Generate new voucher, successful",
"data": {
"VID": "66961331-4c68-40bb-9d2c-e71823762dc5",
"UserID": "User001",
"Points": "1000",
"Value": "5"
}
}

### Vendor Consume Voucher
Method: "POST"

localhost:3000:/api/consumevid

example JSON Input

{
"VID":"4569790e-cdef-4603-a913-6340b7d32e99",
"UserID":"User001",
"MerchantID":"ertertert"
}

example JSON Return

{
"ok": true,
"msg": "[MS-VOUCHERS]: Consume voucher, successful",
"data": {
"VID": "4569790e-cdef-4603-a913-6340b7d32e99",
"UserID": "User001",
"MerchantID": "ertertert"
}
}