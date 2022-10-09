# Golang-Sql  Golang的SQL数据库接口

## 一、Golang的标准库 sql包

- database/sql
	- sql包提供了保证SQL或类SQL数据库的泛用接口
	- Go官方没有提供数据库驱动，而是为开发数据库驱动定义了一些标准接口(即database/sql)，开发者可以根据定义的接口来开发相应的数据库驱动
	- Go语言中支持MySQL的驱动比较多，如
		- `github.com/go-sql-driver/mysql` 支持 database/sql
		- `github.com/ziutek/mymysql` 支持 database/sql，支持自定义的接口
		- `github.com/Philio/GoMySQL` 不支持 database/sql，支持自定义的接口

### 1. DB注册 `Register` 
```go
	// Register 注册并命名一个数据库，可以在Open函数中使用该命名启用该驱动
	func Register(name string, driver driver.Driver)
```

### 2. DB查询结果 `Result`
```go
	// Scanner 接口会被Rows或Row的Scan方法使用
	type Scanner interface {
		// Scan 方法从数据库驱动获取一个值
		// 参数src的类型保证为如下类型之一：
		//    int64
		//    float64
		//    bool
		//    []byte
		//    string
		//    time.Time
		// 如果不能不丢失信息的保存一个值，应返回错误
		Scan(src interface{}) error
	}

	// NullBool、NullInt64、NullFloat64、NullString 实现了Scanner接口，因此可以作为Rows/Row的Scan方法的参数保存扫描结果
	// NullString代表一个可为NULL的字符串
	type NullString struct {
		String string
		Valid  bool // 如果String不是NULL则Valid为真
	}
	// Scan实现了Scanner接口
	func (ns *NullString) Scan(value interface{}) error
	// Value实现了driver.Valuer接口
	func (ns NullString) Value() (driver.Value, error)
	// 示例 example
	// var s NullString
	// err := db.QueryRow("SELECT name FROM foo WHERE id=?", id).Scan(&s)
	// if s.Valid { // use s.String } else { // NULL value }

	// RawBytes是一个字节切片，保管对内存的引用，为数据库自身所使用
	// 在Scaner接口的Scan方法写入RawBytes数据后，该切片只在限次调用Next、Scan或Close方法之前合法
	type RawBytes []byte

	// Result 是对已执行的SQL命令的总结
	type Result interface {
		// LastInsertId 返回一个数据库生成的回应命令的整数
		// 当插入新行时，一般来自一个"自增"列
		// 不是所有的数据库都支持该功能，该状态的语法也各有不同
		LastInsertId() (int64, error)

		// RowsAffected 返回被update、insert或delete命令影响的行数
		// 不是所有的数据库都支持该功能
		RowsAffected() (int64, error)
	}
```

### 3. DB操作 `DB`
- `DB` 是一个数据库(操作)句柄，代表一个具有零到多个底层连接的连接池，它可以安全的被多个go程同时使用
	- sql包会自动创建和释放连接
	- 它也会维护一个闲置连接的连接池
		- 如果数据库具有单连接状态的概念，该状态只有在事务中被观察时才可信
		- 一旦调用了 `BD.Begin`，返回的Tx会绑定到单个连接
		- 当调用事务Tx的 `Commit` 或 `Rollback` 后，该事务使用的连接会归还到DB的闲置连接池中
		- 连接池的大小可以用 `SetMaxIdleConns` 方法控制

- `Open` 打开一个 `dirverName` 指定的数据库，`dataSourceName` 指定数据源，一般包至少括数据库文件名和(可能的)连接信息
	- 大多数用户会通过数据库特定的连接帮助函数打开数据库，返回一个*DB
	- Go标准库中没有数据库驱动，参见 http://golang.org/s/sqldrivers 获取第三方驱动
	- Open函数可能只是验证其参数，而不创建与数据库的连接
	- 如果要检查数据源的名称是否合法，应调用返回值的Ping方法
	- 返回的DB可以安全的被多个go程同时使用，并会维护自身的闲置连接池；这样一来，Open函数只需调用一次，且很少需要关闭DB
