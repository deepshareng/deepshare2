var needwait = '';
var expectDst = '';
var expectDeeplink = '';
var dstLocation = '';
var deeplinkLocation = '';
var DEBUG = false;
var Params = {
    AppName: '',
    IconUrl: '',
    Pkg: '',
    BundleID: '',
    AppID: '',
    Url: '',
    Match_id: '',
    Download_msg: '',
    Download_title: '',
    Download_btn_text: '',
    Chrome_major: 0,
    Ios_major: 0,
    Redirect_url: '',
    YYB_url: '',
    Scheme: '',
    Host: '',
    DsTag: '',
    UserConf_Bg_WechatAndroidTip_url: '',
    UserConf_Bg_WechatIosTip_url: '',
    AppInsStatus: 0,
    isAndroid: function () {
        return false;
    },
    isIOS: function () {
        return false;
    },
    isWechat: function () {
        return false;
    },
    isQQ: function () {
        return false;
    },
    isQQBrowser: function () {
        return false;
    },
    isWeibo: function () {
        return false;
    },
    isTwitter: function () {
        return false;
    },
    isFacebook: function () {
        return false;
    },
    isFirefox: function () {
        return false;
    },
    isChrome: function () {
        return false;
    },
    isUniversallink: function () {
        return false;
    },
    isDownloadDirectly: function () {
        return false;
    },
    isCannotDeeplink: function () {
        return false;
    },
    isCannotGetWinEvent: function () {
        return false;
    },
    isCannotGoMarket: function () {
        return false;
    },
    isForceUseScheme: function () {
        return false;
    },
    isUC: function () {
        return false;
    },
    isYYBEnableAndroid: function () {
        return false;
    },
    isYYBEnableIosBelow9: function () {
        return false;
    },
    isYYBEnableIosAbove9: function () {
        return false;
    }
};

