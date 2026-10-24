/**

just use postgres - chapter 11
enqueue, dequeue

**/

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

// import (
// 	// rq "github.com/rabbitmq/amqp091-go"
// )

// func main() {
// 	conn, err := rq.Dial("amqp://admin:admin@localhost:5672/")
// 	if err != nil {
// 		panic(err)
// 	}

// 	defer conn.Close()

// 	ch, err := conn.Channel()
// 	if err != nil {
// 		panic(err)
// 	}

// 	defer ch.Close()

// 	ch.ExchangeDeclare()
// }

// func migration(dbSource string) error {
// 	appDB, err := sql.Open("postgres", dbSource)
// 	if err != nil {
// 		return err
// 	}
// 	defer appDB.Close()

// 	driver, err := postgres.WithInstance(appDB, &postgres.Config{})
// 	if err != nil {
// 		return err
// 	}

// 	m, err := migrate.NewWithDatabaseInstance("file://migrations", "postgres", driver)
// 	if err != nil {
// 		return err
// 	}

// 	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
// 		return err
// 	}
// 	return nil
// }

func createSchema(ctx context.Context, conn *sqlx.DB, name string) {
	schemaCreateQuery := fmt.Sprintf(`create schema if not exists %s`, name)
	_, err := conn.ExecContext(ctx, schemaCreateQuery)
	if err != nil {
		panic(err)
	}
}

func listSchema(ctx context.Context, conn *sqlx.DB) {
	var schemas []string
	err := conn.SelectContext(ctx, &schemas, `
        SELECT schema_name
        FROM information_schema.schemata
        WHERE schema_name NOT IN ('pg_catalog', 'information_schema')
          AND schema_name NOT LIKE 'pg_toast%'
        ORDER BY schema_name
    `)
	if err != nil {
		panic(err)
	}

	fmt.Println("Schemas::List::")
	fmt.Println("---------------")
	for _, name := range schemas {
		fmt.Println(name)
	}
	fmt.Println("---------------")
}

func listSchemaTables(ctx context.Context, conn *sqlx.DB, schemaName string) {
	query := `SELECT tablename
		FROM pg_catalog.pg_tables
		WHERE schemaname = $1
		ORDER BY tablename`

	var tables []string
	err := conn.SelectContext(ctx, &tables, query, schemaName)
	if err != nil {
		panic(err)
	}

	fmt.Println("Schema(" + schemaName + ")::TABLE::List::")
	fmt.Println("---------------")
	for _, name := range tables {
		fmt.Println(name)
	}
	fmt.Println("---------------")
}

func createEnum(ctx context.Context, conn *sqlx.DB, schemaName string, name string, values []string) {
	var exists bool
	err := conn.GetContext(ctx, &exists, `
		SELECT EXISTS (
			SELECT 1 FROM pg_type t
			JOIN pg_namespace n ON n.oid = t.typnamespace
			WHERE t.typname = $1 AND n.nspname = $2
		)
	`, name, schemaName)
	if err != nil {
		panic(err)
	}
	if exists {
		return
	}

	quotedValues := make([]string, len(values))
	for i, v := range values {
		quotedValues[i] = pq.QuoteLiteral(v)
	}

	query := fmt.Sprintf(
		"CREATE TYPE %s.%s AS ENUM (%s)",
		pq.QuoteIdentifier(schemaName),
		pq.QuoteIdentifier(name),
		strings.Join(quotedValues, ", "),
	)

	_, err = conn.ExecContext(ctx, query)
	if err != nil {
		panic(err)
	}
}

type EnumInfo struct {
	Name   string         `db:"enum_name"`
	Values pq.StringArray `db:"values"`
}

func listEnums(ctx context.Context, db *sqlx.DB, schema string) {
	var enums []EnumInfo
	err := db.SelectContext(ctx, &enums, `
        SELECT t.typname AS enum_name,
               array_agg(e.enumlabel ORDER BY e.enumsortorder) AS values
        FROM pg_type t
        JOIN pg_enum e ON e.enumtypid = t.oid
        JOIN pg_namespace n ON n.oid = t.typnamespace
        WHERE n.nspname = $1
        GROUP BY t.typname
        ORDER BY t.typname
    `, schema)

	if err != nil {
		panic(err)
	}
	fmt.Println("Schema(" + schema + ")::ENUM::List::")
	fmt.Println("---------------")
	for _, en := range enums {
		fmt.Println(en.Name)
	}
	fmt.Println("---------------")
}

func createTable(ctx context.Context, conn *sqlx.DB, schema string) {
	query := fmt.Sprintf(`
	create table if not exists %s.queue(
		id bigserial primary key,
		message json not null,
		created_at timestamptz default now(),
		status %s.%s not null default %s
	)
	`, pq.QuoteIdentifier(schema), pq.QuoteIdentifier(schema), pq.QuoteIdentifier("status"), pq.QuoteLiteral("new"))

	_, err := conn.ExecContext(ctx, query)
	if err != nil {
		panic(err)
	}
}