```go
	type DB struct { ... }

	func Open(driverName, dataSourceName string) (*DB, error)

	// Driver 方法返回数据库下层驱动
	func (db *DB) Driver() driver.Driver
	// Ping 检查与数据库的连接是否仍有效，如果需要会创建连接
	func (db *DB) Ping() error
	// Close 关闭数据库，释放任何打开的资源
	// 一般不会关闭DB，因为DB句柄通常被多个go程共享，并长期活跃
	func (db *DB) Close() error

	// SetMaxOpenConns 设置与数据库建立连接的最大数目
	// 如果n大于0且小于最大闲置连接数，会将最大闲置连接数减小到匹配最大开启连接数的限制；如果n <= 0，不会限制最大开启连接数，默认为0(无限制)
	func (db *DB) SetMaxOpenConns(n int)
	// SetMaxIdleConns 设置连接池中的最大闲置连接数
	// 如果n大于最大开启连接数，则新的最大闲置连接数会减小到匹配最大开启连接数的限制；如果n <= 0，不会保留闲置连接
	func (db *DB) SetMaxIdleConns(n int)

	// Exec 执行一次命令(包括查询、删除、更新、插入等)，不返回任何执行结果，参数args表示query中的占位参数
	func (db *DB) Exec(query string, args ...interface{}) (Result, error)
	// Query 执行一次查询，返回多行结果(即Rows)，一般用于执行select命令
	func (db *DB) Query(query string, args ...interface{}) (*Rows, error)
	// QueryRow 执行一次查询，并期望返回最多一行结果(即Row)
	// QueryRow 总是返回非nil的值，直到返回值的Scan方法被调用时，才会返回被延迟的错误(如：未找到结果)
	func (db *DB) QueryRow(query string, args ...interface{}) *Row
	// Prepare 创建一个准备好的状态用于之后的查询和命令，返回值可以同时执行多个查询和命令
	func (db *DB) Prepare(query string) (*Stmt, error)
	// Begin 开始一个事务，隔离水平由数据库驱动决定
	func (db *DB) Begin() (*Tx, error)
```

### 4. DB查询 `ROW`
```golang
	// QueryRow 方法返回Row，代表单行查询结果
	type Row struct { ... }

	// Scan 将该行查询结果各列分别保存进dest参数指定的值中
	// 如果该查询匹配多行，Scan会使用第一行结果并丢弃其余各行；如果没有匹配查询的行，Scan会返回ErrNoRows
	func (r *Row) Scan(dest ...interface{}) error

	// Rows 是查询的结果
	// 它的游标指向结果集的第零行，使用Next方法来遍历各行结果
	type Rows struct { ... }

	// Columns 返回列名，如果Rows已经关闭会返回错误
	func (rs *Rows) Columns() ([]string, error)

	// Scan 将当前行各列结果填充进dest指定的各个值中
	func (rs *Rows) Scan(dest ...interface{}) error

	// Next 准备用于Scan方法的下一行结果
	// 如果成功会返回真，如果没有下一行或者出现错误会返回假；Err应该被调用以区分这两种情况
	// 每一次调用Scan方法，甚至包括第一次调用该方法，都必须在前面先调用Next方法
	func (rs *Rows) Next() bool
	// 示例 example
	// rows, err := db.Query("SELECT ...")
	// defer rows.Close()
	// for rows.Next() {
	// 	var id int
	// 	var name string
	// 	err = rows.Scan(&id, &name)
	// }
	// err = rows.Err() // get any error encountered during iteration

	// Close 关闭Rows，阻止对其更多的列举
	// 如果Next方法返回假，Rows会自动关闭，满足检查Err方法结果的条件
	// Close 方法为幂等的(多次调用无效的成功)，不影响Err方法的结果
	func (rs *Rows) Close() error

	// Err 返回可能的、在迭代时出现的错误；Err需在显式或隐式调用Close方法后调用
	func (rs *Rows) Err() error
```

