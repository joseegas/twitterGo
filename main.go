package main

import (
	"context"
	"os"
	"strings"
	"twitterGo/awsgo"
	"twitterGo/bd"
	"twitterGo/handlers"
	"twitterGo/models"
	"twitterGo/secretmanager"

	"github.com/aws/aws-lambda-go/events"
	lambda "github.com/aws/aws-lambda-go/lambda"
)

func main(){
	lambda.Start(EjecutoLambda)
}

func EjecutoLambda(ctx context.Context, request events.APIGatewayProxyRequest)(*events.APIGatewayProxyResponse,error){
	var res *events.APIGatewayProxyResponse
	awsgo.InicializoAWS()

	if !ValidoParametros(){
		res = &events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body: "Error en las variables de entorno. Deben incluir 'SecretName','BucketName','UrlPrefix'",
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}
		return res,nil
	}

	SecretModel,err:=secretmanager.GetSecret(os.Getenv("SecretName"))
	if err!=nil{
		res = &events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body: "Error en la lectura de Secret "+err.Error(),
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}
		return res,nil
	}

	path:=strings.Replace(request.PathParameters["twittergo"],os.Getenv("UrlPrefix"),"",-1)

	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("path"), path)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("method"), request.HTTPMethod)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("user"), SecretModel.Username)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("password"), SecretModel.Password)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("host"), SecretModel.Host)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("database"), SecretModel.Database)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("jwtSign"), SecretModel.JWTSign)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("body"), request.Body)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("bucketName"), os.Getenv("BucketName"))

	//Chequeo a la bd
	err=bd.ConectarBD(awsgo.Ctx)
	if err!=nil{
		res = &events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body: "Error conectando la BD "+err.Error(),
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}
		return res,nil
	}
	
	respApi:=handlers.Manejadores(awsgo.Ctx, request)
	if respApi.CustomResp ==nil{
		res = &events.APIGatewayProxyResponse{
			StatusCode: respApi.Status,
			Body: respApi.Message,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}
		return res,nil
	}else{
		return respApi.CustomResp, nil
	}

}

func ValidoParametros()bool{
	_, traeParametro:=os.LookupEnv("SecretName")
	if !traeParametro {
		return traeParametro
	}

	_, traeParametro=os.LookupEnv("BucketName")
	if !traeParametro {
		return traeParametro
	}

	_, traeParametro=os.LookupEnv("UrlPrefix")
	if !traeParametro {
		return traeParametro
	}

	return true;
}