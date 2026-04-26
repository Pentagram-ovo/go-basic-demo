package service

import (
	"errors"
	"go-forum/dao"
	"go-forum/model"
)

// CreateComment 发布新评论
func CreateComment(userid, postid uint, content string) error {
	var comment model.Comment
	post, err := dao.GetPostByPostId(postid)
	if err != nil {
		return err // 数据库错误直接返回
	}
	if post == nil {
		return errors.New("帖子不存在")
	}
	comment.UserID = userid
	comment.PostID = postid
	comment.Content = content
	return dao.CreateComment(&comment)
}

// GetCommentListByPostID 按帖子id分页查询评论
func GetCommentListByPostID(postid uint, page, size int) ([]model.Comment, int64, error) {
	if page < 1 || size > 20 {
		return nil, 0, errors.New("请调整参数")
	}
	comments, total, err := dao.GetCommentListByPostID(postid, page, size)
	if err != nil {
		return nil, 0, err
	}
	return comments, total, nil
}

// GetCommentByUserID 按用户id查询其所有评论
func GetCommentByUserID(userid uint) ([]model.Comment, error) {
	comments, err := dao.GetCommentByUserId(userid)
	if err != nil {
		return nil, err
	}
	return comments, nil
}