### 5. DB状态 `Stmt`
- `Stmt`是准备好的状态
	- 可以安全的被多个go程同时使用
	- [database/sql: Stmt的使用以及坑](https://studygolang.com/articles/1795)
```go
	type Stmt struct { ... }

	// Exec 使用提供的参数执行准备好的命令状态，返回Result类型的该状态执行结果的总结
	func (s *Stmt) Exec(args ...interface{}) (Result, error)

	// Query 使用提供的参数执行准备好的查询状态，返回Rows类型查询结果
	func (s *Stmt) Query(args ...interface{}) (*Rows, error)

	// QueryRow 使用提供的参数执行准备好的查询状态
	// 如果在执时遇到了错误，该错误会被延迟，直到返回值的Scan方法被调用时才释放，返回值总是非nil的；
	// 如果没有查询到结果，*Row类型返回值的Scan方法会返回ErrNoRows；否则，Scan方法会扫描结果第一行并丢弃其余行
	func (s *Stmt) QueryRow(args ...interface{}) *Row
	// 示例 example
	// var name string
	// err := nameByUseridStmt.QueryRow(id).Scan(&name)

	// Close 关闭状态
	func (s *Stmt) Close() error
```

### 6. DB事物 `Tx`
- `Tx` 代表一个进行中的数据库事务
	- 一次事务必须以对 `Commit` 或 `Rollback` 的调用结束
	- 调用 `Commit` 或 `Rollback` 后，所有对事务的操作都会失败并返回错误值 `ErrTxDone`
```go
	type Tx struct { ... }

	// Exec 执行命令，但不返回结果，例如执行insert和update
	func (tx *Tx) Exec(query string, args ...interface{}) (Result, error)

	// Query执行查询并返回零到多行结果(Rows)，一般执行select命令
	func (tx *Tx) Query(query string, args ...interface{}) (*Rows, error)

	// QueryRow 执行查询并期望返回最多一行结果(Row)
	// QueryRow 总是返回非nil的结果，查询失败的错误会延迟到在调用该结果的Scan方法时释放
	func (tx *Tx) QueryRow(query string, args ...interface{}) *Row

	// Prepare 准备一个专用于该事务的状态
	// 返回的该事务专属状态操作在Tx递交会回滚后不能再使用；要在该事务中使用已存在的状态，参见Tx.Stmt方法
	func (tx *Tx) Prepare(query string) (*Stmt, error)

	// Stmt使用已存在的状态生成一个该事务特定的状态
	func (tx *Tx) Stmt(stmt *Stmt) *Stmt
	// 示例 example
	// updateMoney, err := db.Prepare("UPDATE balance SET money=money+? WHERE id=?")
	// tx, err := db.Begin()
	// res, err := tx.Stmt(updateMoney).Exec(123.45, 98293203)

	// Stmt 使用已存在的状态生成一个该事务特定的状态
	func (tx *Tx) Commit() error

	// Commit 递交事务
	func (tx *Tx) Commit() error

	// Rollback 放弃并回滚事务
	func (tx *Tx) Rollback() error
```

## 二、Golang的标准库 driver包

- database/sql/driver
	- driver包定义了应被数据库驱动实现的接口，这些接口会被sql包使用
	- 绝大多数代码应使用sql包

### 1. ValueConverter
- ValueConverter
	- ValueConverter接口提供了ConvertValue方法，把数据库里的数据类型转换成Value允许的数据类型
	- driver包提供了各种ValueConverter接口的实现，以保证不同驱动之间的实现和转换的一致性
	- ValueConverter接口有如下用途
		- 转换sql包提供的Value类型值到数据库指定列的类型，并保证它的匹配，例如保证某个int64值满足一个表的uint16列
		- 转换数据库提供的值到驱动的Value类型
		- 在扫描时被sql包用于将驱动的Value类型转换为用户的类型
```go
	type ValueConverter interface {
		// ConvertValue将一个值转换为驱动支持的Value类型
		ConvertValue(v interface{}) (Value, error)
	}
```

- ColumnConverter
	- 如果Stmt有自己的列类型，可以实现ColumnConverter接口，返回值可以将任意类型转换为驱动的Value类型
```go
	type ColumnConverter interface {
		// ColumnConverter返回指定列的ValueConverter
		// 如果该列未指定类型，或不应特殊处理，应返回DefaultValueConverter
		ColumnConverter(idx int) ValueConverter
	}
```

- NotNull
	- NotNull实现了ValueConverter接口，不允许nil值，否则会将值交给Converter字段处理
```go
	type NotNull struct {
		Converter ValueConverter
	}

	func (n NotNull) ConvertValue(v interface{}) (Value, error)
```

- Null
	- Null实现了ValueConverter接口，允许nil值，否则会将值交给Converter字段处理
```go
	type Null struct {
		Converter ValueConverter
	}

	func (n Null) ConvertValue(v interface{}) (Value, error)
```

### 2. Variables
- ValueConverter接口转化规则
	- 用于将输入的值转换为对应类型，如 bool、init32、string
	- 布尔类型
		- 不做修改
	- 整数类型
		- 1 为真
		- 0 为假
		- 其余整数会导致错误
	- 字符串和`[]byte`
		- 与 `strconv.ParseBool` 的规则相同
	- 所有其他类型都会导致错误
```go
	// Bool 是ValueConverter接口值，用于将输入的值转换为布尔类型
	var Bool boolType

	// Int32 是一个ValueConverter接口值，用于将值转换为int64类型，会尊重int32类型的限制
	var Int32 int32Type

	// String 是一个ValueConverter接口值，用于将值转换为字符串。
	// 如果值v是字符串或者[]byte类型，不会做修改，如果值v是其它类型，会转换为fmt.Sprintf("%v", v)
	var String stringType

	// DefaultParameterConverter 是ValueConverter接口的默认实现，当一个Stmt没有实现ColumnConverter时，就会使用它
	var DefaultParameterConverter defaultConverter
	// 如果值value满足函数IsValue(value)为真，DefaultParameterConverter直接返回 value
	// 否则，整数类型会被转换为int64，浮点数转换为float64，字符串转换为[]byte，其它类型会导致错误

	// ResultNoRows是预定义的Result类型值，用于当一个DDL命令(如create table)成功时被驱动返回
	// 它的LastInsertId和RowsAffected方法都返回错误
	var ResultNoRows noRows

	// ErrBadConn应被驱动返回，以通知sql包一个driver.Conn处于损坏状态(如服务端之前关闭了连接)，sql包会重启一个新的连接
	var ErrBadConn = errors.New("driver: bad connection")

	// ErrSkip可能会被某些可选接口的方法返回，用于在运行时表明快速方法不可用，sql包应像未实现该接口的情况一样执行
	var ErrSkip = errors.New("driver: skip fast-path; continue as if unimplemented")
```

### 3. Value
- Value
	- Value是驱动必须能处理的值
	- 它要么是nil，要么是如下类型的实例
		- `int64`
		- `float64`
		- `bool`
		- `[]byte`
		- `string`   [*] Rows.Next不会返回该类型值
		- `time.Time` 
```go
	type Value interface{}
```

- Valuer
	- Valuer是提供Value方法的接口
	- 实现了Valuer接口的类型可以将自身转换为驱动支持的Value类型值
```go
	type Valuer interface {
		// Value返回一个驱动支持的Value类型值
		Value() (Value, error)
	}
```

- IsValue
	- IsValue报告v是否是合法的Value类型参数
	- 和IsScanValue不同，IsValue接受字符串类型
```go
	func IsValue(v interface{}) bool
```

- IsScanValue
	- IsScanValue报告v是否是合法的Value扫描类型参数
	- 和IsValue不同，IsScanValue不接受字符串类型
```go
	func IsScanValue(v interface{}) bool
```

### 4. Driver
- Driver 
	- Driver接口必须被数据库驱动实现
	- 注册数据库驱动
	- 打开数据库连接
```go
	type Driver interface {
		// Open返回一个新的与数据库的连接，参数name的格式是驱动特定的
		//
		// Open可能返回一个缓存的连接(之前关闭的连接)，但这么做是不必要的
		// sql包会维护闲置连接池以便有效的重用连接
		//
		// 返回的连接同一时间只会被一个go程使用
		Open(name string) (Conn, error)
	}
```

### 5. Conn
- Conn
	- Conn是与数据库的连接，该连接不会被多线程并行使用，且连接被假定为具有状态的
	- 把一个查询 `query` 传给 `Prepare`，返回 `Stmt` (statement)
	- `Close` 关闭数据库连接
	- `Begin` 返回一个事务 `Tx` (transaction)
```go
	type Conn interface {
		// Prepare返回一个准备好的、绑定到该连接的状态
		Prepare(query string) (Stmt, error)

		// Close作废并停止任何现在准备好的状态和事务，将该连接标注为不再使用
		//
		// 因为sql包维护着一个连接池，只有当闲置连接过剩时才会调用Close方法
		// 驱动的实现中不需要添加自己的连接缓存池
		Close() error

		// Begin开始并返回一个新的事务
		Begin() (Tx, error)
	}
```

- Execer
	- Execer是一个可选的、可能被Conn接口实现的接口
	- 如果一个Conn未实现 Execer接口，sql包的`DB.Exec`会首先准备一个查询，执行状态，然后关闭状态
	- `Query` 可能会返回 `ErrSkip`
```go
	type Execer interface {
		Exec(query string, args []Value) (Result, error)
	}
```

- Queryer
	- Queryer是一个可选的、可能被Conn接口实现的接口
	- 如果一个Conn未实现 Queryer接口，sql包的`DB.Exec`会首先准备一个查询，执行状态，然后关闭状态
	- `Exec` 可能会返回 `ErrSkip`
```go
	type Queryer interface {
		Query(query string, args []Value) (Rows, error)
	}
```

### 6. Stmt
- Stmt
	- Stmt是准备好的状态，它会绑定到一个连接，不应被多go程同时使用
	- `Close` 关闭当前的链接状态
	- `NumInput` 返回当前预留参数的个数
	- `Exec` 执行Prepare准备好的 sql，传入参数执行 update/insert 等操作，返回 `Result` 数据
	- `Query` 执行Prepare准备好的 sql，传入需要的参数执行 select 操作，返回 `Rows` 结果集	
```go
	type Stmt interface {
		// Close关闭Stmt
		//
		// 和Go1.1一样，如果Stmt被任何查询使用中的话，将不会被关闭
		Close() error

		// NumInput返回占位参数的个数。
		//
		// 如果NumInput返回值 >= 0，sql包会提前检查调用者提供的参数个数，
		// 并且会在调用Exec或Query方法前返回数目不对的错误。
		//
		// NumInput可以返回-1，如果驱动占位参数的数量
		// 此时sql包不会提前检查参数个数
		NumInput() int

		// Exec执行查询，而不会返回结果，如insert或update
		Exec(args []Value) (Result, error)

		// Query执行查询并返回结果，如select
		Query(args []Value) (Rows, error)
	}
```

### 7. Tx
- Tx
	- Tx是一次事务
	- `Commit` 提交事务
	- `Rollback` 回滚事务
```go
	type Tx interface {
		Commit() error
		Rollback() error
	}
```

### 8. Result
- Result
	- `LastInsertId` 返回由数据库执行插入操作得到的自增ID号
	- `RowsAffected` 返回操作影响的数据条目数
```go
	type Result interface {
		// LastInsertId返回insert等命令后数据库自动生成的ID
		LastInsertId() (int64, error)

		// RowsAffected返回被查询影响的行数
		RowsAffected() (int64, error)
	}

	// RowsAffected实现了Result接口，用于insert或update操作，这些操作会修改零到多行数据
	type RowsAffected int64
	func (v RowsAffected) LastInsertId() (int64, error)
	func (v RowsAffected) RowsAffected() (int64, error)
```

### 9. Rows
- Rows
	- Rows是执行查询得到的结果的迭代器
	- `Columns` 是查询所需要的表字段
	- `Close` 关闭迭代器
	- `Next` 返回下一条数据，把数据赋值给dest，dest里面的元素必须是 `driver.Value` 的值
		- 如果最后没有数据，Next 函数返回 `io.EOF`
```go
	type Rows interface {
		// Columns返回各列的名称，列的数量可以从切片长度确定
		// 如果某个列的名称未知，对应的条目应为空字符串
		Columns() []string

		// Close关闭Rows
		Close() error

		// 调用Next方法以将下一行数据填充进提供的切片中
		// 提供的切片必须和Columns返回的切片长度相同
		//
		// 切片dest可能被填充同一种驱动Value类型，但字符串除外
		// 所有string值都必须转换为[]byte
		//
		// 当没有更多行时，Next应返回io.EOF
		Next(dest []Value) error
	}
```

- Value
	- `driver.ValueConverter` 把数据库里的数据类型转换成Value允许的数据类型
```go
	// driver.Driver
	package driver // import "database/sql/driver"
	type Driver interface { 
		Open(name string) (Conn, error)                        // func Open(driverName, dataSourceName string) (*DB, error)
	}

	// 注册数据库驱动
	var d = Driver{proto: "tcp", raddr: "127.0.0.1:3306"}
	sql.Register("mysql", &d)                                  // func Register(name string, driver driver.Driver)

	// 打开数据库连接
	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/test?charset=utf8")

	// driver.Conn
	type Conn interface {
		Prepare(query string) (Stmt, error)                    // func (db *DB) Prepare(query string) (*Stmt, error)
		Close() error                                          // func (c *Conn) Close() error
		Begin() (Tx, error)                                    // func (db *DB) Begin() (*Tx, error)
	}

	// driver.Stmt
	type Stmt interface {
		Close() error                                          // func (s *Stmt) Close() error
		NumInput() int
		Exec(args []Value) (Result, error)                     // func (s *Stmt) Exec(args ...interface{}) (Result, error)
		Query(args []Value) (Rows, error)                      // func (s *Stmt) Query(args ...interface{}) (*Rows, error)
	}

	// driver.Tx
	type Tx interface {
		Commit() error                                          // func (tx *Tx) Commit() error
		Rollback() error                                        // func (tx *Tx) Rollback() error
	}

	// driver.Result
	type Result interface {
		LastInsertId() (int64, error)
		RowsAffected() (int64, error)
	}

	// driver.Rows
	type Rows interface {
		Columns() []string                                      // func (rs *Rows) Columns() ([]string, error)
		Close() error                                           // func (rs *Rows) Close() error
		Next(dest []Value) error                                // func (rs *Rows) Next() bool
	}

	// driver.Value
		// nil           // 要么是 nil，要么是下面的任意一种
		// int64 
		// float64 
		// bool 
		// []byte 
		// string 
		// time.Time
	type Value interface{}

	// driver.ValueConverter
	type ValueConverter interface {
		ConvertValue(v interface{}) (Value, error)
	}
```
