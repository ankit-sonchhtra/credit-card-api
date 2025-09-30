db = db.getSiblingDB('credit-card-api');
db.createCollection('accounts');
db.createCollection('users');
db.createCollection('transactions');

print("initialized credit-card-api DB with initial required collections");