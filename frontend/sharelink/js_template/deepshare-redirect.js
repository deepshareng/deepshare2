var winWidth = $(window).width();
var winHeight = $(window).height();

var CONST_DS_TAG = 'ds_tag';
var CONST_DSCOOKIE = 'dscookie';
var CONST_W_DSCOOKIE = 'wcookie';
var deeplinkLocation = '';
var dstLocation = '';
var ResPath = "../../jsserver/";
var ResPathImg = ResPath + "images/";

var LANDING_BG_TOP = 'https://ds-static.fds.so/ds-static/sharelink/landing/landing-bg-top.png';
var DRAGDOWN = 'https://ds-static.fds.so/ds-static/sharelink/landing/dragdown.png';

var dsAction = {
    trackingUrl: '/v2/dsactions/',
    actionJSDeepLink: 'js/deeplink',
    actionJSDst: 'js/dst',
    actionJSUserClick: 'js/userclick',

    destination: {
        dstweixintipandroid: 'dst-weixin-tip-android',
        dstweixintipios: 'dst-weixin-tip-ios',
        dstqqtipandroid: 'dst-qq-tip-android',
        dstqqtipios: 'dst-qq-tip-ios',
        dstweibotipandroid: 'dst-weibo-tip-android',
        dstweibotipios: 'dst-weibo-tip-ios',
        dstfacebooktipandroid: 'dst-facebook-tip-android',
        dstfacebooktipios: 'dst-facebook-tip-ios',
        dsttwittertipandroid: 'dst-twitter-tip-android',
        dsttwittertipios: 'dst-twitter-tip-ios',
        dstcannotdeeplink: 'dst-cannot-deeplink',
        dstucbrowser: 'dst-uc-browser',
        dstios9UniversalLinkLandPage: 'dst-ios9-universallink-landpage',
        dstandroidDirectDownloadLandPage: 'dst-android-direct-download-landpage',
        dstandroidMarketLandPage: 'dst-android-market-landpage',
        dstandroidCannotGoMarketLandPage: 'dst-android-cannot-gomarket-landpage',
        dstplatformNA: 'dst-{Platform}-not-available'
    },
    reportDSJSEvent: function (eventType, dst) {
        var params = {
            action: eventType,
            kvs: {
                "click_id": Params.Match_id,
                "destination": dst,
                "ds_tag": Params.DsTag
            }
        };
        var paramsJson = JSON.stringify(params);
        var url = this.trackingUrl + Params.AppID;
        $.ajax({
            url: url,
            type: 'POST',
            data: paramsJson
        });
    },

    reportDSJSUserClickEvent: function (eventType, btn, choice) {
        var params = {
            action: eventType,
            kvs: {
                "click_id": Params.Match_id,
                "user_btn": btn,
                "user_choice": choice,
                "ds_tag": Params.DsTag
            }
        };
        var paramsJson = JSON.stringify(params);
        var url = this.trackingUrl + Params.AppID;
        $.ajax({
            url: url,
            type: 'POST',
            data: paramsJson
        });
    }
};
var CONST_APP_INS_STATUS = {
    "NotInstall": 0,
    "Installed": 1,
    "Unclear": 2
};

var env = {
    windowLocation: function (loc) {
        window.location = loc;
    },
    windowOpen: function (loc) {
        window.open(loc);
    },
    windowClose: function () {
        window.close();
    },
    windowChangeHistory: function () {
        window.history.replaceState("Object", "Title", "0");
    },
    windowAddEventListener: function (type, listener) {
        window.addEventListener(type, listener);
    },
    windowUrlAddTag: function () {
        if (window.location.search.indexOf(CONST_DS_TAG) < 0) {
            var tag = Math.floor((Math.random() * 1000000));
            var wcookie = Params.Match_id;
            var queryStr = "?" + CONST_DS_TAG + "=" + tag + "&" + CONST_W_DSCOOKIE + "=" + wcookie;
            window.location.search = queryStr;
            if (DEBUG) {
                alert("Url Add Tag:" + queryStr);
            }
        }
    },
    cookieEnabled: function () {
        var isPrivateMode = false;
        try {
            localStorage.test = 2;
        } catch (e) {
            if (DEBUG) {
                alert("private mode");
            }
            isPrivateMode = true;
        }
        return navigator.cookieEnabled && !isPrivateMode;
    }
};

