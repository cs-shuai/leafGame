//如果参数过多，建议通过 object 方式传入
function login(name) {
    console.log("code:" + code)

    initLogin("", name)

    var json = {
        Login:{
            UserName: User.userName,
            UserId: User.userId,
        }
    }
    SendJson(json)
}

function LoginCallBack(data) {
    initLogin(data.GameMessage.Data.UserId, data.GameMessage.Data.UserName)
    $("#name").html(User.userName)
    $.toast("登录成功", 500);
}

// 初始化登录会员信息
function initLogin(userId, userName) {
    if (userId == "") {
        userId = $.cookie('userId' + code);
    } else {
        $.cookie('userId' + code, userId, { expires: 1 });
    }

    if (userName == "") {
        userName = $.cookie('userName' + code);
    } else {
        $.cookie('userName' + code, userName, { expires: 1 });

    }

    if (userId != undefined) {
        User.userId = userId
    }
    if (userName != undefined) {
        User.userName = userName
    }
    console.log("初始化 User")
    console.log(User)
}