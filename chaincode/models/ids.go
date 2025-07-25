package models

import "fmt"

func BuildQueryIdStartsWith(prefix string) string {
	return fmt.Sprintf("{\"selector\": {\"id\": { \"$regex\": \"^(%s-)\" } } }", prefix)
}

func BuildQueryForEntityType(entityName string, selectors string) string {
	return fmt.Sprintf("{\"selector\": {\"id\": { \"$regex\": \"^(%s-)\" }, %s } }", entityName, selectors)
}

func BuildQueryFieldContains(entityName string, fieldName string, substring string) string {
	fieldSelector := BuildContainsSelector(fieldName, substring)
	return BuildQueryForEntityType(entityName, fieldSelector)
}

func BuildContainsSelector(fieldName string, substring string) string {
	return fmt.Sprintf("\"%s\": { \"$regex\": \"%s\" }", fieldName, substring)
}

func FormatKey(entityType string, entityID string) string {
	return fmt.Sprintf("%s-%s", entityType, entityID)
}

func ToProductID(id string) string {
	return FormatKey(PRODUCT_TYPE, id)
}

func ToUserID(id string) string {
	return FormatKey(USER_TYPE, id)
}

func ToReceiptID(id string) string {
	return FormatKey(RECEIPT_TYPE, id)
}

func ToTraderID(id string) string {
	return FormatKey(TRADER_TYPE, id)
}