var renderTemplete = function (templeteId, domElementId, renderParams) {
    var params = renderParams || {};
    var $templeteEle = $('#' + templeteId);
    var $domEle = $('#' + domElementId);
    if ($templeteEle.length > 0 && $domEle.length > 0) {
        var template = $templeteEle.html();
        Mustache.parse(template, ['${', '}']);   // optional, speeds up future uses
        var rendered = Mustache.render(template, renderParams);
        $domEle.html(rendered);
        $domEle.removeClass('hide');
    }
};

var clearTimeoutOnPageUnload = function (redirectTimer) {
    env.windowAddEventListener("pagehide", function () {
        if (DEBUG) {
            alert('window event pagehide');
        }
        clearTimeout(redirectTimer);
        env.windowChangeHistory();
    });
    env.windowAddEventListener("blur", function () {
        if (DEBUG) {
            alert('window event blur');
        }
        clearTimeout(redirectTimer);
        env.windowChangeHistory();
    });
    env.windowAddEventListener("unload", function () {
        if (DEBUG) {
            alert('window event unload');
        }
        clearTimeout(redirectTimer);
        env.windowChangeHistory();
    });
    document.addEventListener("webkitvisibilitychange", function () {
        if (DEBUG) {
            alert('window event webkitvisibilitychange');
        }
        if (document.webkitHidden) {
            clearTimeout(redirectTimer);
            env.windowChangeHistory();
        }
    });
    env.windowAddEventListener("beforeunload", function () {
        if (DEBUG) {
            alert('window event beforeunload');
        }
//        clearTimeout(redirectTimer);
    });
    env.windowAddEventListener("focus", function () {
        if (DEBUG) {
            alert('window event focus');
        }
        //focus event is dangerous, it shows at least Firefox and Xiaomi Browser will receive a focus event when it try to deeplink.
        //So do not use this event to clear timer.
        //if (isFirefox()) {
        //    return;
        //}
        //clearTimeout(redirectTimer);
        //env.windowChangeHistory();
    });
    env.windowAddEventListener("focusout", function () {
        if (DEBUG) {
            alert('window event focusout');
        }
        clearTimeout(redirectTimer);
        env.windowChangeHistory();
    });
};


var gotoTip = function (type, dst, userImgInfo) {
    var dom = "";
    var templateId = '';
    var renderParams = {};
    if (userImgInfo !== ""){
        templateId = 'weixinTipUserConfigTemplate';
        renderParams = { 
            'userconfigBg': userImgInfo,
        };
    }else{
        var bgImgUrl, iconsImgUrl;
        var openbrowserStep1, openbrowserStep2, openbrowserStep3, openbrowserStep4;
        templateId = 'weixinTipTemplate';
        if (type == 'ios' || type == 'ios_right_down') {
            bgImgUrl = 'https://ds-static.fds.so/ds-static/sharelink/openbrowser_ios/bg.png';
            iconsImgUrl = 'https://ds-static.fds.so/ds-static/sharelink/openbrowser_ios/icons.png';
            if (type == 'ios_right_down') {
                openbrowserStep1 = ResiOSOpenbrowserStep1RightDown;
            } else {
                openbrowserStep1 = ResiOSOpenbrowserStep1;
            }
            openbrowserStep2 = ResiOSOpenbrowserStep2;
            openbrowserStep3 = ResiOSOpenbrowserStep3;
            openbrowserStep4 = ResiOSOpenbrowserStep4;
            
        }else{
            //type == android
            bgImgUrl = 'https://ds-static.fds.so/ds-static/sharelink/openbrowser_android/bg.png';
            iconsImgUrl = 'https://ds-static.fds.so/ds-static/sharelink/openbrowser_android/icons.png';
            openbrowserStep1 = ResAndroidOpenbrowserStep1;
            openbrowserStep2 = ResAndroidOpenbrowserStep2;
            openbrowserStep3 = ResAndroidOpenbrowserStep3;
            openbrowserStep4 = ResAndroidOpenbrowserStep4;
        }

        renderParams = { 
            'iconUrl':Params.IconUrl,
            'resOpenbrowserMsg': ResOpenbrowserMsg, 
            'bgImgUrl': bgImgUrl,
            'iconsImgUrl': iconsImgUrl,
            'openbrowserStep1': openbrowserStep1,
            'openbrowserStep2': openbrowserStep2,
            'openbrowserStep3': openbrowserStep3,
            'openbrowserStep4': openbrowserStep4,
            'downloadUrl': Params.Url,
        };
    }
    renderTemplete(templateId, 'gotoTip', renderParams);
    //$(".image-tip").css("height", winHeight);
    $(".image-tip").show();
    // loadCanvas("div_weixintip", IOS);
    env.windowUrlAddTag();
    dstLocation = dst;
    dsAction.reportDSJSEvent(dsAction.actionJSDst, dstLocation);
};

