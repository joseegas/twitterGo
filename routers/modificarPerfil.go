package routers

import (
	"context"
	"encoding/json"
	"twitterGo/bd"
	"twitterGo/models"
)

func ModificarPerfil(ctx context.Context, claim models.Claim) models.RespApi{
	var r models.RespApi
	r.Status=400

	var t models.Usuario

	body:= ctx.Value(models.Key("body")).(string)
	err:=json.Unmarshal([]byte(body), &t)
	if err!=nil{
		r.Message="Datos Incorrecto "+err.Error()
	}

	status,err:=bd.ModificoRegistro(t,claim.ID.Hex())
	if err!=nil{
		r.Message="Ocurrió un error al intentar modificar el registro. "+err.Error()
		return r
	}

	if !status{
		r.Message="No se ha logrado modificar el registro del usuario. "+err.Error()
		return r
	}

	r.Status=200
	r.Message="Modificación de perfil ok"
	return r
}