func createEnqueueFunction(ctx context.Context, conn *sqlx.DB, channel string) {
	query := fmt.Sprintf(`
	create or replace function mqueue.enqueue(msg json)
	returns BIGINT as $$
	DECLARE
    	new_id BIGINT;
	begin
		insert into mqueue.queue(message) values(msg)
		RETURNING id INTO new_id;

		PERFORM pg_notify(%s, new_id::text);

		RETURN new_id;
	end;
	$$ language plpgsql;
	`, pq.QuoteLiteral(channel))

	_, err := conn.ExecContext(ctx, query)
	if err != nil {
		panic(err)
	}
}

func enqueue(ctx context.Context, conn *sqlx.DB, msg any) {
	payload, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}

	var id int64
	err = conn.GetContext(ctx, &id, `SELECT mqueue.enqueue($1)`, payload)
	if err != nil {
		panic(err)
	}
	fmt.Println("QueueId: ", id)
}

func createDequeueFunction(ctx context.Context, conn *sqlx.DB) {
	query := `
	create or replace function mqueue.dequeue(cnt int)
	returns table(
		msgid bigint,
		message json
	) as $$
	begin
		return query
		with new_messages as (
			select id from mqueue.queue
			where status = 'new' order by created_at
			for update skip locked
			limit cnt
		)
		update mqueue.queue q
		set status='processing'
		from new_messages where q.id = new_messages.id
		returning q.id, q.message;
	end;
	$$ language plpgsql;
	`

	_, err := conn.ExecContext(ctx, query)
	if err != nil {
		panic(err)
	}
}

type QueueRow struct {
	MsgID   int64           `db:"msgid"`
	Message json.RawMessage `db:"message"`
}

type Payload struct {
	Name string `json:"name"`
	Dept string `json:"dept"`
}

func dequeue(ctx context.Context, conn *sqlx.DB, cnt int) {
	var result []QueueRow
	err := conn.SelectContext(ctx, &result, `select * from mqueue.dequeue($1)`, cnt)
	if err != nil {
		panic(err)
	}

	for _, r := range result {
		var p Payload
		if err := json.Unmarshal(r.Message, &p); err != nil {
			fmt.Println("bad payload for msg", r.MsgID, err)
			continue
		}
		fmt.Println("MSG:", r.MsgID, p.Name, p.Dept)
	}
}

type FunctionInfo struct {
	FunctionName string `db:"function_name"`
	Arguments    string `db:"arguments"`
	ReturnType   string `db:"return_type"`
}

func listFunctions(ctx context.Context, db *sqlx.DB, schema string) {
	schema = "mqueue"
	var fns []FunctionInfo
	err := db.SelectContext(ctx, &fns, `
        SELECT p.proname AS function_name,
               pg_get_function_arguments(p.oid) AS arguments,
               pg_get_function_result(p.oid) AS return_type
        FROM pg_catalog.pg_proc p
        JOIN pg_catalog.pg_namespace n ON n.oid = p.pronamespace
        WHERE n.nspname = $1
        ORDER BY p.proname
    `, schema)
	if err != nil {
		panic(err)
	}

	fmt.Println("Schema(" + schema + ")::FUNC::List::")
	fmt.Println("---------------")
	for _, en := range fns {
		fmt.Println(en.FunctionName)
	}
	fmt.Println("---------------")
}

func main() {
	conn, err := sqlx.Connect("postgres", "postgres://mq_user:mq_pass@localhost:5432/mq_db?sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := conn.PingContext(ctx); err != nil {
		panic(err)
	}

	listner := pq.NewListener("postgres://mq_user:mq_pass@localhost:5432/mq_db?sslmode=disable", time.Second, 15*time.Second, nil)

	defer listner.Close()

	if err := listner.Listen("queue"); err != nil {
		panic(err)
	}

	if err := listner.Ping(); err != nil {
		panic(err)
	}

	go func(ctx context.Context, listner *pq.Listener, conn *sqlx.DB) {
		select {
		case <-ctx.Done():
			return
		case msg := <-listner.Notify:
			fmt.Println("GET ", msg.Channel, msg.Extra)
			dequeue(ctx, conn, 15)
		}
	}(ctx, listner, conn)

	createSchema(ctx, conn, "mqueue")
	listSchema(ctx, conn)
	listSchemaTables(ctx, conn, "mqueue")
	createEnum(ctx, conn, "mqueue", "status", []string{"new", "processing", "completed"})
	listEnums(ctx, conn, "mqueue")
	createTable(ctx, conn, "mqueue")
	createEnqueueFunction(ctx, conn, "queue")
	createDequeueFunction(ctx, conn)
	listFunctions(ctx, conn, "mqueue")
	enqueue(ctx, conn, map[string]string{
		"name": "labib",
		"dept": "cse",
	})

	time.Sleep(5 * time.Second)
}