var gotoTencentProduct = function (dstIos, dstAndroid){
    if (Params.isIOS()) {
        if (DEBUG) {
            alert("isIOS");
        }
        if ((Params.Ios_major < 9 && Params.isYYBEnableIosBelow9()) ||
            (Params.Ios_major >= 9 && Params.isYYBEnableIosAbove9()) ) {
            if (DEBUG) {
                alert(Params.YYB_url)
            }
            gotoUrl(Params.YYB_url);
        }else{
            gotoTip('ios', dstIos, Params.UserConf_Bg_WechatIosTip_url);
        }
    } else if (Params.isAndroid()) {
        if (DEBUG) {
            alert("isAndroid");
        }
        if (Params.isYYBEnableAndroid()) {
            if (DEBUG) {
                alert(Params.YYB_url)
            }
            gotoUrl(Params.YYB_url);
        } else {
            gotoTip('android', dstAndroid, Params.UserConf_Bg_WechatAndroidTip_url);
        }
    }
};

var gotoCannotDeeplink = function () {
    if (DEBUG) {
        alert('cannot deeplink');
    }
    if (Params.isDownloadDirectly()) { 
        renderTemplete('gotoCannotDeeplinkWithDownloadBtnTemplate' , 'gotoCannotDeeplink', {
            'bgTopImgUrl': ResPathImg + 'oops/oops-bg-top.png',
            'iconUrl': ResPathImg + 'oops/oops-icon.png',
            'bgBottomImgUrl': ResPathImg + 'oops/oops-bg-jump.png',
            'resOopsMsg': ResOopsMsg,
            'resOopsTips': ResOopsTips,
            'resDownloadAPK':ResDownloadAPK,
        });

        $('#btnGotoAndroidDownload').click(function () {
            dsAction.reportDSJSUserClickEvent(dsAction.actionJSUserClick, "gotoAndroidDirectDownload", "yes");
            gotoUrl(Params.Url);
        });
    } else {
        renderTemplete('gotoCannotDeeplinkWithMarketBtnTemplate' , 'gotoCannotDeeplink', {
            'bgTopImgUrl': ResPathImg + 'oops/oops-bg-top.png',
            'iconUrl': ResPathImg + 'oops/oops-icon.png',
            'bgBottomImgUrl': ResPathImg + 'oops/oops-bg-jump.png',
            'resOopsMsg': ResOopsMsg,
            'resOopsTips': ResOopsTips,
            'resGotoAppStore':ResGotoAppStoreDownload,
        });

        $('#btnGotoAndroidMarket').click(function () {
            dsAction.reportDSJSUserClickEvent(dsAction.actionJSUserClick, "gotoAndroidMarket", "yes");
            gotoAndroidMarket();
        });
    }

    dstLocation = dsAction.destination.dstcannotdeeplink;
    dsAction.reportDSJSEvent(dsAction.actionJSDst, dstLocation);
};

