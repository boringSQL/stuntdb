package facts

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"slices"
)

var plannerSettings = []string{
	"constraint_exclusion",
	"cpu_index_tuple_cost",
	"cpu_operator_cost",
	"cpu_tuple_cost",
	"cursor_tuple_fraction",
	"default_statistics_target",
	"effective_cache_size",
	"enable_async_append",
	"enable_bitmapscan",
	"enable_distinct_reordering",
	"enable_eager_aggregate",
	"enable_gathermerge",
	"enable_group_by_reordering",
	"enable_hashagg",
	"enable_hashjoin",
	"enable_incremental_sort",
	"enable_indexonlyscan",
	"enable_indexscan",
	"enable_material",
	"enable_memoize",
	"enable_mergejoin",
	"enable_nestloop",
	"enable_parallel_append",
	"enable_parallel_hash",
	"enable_partition_pruning",
	"enable_partitionwise_aggregate",
	"enable_partitionwise_join",
	"enable_presorted_aggregate",
	"enable_self_join_elimination",
	"enable_seqscan",
	"enable_sort",
	"enable_tidscan",
	"from_collapse_limit",
	"geqo",
	"geqo_effort",
	"geqo_generations",
	"geqo_pool_size",
	"geqo_seed",
	"geqo_selection_bias",
	"geqo_threshold",
	"hash_mem_multiplier",
	"jit",
	"jit_above_cost",
	"jit_inline_above_cost",
	"jit_optimize_above_cost",
	"join_collapse_limit",
	"max_parallel_workers",
	"max_parallel_workers_per_gather",
	"min_eager_agg_group_size",
	"min_parallel_index_scan_size",
	"min_parallel_table_scan_size",
	"parallel_setup_cost",
	"parallel_tuple_cost",
	"plan_cache_mode",
	"random_page_cost",
	"recursive_worktable_factor",
	"seq_page_cost",
	"work_mem",
}

var extraSettings = []string{
	"work_mem",
	"hash_mem_multiplier",
	"max_parallel_workers",
	"max_parallel_workers_per_gather",
}

const settingsQuery = `
select name, setting
from pg_settings
where context = 'user'
  and (category like 'Query Tuning%' or name = any($1))
order by name`

func CaptureSettings(db *sql.DB) (map[string]string, error) {
	rows, err := db.Query(settingsQuery, extraSettings)
	if err != nil {
		return nil, fmt.Errorf("reading settings: %w", err)
	}
	defer rows.Close()

	settings := map[string]string{}
	for rows.Next() {
		var name, value string
		if err := rows.Scan(&name, &value); err != nil {
			return nil, fmt.Errorf("reading settings: %w", err)
		}
		settings[name] = value
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("reading settings: %w", err)
	}
	return settings, nil
}

func LoadSettings(data []byte) (map[string]string, error) {
	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, err
	}
	settings := make(map[string]string, len(raw))
	for name, value := range raw {
		if !slices.Contains(plannerSettings, name) {
			return nil, fmt.Errorf("%q is not a planner setting", name)
		}
		text, ok := value.(string)
		if !ok {
			return nil, fmt.Errorf("%q must be a string", name)
		}
		settings[name] = text
	}
	return settings, nil
}

func MarshalSettings(settings map[string]string) ([]byte, error) {
	data, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return nil, err
	}
	return append(data, '\n'), nil
}
