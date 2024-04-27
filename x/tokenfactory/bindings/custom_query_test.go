package bindings_test

//func TestQueryFullDenom(t *testing.T) {
//	actor := RandomAccountAddress()
//	junoapp, ctx := SetupCustomApp(t, actor)
//
//	reflect := instantiateReflectContract(t, ctx, junoapp, actor)
//	require.NotEmpty(t, reflect)
//
//	// query full denom
//	query := bindings.TokenFactoryQuery{
//		FullDenom: &bindings.FullDenom{
//			CreatorAddr: reflect.String(),
//			Subdenom:    "ustart",
//		},
//	}
//	resp := bindings.FullDenomResponse{}
//	queryCustom(t, ctx, junoapp, reflect, query, &resp)
//
//	expected := fmt.Sprintf("factory/%s/ustart", reflect.String())
//	require.EqualValues(t, expected, resp.Denom)
//}
//
//type ReflectQuery struct {
//	Chain *ChainRequest `json:"chain,omitempty"`
//}
//
//type ChainRequest struct {
//	Request wasmvmtypes.QueryRequest `json:"request"`
//}
//
//type ChainResponse struct {
//	Data []byte `json:"data"`
//}
//
//func queryCustom(t *testing.T, ctx sdk.Context, junoapp *app.App, contract sdk.AccAddress, request bindings.TokenFactoryQuery, response interface{}) {
//	msgBz, err := json.Marshal(request)
//	require.NoError(t, err)
//	fmt.Println("queryCustom1", string(msgBz))
//
//	query := ReflectQuery{
//		Chain: &ChainRequest{
//			Request: wasmvmtypes.QueryRequest{Custom: msgBz},
//		},
//	}
//	queryBz, err := json.Marshal(query)
//	require.NoError(t, err)
//	fmt.Println("queryCustom2", string(queryBz))
//
//	resBz, err := junoapp.AppKeepers.WasmKeeper.QuerySmart(ctx, contract, queryBz)
//	require.NoError(t, err)
//	var resp ChainResponse
//	err = json.Unmarshal(resBz, &resp)
//	require.NoError(t, err)
//	err = json.Unmarshal(resp.Data, response)
//	require.NoError(t, err)
//}