var gotoAndroidNewInstall = function () {
    if (Params.isDownloadDirectly()) {
        gotoAndroidDownloadLandingPage();
    } else {
        if (Params.isCannotGoMarket()) {
            gotoAndroidCannotGoMarketLandingPage();
        } else if (Params.isCannotGetWinEvent() || Params.isUC()) {
            gotoAndroidMarketLandingPage();
        } else {
            gotoAndroidMarket();
        }
    }
};

var gotoIOSLandingPage = function () {
    dstLocation = dsAction.destination.dstios9UniversalLinkLandPage;
    if (DEBUG) {
        alert(dstLocation)
    }

    renderTemplete('gotoLandingpageTemplate', 'gotoLandingpage', {
        'bgTop': LANDING_BG_TOP,
        'appName':Params.AppName, 
        'iconUrl':Params.IconUrl, 
        'downloadTitle':Params.Download_title,
        'downloadMsg':Params.Download_msg, 
        'btnLandingpageText':ResGotoAppStoreDownload, 
        'elementType':'button',
        'dragdownIcon': DRAGDOWN,
        'dragdownTip': ResDragdownTip,
        'dragdown-display': '',
    });

    dsAction.reportDSJSEvent(dsAction.actionJSDst, dstLocation);
    $('#btnGotoLandingPage').click(function () {
        dsAction.reportDSJSUserClickEvent(dsAction.actionJSUserClick, "gotoIosAppStore", "yes");
        gotoUrl(Params.Url);
    });
};

var gotoAndroidMarketLandingPage = function () {
    dstLocation = dsAction.destination.dstandroidMarketLandPage;
    if (DEBUG) {
        alert(dstLocation)
    }

    renderTemplete('gotoLandingpageTemplate', 'gotoLandingpage', {
        'bgTop': LANDING_BG_TOP,
        'appName':Params.AppName, 
        'iconUrl':Params.IconUrl, 
        'downloadTitle':Params.Download_title, 
        'downloadMsg':Params.Download_msg, 
        'btnLandingpageText':ResGotoAppStoreDownload, 
        'elementType':'button',
        'dragdownIcon': '',
        'dragdownTip': '',
        'dragdown-display': 'hide',
    });

    dsAction.reportDSJSEvent(dsAction.actionJSDst, dstLocation);
    $('#btnGotoLandingPage').click(function () {
        dsAction.reportDSJSUserClickEvent(dsAction.actionJSUserClick, "gotoAndroidMarket", "yes");
        gotoAndroidMarket();
    });
};

var gotoAndroidCannotGoMarketLandingPage = function () {
    dstLocation = dsAction.destination.dstandroidCannotGoMarketLandPage;
    if (DEBUG) {
        alert(dstLocation);
    }

    renderTemplete('gotoLandingpageTemplate', 'gotoLandingpage', {
        'bgTop': LANDING_BG_TOP,
        'appName':Params.AppName, 
        'iconUrl':Params.IconUrl, 
        'downloadTitle':Params.Download_title, 
        'downloadMsg':Params.Download_msg, 
        'btnLandingpageText':ResPleaseOpenAppStore, 
        'elementType':'p',
        'dragdownIcon':  '',
        'dragdownTip': '',
        'dragdown-display': 'hide',
    });

    dsAction.reportDSJSEvent(dsAction.actionJSDst, dstLocation);
};

