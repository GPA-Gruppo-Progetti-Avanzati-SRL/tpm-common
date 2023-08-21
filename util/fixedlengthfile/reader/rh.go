package reader

import (
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-common/util/fixedlengthfile"
)

var RHDefinition = fixedlengthfile.FixedLengthRecordDefinition{
	Id:                  "RH",
	PrefixDiscriminator: " RH",
	Fields: []fixedlengthfile.FixedLengthFieldDefinition{
		{Id: "start-filler", Name: "start-filler", Length: 1},
		{Id: "record-type", Name: "record-type", Length: 2},
		{Id: "sender", Name: "sender", Length: 5},
		{Id: "recipient", Name: "recipient", Length: 5},
		{Id: "creation-date", Name: "creation-date", Length: 6},
		{Id: "support-name", Name: "support-name", Length: 20},
		{Id: "filler-2", Name: "filler-2", Length: 76},
		{Id: "field-na", Name: "field-na", Length: 5},
	},
}

var RHEFDefinition = fixedlengthfile.FixedLengthRecordDefinition{
	Id:                  "RH-EF",
	PrefixDiscriminator: " EF",
	Fields: []fixedlengthfile.FixedLengthFieldDefinition{
		{Id: "start-filler", Name: "start-filler", Length: 1},
		{Id: "record-type", Name: "record-type", Length: 2},
		{Id: "sender", Name: "sender", Length: 5},
		{Id: "recipient", Name: "recipient", Length: 5},
		{Id: "creation-date", Name: "creation-date", Length: 6},
		{Id: "support-name", Name: "support-name", Length: 20},
		{Id: "filler-2", Name: "filler-2", Length: 6},
		{Id: "no-statements", Name: "no-statements", Length: 7},
		{Id: "filler-3", Name: "filler-3", Length: 30},
		{Id: "no-records", Name: "no-records", Length: 7},
		{Id: "filler-4", Name: "filler-4", Length: 25},
		{Id: "field-na", Name: "field-na", Length: 6},
	},
}

var RH61Definition = fixedlengthfile.FixedLengthRecordDefinition{
	Id:                  "RH-61",
	PrefixDiscriminator: " 61",
	Fields: []fixedlengthfile.FixedLengthFieldDefinition{
		{Id: "start-filler", Name: "start-filler", Length: 1},
		{Id: "record-type", Name: "record-type", Length: 2},
		{Id: "progr-number", Name: "progr-number", Length: 7},
		{Id: "filler-2", Name: "filler-2", Length: 13},
		{Id: "orig-bank-abi", Name: "orig-bank-abi", Length: 5},
		{Id: "reason", Name: "reason", Length: 5},
		{Id: "description", Name: "description", Length: 16},
		{Id: "account-type", Name: "account-type", Length: 2},
		// {Id: "bank-details", Name: "bank-details", Length: 23},
		{Id: "cin", Name: "cin", Length: 1},
		{Id: "bank-abi", Name: "bank-abi", Length: 5},
		{Id: "bank-cab", Name: "bank-cab", Length: 5},
		{Id: "current-account-code", Name: "current-account-code", Length: 12},
		{Id: "currency-code", Name: "currency-code", Length: 3},
		{Id: "accounting-date", Name: "accounting-date", Length: 6},
		{Id: "sign", Name: "sign", Length: 1},
		{Id: "opening-balance", Name: "opening-balance", Length: 15},
		// {Id: "more-iban-details", Name: "more-iban-details", Length: 4},
		{Id: "country-code", Name: "country-code", Length: 2},
		{Id: "check-digit", Name: "check-digit", Length: 2},
		{Id: "filler-3", Name: "filler-3", Length: 17},
	},
}

