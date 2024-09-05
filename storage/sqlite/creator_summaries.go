package sqlite

import (
	"database/sql"

	"github.com/RaghavSood/collectibles/types"
)

func (d *SqliteBackend) CreatorSummaries() ([]types.CreatorSummary, error) {
	query := `SELECT name, slug, series_count, item_count, tvl, unfunded, redeemed, unredeemed FROM creator_summary ORDER BY tvl desc;`

	rows, err := d.db.Query(query)
	if err != nil {
		return nil, err
	}

	return scanCreatorSummaries(rows)
}

func (d *SqliteBackend) ScamCreatorSummaries() ([]types.CreatorSummary, error) {
	query := `SELECT name, slug, series_count, item_count, tvl, unfunded, redeemed, unredeemed FROM creator_summary WHERE slug IN (SELECT flag_key FROM flags WHERE flag_scope = 'creators' AND flag_type = 'scam') ORDER BY tvl desc;`

	rows, err := d.db.Query(query)
	if err != nil {
		return nil, err
	}

	return scanCreatorSummaries(rows)
}

func (d *SqliteBackend) CreatorSummary(creatorSlug string) (*types.CreatorSummary, error) {
	query := `SELECT name, slug, series_count, item_count, tvl, unfunded, redeemed, unredeemed FROM creator_summary WHERE slug = ?;`

	row := d.db.QueryRow(query, creatorSlug)

	var summary types.CreatorSummary
	err := row.Scan(&summary.Name, &summary.Slug, &summary.SeriesCount, &summary.ItemCount, &summary.TotalValue, &summary.Unfunded, &summary.Redeemed, &summary.Unredeemed)
	if err != nil {
		return nil, err
	}

	return &summary, nil
}

func scanCreatorSummaries(rows *sql.Rows) ([]types.CreatorSummary, error) {
	var summaries []types.CreatorSummary
	for rows.Next() {
		var summary types.CreatorSummary
		err := rows.Scan(&summary.Name, &summary.Slug, &summary.SeriesCount, &summary.ItemCount, &summary.TotalValue, &summary.Unfunded, &summary.Redeemed, &summary.Unredeemed)
		if err != nil {
			return nil, err
		}
		summaries = append(summaries, summary)
	}
	return summaries, nil
}
