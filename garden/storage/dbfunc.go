package dbfunc

import (
	"context"
	plantsfunc "garden/plants"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresDB struct {
	pl  *pgxpool.Pool
	ctx context.Context
}

func NewPostgresDB(context context.Context, pool *pgxpool.Pool) *PostgresDB {
	return &PostgresDB{pl: pool, ctx: context}
}

func (db *PostgresDB) GetAllPlants() ([]plantsfunc.Plant, error) {
	rows, err := db.pl.Query(db.ctx, "SELECT * FROM plants ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var plants []plantsfunc.Plant
	for rows.Next() {
		var p plantsfunc.Plant
		err := rows.Scan(
			&p.ID,
			&p.Name,
			&p.Species,
			&p.Health,
			&p.Water,
			&p.Fertilizer,
			&p.PlantedAt,
			&p.LastWatered,
			&p.LastFertilized,
			&p.Stage,
			&p.Status,
		)
		if err != nil {
			return nil, err
		}
		plants = append(plants, p)
	}
	return plants, nil
}

func (db *PostgresDB) PlantExists(id int) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM plants WHERE id = $1)`
	err := db.pl.QueryRow(db.ctx, query, id).Scan(&exists)
	return exists, err
}

func (db *PostgresDB) PlantExistsByName(name string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM plants WHERE name = $1)`
	err := db.pl.QueryRow(db.ctx, query, name).Scan(&exists)
	return exists, err
}

func (db *PostgresDB) UpdatePlant(plant *plantsfunc.Plant) error {
	_, err := db.pl.Exec(db.ctx,
		`UPDATE plants SET 
			health = $1, 
			water = $2, 
			fertilizer = $3,
			last_watered = $4,
			last_fertilized = $5,
			stage = $6,
			status = $7
		WHERE id = $8`,
		plant.Health,
		plant.Water,
		plant.Fertilizer,
		plant.LastWatered,
		plant.LastFertilized,
		plant.Stage,
		plant.Status,
		plant.ID,
	)
	return err
}

func (db *PostgresDB) WaterPlant(id int) error {
	query := `UPDATE plants SET
        water = LEAST(water + 30, 100),
        last_watered = NOW(),
        status = 'Watered'
        WHERE id = $1`
	_, err := db.pl.Exec(db.ctx, query, id)
	return err
}

func (db *PostgresDB) FertilizePlant(id int) error {
	query := `UPDATE plants SET 
        fertilizer = LEAST(fertilizer + 25, 100),
        last_fertilized = NOW(),
        status = 'Fertilized'
        WHERE id = $1`
	_, err := db.pl.Exec(db.ctx, query, id)
	return err
}

func (db *PostgresDB) DeletePlant(id int) error {
	query := `DELETE FROM plants WHERE id = $1 RETURNING id`
	var deletedID int
	err := db.pl.QueryRow(db.ctx, query, id).Scan(&deletedID)
	return err
}

func (db *PostgresDB) CreatePlantWithParams(name, species string) error {
	query := `INSERT INTO plants (
        name, species, health, water, fertilizer, 
        planted_at, last_watered, last_fertilized, stage, status
    ) VALUES ($1, $2, 100, 100, 100, NOW(), NOW(), NOW(), 'seedling', 'healthy')`

	_, err := db.pl.Exec(db.ctx, query, name, species)
	return err
}
