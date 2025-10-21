# Toolbox 使用者手冊

這是一份針對 `toolbox.exe` 這個自製工具箱的官方使用者手冊。本文件詳細說明了所有可用的指令、參數及其使用方法。

## 指令總覽

| 指令 | 功能描述 |
| --- | --- |
| `help` | 顯示主說明或特定指令的詳細說明。 |
| `copy` | 複製檔案或資料夾。 |
| `remove` | 刪除檔案或資料夾。 |
| `clear` | 清空指定的資料夾。 |
| `randomString` | 產生指定長度和數量的隨機字串。 |
| `mssqlCompareTableData` | 比較兩個 MSSQL 資料表的資料是否完全相同。 |
| `mssqlGetDatabaseSchema` | 取得指定 MSSQL 資料庫的 Schema，並可匯出成 JSON 或 Markdown 格式。 |

---

## `copy`

### 功能

複製檔案或資料夾。這個指令可以處理以下幾種情況：

- 檔案到檔案
- 檔案到資料夾
- 資料夾到資料夾

**注意：** 不支援將資料夾複製到檔案。

### 用法

```bash
toolbox.exe copy <source> <destination>
```

### 參數

| 參數 | 類型 | 描述 |
| --- | --- | --- |
| `<source>` | `string` | 指定要複製的來源檔案或資料夾路徑。 |
| `<destination>` | `string` | 指定要複製到的目的地檔案或資料夾路徑。 |

### 範例

**1. 複製檔案到另一個檔案**

```bash
toolbox.exe copy file1.txt file2.txt
```

**2. 複製檔案到資料夾**

```bash
toolbox.exe copy file1.txt ./my_folder/
```

**3. 複製資料夾到另一個資料夾**

```bash
toolbox.exe copy ./folder1/ ./folder2/
```

---

## `remove`

### 功能

刪除指定的檔案或資料夾。這個指令會遞迴地刪除所有內容，類似於 `rm -rf`。

### 用法

```bash
toolbox.exe remove <target>
```

### 參數

| 參數 | 類型 | 描述 |
| --- | --- | --- |
| `<target>` | `string` | 指定要刪除的檔案或資料夾路徑。 |

### 範例

**1. 刪除檔案**

```bash
toolbox.exe remove file_to_delete.txt
```

**2. 刪除資料夾**

```bash
toolbox.exe remove ./folder_to_delete/
```

---

## `clear`

### 功能

清空指定的資料夾，也就是刪除資料夾內的所有檔案和子資料夾，但保留該資料夾本身。

### 用法

```bash
toolbox.exe clear <target>
```

### 參數

| 參數 | 類型 | 描述 |
| --- | --- | --- |
| `<target>` | `string` | 指定要清空的資料夾路徑。 |

### 範例

```bash
toolbox.exe clear ./logs/
```

---

## `randomString`

### 功能

產生隨機字串。使用者可以自訂字串的來源字元、長度，以及要產生的數量。

### 用法

```bash
toolbox.exe randomString [flags]
```

### 參數

| 參數 | 類型 | 預設值 | 描述 |
| --- | --- | --- | --- |
| `--src` | `string` | `0123...xyz` | 指定用來產生隨機字串的字元集合。 |
| `--length` | `int` | `8` | 指定每個隨機字串的長度。 |
| `--count` | `int` | `1` | 指定要產生的隨機字串數量。 |

### 範例

**1. 產生一個長度為 12 的隨機字串**

```bash
toolbox.exe randomString --length 12
```

**2. 產生 5 個長度為 6 且只包含小寫英文字母的隨機字串**

```bash
toolbox.exe randomString --src "abcdefghijklmnopqrstuvwxyz" --length 6 --count 5
```

---

## `mssqlCompareTableData`

### 功能

比較兩個 MSSQL 資料表的資料是否完全相等。此指令只會比較兩個資料表中共同擁有的欄位。

### 用法

```bash
toolbox.exe mssqlCompareTableData <setting_file_path>
```

### 參數

| 參數 | 類型 | 描述 |
| --- | --- | --- |
| `<setting_file_path>` | `string` | 指定一個 JSON 格式的設定檔，其中包含了兩個資料庫的連線資訊以及要比較的資料表清單。 |

### 設定檔範例

```json
{
    "ipL": "127.0.0.1",
    "portL": "1433",
    "accL": "sa",
    "pwL": "<YourStrong@Passw0rd>",
    "dbnameL": "Database1",
    "ipR": "127.0.0.1",
    "portR": "1433",
    "accR": "sa",
    "pwR": "<YourStrong@Passw0rd>",
    "dbnameR": "Database2",
    "tables": ["TableA", "TableB", "TableC"]
}
```

---

## `mssqlGetDatabaseSchema`

### 功能

取得指定 MSSQL 資料庫的 Schema 資料，並可將結果匯出成 JSON 或 Markdown 檔案。

### 用法

```bash
toolbox.exe mssqlGetDatabaseSchema [flags]
```

### 參數

| 參數 | 類型 | 預設值 | 描述 |
| --- | --- | --- | --- |
| `--ip` | `string` | | **(必填)** 資料庫伺服器的 IP 位址。 |
| `--port` | `string` | | **(必填)** 資料庫伺服器的埠號。 |
| `--acc` | `string` | | **(必填)** 登入帳號。 |
| `--pw` | `string` | | **(必填)** 登入密碼。 |
| `--dbname` | `string` | | **(必填)** 要存取的資料庫名稱。 |
| `--json` | `bool` | `true` | 是否產生 JSON 格式的 Schema 檔案。 |
| `--md` | `bool` | `false` | 是否產生 Markdown 格式的 Schema 檔案。 |

### 範例

```bash
toolbox.exe mssqlGetDatabaseSchema --ip "127.0.0.1" --port "1433" --acc "sa" --pw "<YourStrong@Passw0rd>" --dbname "MyDatabase" --md true
```
