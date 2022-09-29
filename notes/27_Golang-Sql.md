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
	- driver包定义了应被数据库驱动实现的接口，这些接口会被sql包使用，绝大多数代码应使用sql包绝大多数代码应使用sql包

### 1. database/sql/driver
- `driver.Driver` 
	- 注册数据库驱动
	- 打开数据库连接

- Conn
	- 把一个查询 query传给Prepare，返回 Stmt(statement)
	- Close关闭数据库连接
	- Begin返回一个事务 Tx(transaction)

- Stmt
	- Close关闭当前的链接状态
	- NumInput返回当前预留参数的个数
	- Exec执行Prepare准备好的 sql，传入参数执行 update/insert 等操作，返回 Result 数据
	- Query执行Prepare准备好的 sql，传入需要的参数执行 select 操作，返回 Rows 结果集	

- Tx
	- Commit提交事务
	- Rollback回滚事务

- Result
	- LastInsertId返回由数据库执行插入操作得到的自增ID号
	- RowsAffected返回操作影响的数据条目数
	- RowsAffected
		- RowsAffected是int64的别名，它实现了Result接口
			- `type RowsAffected int64`
			- `func (RowsAffected) LastInsertId() (int64, error)`
			- `func (v RowsAffected) RowsAffected() (int64, error)`

- Rows
	- Columns是查询所需要的表字段
	- Close关闭迭代器
	- Next返回下一条数据，把数据赋值给dest，dest里面的元素必须是driver.Value的值
		- 如果最后没有数据，Next 函数返回 `io.EOF`
	
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
		Columns() []string                             // func (rs *Rows) Columns() ([]string, error)
		Close() error                                  // func (rs *Rows) Close() error
		Next(dest []Value) error                       // func (rs *Rows) Next() bool
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