var RH62Definition = fixedlengthfile.FixedLengthRecordDefinition{
	Id:                  "RH-62",
	PrefixDiscriminator: " 62",
	Fields: []fixedlengthfile.FixedLengthFieldDefinition{
		{Id: "start-filler", Name: "start-filler", Length: 1},
		{Id: "record-type", Name: "record-type", Length: 2},
		{Id: "progr-number", Name: "progr-number", Length: 7},
		{Id: "movmnt-progr-number", Name: "movmnt-progr-number", Length: 3},
		{Id: "value-date", Name: "value-date", Length: 6},
		{Id: "accounting-date", Name: "accounting-date", Length: 6},
		{Id: "movmnt-sign", Name: "movmnt-sign", Length: 1},
		{Id: "movmntamount", Name: "movmntamount", Length: 15},
		{Id: "cbi-reason", Name: "cbi-reason", Length: 2},
		{Id: "internal-reason", Name: "internal-reason", Length: 2},
		{Id: "cheque-number", Name: "cheque-number", Length: 16},
		{Id: "bank-ref", Name: "bank-ref", Length: 16},
		{Id: "cust-ref-type", Name: "cust-ref-type", Length: 9},
		{Id: "cust-ref-movmnt-descr", Name: "cust-ref-movmnt-descr", Length: 34},
	},
}

var RH63Definition_KKK = fixedlengthfile.FixedLengthRecordDefinition{
	Id:                  "RH-63-KKK",
	PrefixDiscriminator: " 63**********KKK",
	Fields: []fixedlengthfile.FixedLengthFieldDefinition{
		{Id: "start-filler", Name: "start-filler", Length: 1},
		{Id: "record-type", Name: "record-type", Length: 2},
		{Id: "progr-number", Name: "progr-number", Length: 7},
		{Id: "movmnt-progr-number", Name: "movmnt-progr-number", Length: 3},
		{Id: "structure-flag", Name: "structure-flag", Length: 3},
		{Id: "type-identifier", Name: "type-identifier", Length: 23},
		{Id: "filler-2", Name: "filler-2", Length: 81},
	},
}

var RH63Definition_YYY = fixedlengthfile.FixedLengthRecordDefinition{
	Id:                  "RH-63-YYY",
	PrefixDiscriminator: " 63**********YYY",
	Fields: []fixedlengthfile.FixedLengthFieldDefinition{
		{Id: "start-filler", Name: "start-filler", Length: 1},
		{Id: "record-type", Name: "record-type", Length: 2},
		{Id: "progr-number", Name: "progr-number", Length: 7},
		{Id: "movmnt-progr-number", Name: "movmnt-progr-number", Length: 3},
		{Id: "structure-flag", Name: "structure-flag", Length: 3},
		{Id: "order-date", Name: "order-date", Length: 8},
		{Id: "ordering-prty--taxpayer-code", Name: "ordering-prty--taxpayer-code", Length: 16},
		{Id: "ordering-prty-descr", Name: "ordering-prty-descr", Length: 40},
		{Id: "country", Name: "country", Length: 40},
	},
}

var RH63Definition_YY2 = fixedlengthfile.FixedLengthRecordDefinition{
	Id:                  "RH-63-YY2",
	PrefixDiscriminator: " 63**********YY2",
	Fields: []fixedlengthfile.FixedLengthFieldDefinition{
		{Id: "start-filler", Name: "start-filler", Length: 1},
		{Id: "record-type", Name: "record-type", Length: 2},
		{Id: "progr-number", Name: "progr-number", Length: 7},
		{Id: "movmnt-progr-number", Name: "movmnt-progr-number", Length: 3},
		{Id: "structure-flag", Name: "structure-flag", Length: 3},
		{Id: "ordering-prty-addr", Name: "ordering-prty-addr", Length: 50},
		{Id: "ordering-prty-iban", Name: "ordering-prty-iban", Length: 34},
		{Id: "filler-2", Name: "filler-2", Length: 20},
	},
}

var RH63Definition_ZZ1 = fixedlengthfile.FixedLengthRecordDefinition{
	Id:                  "RH-63-ZZ1",
	PrefixDiscriminator: " 63**********ZZ1",
	Fields: []fixedlengthfile.FixedLengthFieldDefinition{
		{Id: "start-filler", Name: "start-filler", Length: 1},
		{Id: "record-type", Name: "record-type", Length: 2},
		{Id: "progr-number", Name: "progr-number", Length: 7},
		{Id: "movmnt-progr-number", Name: "movmnt-progr-number", Length: 3},
		{Id: "structure-flag", Name: "structure-flag", Length: 3},
		{Id: "orig-amnt", Name: "orig-amnt", Length: 18},
		{Id: "orig-amnt-currency-code", Name: "orig-amnt-currency-code", Length: 3},
		{Id: "paid-amnt", Name: "paid-amnt", Length: 18},
		{Id: "paid-amnt-currency-code", Name: "paid-amnt-currency-code", Length: 3},
		{Id: "trx-amnt", Name: "trx-amnt", Length: 18},
		{Id: "trx-amnt-currency-code", Name: "trx-amnt-currency-code", Length: 3},
		{Id: "exchg-rate", Name: "exchg-rate", Length: 12},
		{Id: "commission-amnt", Name: "commission-amnt", Length: 13},
		{Id: "commission-fees-amnt", Name: "commission-fees-amnt", Length: 13},
		{Id: "country-code", Name: "country-code", Length: 3},
	},
}

