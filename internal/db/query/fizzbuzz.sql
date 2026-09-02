-- name: UpsertFizzBuzzStat :exec
INSERT INTO fizzbuzz_stats (int1, int2, "limit", str1, str2, hits)
VALUES ($1, $2, $3, $4, $5, 1)
ON CONFLICT (int1, int2, "limit", str1, str2)
DO UPDATE SET
    hits       = fizzbuzz_stats.hits + 1,
    updated_at = NOW();

-- name: GetTopFizzBuzzStat :one
SELECT int1, int2, "limit", str1, str2, hits
FROM fizzbuzz_stats
ORDER BY hits DESC
LIMIT 1;
