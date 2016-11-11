describe('deepshare', function() {
    beforeEach(function () {
        needwait = false;
        expectDst = '';
        expectDeeplink = '';
        dstLocation = '';
        deeplinkLocation = '';
        DEBUG = true;
        Params = {
            // In app data
            appId: 'appid1', 
            downloadMsg: '',
            downloadTitle: '',
            downloadBtnText: '',
            // Param
            app_name: '她理财',
            icon_url: 'http://file.market.xiaomi.com/thumbnail/PNG/l114/AppStore/06dc64e1221c8a9ebc7cd8d96b228ecc517439b95',
            pkg: 'pkg',
            bundle_id: 'bundleid',
            app_id: 'appid1',
            url: '',
            match_id: '',
            chrome_major: 0,
            ios_major: 0,
            redirect_url: '',
            yyb_url: '',
            scheme: '',
            host: '',
            ds_tag: '',
            app_ins_status: 0,
            is_android: false,
            is_ios: false,
            is_wechat: false,
            is_qq: false,
            is_qq_browser: false,
            is_weibo: false,
            is_firefox: false,
            is_universal_link: false,
            is_download_directly: false,
            cannot_deeplink: false,
            cannot_get_win_event: false,
            cannot_go_market: false,
            force_use_scheme: false,
            is_uc: false,
        };

        deepshare = new DeepShare('appid1');
        deepshare.BindParams(Params);

        jasmine.clock().install();
        spyOn(deepshare._env, "windowLocation").and.callFake(function () {
        });
        spyOn(deepshare._env, "windowOpen").and.callFake(function () {
        });
        spyOn(deepshare._env, "windowClose").and.callFake(function () {
        });
        spyOn(deepshare._env, "windowChangeHistory").and.callFake(function () {
        });
        spyOn(deepshare._env, "windowAddEventListener").and.callFake(function () {
        });
        spyOn(deepshare._env, "windowUrlAddTag").and.callFake(function () {
        });
    });

    afterEach(function () {
        jasmine.clock().uninstall();
    });


    it("should go to weixin tip in android", function () {
        expectDst = deepshare._DSAction.destination.dstweixintipandroid;
        expectDeeplink = '';

        deepshare._Params.is_android = true;
        deepshare._Params.is_wechat = true;

        deepshare.Start();

        //expect(deepshare._dstLocation).toEqual(expectDst);
        expect(deepshare._dstLocation).toEqual(expectDst);
        expect(deepshare._deeplinkLocation).toEqual(expectDeeplink);
    });

    it("should go to weixi tip in ios", function () {
        expectDst = deepshare._DSAction.destination.dstweixintipios;
        expectDeeplink = '';

        deepshare._Params.is_ios = true;
        deepshare._Params.is_wechat = true;

        deepshare.Start();

        expect(deepshare._dstLocation).toEqual(expectDst);
        expect(deepshare._deeplinkLocation).toEqual(expectDeeplink);
    });

    it("should go to weixin tip in ios with yyb url", function () {
        expectDst = deepshare._DSAction.destination.dstweixintipios;
        expectDeeplink = '';

        deepshare._Params.yyb_url = 'http://xxx.com';
        deepshare._Params.is_ios = true;
        deepshare._Params.is_wechat = true;

        deepshare.Start();

        expect(deepshare._dstLocation).toEqual(expectDst);
        expect(deepshare._deeplinkLocation).toEqual(expectDeeplink);
    });

    it("should go to weixin tip in ios9 with yyb url", function () {
        expectDst = deepshare._DSAction.destination.dstweixintipios;
        expectDeeplink = '';

        deepshare._Params.yyb_url = 'http://xxx.com';
        deepshare._Params.ios_major = 9;
        deepshare._Params.is_ios = true;
        deepshare._Params.is_wechat = true;

        deepshare.Start();

        expect(deepshare._dstLocation).toEqual(expectDst);
        expect(deepshare._deeplinkLocation).toEqual(expectDeeplink);
    });

    it("should go to qq tip in ios with yyb url", function () {
        expectDst = deepshare._DSAction.destination.dstqqtipios;
        expectDeeplink = '';

        deepshare._Params.yyb_url = 'http://xxx.com';
        deepshare._Params.is_ios = true;
        deepshare._Params.is_qq = true;

        deepshare.Start();

        expect(deepshare._dstLocation).toEqual(expectDst);
        expect(deepshare._deeplinkLocation).toEqual(expectDeeplink);
    });

    it("should go to qq tip in android without yyb url", function () {
        expectDst = deepshare._DSAction.destination.dstqqtipandroid;
        expectDeeplink = '';

        deepshare._Params.is_android = true;

        deepshare._Params.is_qq = true;

        deepshare.Start();

        expect(deepshare._dstLocation).toEqual(expectDst);
        expect(deepshare._deeplinkLocation).toEqual(expectDeeplink);
    });

    it("should go to yyb url in qq browser with yyb url", function () {
        expectDst = 'http://xxx.com';
        expectDeeplink = '';

        deepshare._Params.yyb_url = 'http://xxx.com';
        deepshare._Params.is_qq_browser = true;
        deepshare._Params.is_android = true;
        deepshare.Start();

        expect(deepshare._dstLocation).toEqual(expectDst);
        expect(deepshare._deeplinkLocation).toEqual(expectDeeplink);
    });

    it("should demostrate cannot deeplink in qqbrowser without yyb url", function () {
        expectDst = deepshare._DSAction.destination.dstcannotdeeplink;
        expectDeeplink = '';

        deepshare._Params.is_qq_browser = true;
        deepshare._Params.is_android = true;
        deepshare.Start();

        expect(deepshare._dstLocation).toEqual(expectDst);
        expect(deepshare._deeplinkLocation).toEqual(expectDeeplink);
    });

    it("testQQIos", function () {
        expectDst = deepshare._DSAction.destination.dstqqtipios;
        expectDeeplink = '';

        deepshare._Params.is_ios = true;
        deepshare._Params.is_qq = true;

        deepshare.Start();

        expect(deepshare._dstLocation).toEqual(expectDst);
        expect(deepshare._deeplinkLocation).toEqual(expectDeeplink);

    });

    it("testWeiboIos", function () {
        expectDst = deepshare._DSAction.destination.dstweibotipios;
        expectDeeplink = '';

        deepshare._Params.is_ios = true;
        deepshare._Params.is_weibo = true;

        deepshare.Start();

        expect(deepshare._dstLocation).toEqual(expectDst);
        expect(deepshare._deeplinkLocation).toEqual(expectDeeplink);
    });
    it("testWeiboAndroid", function () {
        expectDst = deepshare._DSAction.destination.dstweibotipandroid;
        expectDeeplink = '';

        deepshare._Params.is_android = true;
        deepshare._Params.is_weibo = true;

        deepshare.Start();

        expect(deepshare._dstLocation).toEqual(expectDst);
        expect(deepshare._deeplinkLocation).toEqual(expectDeeplink);
    });

    it("testIos", function () {
        needwait = true;
        expectDst = 'https://itunes.apple.com/cn/app/yi-da/id994848419?mt=8';
        expectDeeplink = 'yeeda://?click_id=123';

        deepshare._Params.url = 'https://itunes.apple.com/cn/app/yi-da/id994848419?mt=8';
        deepshare._Params.match_id = '123';
        deepshare._Params.scheme = 'yeeda';

        deepshare._Params.is_ios = true;

        deepshare.Start();
        jasmine.clock().tick(2001);
        expect(expectDst).toEqual(deepshare._dstLocation);
        expect(expectDeeplink).toEqual(deepshare._deeplinkLocation);
    });

    it("testIos9UniversallinkNoCookie", function () {
        expectDst = deepshare._DSAction.destination.dstios9UniversalLinkLandPage;
        expectDeeplink = '';

        deepshare._Params.ios_major = 9;

        deepshare._Params.is_ios = true;
        deepshare._Params.is_universal_link = true;
        spyOn(deepshare._env, "cookieEnabled").and.callFake(function () {
            return false;
        });
        deepshare.Start();

        expect(expectDst).toEqual(deepshare._dstLocation);
        expect(expectDeeplink).toEqual(deepshare._deeplinkLocation);
    });

    it("testIos9UniversallinkCookieEnableNoInstall", function () {
        expectDst = 'http://xxx.com';
        expectDeeplink = '';

        deepshare._Params.ios_major = 9;
        deepshare._Params.url = 'http://xxx.com';

        deepshare._Params.is_ios = true;
        deepshare._Params.is_universal_link = true;
        spyOn(deepshare._env, "cookieEnabled").and.callFake(function () {
            return true;
        });
        app_ins_status = 0;
        deepshare.Start();

        expect(expectDst).toEqual(deepshare._dstLocation);
        expect(expectDeeplink).toEqual(deepshare._deeplinkLocation);
    });

    it("testIos9UniversallinkCookieEnableInstalled", function () {
        expectDst = deepshare._DSAction.destination.dstios9UniversalLinkLandPage;
        expectDeeplink = '';

        deepshare._Params.ios_major = 9;
        deepshare._Params.url = 'http://xxx.com';

        deepshare._Params.is_ios = true;
        deepshare._Params.is_universal_link = true;
        spyOn(deepshare._env, "cookieEnabled").and.callFake(function () {
            return true;
        });
        deepshare._Params.app_ins_status = 1;
        deepshare.Start();

        expect(expectDst).toEqual(deepshare._dstLocation);
        expect(expectDeeplink).toEqual(deepshare._deeplinkLocation);
    });

    it("testIos9UniversallinkCookieEnableInstallNoClear", function () {
        expectDst = deepshare._DSAction.destination.dstios9UniversalLinkLandPage;
        expectDeeplink = '';

        deepshare._Params.ios_major = 9;
        deepshare._Params.url = 'http://xxx.com';

        deepshare._Params.is_ios = true;
        deepshare._Params.is_universal_link = true;
        spyOn(deepshare._env, "cookieEnabled").and.callFake(function () {
            return true;
        });
        deepshare._Params.app_ins_status = 2;
        deepshare.Start();

        expect(expectDst).toEqual(deepshare._dstLocation);
        expect(expectDeeplink).toEqual(deepshare._deeplinkLocation);
    });

    it("testIos9Chrome", function () {
        expectDst = 'http://xxx.com';
        expectDeeplink = 'yeeda://?click_id=123';

        deepshare._Params.match_id = '123';
        deepshare._Params.scheme = 'yeeda';
        deepshare._Params.host = 'xxx';
        deepshare._Params.ios_major = 9;
        deepshare._Params.url = 'http://xxx.com';
        deepshare._Params.is_ios = true;
        deepshare._Params.is_chrome = true;
        deepshare._Params.chrome_major = 20;

        deepshare.Start();
        expect(expectDst).toEqual(deepshare._dstLocation);
        expect(expectDeeplink).toEqual(deepshare._deeplinkLocation);
    });

    it("testIos9Safari", function () {
        expectDst = 'http://xxx.com';
        expectDeeplink = 'yeeda://?click_id=123';

        deepshare._Params.match_id = '123';
        deepshare._Params.scheme = 'yeeda';
        deepshare._Params.host = 'xxx';
        deepshare._Params.ios_major = 9;
        deepshare._Params.url = 'http://xxx.com';
        deepshare._Params.is_ios = true;

        deepshare.Start();
        jasmine.clock().tick(501);
        expect(expectDst).toEqual(deepshare._dstLocation);
        expect(expectDeeplink).toEqual(deepshare._deeplinkLocation);
    });

    it("testAndroidFirefox", function () {
        expectDst = 'market://details?id=com.yeeda';
        expectDeeplink = 'deepshare://com.yeeda?click_id=123';

        deepshare._Params.match_id = '123';
        deepshare._Params.scheme = 'deepshare';
        deepshare._Params.host = 'com.yeeda';
        deepshare._Params.pkg = 'com.yeeda';
        deepshare._Params.is_android = true;
        deepshare._Params.is_firefox = true;

        deepshare.Start();
        jasmine.clock().tick(2001);
        expect(expectDst).toEqual(deepshare._dstLocation);
        expect(expectDeeplink).toEqual(deepshare._deeplinkLocation);
    });

    it("testAndroidChrome24", function () {
        expectDst = 'market://details?id=com.yeeda';
        expectDeeplink = 'deepshare://com.yeeda?click_id=123';

        deepshare._Params.match_id = '123';
        deepshare._Params.scheme = 'deepshare';
        deepshare._Params.host = 'com.yeeda';
        deepshare._Params.pkg = 'com.yeeda';
        deepshare._Params.chrome_major = 24;
        deepshare._Params.is_android = true;

        deepshare.Start();

        jasmine.clock().tick(2001);
        expect(expectDst).toEqual(deepshare._dstLocation);
        expect(expectDeeplink).toEqual(deepshare._deeplinkLocation);
    });

    it("testAndroidChrome36", function () {
        expectDst = 'market://details?id=com.yeeda';
        expectDeeplink = 'intent://com.yeeda?click_id=123#Intent;scheme=deepshare;package=com.yeeda;S.browser_fallback_url=http://xxx.com;end';

        deepshare._Params.match_id = '123';
        deepshare._Params.scheme = 'deepshare';
        deepshare._Params.host = 'com.yeeda';
        deepshare._Params.pkg = 'com.yeeda';
        deepshare._Params.url = 'http://xxx.com';
        deepshare._Params.chrome_major = 36;
        deepshare._Params.is_android = true;

        deepshare.Start();

        jasmine.clock().tick(2001);
        expect(expectDst).toEqual(deepshare._dstLocation);
        expect(expectDeeplink).toEqual(deepshare._deeplinkLocation);
    });

    it("testDefaultBrowser", function () {
        expectDst = 'market://details?id=com.yeeda';
        expectDeeplink = 'deepshare://com.yeeda?click_id=123';

        deepshare._Params.match_id = '123';
        deepshare._Params.scheme = 'deepshare';
        deepshare._Params.host = 'com.yeeda';
        deepshare._Params.pkg = 'com.yeeda';
        deepshare._Params.is_android = true;

        deepshare.Start();

        jasmine.clock().tick(2001);
        expect(expectDst).toEqual(deepshare._dstLocation);
        expect(expectDeeplink).toEqual(deepshare._deeplinkLocation);
    });

    it("testDefaultBrowserYYB", function () {
        expectDst = 'market://details?id=com.yeeda';
        expectDeeplink = 'deepshare://com.yeeda?click_id=123';

        deepshare._Params.match_id = '123';
        deepshare._Params.scheme = 'deepshare';
        deepshare._Params.host = 'com.yeeda';
        deepshare._Params.pkg = 'com.yeeda';
        deepshare._Params.yyb_url = 'http://xxx.com';
        deepshare._Params.is_android = true;

        deepshare.Start();

        jasmine.clock().tick(2001);
        expect(expectDst).toEqual(deepshare._dstLocation);
        expect(expectDeeplink).toEqual(deepshare._deeplinkLocation);
    });

    it("testAndroidDirectDownload", function () {
        expectDst = deepshare._DSAction.destination.dstandroidDirectDownloadLandPage;
        expectDeeplink = 'deepshare://com.yeeda?click_id=123';

        deepshare._Params.match_id = '123';
        deepshare._Params.scheme = 'deepshare';
        deepshare._Params.host = 'com.yeeda';
        deepshare._Params.is_android = true;
        deepshare._Params.is_download_directly = true;

        deepshare.Start();

        jasmine.clock().tick(2001);
        expect(expectDst).toEqual(deepshare._dstLocation);
        expect(expectDeeplink).toEqual(deepshare._deeplinkLocation);
    });

    it("testCannotDeeplink", function () {
        expectDst = deepshare._DSAction.destination.dstcannotdeeplink;
        expectDeeplink = 'deepshare://com.yeeda?click_id=123';

        deepshare._Params.match_id = '123';
        deepshare._Params.scheme = 'deepshare';
        deepshare._Params.host = 'com.yeeda';
        deepshare._Params.is_android = true;
        deepshare._Params.cannot_deeplink = true;

        deepshare.Start();

        jasmine.clock().tick(2001);
        expect(expectDst).toEqual(deepshare._dstLocation);
        expect(expectDeeplink).toEqual(deepshare._deeplinkLocation);
    });

    /*
    it("testIsUc", function () {
        expectDst1 = deepshare._DSAction.destination.dstucbrowser;
        expectDst2 = deepshare._DSAction.destination.dstandroidMarketLandPage;
        expectDeeplink = 'deepshare://com.yeeda?click_id=123';

        deepshare._Params.match_id = '123';
        deepshare._Params.scheme = 'deepshare';
        deepshare._Params.host = 'com.yeeda';
        deepshare._Params.pkg = 'com.yeeda';
        deepshare._Params.is_android = function () {
            return true;
        };
        deepshare._Params.isUC = function () {
            return true;
        };
        deepshare.Start();
        expect(expectDst2).toEqual(deepshare._dstLocation);
        jasmine.clock().tick(6001);
        expect(expectDeeplink).toEqual(deepshare._deeplinkLocation);
    });
    */

    it("testCannotGoMarket", function () {
        expectDst = deepshare._DSAction.destination.dstandroidCannotGoMarketLandPage;
        expectDeeplink = 'deepshare://com.yeeda?click_id=123';

        deepshare._Params.match_id = '123';
        deepshare._Params.scheme = 'deepshare';
        deepshare._Params.host = 'com.yeeda';
        deepshare._Params.is_android = true;
        deepshare._Params.cannot_go_market = true;

        deepshare.Start();

        jasmine.clock().tick(2001);
        expect(expectDst).toEqual(deepshare._dstLocation);
        expect(expectDeeplink).toEqual(deepshare._deeplinkLocation);
    });

    it("testForceUseScheme", function () {
        expectDst = 'market://details?id=com.yeeda';
        expectDeeplink = 'deepshare://com.yeeda?click_id=123';

        deepshare._Params.match_id = '123';
        deepshare._Params.scheme = 'deepshare';
        deepshare._Params.host = 'com.yeeda';
        deepshare._Params.pkg = 'com.yeeda';
        deepshare._Params.chrome_major = 36;
        deepshare._Params.is_android = true;
        deepshare._Params.force_use_scheme = true;

        deepshare.Start();

        jasmine.clock().tick(2001);
        expect(expectDst).toEqual(deepshare._dstLocation);
        expect(expectDeeplink).toEqual(deepshare._deeplinkLocation);
    });

    it("testCannotGoMarket", function () {
        expectDst = deepshare._DSAction.destination.dstandroidCannotGoMarketLandPage;
        expectDeeplink = 'deepshare://com.yeeda?click_id=123';

        deepshare._Params.match_id = '123';
        deepshare._Params.scheme = 'deepshare';
        deepshare._Params.host = 'com.yeeda';
        deepshare._Params.chrome_major = 15;
        deepshare._Params.is_android = true;
        deepshare._Params.cannot_go_market = true;

        deepshare.Start();

        jasmine.clock().tick(2001);
        expect(expectDst).toEqual(deepshare._dstLocation);
        expect(expectDeeplink).toEqual(deepshare._deeplinkLocation);
    });

    it("testIosNA", function () {
        expectDst = 'dst-ios-not-available'

        deepshare._Params.bundle_id = '';
        deepshare._Params.is_ios = true;

        deepshare.Start();
        expect(expectDst).toEqual(deepshare._dstLocation);
    });

    it("testAndroidNA", function () {
        expectDst = 'dst-android-not-available'

        deepshare._Params.pkg = '';
        deepshare._Params.is_android = true;

        deepshare.Start();
        expect(expectDst).toEqual(deepshare._dstLocation);
    });
});
