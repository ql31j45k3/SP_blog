# configs

# 欄位解說
    system.log.path
        log 檔案路徑，依照時間做切分檔案邏輯與建立一個軟連結的固定檔案名稱

# 欄位預設值
    當設定檔 key 不存在未設定，程式取以下 default 值

    system.log.level: "warn"
    system.log.path: /var/log/ostrich/def_log

    gin.mode: "debug"

    gorm.log.mode: "silent"

# 動態調整欄位
    以下欄位不需重啟服務，程式會重新取得新的值

    system.log.level

# 欄位固定參數值
    system.log.level
        panic、fatal、error、warn、info、debug、trace
    gin.mode
        debug、release、test
    gorm.log.mode
        silent、error、warn、info, default silent

# 欄位單位
    
