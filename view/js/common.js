// 长链地址
var wsUri = "ws://172.16.1.74:8889/";
// 用户code
var code = getUrlParam("code");

// 设置用户code
if (code == null) {
    code = _getRandomString(6);
    window.location.href = UrlUpdateParams(window.location.href, "code", code);
}

function Message(data) {
    console.log(data);
    var div = document.getElementById("div");
    if (data.GameMessage.Msg !== "") {
        div.innerHTML = div.innerHTML + "<br>" + data.GameMessage.Msg
        div.scrollTop = div.scrollHeight;
    }

    console.log("回调类型: " + data.GameMessage.Type)
    console.log("回调数据: " + data)
    switch (data.GameMessage.Type) {
        case "Login": LoginCallBack(data); break;
        case "Room": rRoomCallBack(data); break;
        case "Prepare": rPrepareCallBack(data); break;
        case "Start": rStartCallBack(data); break;
        case "StartGame": rStartGameCallBack(data); break;
        case "Vote": rVoteCallBack(data); break;
        case "Over": rOverCallBack(data); break;
        case "VoteSuccess": rVoteSuccessCallBack(data); break;
    }
}

// 获取长度为len的随机字符串
function _getRandomString(len) {
    len = len || 32;
    var $chars = 'ABCDEFGHJKMNPQRSTWXYZabcdefhijkmnprstwxyz2345678'; // 默认去掉了容易混淆的字符oOLl,9gq,Vv,Uu,I1
    var maxPos = $chars.length;
    var pwd = '';
    for (i = 0; i < len; i++) {
        pwd += $chars.charAt(Math.floor(Math.random() * maxPos));
    }
    return pwd;
}

// 获取url中的参数
function getUrlParam(name) {
    var reg = new RegExp("(^|&)" + name + "=([^&]*)(&|$)"); //构造一个含有目标参数的正则表达式对象
    var r = window.location.search.substr(1).match(reg);  //匹配目标参数
    if (r != null) return unescape(r[2]); return null; //返回参数值
}


function UrlUpdateParams (url, name, value) {
    var r = url;
    if (r != null && r != 'undefined' && r != "") {
        value = encodeURIComponent(value);
        var reg = new RegExp("(^|)" + name + "=([^&]*)(|$)");
        var tmp = name + "=" + value;
        if (url.match(reg) != null) {
            r = url.replace(eval(reg), tmp);
        }
        else {
            if (url.match("[\?]")) {
                r = url + "&" + tmp;
            } else {
                r = url + "?" + tmp;
            }
        }
    }
    return r;
}


var MAX = 99, MIN = 1;
$('.weui-count__decrease').click(function (e) {
    var $input = $(e.currentTarget).parent().find('.weui-count__number');
    var number = parseInt($input.val() || "0") - 1
    if (number < MIN) number = MIN;
    $input.val(number)
})
$('.weui-count__increase').click(function (e) {
    var $input = $(e.currentTarget).parent().find('.weui-count__number');
    var number = parseInt($input.val() || "0") + 1
    if (number > MAX) number = MAX;
    $input.val(number)
})

var User = {
    "userName" : "",
    "userId": "",
    "roomId": ""
};
var output;
var websocketConnect;
window.addEventListener("load", init, false);
function init() {
    output = document.getElementById("output");
    NewWebSocket();
}

function NewWebSocket() {
    console.log(wsUri)
    $.showLoading("连接中...")
    websocketConnect = new WebSocket(wsUri);
    websocketConnect.onopen = function(evt) {
        onOpen(evt)
        $.hideLoading()
    };
    websocketConnect.onclose = function(evt) {
        onClose(evt)
    };
    websocketConnect.onmessage = function(evt) {
        onMessage(evt)
    };
    websocketConnect.onerror = function(evt) {
        $.hideLoading()
        onError(evt)
    };

}

// 连接
function onOpen(evt) {
    console.log("CONNECTED");
    // login();
    if ($.cookie('userName' + code) == undefined || $.cookie('userName' + code) == "") {
        $.prompt({
            title: '用户名',
            text: '起一个好听的用户名',
            input: User.userName,
            empty: false, // 是否允许为空
            onOK: function (input) {
                //点击确认
                console.log(input)
                login(input)
            },
            onCancel: function () {
                //点击取消
            }
        });
    } else {
        login($.cookie('userName' + code))
    }
}

// 关闭
function onClose(evt) {
    console.log("onClose");
    console.log("DISCONNECTED");
    $.confirm({
        title: '连接失败',
        text: '是否重新连接',
        onOK: function () {
            reconnect()
        }
    });
}

// 发送
function onMessage(evt) {
    console.log(evt);
    var filrReader = new FileReader();
    filrReader.onload = function() {
        var arrayBuffer = this.result;
        var decoder = new TextDecoder('utf-8')
        var json = JSON.parse(decoder.decode(new DataView(arrayBuffer)));
        Message(json)
    };
    filrReader.readAsArrayBuffer(evt.data);
};

// 错误
function onError(evt) {
    console.log("onError");
    $.toast(evt.data, "cancel");
}

// 发送json
function SendJson(json) {
    doSend(JSON.stringify(json));
}

// 发送
function doSend(message) {
    console.log(message);
    websocketConnect.send(message);
}

// 重连
function reconnect() {
    isPrepare = false
    $.showLoading("连接中...")
    var timeReconnect = setTimeout(function(){
        NewWebSocket()
    } , 1000);
}

// 退出房间
$(document).on("click", ".out", function () {
    var json = {
        RoomOut:{
            RoomId: User.roomId,
            UserId: User.userId,
        }
    }
    SendJson(json)
    $("#page1").load("./login.html");
})



// 退出登录
$(document).on("click", ".return", function () {
    var a = $.cookie('userName' + code, "")
    console.log(a)
    var b = $.cookie('userId' + code, "")
    console.log(b)

    initLogin("", "")
    var json = {
        UserOut:{
            UserId: User.userId,
        }
    }
    SendJson(json)
    $("#page1").load("./login.html");
})

