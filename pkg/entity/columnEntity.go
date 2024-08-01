package entity

type ColumnEntity struct {
	UserID          string   `json:"userID" bson:"user_id"`
	TableName       string   `json:"tableName" bson:"table_name"`
	ColumnsSelected []string `json:"columnsSelected" bson:"columns_selected"`
}