var gotoAndroidDownloadLandingPage = function () {
    dstLocation = dsAction.destination.dstandroidDirectDownloadLandPage;
    if (DEBUG) {
        alert(dstLocation);
    }

    renderTemplete('gotoLandingpageTemplate', 'gotoLandingpage', {
        'bgTop': LANDING_BG_TOP,
        'appName':Params.AppName, 
        'iconUrl':Params.IconUrl, 
        'downloadTitle':Params.Download_title, 
        'downloadMsg':Params.Download_msg, 
        'btnLandingpageText':ResDownloadAPK, 
        'elementType':'button',
        'dragdownIcon':  '',
        'dragdownTip': '',
        'dragdown-display': 'hide',
    });

    $('#btnGotoLandingPage').click(function () {
        dsAction.reportDSJSUserClickEvent(dsAction.actionJSUserClick, "gotoAndroidDirectDownload", "yes");
        gotoUrl(Params.Url);
    })
    dsAction.reportDSJSEvent(dsAction.actionJSDst, dstLocation);
};

var gotoUC = function (deeplinkurl) {
    dstLocation = dsAction.destination.dstucbrowser;
    var oriHtml = $("body").html();

    renderTemplete('allowMeDeeplinkTemplate', 'allowMeDeeplink', {
        'bgTopImgUrl': ResPathImg + 'allowme/allowme-bg-top.png',
        'bgBottomImgUrl': ResPathImg + 'allowme/allowme-bg-bottom.png',
        'resAllowMeMsg': ResAllowMeMsg,
        'resAllowMeTips': ResAllowMeTips,
        'resOr': ResOr,
        'resAllowMeThisTime': ResAllowMeThisTime,
        'resAllowMeAlways': ResAllowMeAlways,
    });

    var counter = 6;
    var timeoutCounterFunc = function () {
        //$("body").html(oriHtml);
        counter--;
        $('#textCountDown').html(counter)
        if (counter === 0) {
            $('#textCountDown').html("")
            iframeDeeplinkLaunch(deeplinkurl, 3000, function () {
                $("body").html(oriHtml);
                gotoAndroidNewInstall();
            });
        } else {
            setTimeout(timeoutCounterFunc, 1000);
        }
    };
    timeoutCounterFunc();
};

var gotoAndroidMarket = function () {
    env.windowChangeHistory();
    dstLocation = 'market://details?id=' + Params.Pkg;
    if (DEBUG) {
        alert(dstLocation)
    }
    dsAction.reportDSJSEvent(dsAction.actionJSDst, dstLocation);
    env.windowLocation(dstLocation);
};

//var gotoUrl = function () {
//    env.windowChangeHistory();
//    dstLocation = Url;
//    reportDSJSEvent(actionJSDst, Url);
//    env.windowLocation(Url);
//};

var gotoUrl = function (url) {
    env.windowChangeHistory();
    dstLocation = url;
    dsAction.reportDSJSEvent(dsAction.actionJSDst, url);
    env.windowLocation(url);
};

var gotoDivPlatformNotAvail = function (platform) {
    var resNotAvailable = '';
    if ('ios' == platform) {
        resNotAvailable = ResiOSNotAvailable;
    }else if ('android' == platform) {
        resNotAvailable = ResAndroidNotAvailable;
    }
    renderTemplete('platformNotAvailableTemplate', 'platformNotAvailable', {
        'bgImgUrl': ResPathImg + platform + '_not_avail/notavail.png',
        'resNotAvailable': resNotAvailable, 
    });
    
    dstLocation = dsAction.destination.dstplatformNA
        .replace(/{Platform}/g, platform);
    dsAction.reportDSJSEvent(dsAction.actionJSDst, dstLocation);
};

var deeplinkLaunch = function (deeplink, timeoutTime, timeoutCallback) {
    deeplinkLocation = deeplink;
    if (DEBUG) {
        alert(deeplinkLocation)
    }
    dsAction.reportDSJSEvent(dsAction.actionJSDeepLink, deeplink);
    env.windowLocation(deeplink)
    var timeout = setTimeout(function () {
        timeoutCallback();
    }, timeoutTime);
    clearTimeoutOnPageUnload(timeout);
};

