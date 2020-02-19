// Copyright 2014 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package controller

import (
	"html/template"
	"net/http"
	"strings"

	"github.com/dchest/captcha"
	echo "github.com/labstack/echo/v4"
	"github.com/ue4-community/learnue/modules/goutils"
	"github.com/ue4-community/learnue/modules/logger"
	"github.com/ue4-community/learnue/modules/context"
	"github.com/ue4-community/learnue/modules/echoutils"

	. "github.com/ue4-community/learnue/http"
	"github.com/ue4-community/learnue/http/middleware"
	"github.com/ue4-community/learnue/logic"
	"github.com/ue4-community/learnue/model"
	"github.com/ue4-community/learnue/modules/util"
)

// 在需要评论（喜欢）且要回调的地方注册评论（喜欢）对象
func init() {
	// 注册评论（喜欢）对象
	logic.RegisterCommentObject(model.TypeArticle, logic.ArticleComment{})
	logic.RegisterLikeObject(model.TypeArticle, logic.ArticleLike{})
}

type ArticleController struct{}

// 注册路由
func (self ArticleController) RegisterRoute(g *echo.Group) {
	g.GET("/articles", self.ReadList)
	g.GET("/articles/crawl", self.Crawl)

	g.GET("/articles/:id", self.Detail)

	g.Match([]string{"GET", "POST"}, "/articles/new", self.Create, middleware.NeedLogin(), middleware.Sensivite(), middleware.BalanceCheck(), middleware.PublishNotice(), middleware.CheckCaptcha())
	g.Match([]string{"GET", "POST"}, "/articles/modify", self.Modify, middleware.NeedLogin(), middleware.Sensivite())
}

// ReadList 网友文章列表页
func (ArticleController) ReadList(ctx echo.Context) error {
	limit := 20

	curPage := goutils.MustInt(ctx.QueryParam("p"), 1)
	paginator := logic.NewPaginator(curPage)
	paginator.SetPerPage(limit)
	total := logic.DefaultArticle.Count(context.EchoContext(ctx), "")
	pageHtml := paginator.SetTotal(total).GetPageHtml(ctx.Request().URL.Path)
	pageInfo := template.HTML(pageHtml)

	// TODO: 参考的 topics 的处理方式，但是感觉不应该这样做
	topArticles := logic.DefaultArticle.FindAll(context.EchoContext(ctx), paginator, "id DESC", "top=1")
	unTopArticles := logic.DefaultArticle.FindAll(context.EchoContext(ctx), paginator, "id DESC", "top!=1")
	articles := append(topArticles, unTopArticles...)
	if articles == nil {
		logger.Errorln("article controller: find article error")
		return ctx.Redirect(http.StatusSeeOther, "/articles")
	}

	num := len(articles)
	if num == 0 {
		return render(ctx, "articles/list.html", map[string]interface{}{"articles": articles})
	}

	// 获取当前用户喜欢对象信息
	me, ok := ctx.Get("user").(*model.Me)
	var topLikeFlags map[int]int
	var unTopLikeFlags map[int]int
	likeFlags := map[int]int{}

	if ok {
		topArticlesNum := len(topArticles)
		if topArticlesNum > 0 {
			topLikeFlags, _ = logic.DefaultLike.FindUserLikeObjects(context.EchoContext(ctx), me.Uid, model.TypeArticle, topArticles[0].Id, topArticles[topArticlesNum-1].Id)
			for k, v := range topLikeFlags {
				likeFlags[k] = v
			}
		}

		unTopArticlesNum := len(unTopArticles)
		if unTopArticlesNum > 0 {
			unTopLikeFlags, _ = logic.DefaultLike.FindUserLikeObjects(context.EchoContext(ctx), me.Uid, model.TypeArticle, unTopArticles[0].Id, unTopArticles[unTopArticlesNum-1].Id)
			for k, v := range unTopLikeFlags {
				likeFlags[k] = v
			}
		}
	}

	return render(ctx, "articles/list.html", map[string]interface{}{"articles": articles, "page": pageInfo, "likeflags": likeFlags})
}

