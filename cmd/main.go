package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	// "github.com/aws/aws-sdk-go/service/dynamodb/expression"
	// "log"
	// "net/http"
	//	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

const articleTableName string = "article"
const articleContentTableName string = "contents"
const region = "ap-northeast-1"

func searchArticlesByString(targetString string) {
	dbClientSession := session.Must(session.NewSession())
	ddb := dynamodb.New(dbClientSession, aws.NewConfig().WithRegion(region))

	params := &dynamodb.ScanInput{
		TableName:        aws.String(articleContentTableName),
		FilterExpression: aws.String("contains(content, :search_text)"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":search_text": {
				S: aws.String(targetString),
			},
		},
	}
	res, err := ddb.Scan(params)
	if err != nil {
		panic(err)
	}

	fmt.Println(res)
}

func searchArticles(article_id string, sk_attr string) {
	dbClientSession := session.Must(session.NewSession())
	ddb := dynamodb.New(dbClientSession, aws.NewConfig().WithRegion(region))

	// SKが指定されている場合
	var params *dynamodb.QueryInput
	if sk_attr != "" {
		params = &dynamodb.QueryInput{
			TableName:              aws.String(articleTableName),
			KeyConditionExpression: aws.String("article_id = :id AND SK = :SK"),
			ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
				":id": {
					S: aws.String(article_id),
				},
				":SK": {
					S: aws.String(sk_attr),
				},
			},
		}
	} else {
		params = &dynamodb.QueryInput{
			TableName:              aws.String(articleTableName),
			KeyConditionExpression: aws.String("article_id = :id"),
			ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
				":id": {
					S: aws.String(article_id),
				},
			},
		}
	}

	res, err := ddb.Query(params)
	if err != nil {
		panic(err)
	}

	fmt.Println(res)
}

func sortByPostedDate(createdAt articleCreatedAtInput) {
	dbClientSession := session.Must(session.NewSession())
	ddb := dynamodb.New(dbClientSession, aws.NewConfig().WithRegion(region))

	var params *dynamodb.QueryInput
	if createdAt.date != "" {
		params = &dynamodb.QueryInput{
			TableName: aws.String(articleTableName),
			IndexName: aws.String("created_ym-created_dms-index"),
			// ソートキーをcategories#TRIPなど、1つ1つのcategoryごとにキーを分けているため、このクエリを実行しやすい
			KeyConditionExpression: aws.String("created_ym = :year_month AND begins_with(created_dms, :created_date)"),
			ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
				":year_month": {
					S: aws.String(createdAt.year + "-" + createdAt.month),
				},
				":created_date": {
					S: aws.String(createdAt.date),
				},
			},
		}
	} else {
		params = &dynamodb.QueryInput{
			TableName: aws.String(articleTableName),
			IndexName: aws.String("created_ym-created_dms-index"),
			// ソートキーをcategories#TRIPなど、1つ1つのcategoryごとにキーを分けているため、このクエリを実行しやすい
			KeyConditionExpression: aws.String("created_ym = :year_month"),
			ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
				":year_month": {
					S: aws.String(createdAt.year + "-" + createdAt.month),
				},
			},
		}
	}
	res, err := ddb.Query(params)
	if err != nil {
		panic(err)
	}

	fmt.Println(res)
}

type articleCreatedAtInput struct {
	year  string
	month string
	date  string
}

func main() {
	searchArticles("id1", "title")
	// 月単位での読み取りと、日付指定の読み取りが可能
	// 昇順並び替えのための、何時とか何分とかもレスポンスされる
	createdAtToSearch := articleCreatedAtInput{"2023", "2", "25"}
	sortByPostedDate(createdAtToSearch)

	// searchByTitle()
	// searchArticlesByString("test_string")
	// queryByPkAndSk("id1", "title")
	// reqArticlesByCategories("id1", "TRIP")
	// フリーワード検索
	// カテゴリ別検索
	// 新着順検索
	//

	// // Sessionの作成時に認証情報や環境変数？の設定が読み込まれる
	// // Mustメソッドはセッション作成のエラーハンドリングを行う
	// const attributeToGet string = "content"
	// dbClientSession := session.Must(session.NewSession())
	// // Create a DynamoDB client with additional configuration
	// ddb := dynamodb.New(dbClientSession, aws.NewConfig().WithRegion("ap-northeast-1"))

	// params := &dynamodb.GetItemInput{
	// 	// aws.Stringはポインタ型への変換を行う
	// 	TableName: aws.String("articles"), // テーブル名

	// 	Key: map[string]*dynamodb.AttributeValue{
	// 		// プライマリキー
	// 		"article_id": { // キー名
	// 			S: aws.String("id1"), // 持ってくるキーの値
	// 		},
	// 		// ソートキー
	// 		"categories": {
	// 			S: aws.String("category1"),
	// 		},
	// 	},
	// 	AttributesToGet: []*string{
	// 		aws.String(attributeToGet),
	// 	},
	// 	// 読み取り整合性
	// 	ConsistentRead: aws.Bool(false),

	// 	//返ってくるデータの種類
	// 	ReturnConsumedCapacity: aws.String("TOTAL"),
	// }

	// resp, err := ddb.GetItem(params)

	// if err != nil {
	// 	fmt.Println(err.Error())
	// }

	// //resp.Item[項目名].型 でデータへのポインタを取得
	// fmt.Println(*resp.Item[attributeToGet].S)

	// // handlers
	// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	w.Write("base url health check ok!")
	// })

	// if err := http.ListenAndServe(":8080", nil); err != nil {
	// 	log.Fatal("ListenAndServe:", err)
	// }
}