var chiosDeeplinkLaunch = function (deeplink, timeoutCallback) {
    deeplinkLocation = deeplink;
    var w = null;
    try {
        dsAction.reportDSJSEvent(dsAction.actionJSDeepLink, deeplink);
        w = env.windowOpen(deeplink);
        if (DEBUG) {
            alert('pass');
        }
        env.windowChangeHistory()
    } catch (e) {
        if (DEBUG) {
            alert('exception');
        }
    }
    if (w) {
        env.windowClose();
    } else {
        timeoutCallback();
    }
};

var iframeDeeplinkLaunch = function (deeplink, timeoutTime, timeoutCallback) {
    var hiddenIFrame = document.createElement('iframe');
    hiddenIFrame.style.width = '1px';
    hiddenIFrame.style.height = '1px';
    hiddenIFrame.border = 'none';
    hiddenIFrame.style.display = 'none';
    hiddenIFrame.src = deeplink;
    document.body.appendChild(hiddenIFrame);
    deeplinkLocation = deeplink;
    dsAction.reportDSJSEvent(dsAction.actionJSDeepLink, deeplink);
    var timeout = setTimeout(function () {
        //                 document.getElementById("debug").innerHTML = "Start timer";
        timeoutCallback();
    }, timeoutTime);
    clearTimeoutOnPageUnload(timeout);
};

var isIosNotAvailable = function () {
    return Params.isIOS() && (Params.BundleID === undefined || Params.BundleID === "");
};

var isAndroidNotAvailable = function () {
    return Params.isAndroid() && (Params.Pkg === undefined || Params.Pkg === "");
};

var shouldGotoYYB = function () {
    return Params.YYB_url !== undefined && Params.YYB_url !== "" && (!(Params.isIOS()));
};

