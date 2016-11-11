DEBUG = false;
MOCK_DATA = true;
if (MOCK_DATA === true) {
    Params.AppName = '她理财';
    Params.IconUrl = 'http://file.market.xiaomi.com/thumbnail/PNG/l114/AppStore/06dc64e1221c8a9ebc7cd8d96b228ecc517439b95';
    Params.Pkg = 'com.incn.yida';
    Params.BundleID = 'com.incn.yida';
    Params.AppID = 'appID1';
    Params.Url = 'https://itunes.apple.com/cn/app/yi-da/id994848419?mt=8';
    Params.Match_id = '123';
    Params.Download_msg = '';
    Params.Download_title = '';
    Params.Download_btn_text = '';
    Params.Chrome_major = '30';
    Params.Ios_major = '8';
    Params.Redirect_url = 'http://www.sohu.com';
    Params.YYB_url = 'http://a.app.qq.com/o/simple.jsp?pkgname=com.xianguo.tingguo';
    //Params.YYB_url = '';
    Params.Scheme = 'dsEE26D0331DD04D54';
    Params.Host = 'open';
    Params.AppInsStatus = '2';
    Params.UserConf_Bg_WechatAndroidTip_url = 'http://img0.imgtn.bdimg.com/it/u=938096994,3074232342&fm=21&gp=0.jpg';
    Params.UserConf_Bg_WechatIosTip_url = '';

    Params.isAndroid = function () {
        return false;
    };
    Params.isYYBEnableAndroid = function (){
        return false;
    };
    Params.isYYBEnableIosBelow9 = function (){
        return false;
    };
    Params.isIOS = function () {
        return true;
    };
    Params.isWechat = function () {
        return true;
    };
    Params.isQQ = function () {
        return false;
    };
    Params.isWeibo = function () {
        return false;
    };
    Params.isFirefox = function () {
        return false;
    };
    Params.isChrome = function () {
        return true;
    };
    Params.isUniversallink = function () {
        return true;
    };
    Params.isDownloadDirectly = function () {
        return false;
    };
    Params.isCannotDeeplink = function () {
        return false;
    };
    Params.isCannotGetWinEvent= function () {
        return false;
    };
    Params.isUC = function () {
        return true;
    };
}
