package sqlite

import (
	"database/sql"

	"github.com/RaghavSood/collectibles/types"
)

func (d *SqliteBackend) SeriesSummaries() ([]types.SeriesSummary, error) {
	query := `SELECT slug, name, item_count, tvl, unfunded, redeemed, unredeemed FROM series_summary ORDER BY tvl DESC;`

	rows, err := d.db.Query(query)
	if err != nil {
		return nil, err
	}

	return scanSeriesSummaries(rows)
}

func (d *SqliteBackend) SeriesSummariesByCreator(slug string) ([]types.SeriesSummary, error) {
	query := `SELECT
  ss.slug,
  ss.name,
  ss.item_count,
  ss.tvl,
  ss.unfunded,
  ss.redeemed,
  ss.unredeemed
FROM
  series_summary ss
JOIN
  series_creators sc ON ss.slug = sc.series_slug
WHERE
  sc.creator_slug = ?
ORDER BY
  ss.tvl DESC;`

	rows, err := d.db.Query(query, slug)
	if err != nil {
		return nil, err
	}

	return scanSeriesSummaries(rows)
}

func (d *SqliteBackend) SeriesSummary(slug string) (*types.SeriesSummary, error) {
	query := `SELECT slug, name, item_count, tvl, unfunded, redeemed, unredeemed FROM series_summary WHERE slug = ?;`

	row := d.db.QueryRow(query, slug)
	var summary types.SeriesSummary
	err := row.Scan(&summary.Slug, &summary.Name, &summary.ItemCount, &summary.TotalValue, &summary.Unfunded, &summary.Redeemed, &summary.Unredeemed)
	if err != nil {
		return &types.SeriesSummary{}, err
	}

	return &summary, nil
}

func scanSeriesSummaries(rows *sql.Rows) ([]types.SeriesSummary, error) {
	var summaries []types.SeriesSummary
	for rows.Next() {
		var summary types.SeriesSummary
		err := rows.Scan(&summary.Slug, &summary.Name, &summary.ItemCount, &summary.TotalValue, &summary.Unfunded, &summary.Redeemed, &summary.Unredeemed)
		if err != nil {
			return nil, err
		}
		summaries = append(summaries, summary)
	}
	return summaries, nil
}
