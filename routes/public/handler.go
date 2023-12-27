package public

import (
	"database/sql"
	"strconv"

	"github.com/gofiber/fiber/v2"
	jwt "github.com/golang-jwt/jwt/v5"
	questionDatabase "github.com/returnone-x/server/database/question"
	questionCommentDatabase "github.com/returnone-x/server/database/question/comment"
	questionModal "github.com/returnone-x/server/models/question"
	utils "github.com/returnone-x/server/utils"
)

func GetQuestion(c *fiber.Ctx) error {

	params := c.AllParams()

	if params["id"] == "" {
		return c.Status(400).JSON(utils.InvalidRequest())
	}
	question_data, err := questionDatabase.GetQuestionData(params["id"])

	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(400).JSON(utils.ErrorMessage("Can't find this question", err))
		} else {
			return c.Status(500).JSON(utils.ErrorMessage("Error when get question from database", err))
		}
	}
	up_vote_count, err := questionDatabase.GetQuestionVoteCount(params["id"], 1)

	if err != nil {
		return c.Status(500).JSON(utils.ErrorMessage("Error get up vote count", err))
	}

	down_vote_count, err := questionDatabase.GetQuestionVoteCount(params["id"], 2)

	if err != nil {
		return c.Status(500).JSON(utils.ErrorMessage("Error get down vote count", err))
	}

	vote_count := up_vote_count - down_vote_count
	check_token := c.Locals("access_token_context")
	if check_token == nil {
		return_result := questionModal.ReturnResult{
			Id:            question_data.Id,
			Questioner_id: question_data.Questioner_id,
			Title:         question_data.Title,
			Content:       question_data.Content,
			Tags_name:     question_data.Tags_name,
			Tags_version:  question_data.Tags_version,
			Views:         question_data.Views,
			Vote_count:    vote_count,
			User_vote:     0,
			Create_at:     question_data.Create_at,
			Update_at:     question_data.Update_at,
		}
		return c.Status(200).JSON(fiber.Map{
			"status":  "success",
			"message": "Successful get the question data",
			"data":    return_result,
		})
	}

	token := c.Locals("access_token_context").(*jwt.Token)

	claims := token.Claims.(jwt.MapClaims)
	// get user_id from accessToken cookie
	user_id := claims["user_id"].(string)
	// get user vote data(this question)
	user_vote_result, err := questionDatabase.GetUserQuestionVote(params["id"], user_id)

	if err != nil && err != sql.ErrNoRows {
		return c.Status(500).JSON(utils.ErrorMessage("Error get data", err))
	}

	return_result := questionModal.ReturnResult{
		Id:            question_data.Id,
		Questioner_id: question_data.Questioner_id,
		Title:         question_data.Title,
		Content:       question_data.Content,
		Tags_name:     question_data.Tags_name,
		Tags_version:  question_data.Tags_version,
		Views:         question_data.Views,
		Vote_count:    vote_count,
		User_vote:     user_vote_result.Vote,
		Create_at:     question_data.Create_at,
		Update_at:     question_data.Update_at,
	}
	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Successful get the question data",
		"data":    return_result,
	})
}

func GetQuestionComment(c *fiber.Ctx) error {
	params := c.AllParams()
	limit := c.Query("limit")

	limit_number, err := strconv.Atoi(limit)
	if err != nil {
		return c.Status(400).JSON(utils.RequestValueValid("limit"))
	}

	if limit_number > 100 || limit_number <= 0 {
		return c.Status(400).JSON(utils.RequestValueValid("limit"))
	}

	// if user not send parms data
	if params["question_id"] == "" {
		return c.Status(400).JSON(utils.InvalidRequest())
	}

	result, err := questionCommentDatabase.GetQuestionComments(params["question_id"], limit_number)

	if err != nil {
		return c.Status(500).JSON(utils.ErrorMessage("Error create data", err))
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "successful create new comment data",
		"data":    result,
	})
}
