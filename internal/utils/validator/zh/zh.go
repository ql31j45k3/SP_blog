package zh

func NewField2Name() map[string]string {
	field2Name := make(map[string]string)

	field2Name["title"] = "標籤"
	field2Name["desc"] = "描敘"
	field2Name["content"] = "內容"
	field2Name["status"] = "狀態"

	return field2Name
}
