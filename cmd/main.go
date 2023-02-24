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

func searchFreeWord() {
	dbClientSession := session.Must(session.NewSession())
	// Create a DynamoDB client with additional configuration
	ddb := dynamodb.New(dbClientSession, aws.NewConfig().WithRegion("ap-northeast-1"))

	searchWord := "content"
	params := &dynamodb.ScanInput{
		TableName:        aws.String("articles"),
		FilterExpression: aws.String("contains(content, :search_text)"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":search_text": {
				S: aws.String(searchWord),
			},
		},
	}
	res, err := ddb.Scan(params)
	if err != nil {
		panic(err)
	}

	fmt.Println(res)
}

func queryByPkAndSk(article_id string, sk string) {
	dbClientSession := session.Must(session.NewSession())
	// Create a DynamoDB client with additional configuration
	ddb := dynamodb.New(dbClientSession, aws.NewConfig().WithRegion("ap-northeast-1"))

	params := &dynamodb.QueryInput{
		TableName:              aws.String("articles"),
		KeyConditionExpression: aws.String("article_id = :id AND SK = :SK"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":id": {
				S: aws.String(article_id),
			},
			":SK": {
				S: aws.String(sk),
			},
		},
	}
	res, err := ddb.Query(params)
	if err != nil {
		panic(err)
	}

	fmt.Println(res)
}

func reqArticlesByCategories(article_id string, sk string) {
	dbClientSession := session.Must(session.NewSession())
	// Create a DynamoDB client with additional configuration
	ddb := dynamodb.New(dbClientSession, aws.NewConfig().WithRegion("ap-northeast-1"))

	params := &dynamodb.QueryInput{
		TableName:              aws.String("articles"),
		KeyConditionExpression: aws.String("article_id = :id AND begins_with(SK, :SK)"),
		// FilterExpression:       aws.String("begins_with(categories, )"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":id": {
				S: aws.String(article_id),
			},
			":SK": {
				S: aws.String(sk),
			},
		},
	}
	res, err := ddb.Query(params)
	if err != nil {
		panic(err)
	}

	fmt.Println(res)
}

func main() {
	// searchFreeWord()
	// queryByPkAndSk("id1", "title")
	reqArticlesByCategories("id1", "t")
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