function start() {
    //            document.getElementById("debug").innerHTML =
    //                    "isChrome: " + isChrome() + ";   chromeVersion: " + Chrome_major
    //                    + "<br> isAndroid: " + isAndroid()
    //                    + "<br> isIOS: " + isIOS()
    //                    + "<br> isWechat: " + isWechat()
    //                    + "<br> isWeibo: " + isWeibo()
    //                    + "<br> isFirefox: " + isFirefox();
    if (Params.Download_msg === "") {
        Params.Download_msg = ResDefaultDownloadMsg;
    }
    if (Params.Download_btn_text !== "") {
        ResDownloadAPK = Params.Download_btn_text;
        ResGotoAppStoreDownload = Params.Download_btn_text;
    }
    if (isIosNotAvailable()) {
        gotoDivPlatformNotAvail('ios');
        return;
    } else if (isAndroidNotAvailable()) {
        gotoDivPlatformNotAvail('android');
        return;
    }
    if (Params.isWechat()) {
        //                document.getElementById("debug").innerHTML = "isWechat";
        if (DEBUG) {
            alert('isWeChat');
        }
        //for ios9, check AppInsStatus, if the app is installed or unclear, show wechat tips that shows "open with safari"
        if (Params.isIOS() && Params.Ios_major >= 9 && !Params.Force_download) {
            switch (parseInt(Params.AppInsStatus, 10)) {
                case CONST_APP_INS_STATUS.Installed:
                    gotoTip('ios', dsAction.destination.dstweixintipios, Params.UserConf_Bg_WechatIosTip_url);
                    return;
                case CONST_APP_INS_STATUS.NotInstall:
                    // when uninstall_url is set, redirect to uninstall_url when app is not installed
                    if (Params.Uninstall_url !== "") {
                        gotoUrl(Params.Uninstall_url);
                    } else {
                        gotoTencentProduct(dsAction.destination.dstweixintipios, dsAction.destination.dstweixintipandroid);
                    }
                    return;
                case CONST_APP_INS_STATUS.Unclear:
                    // when uninstall_url is set, redirect to uninstall_url when app is not installed
                    if (Params.Uninstall_url !== "") {
                        gotoUrl(Params.Uninstall_url);
                    } else {
                        /* In webapi, cannot set cookie right
                         * which need client to add: xhrFields: {withCredentials: true}
                        */
                        gotoTencentProduct(dsAction.destination.dstweixintipios, dsAction.destination.dstweixintipandroid);
                        //gotoTip('ios', dsAction.destination.dstweixintipios, Params.UserConf_Bg_WechatIosTip_url);
                    }
                    return;
                default:
                    gotoTip('ios', dsAction.destination.dstweixintipios, Params.UserConf_Bg_WechatIosTip_url);
                    return;
            }
        } else {
            gotoTencentProduct(dsAction.destination.dstweixintipios, dsAction.destination.dstweixintipandroid);
        }

    } else if (Params.isQQ()) {
        if (DEBUG) {
            alert("isQQ");
        }
        gotoTencentProduct(dsAction.destination.dstqqtipios, dsAction.destination.dstqqtipandroid);
    } else if (Params.isWeibo()) {
        if (DEBUG) {
            alert("isWeibo");
        }
        //                document.getElementById("debug").innerHTML = "isWeibo";
        if (Params.isIOS()) {
            if (DEBUG) {
                alert("isIOS");
            }
            gotoTip('ios', dsAction.destination.dstweibotipios, Params.UserConf_Bg_WechatIosTip_url);
            // loadCanvas("div_weixintip", IOS);
        } else if (Params.isAndroid()) {
            if (DEBUG) {
                alert("isAndroid");
            }
            gotoTip('android', dsAction.destination.dstweibotipandroid, Params.UserConf_Bg_WechatAndroidTip_url);
        }
    } else if (Params.isFacebook()) {
        if (DEBUG) {
            alert("isFacebook");
        }
        //                document.getElementById("debug").innerHTML = "isWeibo";
        if (Params.isIOS()) {
            if (DEBUG) {
                alert("isIOS");
            }
            gotoTip('ios_right_down', dsAction.destination.dstfacebooktipios, "");
            // loadCanvas("div_weixintip", IOS);
        } else if (Params.isAndroid()) {
            if (DEBUG) {
                alert("isAndroid");
            }
            gotoTip('android', dsAction.destination.dstfacebooktipandroid, "");
        }
    } else if (Params.isTwitter()) {
        if (DEBUG) {
            alert("isTwitter");
        }
        //                document.getElementById("debug").innerHTML = "isWeibo";
        if (Params.isIOS()) {
            if (DEBUG) {
                alert("isIOS");
            }
            gotoTip('ios_right_down', dsAction.destination.dsttwittertipios, "");
            // loadCanvas("div_weixintip", IOS);
        } else if (Params.isAndroid()) {
            if (DEBUG) {
                alert("isAndroid");
            }
            gotoTip('android', dsAction.destination.dsttwittertipandroid, "");
        }
    } else if (Params.isIOS()) {
        //                document.getElementById("debug").innerHTML = "isIOS";
        if (DEBUG) {
            alert("isIOS");
        }
        var deeplinkurl = Params.Scheme + '://';
        if (Params.Match_id && Params.Match_id.length > 0) {
            deeplinkurl += "?click_id=" + Params.Match_id;
        }
        if (DEBUG) {
            alert(deeplinkurl);
        }
        if (Params.Ios_major < 9) {
            if (DEBUG) {
                alert("IOS Major below 9:" + Params.Ios_major);
            }
            iframeDeeplinkLaunch(deeplinkurl, 2000, function () {
                gotoUrl(Params.Url);
            });
        } else {
            if (DEBUG) {
                alert("IOS Major upper 9:" + Params.Ios_major);
            }
            if (Params.isChrome()) {
                if (DEBUG) {
                    alert("isChrome");
                }
                chiosDeeplinkLaunch(deeplinkurl, function () {
                    gotoUrl(Params.Url);
                });
            } else if (Params.isUniversallink()) {
                //If it is universal link, it mean it is in two situation
                //1.The App is not installed
                //2.The App is installed, but the system prefer the web page
                //So we need to show the landing page to cover both situation
                if (DEBUG) {
                    alert("isUniversallink = true");
                }

                if (env.cookieEnabled()) {
                    if (DEBUG) {
                        alert("cookie Enabled; AppInsStatus:" + Params.AppInsStatus);
                    }
                    switch (parseInt(Params.AppInsStatus, 10)) {
                        case CONST_APP_INS_STATUS.Installed:
                            gotoIOSLandingPage();
                            return;
                        case CONST_APP_INS_STATUS.NotInstall:
                            gotoUrl(Params.Url);
                            return;
                        case CONST_APP_INS_STATUS.Unclear:
                            gotoIOSLandingPage();
                            return;
                        default:
                            gotoIOSLandingPage();
                            return;
                    }
                } else {
                    if (DEBUG) {
                        alert("cookie Not Enabled");
                    }
                    gotoIOSLandingPage();
                }
            } else {
                if (DEBUG) {
                    alert("is safari");
                }
                deeplinkLaunch(deeplinkurl, 500, function () {
                    gotoUrl(Params.Url);
                });
            }

        }

    } else if (Params.isAndroid()) {
        //                document.getElementById("debug").innerHTML = "isAndroid";
        if (DEBUG) {
            alert("isAndroid");
        }
        deeplinkurl = Params.Scheme + '://' + Params.Host;
        if (Params.Match_id && Params.Match_id.length > 0) {
            deeplinkurl += "?click_id=" + Params.Match_id;
        }
        if (DEBUG) {
            alert(deeplinkurl)
        }
        if (Params.isCannotDeeplink()) {
            iframeDeeplinkLaunch(deeplinkurl, 2000, function () {
                gotoCannotDeeplink();
            })
        } else if (Params.isQQBrowser()) {
            if (DEBUG) {
                alert("QQ browser");
            }
            if (Params.isYYBEnableAndroid()) {
                if (DEBUG) {
                    alert(Params.YYB_url)
                }
                gotoUrl(Params.YYB_url);
            } else {
                gotoCannotDeeplink();
            }
        } else if (Params.isUC()) {
            if (DEBUG) {
                alert("UC browser");
            }
            gotoUC(deeplinkurl);
        } else if (Params.isChrome() && Params.Chrome_major >= 25 && !Params.isForceUseScheme()) {
            if (DEBUG) {
                alert("Chrome_major:" + Params.Chrome_major);
            }
            //Extract scheme from deeplink
            //                    document.getElementById("debug").innerHTML = "isAndroid 3";
            var intent = Params.Host;
            if (Params.Match_id && Params.Match_id.length > 0) {
                intent += "?click_id=" + Params.Match_id;
            }
            var pkg = Params.Pkg;
            // When deeplinking on chrome 35+, there is inherent app store fallback logic built into the browser (likely a bug in Chrome)
            // Workaround for Bug https://code.google.com/p/chromium/issues/detail?id=459711&thanks=459711&ts=1424288965.
            var workaroundlink = "intent://" + intent + 
                                 "#Intent;scheme=" + toLowerCase(Params.Scheme) +
                                 ";package=" + pkg + 
                                 ";S.browser_fallback_url=" + Params.Url +
                                 ";end";
            // di.innerHTML = workaroundlink;
            deeplinkLaunch(workaroundlink, 2000, function () {
                gotoAndroidNewInstall();
            });
        } else {
            if (DEBUG) {
                alert("default browser");
            }
            //                     document.getElementById("debug").innerHTML = "isAndroid default; " + Url;
            iframeDeeplinkLaunch(deeplinkurl, 2000, function () {
                gotoAndroidNewInstall();
            });
        }
    }
}

function toLowerCase(str) {
    if (typeof str === "string") {
        return str.toLowerCase();
    } else if ((typeof str === "undefined") || (str === null)) {
        return '';
    } else {
        str = '' + str;
        return str.toLowerCase();
    }
}

window.onload = function () {
    start();
};
