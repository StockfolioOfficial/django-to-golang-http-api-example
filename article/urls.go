package article

import "github.com/stockfolioofficial/django-to-golang-rest-api-example/supporter"

var URLPatterns = supporter.Routes{
	supporter.Path("/articles", articleAsView()),
	supporter.Path("/articles/:article_id", articleDetailAsView()),
}