var RH63Definition_ZZ2 = fixedlengthfile.FixedLengthRecordDefinition{
	Id:                  "RH-63-ZZ2",
	PrefixDiscriminator: " 63**********ZZ2",
	Fields: []fixedlengthfile.FixedLengthFieldDefinition{
		{Id: "start-filler", Name: "start-filler", Length: 1},
		{Id: "record-type", Name: "record-type", Length: 2},
		{Id: "progr-number", Name: "progr-number", Length: 7},
		{Id: "movmnt-progr-number", Name: "movmnt-progr-number", Length: 3},
		{Id: "structure-flag", Name: "structure-flag", Length: 3},
		{Id: "ordering-prty", Name: "ordering-prty", Length: 104},
	},
}

var RH63Definition_ZZ3 = fixedlengthfile.FixedLengthRecordDefinition{
	Id:                  "RH-63-ZZ3",
	PrefixDiscriminator: " 63**********ZZ3",
	Fields: []fixedlengthfile.FixedLengthFieldDefinition{
		{Id: "start-filler", Name: "start-filler", Length: 1},
		{Id: "record-type", Name: "record-type", Length: 2},
		{Id: "progr-number", Name: "progr-number", Length: 7},
		{Id: "movmnt-progr-number", Name: "movmnt-progr-number", Length: 3},
		{Id: "structure-flag", Name: "structure-flag", Length: 3},
		{Id: "payee", Name: "payee", Length: 50},
		{Id: "payment-reason", Name: "payment-reason", Length: 54},
	},
}

var RH63Definition_ID1 = fixedlengthfile.FixedLengthRecordDefinition{
	Id:                  "RH-63-ID1",
	PrefixDiscriminator: " 63**********ID1",
	Fields: []fixedlengthfile.FixedLengthFieldDefinition{
		{Id: "start-filler", Name: "start-filler", Length: 1},
		{Id: "record-type", Name: "record-type", Length: 2},
		{Id: "progr-number", Name: "progr-number", Length: 7},
		{Id: "movmnt-progr-number", Name: "movmnt-progr-number", Length: 3},
		{Id: "structure-flag", Name: "structure-flag", Length: 3},
		{Id: "msg-id", Name: "msg-id", Length: 35},
		{Id: "end-2-end-id", Name: "end-2-end-id", Length: 35},
		{Id: "filler", Name: "filler", Length: 34},
	},
}

var RH63Definition_RI1 = fixedlengthfile.FixedLengthRecordDefinition{
	Id:                  "RH-63-RI1",
	PrefixDiscriminator: " 63**********RI1",
	Fields: []fixedlengthfile.FixedLengthFieldDefinition{
		{Id: "start-filler", Name: "start-filler", Length: 1},
		{Id: "record-type", Name: "record-type", Length: 2},
		{Id: "progr-number", Name: "progr-number", Length: 7},
		{Id: "movmnt-progr-number", Name: "movmnt-progr-number", Length: 3},
		{Id: "structure-flag", Name: "structure-flag", Length: 3},
		{Id: "reconc-data", Name: "reconc-data", Length: 104},
	},
}

var RH63Definition_RI2 = fixedlengthfile.FixedLengthRecordDefinition{
	Id:                  "RH-63-RI2",
	PrefixDiscriminator: " 63**********RI2",
	Fields: []fixedlengthfile.FixedLengthFieldDefinition{
		{Id: "start-filler", Name: "start-filler", Length: 1},
		{Id: "record-type", Name: "record-type", Length: 2},
		{Id: "progr-number", Name: "progr-number", Length: 7},
		{Id: "movmnt-progr-number", Name: "movmnt-progr-number", Length: 3},
		{Id: "structure-flag", Name: "structure-flag", Length: 3},
		{Id: "reconc-data", Name: "reconc-data", Length: 104},
	},
}

