# modules
商業邏輯模組功能，依照模塊區分，一個模塊包含功能所有處理 API -> 邏輯處理 -> 資料操作處理

e.g. 文章模組
```
 - modules
   - article
      api.go (API 註冊)
      use_case.go (邏輯處理)
      model.go (資料模型)
      reposity.go (資料相關操作處理)
      common.go (共用的內部函式)
```