describe("tests for deepshare-redirect.js :", function () {
    beforeEach(function () {
        needwait = false;
        expectDst = '';
        expectDeeplink = '';
        dstLocation = '';
        deeplinkLocation = '';
        DEBUG = false;
        Params = {
            AppName: '她理财',
            IconUrl: 'http://file.market.xiaomi.com/thumbnail/PNG/l114/AppStore/06dc64e1221c8a9ebc7cd8d96b228ecc517439b95',
            Pkg: 'pkg',
            BundleID: 'bundleid',
            AppID: 'appid1',
            Url: '',
            Match_id: '',
            Download_msg: '',
            Download_title: '',
            Download_btn_text: '',
            Chrome_major: 0,
            Ios_major: 0,
            Redirect_url: '',
            YYB_url: '',
            Scheme: '',
            Host: '',
            DsTag: '',
            AppInsStatus: 0,
            UserConf_Bg_WechatAndroidTip_url: '',
            UserConf_Bg_WechatIosTip_url: '',
            isAndroid: function () {
                return false;
            },
            isIOS: function () {
                return false;
            },
            isWechat: function () {
                return false;
            },
            isQQ: function () {
                return false;
            },
            isQQBrowser: function () {
                return false;
            },
            isWeibo: function () {
                return false;
            },
            isTwitter: function () {
                return false;
            },
            isFacebook: function () {
                return false;
            },
            isFirefox: function () {
                return false;
            },
            isChrome: function () {
                return false;
            },
            isUniversallink: function () {
                return false;
            },
            isDownloadDirectly: function () {
                return false;
            },
            isCannotDeeplink: function () {
                return false;
            },
            isCannotGetWinEvent: function () {
                return false;
            },
            isCannotGoMarket: function () {
                return false;
            },
            isForceUseScheme: function () {
                return false;
            },
            isUC: function () {
                return false;
            },
            isYYBEnableAndroid: function () {
                return false;
            },
            isYYBEnableIosBelow9: function () {
                return false;
            },
            isYYBEnableIosAbove9: function () {
                return false;
            }
        };
        jasmine.clock().install();
        spyOn(env, "windowLocation").and.callFake(function () {
        });
        spyOn(env, "windowOpen").and.callFake(function () {
        });
        spyOn(env, "windowClose").and.callFake(function () {
        });
        spyOn(env, "windowChangeHistory").and.callFake(function () {
        });
        spyOn(env, "windowAddEventListener").and.callFake(function () {
        });
        spyOn(env, "windowUrlAddTag").and.callFake(function () {
        });
    });

    afterEach(function () {
        jasmine.clock().uninstall();
    });

    it("testWechatAndroid", function () {
        expectDst = dsAction.destination.dstweixintipandroid;
        expectDeeplink = '';

        Params.isAndroid = function () {
            return true;
        };

        Params.isWechat = function () {
            return true;
        };

        start();

        expect(dstLocation).toEqual(expectDst);
        expect(deeplinkLocation).toEqual(expectDeeplink);
    });

    it("testWeichatIos", function () {
        expectDst = dsAction.destination.dstweixintipios;
        expectDeeplink = '';

        Params.isIOS = function () {
            return true;
        };
        Params.isWechat = function () {
            return true;
        };

        start();

        expect(dstLocation).toEqual(expectDst);
        expect(deeplinkLocation).toEqual(expectDeeplink);

    });

    it("testWechatYYB", function () {
        expectDst = dsAction.destination.dstweixintipios;
        expectDeeplink = '';

        Params.YYB_url = 'http://xxx.com';
        Params.isIOS = function () {
            return true;
        };
        Params.isWechat = function () {
            return true;
        };

        start();

        expect(dstLocation).toEqual(expectDst);
        expect(deeplinkLocation).toEqual(expectDeeplink);
    });

    it("testWechatYYBIos8", function () {
        expectDst = 'http://xxx.com';
        expectDeeplink = '';

        Params.YYB_url = 'http://xxx.com';
        Params.Ios_major = '8';
        Params.isIOS = function () {
            return true;
        };
        Params.isWechat = function () {
            return true;
        };
        Params.isYYBEnableIosBelow9 = function () {
            return true;
        };

        start();

        expect(dstLocation).toEqual(expectDst);
        expect(deeplinkLocation).toEqual(expectDeeplink);
    });

    it("testWechatYYBIos9 enbale", function () {
        expectDst = 'http://xxx.com';
        expectDeeplink = '';

        Params.YYB_url = 'http://xxx.com';
        Params.Ios_major = '9';
        Params.isIOS = function () {
            return true;
        };
        Params.isWechat = function () {
            return true;
        };
        Params.AppInsStatus = 0;
        Params.Uninstall_url = "";
        Params.isYYBEnableIosAbove9 = function () {
            return true;
        };

        start();

        expect(expectDst).toEqual(dstLocation);
        expect(deeplinkLocation).toEqual(expectDeeplink);
    });

    it("testWechatYYBIos9", function () {
        expectDst = dsAction.destination.dstweixintipios;
        expectDeeplink = '';

        Params.YYB_url = 'http://xxx.com';
        Params.Ios_major = '9';
        Params.isIOS = function () {
            return true;
        };
        Params.isWechat = function () {
            return true;
        };
        Params.AppInsStatus = 1;

        start();

        expect(dstLocation).toEqual(expectDst);
        expect(deeplinkLocation).toEqual(expectDeeplink);
    });

    it("testQQYYB", function () {
        expectDst = dsAction.destination.dstqqtipios;
        expectDeeplink = '';

        Params.YYB_url = 'http://xxx.com';
        Params.isIOS = function () {
            return true;
        };
        Params.isQQ = function () {
            return true;
        };

        start();

        expect(dstLocation).toEqual(expectDst);
        expect(deeplinkLocation).toEqual(expectDeeplink);
    });

    it("testQQBrowserYYB", function () {
        expectDst = 'http://xxx.com';
        expectDeeplink = '';

        Params.YYB_url = 'http://xxx.com';
        Params.isQQBrowser = function () {
            return true;
        };
        Params.isAndroid = function () {
            return true;
        };
        Params.isYYBEnableAndroid = function () {
            return true;
        };
        start();

        expect(dstLocation).toEqual(expectDst);
        expect(deeplinkLocation).toEqual(expectDeeplink);
    });

    it("testQQBrowserNoYYB", function () {
        expectDst = 'dst-cannot-deeplink';
        expectDeeplink = '';

        Params.isQQBrowser = function () {
            return true;
        };
        Params.isAndroid = function () {
            return true;
        };
        Params.isYYBEnableAndroid = function () {
            return false;
        };
        start();

        expect(dstLocation).toEqual(expectDst);
        expect(deeplinkLocation).toEqual(expectDeeplink);
    });

    it("testQQAndroid", function () {
        expectDst = dsAction.destination.dstqqtipandroid;
        expectDeeplink = '';

        Params.isAndroid = function () {
            return true;
        };

        Params.isQQ = function () {
            return true;
        };

        start();

        expect(dstLocation).toEqual(expectDst);
        expect(deeplinkLocation).toEqual(expectDeeplink);
    });

    it("testQQIos", function () {
        expectDst = dsAction.destination.dstqqtipios;
        expectDeeplink = '';

        Params.isIOS = function () {
            return true;
        };
        Params.isQQ = function () {
            return true;
        };

        start();

        expect(dstLocation).toEqual(expectDst);
        expect(deeplinkLocation).toEqual(expectDeeplink);

    });

    it("testWeiboIos", function () {
        expectDst = dsAction.destination.dstweibotipios;
        expectDeeplink = '';

        Params.isIOS = function () {
            return true;
        };
        Params.isWeibo = function () {
            return true;
        };

        start();

        expect(dstLocation).toEqual(expectDst);
        expect(deeplinkLocation).toEqual(expectDeeplink);
    });
    it("testWeiboAndroid", function () {
        expectDst = dsAction.destination.dstweibotipandroid;
        expectDeeplink = '';

        Params.isAndroid = function () {
            return true;
        };
        Params.isWeibo = function () {
            return true;
        };

        start();

        expect(dstLocation).toEqual(expectDst);
        expect(deeplinkLocation).toEqual(expectDeeplink);
    });

    it("testIos", function () {
        needwait = true;
        expectDst = 'https://itunes.apple.com/cn/app/yi-da/id994848419?mt=8';
        expectDeeplink = 'yeeda://?click_id=123';

        Params.Url = 'https://itunes.apple.com/cn/app/yi-da/id994848419?mt=8';
        Params.Match_id = '123';
        Params.Scheme = 'yeeda';

        Params.isIOS = function () {
            return true;
        };

        start();
        jasmine.clock().tick(2001);
        expect(expectDst).toEqual(dstLocation);
        expect(expectDeeplink).toEqual(deeplinkLocation);
    });

    it("testIos9UniversallinkNoCookie", function () {
        expectDst = dsAction.destination.dstios9UniversalLinkLandPage;
        expectDeeplink = '';

        Params.Ios_major = '9';

        Params.isIOS = function () {
            return true;
        };
        Params.isUniversallink = function () {
            return true;
        };
        spyOn(env, "cookieEnabled").and.callFake(function () {
            return false;
        });
        start();

        expect(expectDst).toEqual(dstLocation);
        expect(expectDeeplink).toEqual(deeplinkLocation);
    });

    it("testIos9UniversallinkCookieEnableNoInstall", function () {
        expectDst = 'http://xxx.com';
        expectDeeplink = '';

        Params.Ios_major = '9';
        Params.Url = 'http://xxx.com';

        Params.isIOS = function () {
            return true;
        };
        Params.isUniversallink = function () {
            return true;
        };
        spyOn(env, "cookieEnabled").and.callFake(function () {
            return true;
        });
        AppInsStatus = 0;
        start();

        expect(expectDst).toEqual(dstLocation);
        expect(expectDeeplink).toEqual(deeplinkLocation);
    });

    it("testIos9UniversallinkCookieEnableInstalled", function () {
        expectDst = dsAction.destination.dstios9UniversalLinkLandPage;
        expectDeeplink = '';

        Params.Ios_major = '9';
        Params.Url = 'http://xxx.com';

        Params.isIOS = function () {
            return true;
        };
        Params.isUniversallink = function () {
            return true;
        };
        spyOn(env, "cookieEnabled").and.callFake(function () {
            return true;
        });
        Params.AppInsStatus = 1;
        start();

        expect(expectDst).toEqual(dstLocation);
        expect(expectDeeplink).toEqual(deeplinkLocation);
    });

    it("testIos9UniversallinkCookieEnableInstallNoClear", function () {
        expectDst = dsAction.destination.dstios9UniversalLinkLandPage;
        expectDeeplink = '';

        Params.Ios_major = '9';
        Params.Url = 'http://xxx.com';

        Params.isIOS = function () {
            return true;
        };
        Params.isUniversallink = function () {
            return true;
        };
        spyOn(env, "cookieEnabled").and.callFake(function () {
            return true;
        });
        Params.AppInsStatus = 2;
        start();

        expect(expectDst).toEqual(dstLocation);
        expect(expectDeeplink).toEqual(deeplinkLocation);
    });

    it("testIos9Chrome", function () {
        expectDst = 'http://xxx.com';
        expectDeeplink = 'yeeda://?click_id=123';

        Params.Match_id = '123';
        Params.Scheme = 'yeeda';
        Params.Host = 'xxx';
        Params.Ios_major = '9';
        Params.Url = 'http://xxx.com';
        Params.isIOS = function () {
            return true;
        };
        Params.isChrome = function () {
            return true;
        };

        start();
        expect(expectDst).toEqual(dstLocation);
        expect(expectDeeplink).toEqual(deeplinkLocation);
    });

    it("testIos9Safari", function () {
        expectDst = 'http://xxx.com';
        expectDeeplink = 'yeeda://?click_id=123';

        Params.Match_id = '123';
        Params.Scheme = 'yeeda';
        Params.Host = 'xxx';
        Params.Ios_major = '9';
        Params.Url = 'http://xxx.com';
        Params.isIOS = function () {
            return true;
        };

        start();
        jasmine.clock().tick(501);
        expect(expectDst).toEqual(dstLocation);
        expect(expectDeeplink).toEqual(deeplinkLocation);
    });

    it("testAndroidFirefox", function () {
        expectDst = 'market://details?id=com.yeeda';
        expectDeeplink = 'deepshare://com.yeeda?click_id=123';

        Params.Match_id = '123';
        Params.Scheme = 'deepshare';
        Params.Host = 'com.yeeda';
        Params.Pkg = 'com.yeeda';
        Params.isAndroid = function () {
            return true;
        };
        Params.isFirefox = function () {
            return true;
        };

        start();
        jasmine.clock().tick(2001);
        expect(expectDst).toEqual(dstLocation);
        expect(expectDeeplink).toEqual(deeplinkLocation);
    });

    it("testAndroidChrome24", function () {
        expectDst = 'market://details?id=com.yeeda';
        expectDeeplink = 'deepshare://com.yeeda?click_id=123';

        Params.Match_id = '123';
        Params.Scheme = 'deepshare';
        Params.Host = 'com.yeeda';
        Params.Pkg = 'com.yeeda';
        Params.Chrome_major = 24;
        Params.isAndroid = function () {
            return true;
        };
        Params.isChrome = function () {
            return true;
        };
        start();
        jasmine.clock().tick(2001);
        expect(expectDst).toEqual(dstLocation);
        expect(expectDeeplink).toEqual(deeplinkLocation);
    });

    it("testAndroidChrome36", function () {
        expectDst = 'market://details?id=com.yeeda';
        expectDeeplink = 'intent://com.yeeda?click_id=123#Intent;scheme=deepshare;package=com.yeeda;S.browser_fallback_url=http://xxx.com;end';

        Params.Match_id = '123';
        Params.Scheme = 'deepshare';
        Params.Host = 'com.yeeda';
        Params.Pkg = 'com.yeeda';
        Params.Url = 'http://xxx.com';
        Params.Chrome_major = 36;
        Params.isAndroid = function () {
            return true;
        };
        Params.isChrome = function () {
            return true;
        };
        start();
        jasmine.clock().tick(2001);
        expect(expectDst).toEqual(dstLocation);
        expect(expectDeeplink).toEqual(deeplinkLocation);
    });

    it("testDefaultBrowser", function () {
        expectDst = 'market://details?id=com.yeeda';
        expectDeeplink = 'deepshare://com.yeeda?click_id=123';

        Params.Match_id = '123';
        Params.Scheme = 'deepshare';
        Params.Host = 'com.yeeda';
        Params.Pkg = 'com.yeeda';
        Params.isAndroid = function () {
            return true;
        };
        start();
        jasmine.clock().tick(2001);
        expect(expectDst).toEqual(dstLocation);
        expect(expectDeeplink).toEqual(deeplinkLocation);
    });

    it("testDefaultBrowserYYB", function () {
        expectDst = 'market://details?id=com.yeeda';
        expectDeeplink = 'deepshare://com.yeeda?click_id=123';

        Params.Match_id = '123';
        Params.Scheme = 'deepshare';
        Params.Host = 'com.yeeda';
        Params.Pkg = 'com.yeeda';
        Params.YYB_url = 'http://xxx.com';
        Params.isAndroid = function () {
            return true;
        };
        start();
        jasmine.clock().tick(2001);
        expect(expectDst).toEqual(dstLocation);
        expect(expectDeeplink).toEqual(deeplinkLocation);
    });

    it("testAndroidDirectDownload", function () {
        expectDst = dsAction.destination.dstandroidDirectDownloadLandPage;
        expectDeeplink = 'deepshare://com.yeeda?click_id=123';

        Params.Match_id = '123';
        Params.Scheme = 'deepshare';
        Params.Host = 'com.yeeda';
        Params.isAndroid = function () {
            return true;
        };
        Params.isDownloadDirectly = function () {
            return true;
        };
        start();
        jasmine.clock().tick(2001);
        expect(expectDst).toEqual(dstLocation);
        expect(expectDeeplink).toEqual(deeplinkLocation);
    });

    it("testCannotDeeplink", function () {
        expectDst = dsAction.destination.dstcannotdeeplink;
        expectDeeplink = 'deepshare://com.yeeda?click_id=123';

        Params.Match_id = '123';
        Params.Scheme = 'deepshare';
        Params.Host = 'com.yeeda';
        Params.isAndroid = function () {
            return true;
        };
        Params.isCannotDeeplink = function () {
            return true;
        };
        start();
        jasmine.clock().tick(2001);
        expect(expectDst).toEqual(dstLocation);
        expect(expectDeeplink).toEqual(deeplinkLocation);
    });

    it("testIsUc", function () {
        expectDst1 = dsAction.destination.dstucbrowser;
        expectDst2 = dsAction.destination.dstandroidMarketLandPage;
        expectDeeplink = 'deepshare://com.yeeda?click_id=123';

        Params.Match_id = '123';
        Params.Scheme = 'deepshare';
        Params.Host = 'com.yeeda';
        Params.Pkg = 'com.yeeda';
        Params.isAndroid = function () {
            return true;
        };
        Params.isUC = function () {
            return true;
        };
        start();
        expect(expectDst1).toEqual(dstLocation);
        jasmine.clock().tick(6001);
        expect(expectDeeplink).toEqual(deeplinkLocation);
    });

    it("testCannotGoMarket", function () {
        expectDst = dsAction.destination.dstandroidCannotGoMarketLandPage;
        expectDeeplink = 'deepshare://com.yeeda?click_id=123';

        Params.Match_id = '123';
        Params.Scheme = 'deepshare';
        Params.Host = 'com.yeeda';
        Params.isAndroid = function () {
            return true;
        };
        Params.isCannotGoMarket = function () {
            return true;
        };
        start();
        jasmine.clock().tick(2001);
        expect(expectDst).toEqual(dstLocation);
        expect(expectDeeplink).toEqual(deeplinkLocation);
    });

    it("testForceUseScheme", function () {
        expectDst = 'market://details?id=com.yeeda';
        expectDeeplink = 'deepshare://com.yeeda?click_id=123';

        Params.Match_id = '123';
        Params.Scheme = 'deepshare';
        Params.Host = 'com.yeeda';
        Params.Pkg = 'com.yeeda';
        Params.Chrome_major = 36;
        Params.isAndroid = function () {
            return true;
        };
        Params.isChrome = function () {
            return true;
        };
        Params.isForceUseScheme = function () {
            return true;
        };
        start();
        jasmine.clock().tick(2001);
        expect(expectDst).toEqual(dstLocation);
        expect(expectDeeplink).toEqual(deeplinkLocation);
    });

    it("testCannotGoMarket", function () {
        expectDst = dsAction.destination.dstandroidCannotGoMarketLandPage;
        expectDeeplink = 'deepshare://com.yeeda?click_id=123';

        Params.Match_id = '123';
        Params.Scheme = 'deepshare';
        Params.Host = 'com.yeeda';
        Params.Chrome_major = 15;
        Params.isAndroid = function () {
            return true;
        };
        Params.isChrome = function () {
            return true;
        };
        Params.isCannotGoMarket = function () {
            return true;
        };
        start();
        jasmine.clock().tick(2001);
        expect(expectDst).toEqual(dstLocation);
        expect(expectDeeplink).toEqual(deeplinkLocation);
    });

    it("testIosNA", function () {
        expectDst = 'dst-ios-not-available'

        Params.BundleID = '';
        Params.isIOS = function () {
            return true;
        };

        start();
        expect(expectDst).toEqual(dstLocation);
    });

    it("testAndroidNA", function () {
        expectDst = 'dst-android-not-available'

        Params.Pkg = '';
        Params.isAndroid = function () {
            return true;
        };

        start();
        expect(expectDst).toEqual(dstLocation);
    });
});