var RH63Definition_Else = fixedlengthfile.FixedLengthRecordDefinition{
	Id:                  "RH-63-Else",
	PrefixDiscriminator: " 63",
	Fields: []fixedlengthfile.FixedLengthFieldDefinition{
		{Id: "start-filler", Name: "start-filler", Length: 1},
		{Id: "record-type", Name: "record-type", Length: 2},
		{Id: "progr-number", Name: "progr-number", Length: 7},
		{Id: "movmnt-progr-number", Name: "movmnt-progr-number", Length: 3},
		{Id: "descr", Name: "descr", Length: 107},
	},
}

var RH64Definition = fixedlengthfile.FixedLengthRecordDefinition{
	Id:                  "RH-64",
	PrefixDiscriminator: " 64",
	Fields: []fixedlengthfile.FixedLengthFieldDefinition{
		{Id: "start-filler", Name: "start-filler", Length: 1},
		{Id: "record-type", Name: "record-type", Length: 2},
		{Id: "progr-number", Name: "progr-number", Length: 7},
		{Id: "currency-code", Name: "currency-code", Length: 3},
		{Id: "accounting-date", Name: "accounting-date", Length: 6},
		{Id: "accounts-balance-sign", Name: "accounts-balance-sign", Length: 1},
		{Id: "accounts-balance", Name: "accounts-balance", Length: 15},
		{Id: "cash-balance-sign", Name: "cash-balance-sign", Length: 1},
		{Id: "cash-balance", Name: "cash-balance", Length: 15},
		{Id: "filler-2", Name: "filler-2", Length: 54},
		{Id: "filler-3", Name: "filler-3", Length: 15},
	},
}

var RH65Definition = fixedlengthfile.FixedLengthRecordDefinition{
	Id:                  "RH-65",
	PrefixDiscriminator: " 65",
	Fields: []fixedlengthfile.FixedLengthFieldDefinition{
		{Id: "start-filler", Name: "start-filler", Length: 1},
		{Id: "record-type", Name: "record-type", Length: 2},
		{Id: "progr-number", Name: "progr-number", Length: 7},
		// {Id: "first-cash-balance", Name: "first-cash-balance", Length: 22},
		{Id: "first-cash-on-hand-date", Name: "first-cash-on-hand-date", Length: 6},
		{Id: "first-cash-sign", Name: "first-cash-sign", Length: 1},
		{Id: "first-cash-balance", Name: "first-cash-balance", Length: 15},
		// {Id: "second-cash-balance", Name: "second-cash-balance", Length: 22},
		{Id: "second-cash-on-hand-date", Name: "second-cash-on-hand-date", Length: 6},
		{Id: "second-cash-sign", Name: "second-cash-sign", Length: 1},
		{Id: "second-cash-balance", Name: "second-cash-balance", Length: 15},
		// {Id: "third-cash-balance", Name: "third-cash-balance", Length: 22},
		{Id: "third-cash-on-hand-date", Name: "third-cash-on-hand-date", Length: 6},
		{Id: "third-cash-sign", Name: "third-cash-sign", Length: 1},
		{Id: "third-cash-balance", Name: "third-cash-balance", Length: 15},
		// {Id: "fourth-cash-balance", Name: "fourth-cash-balance", Length: 22},
		{Id: "fourth-cash-on-hand-date", Name: "fourth-cash-on-hand-date", Length: 6},
		{Id: "fourth-cash-sign", Name: "fourth-cash-sign", Length: 1},
		{Id: "fourth-cash-balance", Name: "fourth-cash-balance", Length: 15},
		//{Id: "fifth-cash-balance", Name: "fifth-cash-balance", Length: 22},
		{Id: "fifth-cash-on-hand-date", Name: "fifth-cash-on-hand-date", Length: 6},
		{Id: "fifth-cash-sign", Name: "fifth-cash-sign", Length: 1},
		{Id: "fifth-cash-balance", Name: "fifth-cash-balance", Length: 15},
	},
}
