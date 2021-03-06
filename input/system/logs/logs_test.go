package logs_test

import (
	"testing"

	"github.com/kylelemons/godebug/pretty"
	"github.com/pganalyze/collector/input/system/logs"
	"github.com/pganalyze/collector/output/pganalyze_collector"
	"github.com/pganalyze/collector/state"
	uuid "github.com/satori/go.uuid"
)

type testpair struct {
	logLinesIn  []state.LogLine
	logLinesOut []state.LogLine
	samplesOut  []state.PostgresQuerySample
}

var tests = []testpair{
	// Statement duration
	{
		[]state.LogLine{{
			Content: "duration: 3205.800 ms execute a2: SELECT \"servers\".* FROM \"servers\" WHERE \"servers\".\"id\" = $1 LIMIT $2",
		}},
		[]state.LogLine{{
			Query:          "SELECT \"servers\".* FROM \"servers\" WHERE \"servers\".\"id\" = $1 LIMIT $2",
			Classification: pganalyze_collector.LogLineInformation_STATEMENT_DURATION,
			Details:        map[string]interface{}{"duration_ms": 3205.8},
		}},
		[]state.PostgresQuerySample{{
			Query:     "SELECT \"servers\".* FROM \"servers\" WHERE \"servers\".\"id\" = $1 LIMIT $2",
			RuntimeMs: 3205.8,
		}},
	},
	{
		[]state.LogLine{{
			Content: "duration: 3205.800 ms execute a2: SELECT ...[Your log message was truncated]",
		}},
		[]state.LogLine{{
			Classification: pganalyze_collector.LogLineInformation_STATEMENT_DURATION,
			Details:        map[string]interface{}{"truncated": true},
		}},
		nil,
	},
	// Connects/Disconnects
	{
		[]state.LogLine{{
			Content: "connection received: host=172.30.0.165 port=56902",
		}, {
			Content: "connection authorized: user=myuser database=mydb SSL enabled (protocol=TLSv1.2, cipher=ECDHE-RSA-AES256-GCM-SHA384, compression=off)",
		}, {
			Content: "pg_hba.conf rejects connection for host \"172.1.0.1\", user \"myuser\", database \"mydb\", SSL on",
		}, {
			Content:  "no pg_hba.conf entry for host \"8.8.8.8\", user \"postgres\", database \"postgres\", SSL off",
			LogLevel: pganalyze_collector.LogLineInformation_FATAL,
		}, {
			Content:  "password authentication failed for user \"postgres\"",
			LogLevel: pganalyze_collector.LogLineInformation_FATAL,
			UUID:     uuid.UUID{1},
		}, {
			Content:  "Connection matched pg_hba.conf line 4: \"hostssl postgres        postgres        0.0.0.0/0               md5\"",
			LogLevel: pganalyze_collector.LogLineInformation_DETAIL,
		}, {
			Content:  "database \"template0\" is not currently accepting connections",
			LogLevel: pganalyze_collector.LogLineInformation_FATAL,
		}, {
			Content:  "role \"abc\" is not permitted to log in",
			LogLevel: pganalyze_collector.LogLineInformation_FATAL,
		}, {
			Content: "disconnection: session time: 1:53:01.198 user=myuser database=mydb host=172.30.0.165 port=56902",
		}},
		[]state.LogLine{{
			Classification: pganalyze_collector.LogLineInformation_CONNECTION_RECEIVED,
		}, {
			Classification: pganalyze_collector.LogLineInformation_CONNECTION_AUTHORIZED,
		}, {
			Classification: pganalyze_collector.LogLineInformation_CONNECTION_REJECTED,
		}, {
			Classification: pganalyze_collector.LogLineInformation_CONNECTION_REJECTED,
			LogLevel:       pganalyze_collector.LogLineInformation_FATAL,
		}, {
			Classification: pganalyze_collector.LogLineInformation_CONNECTION_REJECTED,
			LogLevel:       pganalyze_collector.LogLineInformation_FATAL,
			UUID:           uuid.UUID{1},
		}, {
			LogLevel:   pganalyze_collector.LogLineInformation_DETAIL,
			ParentUUID: uuid.UUID{1},
		}, {
			Classification: pganalyze_collector.LogLineInformation_CONNECTION_REJECTED,
			LogLevel:       pganalyze_collector.LogLineInformation_FATAL,
		}, {
			Classification: pganalyze_collector.LogLineInformation_CONNECTION_REJECTED,
			LogLevel:       pganalyze_collector.LogLineInformation_FATAL,
		}, {
			Classification: pganalyze_collector.LogLineInformation_CONNECTION_DISCONNECTED,
			Details:        map[string]interface{}{"session_time_secs": 6781.198},
		}},
		nil,
	}, {
		[]state.LogLine{{
			Content: "incomplete startup packet",
		}},
		[]state.LogLine{{
			Classification: pganalyze_collector.LogLineInformation_CONNECTION_CLIENT_FAILED_TO_CONNECT,
		}},
		nil,
	}, {
		[]state.LogLine{{
			Content: "could not receive data from client: Connection reset by peer",
		}},
		[]state.LogLine{{
			Classification: pganalyze_collector.LogLineInformation_CONNECTION_LOST,
		}},
		nil,
	}, {
		[]state.LogLine{{
			Content: "could not send data to client: Broken pipe",
		}},
		[]state.LogLine{{
			Classification: pganalyze_collector.LogLineInformation_CONNECTION_LOST,
		}},
		nil,
	}, {
		[]state.LogLine{{
			Content:  "terminating connection because protocol synchronization was lost",
			LogLevel: pganalyze_collector.LogLineInformation_FATAL,
		}},
		[]state.LogLine{{
			LogLevel:       pganalyze_collector.LogLineInformation_FATAL,
			Classification: pganalyze_collector.LogLineInformation_CONNECTION_LOST,
		}},
		nil,
	}, {
		[]state.LogLine{{
			Content:  "connection to client lost",
			LogLevel: pganalyze_collector.LogLineInformation_FATAL,
		}},
		[]state.LogLine{{
			Classification: pganalyze_collector.LogLineInformation_CONNECTION_LOST,
			LogLevel:       pganalyze_collector.LogLineInformation_FATAL,
		}},
		nil,
	}, {
		[]state.LogLine{{
			Content: "unexpected EOF on client connection",
		}},
		[]state.LogLine{{
			Classification: pganalyze_collector.LogLineInformation_CONNECTION_LOST,
		}},
		nil,
	}, {
		[]state.LogLine{{
			Content:  "unexpected EOF on client connection with an open transaction",
			LogLevel: pganalyze_collector.LogLineInformation_LOG,
		}},
		[]state.LogLine{{
			Classification: pganalyze_collector.LogLineInformation_CONNECTION_LOST_OPEN_TX,
			LogLevel:       pganalyze_collector.LogLineInformation_LOG,
		}},
		nil,
	}, {
		[]state.LogLine{{
			Content:  "remaining connection slots are reserved for non-replication superuser connections",
			LogLevel: pganalyze_collector.LogLineInformation_FATAL,
		}},
		[]state.LogLine{{
			LogLevel:       pganalyze_collector.LogLineInformation_FATAL,
			Classification: pganalyze_collector.LogLineInformation_OUT_OF_CONNECTIONS,
		}},
		nil,
	}, {
		[]state.LogLine{{
			Content:  "too many connections for role \"postgres\"",
			LogLevel: pganalyze_collector.LogLineInformation_FATAL,
		}},
		[]state.LogLine{{
			LogLevel:       pganalyze_collector.LogLineInformation_FATAL,
			Classification: pganalyze_collector.LogLineInformation_TOO_MANY_CONNECTIONS_ROLE,
		}},
		nil,
	}, {
		[]state.LogLine{{
			Content:  "terminating connection due to administrator command",
			LogLevel: pganalyze_collector.LogLineInformation_FATAL,
		}},
		[]state.LogLine{{
			LogLevel:       pganalyze_collector.LogLineInformation_FATAL,
			Classification: pganalyze_collector.LogLineInformation_CONNECTION_TERMINATED,
		}},
		nil,
	},
	// Checkpoints
	{
		[]state.LogLine{{
			Content: "checkpoint starting: xlog",
		}},
		[]state.LogLine{{
			Classification: pganalyze_collector.LogLineInformation_CHECKPOINT_STARTING,
			Details:        map[string]interface{}{"reason": "xlog"},
		}},
		nil,
	}, {
		[]state.LogLine{{
			Content: "checkpoint complete: wrote 111906 buffers (10.9%); 0 transaction log file(s) added, 22 removed, 29 recycled; write=215.895 s, sync=0.014 s, total=216.130 s; sync files=94, longest=0.014 s, average=0.000 s; distance=850730 kB, estimate=910977 kB",
		}},
		[]state.LogLine{{
			Classification: pganalyze_collector.LogLineInformation_CHECKPOINT_COMPLETE,
			Details: map[string]interface{}{
				"bufs_written_pct": 10.9, "write_secs": 215.895, "sync_secs": 0.014,
				"total_secs": 216.130, "longest_secs": 0.014, "average_secs": 0.0,
				"bufs_written": 111906, "segs_added": 0, "segs_removed": 22, "segs_recycled": 29,
				"sync_rels": 94, "distance_kb": 850730, "estimate_kb": 910977,
			},
		}},
		nil,
	}, { // Pre 9.5 syntax (without distance/estimate)
		[]state.LogLine{{
			Content: "checkpoint complete: wrote 15047 buffers (1.4%); 0 transaction log file(s) added, 0 removed, 30 recycled; write=68.980 s, sync=1.542 s, total=70.548 s; sync files=925, longest=0.216 s, average=0.001 s",
		}},
		[]state.LogLine{{
			Classification: pganalyze_collector.LogLineInformation_CHECKPOINT_COMPLETE,
			Details: map[string]interface{}{
				"bufs_written": 15047, "segs_added": 0, "segs_removed": 0, "segs_recycled": 30,
				"sync_rels":        925,
				"bufs_written_pct": 1.4, "write_secs": 68.98, "sync_secs": 1.542, "total_secs": 70.548,
				"longest_secs": 0.216, "average_secs": 0.001},
		}},
		nil,
	}, {
		[]state.LogLine{{
			Content:  "checkpoints are occurring too frequently (18 seconds apart)",
			LogLevel: pganalyze_collector.LogLineInformation_LOG,
			UUID:     uuid.UUID{1},
		}, {
			Content:  "Consider increasing the configuration parameter \"max_wal_size\".",
			LogLevel: pganalyze_collector.LogLineInformation_HINT,
		}},
		[]state.LogLine{{
			Classification: pganalyze_collector.LogLineInformation_CHECKPOINT_TOO_FREQUENT,
			LogLevel:       pganalyze_collector.LogLineInformation_LOG,
			Details: map[string]interface{}{
				"elapsed_secs": 18,
			},
			UUID: uuid.UUID{1},
		}, {
			LogLevel:   pganalyze_collector.LogLineInformation_HINT,
			ParentUUID: uuid.UUID{1},
		}},
		nil,
	},
	{
		[]state.LogLine{{
			Content: "restartpoint starting: shutdown immediate",
		}},
		[]state.LogLine{{
			Classification: pganalyze_collector.LogLineInformation_RESTARTPOINT_STARTING,
			Details:        map[string]interface{}{"reason": "shutdown immediate"},
		}},
		nil,
	}, {
		[]state.LogLine{{
			Content: "restartpoint complete: wrote 693 buffers (0.1%); 0 transaction log file(s) added, 0 removed, 5 recycled; write=0.015 s, sync=0.240 s, total=0.288 s; sync files=74, longest=0.024 s, average=0.003 s; distance=81503 kB, estimate=81503 kB",
		}},
		[]state.LogLine{{
			Classification: pganalyze_collector.LogLineInformation_RESTARTPOINT_COMPLETE,
			Details: map[string]interface{}{
				"bufs_written_pct": 0.1, "write_secs": 0.015, "sync_secs": 0.240,
				"total_secs": 0.288, "longest_secs": 0.024, "average_secs": 0.003,
				"bufs_written": 693, "segs_added": 0, "segs_removed": 0, "segs_recycled": 5,
				"sync_rels": 74, "distance_kb": 81503, "estimate_kb": 81503,
			},
		}},
		nil,
	},
	{
		[]state.LogLine{{
			Content:  "recovery restart point at 4E8/9B13FBB0",
			LogLevel: pganalyze_collector.LogLineInformation_LOG,
			UUID:     uuid.UUID{1},
		}, {
			Content:  "last completed transaction was at log time 2017-05-05 20:17:06.511443+00",
			LogLevel: pganalyze_collector.LogLineInformation_DETAIL,
		}},
		[]state.LogLine{{
			Classification: pganalyze_collector.LogLineInformation_RESTARTPOINT_AT,
			LogLevel:       pganalyze_collector.LogLineInformation_LOG,
			UUID:           uuid.UUID{1},
		}, {
			LogLevel:   pganalyze_collector.LogLineInformation_DETAIL,
			ParentUUID: uuid.UUID{1},
		}},
		nil,
	},
	// WAL/Archiving
	{
		[]state.LogLine{{
			Content: "invalid record length at 4E8/9E0979A8",
		}},
		[]state.LogLine{{
			Classification: pganalyze_collector.LogLineInformation_WAL_INVALID_RECORD_LENGTH,
		}},
		nil,
	},
	{
		[]state.LogLine{{
			Content: "redo starts at 4E8/9B13FBB0",
		}, {
			Content: "redo is not required",
		}},
		[]state.LogLine{{
			Classification: pganalyze_collector.LogLineInformation_WAL_REDO,
		}, {
			Classification: pganalyze_collector.LogLineInformation_WAL_REDO,
		}},
		nil,
	},
	// Lock waits
	{
		[]state.LogLine{{
			Content:  "process 583 acquired AccessExclusiveLock on relation 185044 of database 16384 after 2175.443 ms",
			LogLevel: pganalyze_collector.LogLineInformation_LOG,
			UUID:     uuid.UUID{1},
		}, {
			Content:  "ALTER TABLE x ADD COLUMN y text;",
			LogLevel: pganalyze_collector.LogLineInformation_STATEMENT,
		}},
		[]state.LogLine{{
			Classification: pganalyze_collector.LogLineInformation_LOCK_ACQUIRED,
			LogLevel:       pganalyze_collector.LogLineInformation_LOG,
			Query:          "ALTER TABLE x ADD COLUMN y text;",
			Details: map[string]interface{}{
				"after_ms":  2175.443,
				"lock_mode": "AccessExclusiveLock",
				"lock_type": "relation",
			},
			UUID: uuid.UUID{1},
		}, {
			LogLevel:   pganalyze_collector.LogLineInformation_STATEMENT,
			ParentUUID: uuid.UUID{1},
		}},
		nil,
	},
	{
		[]state.LogLine{{
			Content:  "process 25307 acquired ExclusiveLock on tuple (106,38) of relation 16421 of database 16385 after 1129279.295 ms",
			LogLevel: pganalyze_collector.LogLineInformation_LOG,
		}},
		[]state.LogLine{{
			Classification: pganalyze_collector.LogLineInformation_LOCK_ACQUIRED,
			LogLevel:       pganalyze_collector.LogLineInformation_LOG,
			Details: map[string]interface{}{
				"after_ms":  1129279.295,
				"lock_mode": "ExclusiveLock",
				"lock_type": "tuple",
			},
		}},
		nil,
	},
	{
		[]state.LogLine{{
			Content:  "process 2078 still waiting for ShareLock on transaction 1045207414 after 1000.100 ms",
			LogLevel: pganalyze_collector.LogLineInformation_LOG,
			UUID:     uuid.UUID{1},
		}, {
			Content:  "Process holding the lock: 583. Wait queue: 2078, 456",
			LogLevel: pganalyze_collector.LogLineInformation_DETAIL,
		}, {
			Content:  "INSERT INTO x (y) VALUES (1)",
			LogLevel: pganalyze_collector.LogLineInformation_QUERY,
		}, {
			Content:  "PL/pgSQL function insert_helper(text) line 5 at EXECUTE statement",
			LogLevel: pganalyze_collector.LogLineInformation_CONTEXT,
		}, {
			Content:  "SELECT insert_helper($1)",
			LogLevel: pganalyze_collector.LogLineInformation_STATEMENT,
		}},
		[]state.LogLine{{
			Classification: pganalyze_collector.LogLineInformation_LOCK_WAITING,
			LogLevel:       pganalyze_collector.LogLineInformation_LOG,
			Query:          "SELECT insert_helper($1)",
			Details: map[string]interface{}{
				"lock_holders": []int64{583},
				"lock_waiters": []int64{2078, 456},
				"after_ms":     1000.1,
				"lock_mode":    "ShareLock",
				"lock_type":    "transactionid",
			},
			UUID: uuid.UUID{1},
		}, {
			LogLevel:   pganalyze_collector.LogLineInformation_DETAIL,
			ParentUUID: uuid.UUID{1},
		}, {
			LogLevel:   pganalyze_collector.LogLineInformation_QUERY,
			ParentUUID: uuid.UUID{1},
		}, {
			LogLevel:   pganalyze_collector.LogLineInformation_CONTEXT,
			ParentUUID: uuid.UUID{1},
		}, {
			LogLevel:   pganalyze_collector.LogLineInformation_STATEMENT,
			ParentUUID: uuid.UUID{1},
		}},
		nil,
	},
	{
		[]state.LogLine{{
			Content:  "canceling statement due to lock timeout",
			LogLevel: pganalyze_collector.LogLineInformation_ERROR,
			UUID:     uuid.UUID{1},
		}, {
			Content:  "while updating tuple (24,41) in relation \"mytable\"",
			LogLevel: pganalyze_collector.LogLineInformation_CONTEXT,
		}, {
			Content:  "UPDATE mytable SET y = 2 WHERE x = 1",
			LogLevel: pganalyze_collector.LogLineInformation_STATEMENT,
		}},
		[]state.LogLine{{
			Classification: pganalyze_collector.LogLineInformation_LOCK_TIMEOUT,
			LogLevel:       pganalyze_collector.LogLineInformation_ERROR,
			Query:          "UPDATE mytable SET y = 2 WHERE x = 1",
			UUID:           uuid.UUID{1},
		}, {
			LogLevel:   pganalyze_collector.LogLineInformation_CONTEXT,
			ParentUUID: uuid.UUID{1},
		}, {
			LogLevel:   pganalyze_collector.LogLineInformation_STATEMENT,
			ParentUUID: uuid.UUID{1},
		}},
		nil,
	},
	{
		[]state.LogLine{{
			Content:  "process 2078 avoided deadlock for AccessExclusiveLock on relation 999 by rearranging queue order after 123.456 ms",
			LogLevel: pganalyze_collector.LogLineInformation_LOG,
			UUID:     uuid.UUID{1},
		}, {
			Content:  "Processes holding the lock: 583, 123. Wait queue: 2078",
			LogLevel: pganalyze_collector.LogLineInformation_DETAIL,
		}},
		[]state.LogLine{{
			LogLevel:       pganalyze_collector.LogLineInformation_LOG,
			Classification: pganalyze_collector.LogLineInformation_LOCK_DEADLOCK_AVOIDED,
			Details: map[string]interface{}{
				"lock_holders": []int64{583, 123},
				"lock_waiters": []int64{2078},
				"after_ms":     123.456,
				"lock_mode":    "AccessExclusiveLock",
				"lock_type":    "relation",
			},
			UUID: uuid.UUID{1},
		}, {
			LogLevel:   pganalyze_collector.LogLineInformation_DETAIL,
			ParentUUID: uuid.UUID{1},
		}},
		nil,
	}, {
		[]state.LogLine{{
			Content:  "process 123 detected deadlock while waiting for AccessExclusiveLock on extension of relation 666 of database 123 after 456.000 ms",
			LogLevel: pganalyze_collector.LogLineInformation_LOG,
		}},
		[]state.LogLine{{
			LogLevel:       pganalyze_collector.LogLineInformation_LOG,
			Classification: pganalyze_collector.LogLineInformation_LOCK_DEADLOCK_DETECTED,
			Details: map[string]interface{}{
				"lock_mode": "AccessExclusiveLock",
				"lock_type": "extend",
				"after_ms":  456.0,
			},
		}},
		nil,
	}, {
		[]state.LogLine{{
			Content:  "deadlock detected",
			LogLevel: pganalyze_collector.LogLineInformation_ERROR,
			UUID:     uuid.UUID{1},
		}, {
			Content: "Process 9788 waits for ShareLock on transaction 1035; blocked by process 91." +
				"\nProcess 91 waits for ShareLock on transaction 1045; blocked by process 98.\n" +
				"\nProcess 98: INSERT INTO x (id, name, email) VALUES (1, 'ABC', 'abc@example.com') ON CONFLICT(email) DO UPDATE SET name = excluded.name, /* truncated */" +
				"\nProcess 91: INSERT INTO x (id, name, email) VALUES (1, 'ABC', 'abc@example.com') ON CONFLICT(email) DO UPDATE SET name = excluded.name, /* truncated */",
			LogLevel: pganalyze_collector.LogLineInformation_DETAIL,
		}, {
			Content:  "See server log for query details.",
			LogLevel: pganalyze_collector.LogLineInformation_HINT,
		}, {
			Content:  "while inserting index tuple (1,42) in relation \"x\"",
			LogLevel: pganalyze_collector.LogLineInformation_CONTEXT,
		}, {
			Content:  "INSERT INTO x (id, name, email) VALUES (1, 'ABC', 'abc@example.com') ON CONFLICT(email) DO UPDATE SET name = excluded.name RETURNING id",
			LogLevel: pganalyze_collector.LogLineInformation_STATEMENT,
		}},
		[]state.LogLine{{
			LogLevel:       pganalyze_collector.LogLineInformation_ERROR,
			Classification: pganalyze_collector.LogLineInformation_LOCK_DEADLOCK_DETECTED,
			Query:          "INSERT INTO x (id, name, email) VALUES (1, 'ABC', 'abc@example.com') ON CONFLICT(email) DO UPDATE SET name = excluded.name RETURNING id",
			UUID:           uuid.UUID{1},
		}, {
			LogLevel:   pganalyze_collector.LogLineInformation_DETAIL,
			ParentUUID: uuid.UUID{1},
		}, {
			LogLevel:   pganalyze_collector.LogLineInformation_HINT,
			ParentUUID: uuid.UUID{1},
		}, {
			LogLevel:   pganalyze_collector.LogLineInformation_CONTEXT,
			ParentUUID: uuid.UUID{1},
		}, {
			LogLevel:   pganalyze_collector.LogLineInformation_STATEMENT,
			ParentUUID: uuid.UUID{1},
		}},
		nil,
	}, {
		[]state.LogLine{{
			Content:  "process 663 still waiting for ShareLock on virtual transaction 2/7 after 1000.123 ms",
			LogLevel: pganalyze_collector.LogLineInformation_LOG,
		}},
		[]state.LogLine{{
			LogLevel:       pganalyze_collector.LogLineInformation_LOG,
			Classification: pganalyze_collector.LogLineInformation_LOCK_WAITING,
			Details: map[string]interface{}{
				"lock_mode": "ShareLock",
				"lock_type": "virtualxid",
				"after_ms":  1000.123,
			},
		}},
		nil,
	}, {
		[]state.LogLine{{
			Content:  "process 663 still waiting for ExclusiveLock on advisory lock [233136,1,2,2] after 1000.365 ms",
			LogLevel: pganalyze_collector.LogLineInformation_LOG,
			UUID:     uuid.UUID{1},
		}, {
			Content:  "Process holding the lock: 660. Wait queue: 663.",
			LogLevel: pganalyze_collector.LogLineInformation_DETAIL,
		}, {
			Content:  "SELECT pg_advisory_lock(1, 2);",
			LogLevel: pganalyze_collector.LogLineInformation_STATEMENT,
		}},
		[]state.LogLine{{
			LogLevel:       pganalyze_collector.LogLineInformation_LOG,
			Classification: pganalyze_collector.LogLineInformation_LOCK_WAITING,
			Details: map[string]interface{}{
				"lock_mode":    "ExclusiveLock",
				"lock_type":    "advisory",
				"lock_holders": []int64{660},
				"lock_waiters": []int64{663},
				"after_ms":     1000.365,
			},
			Query: "SELECT pg_advisory_lock(1, 2);",
			UUID:  uuid.UUID{1},
		}, {
			LogLevel:   pganalyze_collector.LogLineInformation_DETAIL,
			ParentUUID: uuid.UUID{1},
		}, {
			LogLevel:   pganalyze_collector.LogLineInformation_STATEMENT,
			ParentUUID: uuid.UUID{1},
		}},
		nil,
	},
	// Autovacuum
	{
		[]state.LogLine{{
			Content:  "canceling autovacuum task",
			LogLevel: pganalyze_collector.LogLineInformation_ERROR,
			UUID:     uuid.UUID{1},
		}, {
			Content:  "automatic analyze of table \"dbname.schemaname.tablename\"",
			LogLevel: pganalyze_collector.LogLineInformation_CONTEXT,
		}},
		[]state.LogLine{{
			LogLevel:       pganalyze_collector.LogLineInformation_ERROR,
			Classification: pganalyze_collector.LogLineInformation_AUTOVACUUM_CANCEL,
			UUID:           uuid.UUID{1},
		}, {
			LogLevel:   pganalyze_collector.LogLineInformation_CONTEXT,
			ParentUUID: uuid.UUID{1},
		}},
		nil,
	},
	{
		[]state.LogLine{{
			Content:  "database \"template1\" must be vacuumed within 938860 transactions",
			LogLevel: pganalyze_collector.LogLineInformation_WARNING,
			UUID:     uuid.UUID{1},
		}, {
			Content:  "To avoid a database shutdown, execute a full-database VACUUM in \"template1\".",
			LogLevel: pganalyze_collector.LogLineInformation_HINT,
		}},
		[]state.LogLine{{
			LogLevel:       pganalyze_collector.LogLineInformation_WARNING,
			Classification: pganalyze_collector.LogLineInformation_TXID_WRAPAROUND_WARNING,
			Details: map[string]interface{}{
				"database_name":  "template1",
				"remaining_xids": 938860,
			},
			UUID: uuid.UUID{1},
		}, {
			LogLevel:   pganalyze_collector.LogLineInformation_HINT,
			ParentUUID: uuid.UUID{1},
		}},
		nil,
	},
	{
		[]state.LogLine{{
			Content: "database with OID 10 must be vacuumed within 100 transactions",
		}},
		[]state.LogLine{{
			Classification: pganalyze_collector.LogLineInformation_TXID_WRAPAROUND_WARNING,
			Details: map[string]interface{}{
				"database_oid":   10,
				"remaining_xids": 100,
			},
		}},
		nil,
	},
	{
		[]state.LogLine{{
			Content:  "database is not accepting commands to avoid wraparound data loss in database \"mydb\"",
			LogLevel: pganalyze_collector.LogLineInformation_ERROR,
			UUID:     uuid.UUID{1},
		}, {
			Content:  "Stop the postmaster and use a standalone backend to vacuum that database. You might also need to commit or roll back old prepared transactions.",
			LogLevel: pganalyze_collector.LogLineInformation_HINT,
		}},
		[]state.LogLine{{
			Classification: pganalyze_collector.LogLineInformation_TXID_WRAPAROUND_ERROR,
			LogLevel:       pganalyze_collector.LogLineInformation_ERROR,
			Details: map[string]interface{}{
				"database_name": "mydb",
			},
			UUID: uuid.UUID{1},
		}, {
			LogLevel:   pganalyze_collector.LogLineInformation_HINT,
			ParentUUID: uuid.UUID{1},
		}},
		nil,
	},
	{
		[]state.LogLine{{
			Content: "database is not accepting commands to avoid wraparound data loss in database with OID 16384",
		}},
		[]state.LogLine{{
			Classification: pganalyze_collector.LogLineInformation_TXID_WRAPAROUND_ERROR,
			Details: map[string]interface{}{
				"database_oid": 16384,
			},
		}},
		nil,
	},
	{
		[]state.LogLine{{
			Content: "autovacuum launcher started",
		}},
		[]state.LogLine{{
			Classification: pganalyze_collector.LogLineInformation_AUTOVACUUM_LAUNCHER_STARTED,
		}},
		nil,
	},
	{
		[]state.LogLine{{
			Content: "autovacuum launcher shutting down",
		}},
		[]state.LogLine{{
			Classification: pganalyze_collector.LogLineInformation_AUTOVACUUM_LAUNCHER_SHUTTING_DOWN,
		}},
		nil,
	},
	{
		[]state.LogLine{{
			Content: "automatic vacuum of table \"mydb.public.vac_test\": index scans: 1" +
				"\n pages: 0 removed, 1 remain, 0 skipped due to pins, 0 skipped frozen" +
				"\n tuples: 3 removed, 6 remain, 0 are dead but not yet removable" +
				"\n buffer usage: 70 hits, 4 misses, 4 dirtied" +
				"\n avg read rate: 62.877 MB/s, avg write rate: 62.877 MB/s" +
				"\n system usage: CPU 0.00s/0.00u sec elapsed 0.00 sec",
			LogLevel: pganalyze_collector.LogLineInformation_LOG,
		}},
		[]state.LogLine{{
			Classification: pganalyze_collector.LogLineInformation_AUTOVACUUM_COMPLETED,
			LogLevel:       pganalyze_collector.LogLineInformation_LOG,
			Details: map[string]interface{}{
				"num_index_scans":     1,
				"pages_removed":       0,
				"rel_pages":           1,
				"pinskipped_pages":    0,
				"frozenskipped_pages": 0,
				"tuples_deleted":      3,
				"new_rel_tuples":      6,
				"new_dead_tuples":     0,
				"vacuum_page_hit":     70,
				"vacuum_page_miss":    4,
				"vacuum_page_dirty":   4,
				"read_rate_mb":        62.877,
				"write_rate_mb":       62.877,
				"rusage_kernel":       0,
				"rusage_user":         0,
				"elapsed_secs":        0,
			},
		}},
		nil,
	},
	{
		[]state.LogLine{{
			Content: "automatic vacuum of table \"postgres.public.pgbench_branches\": index scans: 1" +
				"\npages: 0 removed, 12 remain" +
				"\ntuples: 423 removed, 107 remain, 3 are dead but not yet removable" +
				"\nbuffer usage: 52 hits, 1 misses, 1 dirtied" +
				"\navg read rate: 7.455 MB/s, avg write rate: 7.455 MB/s" +
				"\nsystem usage: CPU 0.00s/0.00u sec elapsed 0.00 sec",
			LogLevel: pganalyze_collector.LogLineInformation_LOG,
		}},
		[]state.LogLine{{
			Classification: pganalyze_collector.LogLineInformation_AUTOVACUUM_COMPLETED,
			LogLevel:       pganalyze_collector.LogLineInformation_LOG,
			Details: map[string]interface{}{
				"num_index_scans":   1,
				"pages_removed":     0,
				"rel_pages":         12,
				"tuples_deleted":    423,
				"new_rel_tuples":    107,
				"new_dead_tuples":   3,
				"vacuum_page_hit":   52,
				"vacuum_page_miss":  1,
				"vacuum_page_dirty": 1,
				"read_rate_mb":      7.455,
				"write_rate_mb":     7.455,
				"rusage_kernel":     0,
				"rusage_user":       0,
				"elapsed_secs":      0,
			},
		}},
		nil,
	},
	{
		[]state.LogLine{{
			Content: "automatic vacuum of table \"my_db.public.my_dimension\": index scans: 1" +
				"\n  pages: 0 removed, 29457 remain" +
				"\n  tuples: 3454 removed, 429481 remain, 0 are dead but not yet removable" +
				"\n  buffer usage: 64215 hits, 8056 misses, 22588 dirtied" +
				"\n  avg read rate: 1.018 MB/s, avg write rate: 2.855 MB/s" +
				"\n  system usage: CPU 0.10s/0.88u sec elapsed 61.80 seconds",
			LogLevel: pganalyze_collector.LogLineInformation_LOG,
		}},
		[]state.LogLine{{
			Classification: pganalyze_collector.LogLineInformation_AUTOVACUUM_COMPLETED,
			LogLevel:       pganalyze_collector.LogLineInformation_LOG,
			Details: map[string]interface{}{
				"num_index_scans":   1,
				"pages_removed":     0,
				"rel_pages":         29457,
				"tuples_deleted":    3454,
				"new_rel_tuples":    429481,
				"new_dead_tuples":   0,
				"vacuum_page_hit":   64215,
				"vacuum_page_miss":  8056,
				"vacuum_page_dirty": 22588,
				"read_rate_mb":      1.018,
				"write_rate_mb":     2.855,
				"rusage_kernel":     0.10,
				"rusage_user":       0.88,
				"elapsed_secs":      61.80,
			},
		}},
		nil,
	},
	{
		[]state.LogLine{{
			Content:  "automatic analyze of table \"postgres.public.pgbench_branches\" system usage: CPU 1.02s/2.08u sec elapsed 108.25 sec",
			LogLevel: pganalyze_collector.LogLineInformation_LOG,
		}},
		[]state.LogLine{{
			Classification: pganalyze_collector.LogLineInformation_AUTOANALYZE_COMPLETED,
			LogLevel:       pganalyze_collector.LogLineInformation_LOG,
			Details: map[string]interface{}{
				"rusage_kernel": 1.02,
				"rusage_user":   2.08,
				"elapsed_secs":  108.25,
			},
		}},
		nil,
	},
	// Statement cancellation (other than lock timeout)
	{
		[]state.LogLine{{
			Content:  "canceling statement due to user request",
			LogLevel: pganalyze_collector.LogLineInformation_ERROR,
			UUID:     uuid.UUID{1},
		}, {
			Content:  "SELECT 1",
			LogLevel: pganalyze_collector.LogLineInformation_STATEMENT,
		}},
		[]state.LogLine{{
			LogLevel:       pganalyze_collector.LogLineInformation_ERROR,
			Classification: pganalyze_collector.LogLineInformation_STATEMENT_CANCELED_USER,
			Query:          "SELECT 1",
			UUID:           uuid.UUID{1},
		}, {
			LogLevel:   pganalyze_collector.LogLineInformation_STATEMENT,
			ParentUUID: uuid.UUID{1},
		}},
		nil,
	}, {
		[]state.LogLine{{
			Content:  "canceling statement due to statement timeout",
			LogLevel: pganalyze_collector.LogLineInformation_ERROR,
			UUID:     uuid.UUID{1},
		}, {
			Content:  "SELECT 1",
			LogLevel: pganalyze_collector.LogLineInformation_STATEMENT,
		}},
		[]state.LogLine{{
			LogLevel:       pganalyze_collector.LogLineInformation_ERROR,
			Classification: pganalyze_collector.LogLineInformation_STATEMENT_CANCELED_TIMEOUT,
			Query:          "SELECT 1",
			UUID:           uuid.UUID{1},
		}, {
			LogLevel:   pganalyze_collector.LogLineInformation_STATEMENT,
			ParentUUID: uuid.UUID{1},
		}},
		nil,
	},
	// Server events
	{
		[]state.LogLine{{
			Content:  "server process (PID 660) was terminated by signal 6: Aborted",
			LogLevel: pganalyze_collector.LogLineInformation_LOG,
			UUID:     uuid.UUID{1},
		}, {
			Content:  "Failed process was running: SELECT pg_advisory_lock(1, 2);",
			LogLevel: pganalyze_collector.LogLineInformation_DETAIL,
		}, {
			Content:  "terminating any other active server processes",
			LogLevel: pganalyze_collector.LogLineInformation_LOG,
		}, {
			Content:  "terminating connection because of crash of another server process",
			LogLevel: pganalyze_collector.LogLineInformation_WARNING,
			UUID:     uuid.UUID{2},
		}, {
			Content:  "The postmaster has commanded this server process to roll back the current transaction and exit, because another server process exited abnormally and possibly corrupted shared memory.",
			LogLevel: pganalyze_collector.LogLineInformation_DETAIL,
		}, {
			Content:  "In a moment you should be able to reconnect to the database and repeat your command.",
			LogLevel: pganalyze_collector.LogLineInformation_HINT,
		}, {
			Content:  "all server processes terminated; reinitializing",
			LogLevel: pganalyze_collector.LogLineInformation_LOG,
		}},
		[]state.LogLine{{
			Classification: pganalyze_collector.LogLineInformation_SERVER_CRASHED,
			LogLevel:       pganalyze_collector.LogLineInformation_LOG,
			UUID:           uuid.UUID{1},
		}, {
			LogLevel:   pganalyze_collector.LogLineInformation_DETAIL,
			ParentUUID: uuid.UUID{1},
		}, {
			Classification: pganalyze_collector.LogLineInformation_SERVER_CRASHED,
			LogLevel:       pganalyze_collector.LogLineInformation_LOG,
		}, {
			Classification: pganalyze_collector.LogLineInformation_SERVER_CRASHED,
			LogLevel:       pganalyze_collector.LogLineInformation_WARNING,
			UUID:           uuid.UUID{2},
		}, {
			LogLevel:   pganalyze_collector.LogLineInformation_DETAIL,
			ParentUUID: uuid.UUID{2},
		}, {
			LogLevel:   pganalyze_collector.LogLineInformation_HINT,
			ParentUUID: uuid.UUID{2},
		}, {
			Classification: pganalyze_collector.LogLineInformation_SERVER_CRASHED,
			LogLevel:       pganalyze_collector.LogLineInformation_LOG,
		}},
		nil,
	},
	{
		[]state.LogLine{{
			Content: "database system was shut down in recovery at 2017-05-05 20:17:07 UTC",
		}, {
			Content: "entering standby mode",
		}, {
			Content: "database system is ready to accept read only connections",
		}, {
			Content: "database system was shut down at 2017-05-03 23:23:37 UTC",
		}, {
			Content: "MultiXact member wraparound protections are now enabled",
		}, {
			Content: "database system is ready to accept connections",
		}},
		[]state.LogLine{{
			Classification: pganalyze_collector.LogLineInformation_SERVER_START,
		}, {
			Classification: pganalyze_collector.LogLineInformation_SERVER_START,
		}, {
			Classification: pganalyze_collector.LogLineInformation_SERVER_START,
		}, {
			Classification: pganalyze_collector.LogLineInformation_SERVER_START,
		}, {
			Classification: pganalyze_collector.LogLineInformation_SERVER_START,
		}, {
			Classification: pganalyze_collector.LogLineInformation_SERVER_START,
		}},
		nil,
	},
	{
		[]state.LogLine{{
			Content: "database system was interrupted; last known up at 2017-05-07 22:33:02 UTC",
		}, {
			Content: "database system was not properly shut down; automatic recovery in progress",
		}},
		[]state.LogLine{{
			Classification: pganalyze_collector.LogLineInformation_SERVER_START_RECOVERING,
		}, {
			Classification: pganalyze_collector.LogLineInformation_SERVER_START_RECOVERING,
		}},
		nil,
	}, {
		[]state.LogLine{{
			Content: "received smart shutdown request",
		}, {
			Content: "received fast shutdown request",
		}, {
			Content: "aborting any active transactions",
		}, {
			Content: "shutting down",
		}, {
			Content: "database system is shut down",
		}},
		[]state.LogLine{{
			Classification: pganalyze_collector.LogLineInformation_SERVER_SHUTDOWN,
		}, {
			Classification: pganalyze_collector.LogLineInformation_SERVER_SHUTDOWN,
		}, {
			Classification: pganalyze_collector.LogLineInformation_SERVER_SHUTDOWN,
		}, {
			Classification: pganalyze_collector.LogLineInformation_SERVER_SHUTDOWN,
		}, {
			Classification: pganalyze_collector.LogLineInformation_SERVER_SHUTDOWN,
		}},
		nil,
	}, {
		[]state.LogLine{{
			Content:  "temporary file: path \"base/pgsql_tmp/pgsql_tmp15967.0\", size 200204288",
			LogLevel: pganalyze_collector.LogLineInformation_LOG,
			UUID:     uuid.UUID{1},
		}, {
			Content:  "alter table pgbench_accounts add primary key (aid)",
			LogLevel: pganalyze_collector.LogLineInformation_STATEMENT,
		}},
		[]state.LogLine{{
			LogLevel:       pganalyze_collector.LogLineInformation_LOG,
			Classification: pganalyze_collector.LogLineInformation_SERVER_TEMP_FILE_CREATED,
			UUID:           uuid.UUID{1},
			Query:          "alter table pgbench_accounts add primary key (aid)",
			Details: map[string]interface{}{
				"file": "base/pgsql_tmp/pgsql_tmp15967.0",
				"size": 200204288,
			},
		}, {
			LogLevel:   pganalyze_collector.LogLineInformation_STATEMENT,
			ParentUUID: uuid.UUID{1},
		}},
		nil,
	}, {
		[]state.LogLine{{
			Content: "could not open usermap file \"/var/lib/pgsql/9.5/data/pg_ident.conf\": No such file or directory",
		}, {
			Content: "invalid byte sequence for encoding \"UTF8\": 0xd0 0x2e",
		}, {
			Content: "could not link file \"pg_xlog/xlogtemp.26115\" to \"pg_xlog/000000010000021B000000C5\": File exists",
		}, {
			Content: "unexpected pageaddr 2D5/12000000 in log segment 00000001000002D500000022, offset 0",
		}},
		[]state.LogLine{{
			Classification: pganalyze_collector.LogLineInformation_SERVER_MISC,
		}, {
			Classification: pganalyze_collector.LogLineInformation_SERVER_MISC,
		}, {
			Classification: pganalyze_collector.LogLineInformation_SERVER_MISC,
		}, {
			Classification: pganalyze_collector.LogLineInformation_SERVER_MISC,
		}},
		nil,
	}, {
		[]state.LogLine{{
			Content:  "out of memory",
			LogLevel: pganalyze_collector.LogLineInformation_ERROR,
			UUID:     uuid.UUID{1},
		}, {
			Content:  "Failed on request of size 324589128.",
			LogLevel: pganalyze_collector.LogLineInformation_DETAIL,
		}, {
			Content:  "SELECT 123",
			LogLevel: pganalyze_collector.LogLineInformation_STATEMENT,
		}},
		[]state.LogLine{{
			LogLevel:       pganalyze_collector.LogLineInformation_ERROR,
			Classification: pganalyze_collector.LogLineInformation_SERVER_OUT_OF_MEMORY,
			UUID:           uuid.UUID{1},
			Query:          "SELECT 123",
		}, {
			LogLevel:   pganalyze_collector.LogLineInformation_DETAIL,
			ParentUUID: uuid.UUID{1},
		}, {
			LogLevel:   pganalyze_collector.LogLineInformation_STATEMENT,
			ParentUUID: uuid.UUID{1},
		}},
		nil,
	}, {
		[]state.LogLine{{
			Content:  "server process (PID 123) was terminated by signal 9: Killed",
			LogLevel: pganalyze_collector.LogLineInformation_LOG,
		}},
		[]state.LogLine{{
			LogLevel:       pganalyze_collector.LogLineInformation_LOG,
			Classification: pganalyze_collector.LogLineInformation_SERVER_OUT_OF_MEMORY,
		}},
		nil,
	}, {
		[]state.LogLine{{
			Content:  "page verification failed, calculated checksum 20919 but expected 15254",
			LogLevel: pganalyze_collector.LogLineInformation_WARNING,
		}, {
			Content:  "invalid page in block 335458 of relation base/16385/99454",
			LogLevel: pganalyze_collector.LogLineInformation_ERROR,
			UUID:     uuid.UUID{1},
		}, {
			Content:  "SELECT 1",
			LogLevel: pganalyze_collector.LogLineInformation_STATEMENT,
		}},
		[]state.LogLine{{
			LogLevel:       pganalyze_collector.LogLineInformation_WARNING,
			Classification: pganalyze_collector.LogLineInformation_SERVER_INVALID_CHECKSUM,
		}, {
			LogLevel:       pganalyze_collector.LogLineInformation_ERROR,
			Classification: pganalyze_collector.LogLineInformation_SERVER_INVALID_CHECKSUM,
			UUID:           uuid.UUID{1},
			Query:          "SELECT 1",
			Details: map[string]interface{}{
				"block": 335458,
				"file":  "base/16385/99454",
			},
		}, {
			LogLevel:   pganalyze_collector.LogLineInformation_STATEMENT,
			ParentUUID: uuid.UUID{1},
		}},
		nil,
	},
	{
		[]state.LogLine{{
			Content: "received SIGHUP, reloading configuration files",
		}, {
			Content: "parameter \"log_autovacuum_min_duration\" changed to \"0\"",
		}, {
			Content: "parameter \"shared_preload_libraries\" cannot be changed without restarting the server",
		}, {
			Content: "configuration file \"/var/lib/postgresql/data/postgresql.auto.conf\" contains errors; unaffected changes were applied",
		}},
		[]state.LogLine{{
			Classification: pganalyze_collector.LogLineInformation_SERVER_RELOAD,
		}, {
			Classification: pganalyze_collector.LogLineInformation_SERVER_RELOAD,
		}, {
			Classification: pganalyze_collector.LogLineInformation_SERVER_RELOAD,
		}, {
			Classification: pganalyze_collector.LogLineInformation_SERVER_RELOAD,
		}},
		nil,
	},
	// Standby
	{
		[]state.LogLine{{
			Content:  "restored log file \"00000006000004E80000009C\" from archive",
			LogLevel: pganalyze_collector.LogLineInformation_LOG,
		}},
		[]state.LogLine{{
			Classification: pganalyze_collector.LogLineInformation_STANDBY_RESTORED_WAL_FROM_ARCHIVE,
			LogLevel:       pganalyze_collector.LogLineInformation_LOG,
		}},
		nil,
	},
	{
		[]state.LogLine{{
			Content:  "started streaming WAL from primary at 4E8/9E000000 on timeline 6",
			LogLevel: pganalyze_collector.LogLineInformation_LOG,
		}},
		[]state.LogLine{{
			Classification: pganalyze_collector.LogLineInformation_STANDBY_STARTED_STREAMING,
			LogLevel:       pganalyze_collector.LogLineInformation_LOG,
		}},
		nil,
	},
	{
		[]state.LogLine{{
			Content:  "could not receive data from WAL stream: SSL error: sslv3 alert unexpected message",
			LogLevel: pganalyze_collector.LogLineInformation_FATAL,
		}},
		[]state.LogLine{{
			Classification: pganalyze_collector.LogLineInformation_STANDBY_STREAMING_INTERRUPTED,
			LogLevel:       pganalyze_collector.LogLineInformation_FATAL,
		}},
		nil,
	}, {
		[]state.LogLine{{
			Content:  "terminating walreceiver process due to administrator command",
			LogLevel: pganalyze_collector.LogLineInformation_FATAL,
		}},
		[]state.LogLine{{
			Classification: pganalyze_collector.LogLineInformation_STANDBY_STOPPED_STREAMING,
			LogLevel:       pganalyze_collector.LogLineInformation_FATAL,
		}},
		nil,
	}, {
		[]state.LogLine{{
			Content:  "consistent recovery state reached at 4E8/9E0979A8",
			LogLevel: pganalyze_collector.LogLineInformation_LOG,
		}},
		[]state.LogLine{{
			Classification: pganalyze_collector.LogLineInformation_STANDBY_CONSISTENT_RECOVERY_STATE,
			LogLevel:       pganalyze_collector.LogLineInformation_LOG,
		}},
		nil,
	}, {
		[]state.LogLine{{
			Content:  "canceling statement due to conflict with recovery",
			LogLevel: pganalyze_collector.LogLineInformation_ERROR,
			UUID:     uuid.UUID{1},
		}, {
			Content:  "User query might have needed to see row versions that must be removed.",
			LogLevel: pganalyze_collector.LogLineInformation_DETAIL,
		}, {
			Content:  "SELECT 1",
			LogLevel: pganalyze_collector.LogLineInformation_STATEMENT,
		}},
		[]state.LogLine{{
			LogLevel:       pganalyze_collector.LogLineInformation_ERROR,
			Classification: pganalyze_collector.LogLineInformation_STANDBY_STATEMENT_CANCELED,
			Query:          "SELECT 1",
			UUID:           uuid.UUID{1},
		}, {
			LogLevel:   pganalyze_collector.LogLineInformation_DETAIL,
			ParentUUID: uuid.UUID{1},
		}, {
			LogLevel:   pganalyze_collector.LogLineInformation_STATEMENT,
			ParentUUID: uuid.UUID{1},
		}},
		nil,
	},
	{
		[]state.LogLine{{
			Content:  "according to history file, WAL location 2D5/22000000 belongs to timeline 3, but previous recovered WAL file came from timeline 4",
			LogLevel: pganalyze_collector.LogLineInformation_FATAL,
		}},
		[]state.LogLine{{
			LogLevel:       pganalyze_collector.LogLineInformation_FATAL,
			Classification: pganalyze_collector.LogLineInformation_STANDBY_INVALID_TIMELINE,
		}},
		nil,
	},
	// Constraint violations
	{
		[]state.LogLine{{
			Content:  "duplicate key value violates unique constraint \"test_constraint\"",
			LogLevel: pganalyze_collector.LogLineInformation_ERROR,
			UUID:     uuid.UUID{1},
		}, {
			Content:  "Key (b, c)=(12345, mysecretdata) already exists.",
			LogLevel: pganalyze_collector.LogLineInformation_DETAIL,
		}, {
			Content:  "INSERT INTO a (b, c) VALUES ($1,$2) RETURNING id",
			LogLevel: pganalyze_collector.LogLineInformation_STATEMENT,
		}},
		[]state.LogLine{{
			LogLevel:       pganalyze_collector.LogLineInformation_ERROR,
			Classification: pganalyze_collector.LogLineInformation_UNIQUE_CONSTRAINT_VIOLATION,
			UUID:           uuid.UUID{1},
			Query:          "INSERT INTO a (b, c) VALUES ($1,$2) RETURNING id",
		}, {
			LogLevel:   pganalyze_collector.LogLineInformation_DETAIL,
			ParentUUID: uuid.UUID{1},
		}, {
			LogLevel:   pganalyze_collector.LogLineInformation_STATEMENT,
			ParentUUID: uuid.UUID{1},
		}},
		nil,
	}, {
		[]state.LogLine{{
			Content:  "insert or update on table \"weather\" violates foreign key constraint \"weather_city_fkey\"",
			LogLevel: pganalyze_collector.LogLineInformation_ERROR,
			UUID:     uuid.UUID{1},
		}, {
			Content:  "Key (city)=(Berkeley) is not present in table \"cities\".",
			LogLevel: pganalyze_collector.LogLineInformation_DETAIL,
		}, {
			Content:  "INSERT INTO weather VALUES ('Berkeley', 45, 53, 0.0, '1994-11-28');",
			LogLevel: pganalyze_collector.LogLineInformation_STATEMENT,
		}},
		[]state.LogLine{{
			LogLevel:       pganalyze_collector.LogLineInformation_ERROR,
			Classification: pganalyze_collector.LogLineInformation_FOREIGN_KEY_CONSTRAINT_VIOLATION,
			UUID:           uuid.UUID{1},
			Query:          "INSERT INTO weather VALUES ('Berkeley', 45, 53, 0.0, '1994-11-28');",
		}, {
			LogLevel:   pganalyze_collector.LogLineInformation_DETAIL,
			ParentUUID: uuid.UUID{1},
		}, {
			LogLevel:   pganalyze_collector.LogLineInformation_STATEMENT,
			ParentUUID: uuid.UUID{1},
		}},
		nil,
	},
	{
		[]state.LogLine{{
			Content:  "update or delete on table \"test\" violates foreign key constraint \"test_fkey\" on table \"othertest\"",
			LogLevel: pganalyze_collector.LogLineInformation_ERROR,
			UUID:     uuid.UUID{1},
		}, {
			Content:  "Key (id)=(123) is still referenced from table \"othertest\".",
			LogLevel: pganalyze_collector.LogLineInformation_DETAIL,
		}, {
			Content:  "DELETE FROM test WHERE id = 123",
			LogLevel: pganalyze_collector.LogLineInformation_STATEMENT,
		}},
		[]state.LogLine{{
			LogLevel:       pganalyze_collector.LogLineInformation_ERROR,
			Classification: pganalyze_collector.LogLineInformation_FOREIGN_KEY_CONSTRAINT_VIOLATION,
			UUID:           uuid.UUID{1},
			Query:          "DELETE FROM test WHERE id = 123",
		}, {
			LogLevel:   pganalyze_collector.LogLineInformation_DETAIL,
			ParentUUID: uuid.UUID{1},
		}, {
			LogLevel:   pganalyze_collector.LogLineInformation_STATEMENT,
			ParentUUID: uuid.UUID{1},
		}},
		nil,
	},
	{
		[]state.LogLine{{
			Content:  "null value in column \"mycolumn\" violates not-null constraint",
			LogLevel: pganalyze_collector.LogLineInformation_ERROR,
			UUID:     uuid.UUID{1},
		}, {
			Content:  "Failing row contains (null).",
			LogLevel: pganalyze_collector.LogLineInformation_DETAIL,
		}, {
			Content:  "INSERT INTO \"test\" (\"mycolumn\") VALUES ($1) RETURNING \"id\"",
			LogLevel: pganalyze_collector.LogLineInformation_STATEMENT,
		}},
		[]state.LogLine{{
			LogLevel:       pganalyze_collector.LogLineInformation_ERROR,
			Classification: pganalyze_collector.LogLineInformation_NOT_NULL_CONSTRAINT_VIOLATION,
			UUID:           uuid.UUID{1},
			Query:          "INSERT INTO \"test\" (\"mycolumn\") VALUES ($1) RETURNING \"id\"",
		}, {
			LogLevel:   pganalyze_collector.LogLineInformation_DETAIL,
			ParentUUID: uuid.UUID{1},
		}, {
			LogLevel:   pganalyze_collector.LogLineInformation_STATEMENT,
			ParentUUID: uuid.UUID{1},
		}},
		nil,
	},
	{
		[]state.LogLine{{
			Content:  "new row for relation \"test\" violates check constraint \"positive_value_check\"",
			LogLevel: pganalyze_collector.LogLineInformation_ERROR,
			UUID:     uuid.UUID{1},
		}, {
			Content:  "Failing row contains (-123).",
			LogLevel: pganalyze_collector.LogLineInformation_DETAIL,
		}, {
			Content:  "check constraint \"valid_tag\" is violated by some row",
			LogLevel: pganalyze_collector.LogLineInformation_ERROR,
		}, {
			Content:  "column \"mycolumn\" of table \"test\" contains values that violate the new constraint",
			LogLevel: pganalyze_collector.LogLineInformation_ERROR,
		}, {
			Content:  "value for domain mydomain violates check constraint \"mydomain_check\"",
			LogLevel: pganalyze_collector.LogLineInformation_ERROR,
		}},
		[]state.LogLine{{
			LogLevel:       pganalyze_collector.LogLineInformation_ERROR,
			Classification: pganalyze_collector.LogLineInformation_CHECK_CONSTRAINT_VIOLATION,
			UUID:           uuid.UUID{1},
		}, {
			LogLevel:   pganalyze_collector.LogLineInformation_DETAIL,
			ParentUUID: uuid.UUID{1},
		}, {
			LogLevel:       pganalyze_collector.LogLineInformation_ERROR,
			Classification: pganalyze_collector.LogLineInformation_CHECK_CONSTRAINT_VIOLATION,
		}, {
			LogLevel:       pganalyze_collector.LogLineInformation_ERROR,
			Classification: pganalyze_collector.LogLineInformation_CHECK_CONSTRAINT_VIOLATION,
		}, {
			LogLevel:       pganalyze_collector.LogLineInformation_ERROR,
			Classification: pganalyze_collector.LogLineInformation_CHECK_CONSTRAINT_VIOLATION,
		}},
		nil,
	},
	{
		[]state.LogLine{{
			Content:  "conflicting key value violates exclusion constraint \"reservation_during_excl\"",
			LogLevel: pganalyze_collector.LogLineInformation_ERROR,
			UUID:     uuid.UUID{1},
		}, {
			Content:  "Key (during)=([\"2010-01-01 14:45:00\",\"2010-01-01 15:45:00\")) conflicts with existing key (during)=([\"2010-01-01 11:30:00\",\"2010-01-01 15:00:00\")).",
			LogLevel: pganalyze_collector.LogLineInformation_DETAIL,
		}, {
			Content:  "INSERT INTO reservation VALUES ('[2010-01-01 14:45, 2010-01-01 15:45)');",
			LogLevel: pganalyze_collector.LogLineInformation_STATEMENT,
		}},
		[]state.LogLine{{
			LogLevel:       pganalyze_collector.LogLineInformation_ERROR,
			Classification: pganalyze_collector.LogLineInformation_EXCLUSION_CONSTRAINT_VIOLATION,
			UUID:           uuid.UUID{1},
			Query:          "INSERT INTO reservation VALUES ('[2010-01-01 14:45, 2010-01-01 15:45)');",
		}, {
			LogLevel:   pganalyze_collector.LogLineInformation_DETAIL,
			ParentUUID: uuid.UUID{1},
		}, {
			LogLevel:   pganalyze_collector.LogLineInformation_STATEMENT,
			ParentUUID: uuid.UUID{1},
		}},
		nil,
	},
}

func TestAnalyzeLogLines(t *testing.T) {
	for _, pair := range tests {
		l, s := logs.AnalyzeLogLines(pair.logLinesIn)

		cfg := pretty.CompareConfig
		cfg.SkipZeroFields = true

		if diff := cfg.Compare(pair.logLinesOut, l); diff != "" {
			t.Errorf("For %v: log lines diff: (-got +want)\n%s", pair.logLinesIn, diff)
		}
		if diff := cfg.Compare(pair.samplesOut, s); diff != "" {
			t.Errorf("For %v: query samples diff: (-got +want)\n%s", pair.samplesOut, diff)
		}
	}
}
