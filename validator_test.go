package eos

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type ruleTestCase struct {
	name          string
	value         interface{}
	expectedError string
}

func TestEOSBlockNumRule(t *testing.T) {
	tag := "eos_block_num"
	validator := func(field string, value interface{}) error {
		return EOSBlockNumRule(field, tag, "", value)
	}

	tests := []ruleTestCase{
		{"should be a string", true, "The test field must be a string"},
		{"should not contains invalid characters", "!", "The test field must be a valid EOS block num"},

		{"valid block num", "10", ""},
	}

	runRuleTestCases(t, tag, tests, validator)
}

func TestEOSNameRule(t *testing.T) {
	tag := "eos_name"
	validator := func(field string, value interface{}) error {
		return EOSNameRule(field, tag, "", value)
	}

	tests := []ruleTestCase{
		{"should be a string", true, "The test field is not a known type for an EOS name"},
		{"should not contains invalid characters", "6", "The test field must be a valid EOS name"},
		{"should not be longer than 13", "abcdefghigklma", "The test field must be a valid EOS name"},

		{"valid empty", "", ""},
		{"valid single", "e", ""},
		{"valid limit", "5", ""},
		{"valid with dots and 13 chars", "eosio.tokenfl", ""},
		{"valid eos.Name", Name("eosio"), ""},
		{"valid eos.PermissionName", PermissionName("eosio"), ""},
		{"valid eos.ActionName", ActionName("eosio"), ""},
		{"valid eos.AccountName", AccountName("eosio"), ""},
		{"valid eos.TableName", TableName("eosio"), ""},
	}

	runRuleTestCases(t, tag, tests, validator)
}

func TestEOSExtendedNameRule(t *testing.T) {
	tag := "eos_extended_name"
	validator := func(field string, value interface{}) error {
		return EOSExtendedNameRule(field, tag, "", value)
	}

	tests := []ruleTestCase{
		{"should be a string", true, "The test field is not a known type for an EOS name"},
		{"should not contains invalid characters", "6", "The test field must be a valid EOS name"},
		{"should not be longer than 13", "abcdefghigklma", "The test field must be a valid EOS name"},

		{"valid empty", "", ""},
		{"valid single", "e", ""},
		{"valid limit", "5", ""},
		{"valid with dots and 13 chars", "eosio.tokenfl", ""},
		{"valid with whem symbol", "4,EOS", ""},
		{"valid with whem symbol code", "EOS", ""},

		{"valid eos.Name", Name("eosio"), ""},
		{"valid eos.PermissionName", PermissionName("eosio"), ""},
		{"valid eos.ActionName", ActionName("eosio"), ""},
		{"valid eos.AccountName", AccountName("eosio"), ""},
		{"valid eos.TableName", TableName("eosio"), ""},
	}

	runRuleTestCases(t, tag, tests, validator)
}

func TestEOSNamesListRule(t *testing.T) {
	tag := "eos_names_list"
	rule := EOSNamesListRuleFactory("|", 2)
	validator := func(field string, value interface{}) error {
		return rule(field, tag, "", value)
	}

	tests := []ruleTestCase{
		{"should be a string", true, "The test field must be a string"},
		{"should have at least 1 element", "", "The test field must have at least 1 element"},
		{"should have at max macCount element", "eos|eos|eos", "The test field must have at most 2 elements"},
		{"should fail on single error", "6", "The test[0] field must be a valid EOS name"},
		{"should fail if any element error", "ab|6", "The test[1] field must be a valid EOS name"},

		{"valid single", "ab", ""},
		{"valid multiple", "ded|eos", ""},
	}

	runRuleTestCases(t, tag, tests, validator)
}

func TestEOSExtendedNamesListRule(t *testing.T) {
	tag := "eos_extended_names_list"
	rule := EOSExtendedNamesListRuleFactory("|", 3)
	validator := func(field string, value interface{}) error {
		return rule(field, tag, "", value)
	}

	tests := []ruleTestCase{
		{"should be a string", true, "The test field must be a string"},
		{"should have at least 1 element", "", "The test field must have at least 1 element"},
		{"should have at max macCount element", "eos|eos|eos|eos", "The test field must have at most 3 elements"},
		{"should fail on single error", "6", "The test[0] field must be a valid EOS name"},
		{"should fail if any element error", "ab|6", "The test[1] field must be a valid EOS name"},

		{"valid single", "ab", ""},
		{"valid multiple", "ded|eos", ""},
		{"valid multiple symbol", "ded|eos|EOS", ""},
		{"valid multiple symbol code", "ded|4,EOS", ""},
		{"valid multiple mixed", "ded|EOS|4,EOS", ""},
	}

	runRuleTestCases(t, tag, tests, validator)
}

func TestEOSTrxIDRule(t *testing.T) {
	tag := "eos_trx_id"
	validator := func(field string, value interface{}) error {
		return EOSTrxIDRule(field, tag, "", value)
	}

	tests := []ruleTestCase{
		{"should be a string", true, "The test field must be a string"},
		{"should contains something", "", "The test field must be a valid hexadecimal"},
		{"should contains a least two characters", "a", "The test field must be a valid hexadecimal"},
		{"should not contains invalid characters", "az", "The test field must be a valid hexadecimal"},
		{"should be a multple of 2", "ab01020", "The test field must be a valid hexadecimal"},
		{"should be long enough", "d8fe02221408fbcc221d1207c1b8cc67e0d9b3ca1c6005a36ea10428dd7fd1", "The test field must have exactly 64 characters"},

		{"valid", "d8fe02221408fbcc221d1207c1b8cc67e0d9b3ca1c6005a36ea10428dd7fd148", ""},
		{"valid", "D8FE02221408FBCC221D1207C1B8CC67E0D9B3CA1C6005A36EA10428DD7FD148", ""},
	}

	runRuleTestCases(t, tag, tests, validator)
}

func runRuleTestCases(t *testing.T, tag string, tests []ruleTestCase, validator func(field string, value interface{}) error) {
	for _, test := range tests {
		t.Run(fmt.Sprintf("%s_%s", tag, test.name), func(t *testing.T) {
			err := validator("test", test.value)

			if test.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.Equal(t, errors.New(test.expectedError), err)
			}
		})
	}
}
