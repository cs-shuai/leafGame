// 词组
var keywrod = ""
// 投票的会员编码
var voteUserId = ""
// 投票时间函数Id
var voteTime = ""

function initGame(newKeywrod, createUserId) {
    keywrod = newKeywrod
    voteUserId = ""
    voteTime = ""
    if (createUserId == User.userId) {
        isHomeowner = true
    }

    // 房主
    $("#page1").load("./game.html");
}

// 开始游戏回调
function rStartGameCallBack(data) {
    initGame(data.GameMessage.Data.Keyword, data.GameMessage.Data.CreateUserId)
}

// 开始投票回调
function rVoteCallBack(data) {
    var Round = data.GameMessage.Data.Round;
    var text = "<span class='vote-countdown'> "+ data.GameMessage.Data.VoteTime +"</span> <br>";
    $.each(data.GameMessage.Data.SurvivalUserList, function (i, v) {
        text += "<a href=\"javascript:;\" userId=\"" + i + "\" class=\"weui-btn weui-btn_mini weui-btn_default btn_tp\">"+v.Name+"</a>";
    })
    $.alert({
        title: '第' + Round + '回合投票',
        text: text,
        onOK: function () {
        //点击确认
            SendVote()
        }
    });

    // 重置时间函数
    if (voteTime != "") {
        clearInterval(voteTime)
    }

    // 设置投票倒计时
    voteTime = setInterval(function(){
        var t = $('.vote-countdown').html();
        t--
        $('.vote-countdown').html(t);
        if (t <= 0) {
            SendVote()
            clearInterval(t)
        }
    }, 1000)
}

// 结束本机游戏回调
function rOverCallBack(data) {
    console.log("游戏结束")
    $.alert({
        text: data.GameMessage.Msg,
        onOK: function () {
            if (data.GameMessage.Data.Stage === "Over") {
                initRoom(User.roomId, "")
            }
        }
    });
}

// 投票成功回调
function rVoteSuccessCallBack(data) {
    console.log("投票成功")
    voteUserId = ""

    $.toast(data.GameMessage.Msg, 500);
}

$(document).on("click", ".btn_tp", function () {
    voteUserId = $(this).attr("userId")
    $(".btn_tp").each(function () {
        $(this).removeClass("weui-btn_primary")
        $(this).addClass("weui-btn_default")
    })
    $(this).addClass("weui-btn_primary")
    $(this).removeClass("weui-btn_default")
})

// 投票
function Vote() {
    var roomId = $("#RoomId").val();
    var json = {
        Game: {
            Stage: "Vote",
            RoomId: roomId,
        }
    }
    SendJson(json)
}

function SendVote() {
    // 清空投票倒计时 和 关闭弹窗
    clearInterval(voteTime)
    $.closeModal();

    var json = {
        Vote: {
            UserId: User.userId,
            RoomId: User.roomId,
            VotePlayerNumber: voteUserId,
        }
    }
    SendJson(json)
}

// 查看词组
function lookKeyword() {
    $.alert(keywrod);
}