// Detail 文章详细页
func (ArticleController) Detail(ctx echo.Context) error {
	article, prevNext, err := logic.DefaultArticle.FindByIdAndPreNext(context.EchoContext(ctx), goutils.MustInt(ctx.Param("id")))
	if err != nil {
		return ctx.Redirect(http.StatusSeeOther, "/articles")
	}

	if article == nil || article.Id == 0 || article.Status == model.ArticleStatusOffline {
		return ctx.Redirect(http.StatusSeeOther, "/articles")
	}

	articleGCTT := logic.DefaultArticle.FindArticleGCTT(context.EchoContext(ctx), article)
	data := map[string]interface{}{
		"activeArticles": "active",
		"article":        article,
		"article_gctt":   articleGCTT,
		"prev":           prevNext[0],
		"next":           prevNext[1],
	}

	me, ok := ctx.Get("user").(*model.Me)
	if ok {
		data["likeflag"] = logic.DefaultLike.HadLike(context.EchoContext(ctx), me.Uid, article.Id, model.TypeArticle)
		data["hadcollect"] = logic.DefaultFavorite.HadFavorite(context.EchoContext(ctx), me.Uid, article.Id, model.TypeArticle)

		logic.Views.Incr(Request(ctx), model.TypeArticle, article.Id, me.Uid)

		if !article.IsSelf || me.Uid != article.User.Uid {
			go logic.DefaultViewRecord.Record(article.Id, model.TypeArticle, me.Uid)
		}

		if me.IsRoot || (article.IsSelf && me.Uid == article.User.Uid) {
			data["view_user_num"] = logic.DefaultViewRecord.FindUserNum(context.EchoContext(ctx), article.Id, model.TypeArticle)
			data["view_source"] = logic.DefaultViewSource.FindOne(context.EchoContext(ctx), article.Id, model.TypeArticle)
		}
	} else {
		logic.Views.Incr(Request(ctx), model.TypeArticle, article.Id)
	}

	// 为了阅读数即时看到
	article.Viewnum++

	data["subjects"] = logic.DefaultSubject.FindArticleSubjects(context.EchoContext(ctx), article.Id)

	return render(ctx, "articles/detail.html,common/comment.html", data)
}

// Create 发布新文章
func (ArticleController) Create(ctx echo.Context) error {
	me := ctx.Get("user").(*model.Me)

	title := ctx.FormValue("title")
	if title == "" || ctx.Request().Method != "POST" {
		data := map[string]interface{}{"activeArticles": "active"}
		if logic.NeedCaptcha(me) {
			data["captchaId"] = captcha.NewLen(util.CaptchaLen)
		}
		return render(ctx, "articles/new.html", data)
	}

	if ctx.FormValue("content") == "" {
		return fail(ctx, 1, "内容不能为空")
	}

	forms, _ := ctx.FormParams()
	id, err := logic.DefaultArticle.Publish(echoutils.WrapEchoContext(ctx), me, forms)
	if err != nil {
		return fail(ctx, 2, "内部服务错误")
	}

	return success(ctx, map[string]interface{}{"id": id})
}

// Modify 修改文章
func (ArticleController) Modify(ctx echo.Context) error {
	id := ctx.FormValue("id")
	article, err := logic.DefaultArticle.FindById(context.EchoContext(ctx), id)

	if ctx.Request().Method != "POST" {
		if err != nil {
			return ctx.Redirect(http.StatusSeeOther, "/articles/"+id)
		}

		return render(ctx, "articles/new.html", map[string]interface{}{
			"article":        article,
			"activeArticles": "active",
		})
	}

	if id == "" || ctx.FormValue("content") == "" {
		return fail(ctx, 1, "内容不能为空")
	}

	if err != nil {
		return fail(ctx, 2, "文章不存在")
	}

	me := ctx.Get("user").(*model.Me)
	if !logic.CanEdit(me, article) {
		return fail(ctx, 3, "没有修改权限")
	}

	forms, _ := ctx.FormParams()
	errMsg, err := logic.DefaultArticle.Modify(echoutils.WrapEchoContext(ctx), me, forms)
	if err != nil {
		return fail(ctx, 4, errMsg)
	}

	return success(ctx, map[string]interface{}{"id": article.Id})
}

func (ArticleController) Crawl(ctx echo.Context) error {
	strUrl := ctx.QueryParam("url")

	var (
		errMsg string
		err    error
	)
	strUrl = strings.TrimSpace(strUrl)
	_, err = logic.DefaultArticle.ParseArticle(context.EchoContext(ctx), strUrl, false)
	if err != nil {
		errMsg = err.Error()
	}

	if errMsg != "" {
		return fail(ctx, 1, errMsg)
	}
	return success(ctx, nil)
}
