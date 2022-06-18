# Voucher System API

### Adding Sponsored Fund
Method: "POST"

localhost:3000/api/sponsor/{sponsorID}

example JSON Input

{
"SponsorNameOrUserID":"Seng Siong",
"Amount":"2000"
}


### List Deposit/Withdrawal from MasterFund table
Method: "POST"

localhost:3000/api/masterfund

example JSON Return

{"Mfund_ID":24,"TransactionType":"Deposit","SponsorIDOrVID":"NTUC01","SponsorNameOrUserID":"NTUC","TransactionDate":"2022-06-16 16:58:10","Amount":"1000","BalancedFund":"1000"}
{"Mfund_ID":25,"TransactionType":"Deposit","SponsorIDOrVID":"SS01","SponsorNameOrUserID":"Seng Siong","TransactionDate":"2022-06-16 16:58:52","Amount":"2000","BalancedFund":"3000"}
{"Mfund_ID":26,"TransactionType":"Deposit","SponsorIDOrVID":"SS02","SponsorNameOrUserID":"Seng Siong","TransactionDate":"2022-06-16 21:18:34","Amount":"2000","BalancedFund":"5000"}


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

{"VID":"4c04afc3-309e-4016-a753-c031d0e686af","UserID":"User001","Points":"1000","Value":"5"}

### Vendor Consume Voucher
Method: "POST"

localhost:3000:/api/consumevid

example JSON Input

{
"VID":"9f38ffcb-ec90-4f20-87cb-1d3566b16392",
"UserID":"User001",
"MerchantID":"ertertert"
}

example JSON Return

{
"Status": "202",
"Message": "Successfully Consumed",
"VID": "9f38ffcb-ec90-4f20-87cb-1d3566b16392",
"UserID": "User001",
"MerchantID": "ertertert"
}