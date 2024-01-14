package db

import (
	"context"
	"database/sql"
	"log"

	"github.com/duduBTW/steam-ssr/data"
)

func newNullInt(i int) sql.NullInt32 {
	if i == 0 {
		return sql.NullInt32{}
	}

	return sql.NullInt32{
		Valid: true,
		Int32: int32(i),
	}
}

func CreateCart(userId int, db *sql.DB) int {
	query := `
		INSERT INTO Cart
			(UserId)
		VALUES
			($1)
		RETURNING
			CartId
	`

	var CartId int
	err := db.QueryRow(query, newNullInt(userId)).Scan(&CartId)
	if err != nil {
		log.Fatal(err)
	}

	return CartId
}

func InsertGameToCart(gameId, cartId int, db *sql.DB) int {
	query := `
		INSERT INTO CartGame
			(GameId, CartId)
		VALUES
			($2, $3)
		RETURNING
			CartGameId
	`

	var CartGameId int
	err := db.QueryRow(query, gameId, cartId).Scan(&CartGameId)
	if err != nil {
		log.Fatal(err)
	}

	return CartGameId
}

func GetCartWithGamesById(cartId, userId int, db *sql.DB) []data.CartGame {
	var cartList []data.CartGame

	query := `
		SELECT 
			Cart.CartId, 
			CartGame.CartGameId, 
			IFNULL(Cart.UserId, 0), 
			Game.GameId,
			Game.Title,
			Game.CoverImage,
			Game.Price
		FROM 
			CartGame 
		JOIN
			Cart ON CartGame.CartId = Cart.CartId
		JOIN
			Game ON Game.GameId = CartGame.GameId
		WHERE
			Cart.UserId = $2 OR Cart.CartId = $1 
	`

	rows, err := db.Query(query, userId, cartId)
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	for rows.Next() {
		var cart data.CartGame

		err := rows.Scan(
			&cart.CartId,
			&cart.CartGameId,
			&cart.UserId,
			&cart.GameId,
			&cart.Title,
			&cart.CoverImage,
			&cart.Price,
		)
		if err != nil {
			log.Fatal(err)
		}

		cartList = append(cartList, cart)
	}

	return cartList
}

func CountCartEntries(cartId, userId int, db *sql.DB) int {
	query := `
		SELECT 
			COUNT(Cart.CartId)
		FROM 
			CartGame 
		JOIN
			Cart ON CartGame.CartId = Cart.CartId
		JOIN
			Game ON Game.GameId = CartGame.GameId
		WHERE
			Cart.UserId = $1 OR Cart.CartId = $2 
	`

	var count int
	db.QueryRow(query, userId, cartId).Scan(&count)
	return count
}

func VinculateUserToCart(userId, cartId int, db *sql.DB) {
	query := `
		UPDATE 
			Cart
		SET 
			UserId = $1
		WHERE 
			CartId = $2
	`

	_, err := db.ExecContext(context.TODO(), query, userId, cartId)
	if err != nil {
		log.Fatalf("Unable to update cart to vinculate user %s", err)
	}
}

func UserCartId(userId int, db *sql.DB) int {
	query := `
		SELECT 
			CartId 
		FROM 
			Cart 
		WHERE 
			UserId = $1
	`

	var cartId int
	err := db.QueryRow(query, userId).Scan(&cartId)
	if err != nil && err != sql.ErrNoRows {
		log.Fatal(err)
	}

	return cartId
}

func UserCartGameId(userId, cartGameId int, db *sql.DB) int {
	query := `
		SELECT 
			Cart.CartId 
		FROM 
			Cart
		JOIN
			CartGame ON Cart.CartId = CartGame.CartId  
		WHERE 
			Cart.UserId = $1 AND CartGame.CartGameId = $2
	`

	var cartId int
	err := db.QueryRow(query, userId, cartGameId).Scan(&cartId)
	if err != nil && err != sql.ErrNoRows {
		log.Fatal(err)
	}

	return cartId
}

func RemoveCartGame(cartGameId int, db *sql.DB) {
	query := `
		DELETE FROM 
			CartGame 
		WHERE 
			CartGameId = $1
	`

	db.ExecContext(context.TODO(), query, cartGameId)
}

func ClearCartByUserId(userId int, db *sql.DB) {
	cartGameQuery := `
		DELETE FROM CartGame
		WHERE 
			CartId IN (SELECT CartId FROM Cart WHERE UserId = $1)	
	`
	cartQuery := `
		DELETE FROM 
			Cart 
		WHERE 
			UserId = $1
	`

	db.ExecContext(context.TODO(), cartGameQuery, userId)
	db.ExecContext(context.TODO(), cartQuery, userId)
}

func ClearCartByCartId(cartId int, db *sql.DB) {
	cartGameQuery := `
		DELETE FROM 
			CartGame 
		WHERE 
			CartId = $1
	`

	cartQuery := `
		DELETE FROM 
			Cart 
		WHERE 
			CartId = $1
	`

	db.ExecContext(context.TODO(), cartGameQuery, cartId)
	db.ExecContext(context.TODO(), cartQuery, cartId)
}

func IsGameOnCart(gameId string, userId, cartId int, db *sql.DB) bool {
	query := `
		SELECT 
			COUNT(Cart.CartId)
		FROM 
			CartGame 
		JOIN
			Cart ON CartGame.CartId = Cart.CartId
		JOIN
			Game ON Game.GameId = CartGame.GameId
		WHERE
			(Cart.UserId = $1 OR Cart.CartId = $2) AND CartGame.GameId = $3
	`

	var count int
	db.QueryRow(query, userId, cartId, gameId).Scan(&count)
	return count > 0
}
