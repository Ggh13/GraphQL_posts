package unitTest

import (
	"context"
	"fmt"
	"qraphQL_posts/api/graph"
	"qraphQL_posts/api/graph/model"
	"qraphQL_posts/pkg/logger"
	"reflect"
)

func MainTest(ctx context.Context, resolver graph.Resolver) {

	var numberOfTest int
	numberOfTest = 1
	//Test 1
	var NewComment model.NewComment
	NewComment = model.NewComment{
		UserID:          1,
		ParentIDComment: -1,
		Content:         "Если перед тем как посадить тыкву добавить народное средство под названием... читать дальше",
		PostID:          1,
	}

	comment, err := resolver.Mutation().CreateComment(ctx, NewComment)
	if err == nil || comment != nil {
		logger.GetLoggerFromCtx(ctx).Info(ctx, fmt.Sprintf("@@@@@@ Test %v failed @@@@@@ %v", numberOfTest, err))
	} else {
		logger.GetLoggerFromCtx(ctx).Info(ctx, fmt.Sprintf("Test %v done", numberOfTest))
	}
	numberOfTest += 1

	//Test 2
	var NewPost model.NewPost
	NewPost = model.NewPost{
		UserID:      1,
		Content:     "Сажаем тыкву с Копатычем",
		Commentable: false,
	}
	post, err := resolver.Mutation().CreatePost(ctx, NewPost)
	if err == nil || post != nil {
		logger.GetLoggerFromCtx(ctx).Info(ctx, fmt.Sprintf("@@@@@@ Test %v failed @@@@@@ %v", numberOfTest, err))
	} else {
		logger.GetLoggerFromCtx(ctx).Info(ctx, fmt.Sprintf("Test %v done", numberOfTest))
	}
	numberOfTest += 1

	//Test 3
	var NewUser model.NewUser
	NewUser = model.NewUser{
		Name:    "Копатыч",
		Surname: "Пчолыч",
	}
	user, err := resolver.Mutation().CreateUser(ctx, NewUser)
	if err != nil || user == nil {
		logger.GetLoggerFromCtx(ctx).Info(ctx, fmt.Sprintf("@@@@@@ Test %v failed @@@@@@ %v", numberOfTest, err))
	} else {
		logger.GetLoggerFromCtx(ctx).Info(ctx, fmt.Sprintf("Test %v done. User created with id %v", numberOfTest, user.ID))
	}
	numberOfTest += 1

	//Test 4
	var DonePost model.Post
	DonePost = model.Post{
		ID:          1,
		User:        user,
		Content:     NewPost.Content,
		Commentable: NewPost.Commentable,
		Comments:    []*model.Comment{},
	}
	post, err = resolver.Mutation().CreatePost(ctx, NewPost)
	if err != nil || post == nil {
		logger.GetLoggerFromCtx(ctx).Info(ctx, fmt.Sprintf("@@@@@@ Test %v failed @@@@@@ %v", numberOfTest, err))
	} else {
		if !reflect.DeepEqual(*post, DonePost) {
			//logger.GetLoggerFromCtx(ctx).Info(ctx, fmt.Sprintf("%v %v %v", post.ID, post.Comments, post.Commentable))
			logger.GetLoggerFromCtx(ctx).Info(ctx, fmt.Sprintf("@@@@@@ Test %v failed @@@@@@ They are not equal", numberOfTest))
		} else {
			logger.GetLoggerFromCtx(ctx).Info(ctx, fmt.Sprintf("Test %v done", numberOfTest))
		}
	}
	numberOfTest += 1

	//Test 5
	var DoneComment model.Comment
	DoneComment = model.Comment{
		ID:              1,
		User:            user,
		ParentIDComment: -1,
		Content:         NewComment.Content,
		PostID:          1,
		Comments:        []*model.Comment{},
	}
	comment, err = resolver.Mutation().CreateComment(ctx, NewComment)
	if err == nil || comment != nil {
		logger.GetLoggerFromCtx(ctx).Info(ctx, fmt.Sprintf("@@@@@@ Test %v failed @@@@@@ %v", numberOfTest, err))
	} else {
		if comment == nil {
			logger.GetLoggerFromCtx(ctx).Info(ctx, fmt.Sprintf("Test %v done", numberOfTest))
		} else if !reflect.DeepEqual(*comment, DoneComment) {
			//logger.GetLoggerFromCtx(ctx).Info(ctx, fmt.Sprintf("%v %v %v", post.ID, post.Comments, post.Commentable))
			logger.GetLoggerFromCtx(ctx).Info(ctx, fmt.Sprintf("@@@@@@ Test %v failed @@@@@@ They are not equal", numberOfTest))
		} else {
			logger.GetLoggerFromCtx(ctx).Info(ctx, fmt.Sprintf("Test %v done", numberOfTest))
		}
	}
	numberOfTest += 1

	//Test 6
	var UpdPost model.UpdatePost
	UpdPost = model.UpdatePost{
		ID:          1,
		Commentable: true,
	}

	DonePost = model.Post{
		ID:          1,
		User:        user,
		Content:     NewPost.Content,
		Commentable: true,
		Comments:    []*model.Comment{},
	}

	post, err = resolver.Mutation().PostUpdate(ctx, UpdPost)
	if err != nil || post == nil {
		logger.GetLoggerFromCtx(ctx).Info(ctx, fmt.Sprintf("@@@@@@ Test %v failed @@@@@@ %v", numberOfTest, err))
	} else {
		if DonePost.Commentable != post.Commentable {
			//logger.GetLoggerFromCtx(ctx).Info(ctx, fmt.Sprintf("%v %v %v", post.ID, post.Comments, post.Commentable))
			logger.GetLoggerFromCtx(ctx).Info(ctx, fmt.Sprintf("@@@@@@ Test %v failed @@@@@@ They are not equal", numberOfTest))
		} else {
			logger.GetLoggerFromCtx(ctx).Info(ctx, fmt.Sprintf("Test %v done", numberOfTest))
		}
	}
	numberOfTest += 1

	//Test 7
	comment, err = resolver.Mutation().CreateComment(ctx, NewComment)
	if err != nil || comment == nil {
		logger.GetLoggerFromCtx(ctx).Info(ctx, fmt.Sprintf("@@@@@@ Test %v failed @@@@@@ %v", numberOfTest, err))
	} else {
		if !CheckEqualComment(*comment, DoneComment, false) {
			//logger.GetLoggerFromCtx(ctx).Info(ctx, fmt.Sprintf("%v %v %v", post.ID, post.Comments, post.Commentable))
			logger.GetLoggerFromCtx(ctx).Info(ctx, fmt.Sprintf("@@@@@@ Test %v failed @@@@@@ They are not equal", numberOfTest))
		} else {
			logger.GetLoggerFromCtx(ctx).Info(ctx, fmt.Sprintf("Test %v done", numberOfTest))
		}
	}
	numberOfTest += 1

}

func CheckEqualComment(com1 model.Comment, com2 model.Comment, withComments bool) bool {
	usl1 := (com1.ID == com2.ID && com1.Content == com2.Content && com1.PostID == com2.ID && com1.ParentIDComment == com2.ParentIDComment)
	isEqual := reflect.DeepEqual(com1.Comments, com2.Comments)
	if withComments {
		return usl1 && isEqual
	}
	return usl1
}
