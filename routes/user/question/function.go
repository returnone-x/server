package question

import (
	"database/sql"
	"github.com/gofiber/fiber/v2"
	jwt "github.com/golang-jwt/jwt/v5"
	questionDatabase "github.com/returnone-x/server/database/question"
	utils "github.com/returnone-x/server/utils"
)

func Vote(c *fiber.Ctx, vote int) error {
	// get params
	params := c.AllParams()

	// if user not send parms data
	if params["question_id"] == "" {
		return c.Status(400).JSON(utils.InvalidRequest())
	}

	// get question data from database for check does this question exist
	_, err := questionDatabase.GetQuestionData(params["question_id"])

	// handle error
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(400).JSON(utils.ErrorMessage("Can't find this question", err))
		} else {
			return c.Status(500).JSON(utils.ErrorMessage("Error when get question from database", err))
		}
	}

	// get user_id from accessToken cookie
	token := c.Locals("access_token_context").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	user_id := claims["user_id"].(string)

	// get user vote data(this question)
	user_vote_result, err := questionDatabase.GetUserQuestionVote(params["question_id"], user_id)

	// handle error (exclude didnt find row error)
	if err != nil && err != sql.ErrNoRows {
		return c.Status(500).JSON(utils.ErrorMessage("Error get data", err))
	}

	// if the top ticket result is the same as what is being updated/created now
	if user_vote_result.Vote == vote {
		return c.Status(400).JSON(utils.ErrorMessage("You already voted", err))
	
	// else if got the vote data
	} else if err != sql.ErrNoRows {

		// update question vote
		update_result, err := questionDatabase.UpdateQuestionVote(params["question_id"], user_id, vote)

		// handle error
		if err != nil {
			return c.Status(500).JSON(utils.ErrorMessage("Error update data", err))
		}

		// check does it really update
		row_affected, _ := update_result.RowsAffected()

		// if not update return server error or return 200 status code
		if row_affected == 0 {
			return c.Status(500).JSON(utils.ErrorMessage("Error update data", err))
		} else {
			return c.Status(200).JSON(fiber.Map{
				"status":  "success",
				"message": "Successful update vote",
				"data": fiber.Map{
					"question_id": params["question_id"],
					"voter_id":    user_id,
					"vote":        vote,
				},
			})
		}

	}

	// create question vote
	result, err := questionDatabase.NewQuestionVote(params["question_id"], user_id, vote)

	if err != nil {
		return c.Status(500).JSON(utils.ErrorMessage("Error create data", err))
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Successful create vote",
		"data":    result,
	})
}
