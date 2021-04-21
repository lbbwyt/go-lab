package blue_king_client

//blue king API调用说明

//请求API时，需提供APP应用认证信息和用户认证信息，用以对APP应用和当前用户进行认证。

//使用API前，需先根据已有应用或创建一个应用，获取应用ID和安全密钥，作为应用认证信息。

//应用ID: 应用唯一标识，创建应用时由创建者指定，可在应用基本信息中通过 应用 ID 获取
//安全密钥: 应用密钥，创建应用后由平台生成，可在应用基本信息中通过 安全密钥 获取

//用户认证
//用户认证，通过验证用户登录态实现。用户登录态bk_token，在用户登录后，存储在浏览器的Cookies中。
//
//调用API时，若无法提供用户登录态，可将应用ID添加到应用免登录态验证白名单中，则此应用请求API时，提供当前用户的bk_username即可。
//
//functioncontroller: 通用白名单管理，通过指定不同的功能code，维护不同的白名单
//user_auth::skip_user_auth: "应用免登录态验证白名单" 功能code，添加此功能code，然后将应用ID添加到"功能测试白名单"中

//测试应用信息
//应用Id: felix
//secret: d137926b-e0ed-40c9-81b3-41b53b102384
//bk_token: J26-ynIUcS-Zk2ofjY5ueKQfQkA1zk58Uv73VDfslfE

//调用示例：
//http://bk-paas.cn/api/c/compapi/v2/cc/search_business/
//
//{
//	"bk_app_code": "felix",
//	"bk_app_secret": "d137926b-e0ed-40c9-81b3-41b53b102384",
//	"bk_token": "J26-ynIUcS-Zk2ofjY5ueKQfQkA1zk58Uv73VDfslfE",
//	"page": {
//			"start": 0,
//			"limit": 10,
//			"sort": ""
//			}
//}
